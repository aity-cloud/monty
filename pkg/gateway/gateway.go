package gateway

import (
	"context"
	"crypto"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"slices"

	"log/slog"

	bootstrapv1 "github.com/aity-cloud/monty/pkg/apis/bootstrap/v1"
	bootstrapv2 "github.com/aity-cloud/monty/pkg/apis/bootstrap/v2"
	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	streamv1 "github.com/aity-cloud/monty/pkg/apis/stream/v1"
	"github.com/aity-cloud/monty/pkg/auth/challenges"
	"github.com/aity-cloud/monty/pkg/auth/cluster"
	authv1 "github.com/aity-cloud/monty/pkg/auth/cluster/v1"
	authv2 "github.com/aity-cloud/monty/pkg/auth/cluster/v2"
	"github.com/aity-cloud/monty/pkg/auth/session"
	"github.com/aity-cloud/monty/pkg/bootstrap"
	"github.com/aity-cloud/monty/pkg/capabilities"
	"github.com/aity-cloud/monty/pkg/config"
	cfgmeta "github.com/aity-cloud/monty/pkg/config/meta"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/health"
	"github.com/aity-cloud/monty/pkg/keyring"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/hooks"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/plugins/types"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/update"
	k8sserver "github.com/aity-cloud/monty/pkg/update/kubernetes/server"
	patchserver "github.com/aity-cloud/monty/pkg/update/patch/server"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"golang.org/x/mod/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Gateway struct {
	GatewayOptions
	config        *config.GatewayConfig
	tlsConfig     *tls.Config
	logger        *slog.Logger
	httpServer    *GatewayHTTPServer
	grpcServer    *GatewayGRPCServer
	statusQuerier health.HealthStatusQuerier

	storageBackend  storage.Backend
	capBackendStore capabilities.BackendStore
}

type GatewayOptions struct {
	lifecycler          config.Lifecycler
	extraUpdateHandlers []update.UpdateTypeHandler
}

type GatewayOption func(*GatewayOptions)

func (o *GatewayOptions) apply(opts ...GatewayOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithLifecycler(lc config.Lifecycler) GatewayOption {
	return func(o *GatewayOptions) {
		o.lifecycler = lc
	}
}

func WithExtraUpdateHandlers(handlers ...update.UpdateTypeHandler) GatewayOption {
	return func(o *GatewayOptions) {
		o.extraUpdateHandlers = append(o.extraUpdateHandlers, handlers...)
	}
}

func NewGateway(ctx context.Context, conf *config.GatewayConfig, pl plugins.LoaderInterface, opts ...GatewayOption) *Gateway {
	options := GatewayOptions{
		lifecycler: config.NewUnavailableLifecycler(cfgmeta.ObjectList{conf}),
	}
	options.apply(opts...)

	lg := logger.New().WithGroup("gateway")
	conf.Spec.SetDefaults()

	storageBackend, err := machinery.ConfigureStorageBackend(ctx, &conf.Spec.Storage)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to configure storage backend")
		panic("failed to configure storage backend")
	}

	// configure the server-side installer template with the external hostname
	// and grpc port which agents will connect to
	_, port, err := net.SplitHostPort(conf.Spec.GRPCListenAddress)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to parse listen address")
		panic("failed to parse listen address")
	}
	capBackendStore := capabilities.NewBackendStore(capabilities.ServerInstallerTemplateSpec{
		Address: conf.Spec.Hostname,
		Port:    fmt.Sprint(port),
	}, lg)

	// add capabilities from plugins
	pl.Hook(hooks.OnLoadM(func(p types.CapabilityBackendPlugin, md meta.PluginMeta) {
		info, err := p.Info(ctx, &emptypb.Empty{})
		if err != nil {
			lg.With(
				"plugin", md.Module,
			).Error("failed to get capability info")
			return
		}
		if err := capBackendStore.Add(info.Name, p); err != nil {
			lg.With(
				"plugin", md.Module,
				logger.Err(err),
			).Error("failed to add capability backend")
		}
		lg.With(
			"plugin", md.Module,
			"capability", info.Name,
		).Info("added capability backend")
	}))

	// serve system plugin kv stores
	pl.Hook(hooks.OnLoadM(func(p types.SystemPlugin, md meta.PluginMeta) {
		ns := md.Module
		if err := module.CheckPath(ns); err != nil {
			lg.With(
				"namespace", ns,
				logger.Err(err),
			).Warn("system plugin module name is invalid")
			return
		}
		store := storageBackend.KeyValueStore(ns)
		go p.ServeKeyValueStore(store)
	}))

	// serve caching provider for plugin RPCs
	pl.Hook(hooks.OnLoadM(func(p types.SystemPlugin, md meta.PluginMeta) {
		ns := md.Module
		if err := module.CheckPath(ns); err != nil {
			lg.With(
				"namespace", ns,
				logger.Err(err),
			).Warn("system plugin module name is invalid")
			return
		}
		go p.ServeCachingProvider()
	}))

	// set up http server
	httpServer := NewHTTPServer(ctx, &conf.Spec, lg, pl)

	// Set up cluster auth
	ephemeralKeys, err := machinery.LoadEphemeralKeys(afero.Afero{
		Fs: afero.NewOsFs(),
	}, conf.Spec.Keyring.EphemeralKeyDirs...)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to load ephemeral keys")
		panic("failed to load ephemeral keys")
	}

	v1Verifier := challenges.NewKeyringVerifier(storageBackend, authv1.DomainString, lg.WithGroup("authv1"))
	v2Verifier := challenges.NewKeyringVerifier(storageBackend, authv2.DomainString, lg.WithGroup("authv2"))

	sessionAttrChallenge, err := session.NewServerChallenge(
		keyring.New(lo.ToAnySlice(ephemeralKeys)...))
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to configure authentication")
		panic("failed to configure authentication")
	}

	v1Challenge := authv1.NewServerChallenge(streamv1.Stream_Connect_FullMethodName, v1Verifier, lg.WithGroup("authv1"))
	v2Challenge := authv2.NewServerChallenge(v2Verifier, lg.WithGroup("authv2"))

	clusterAuth := cluster.StreamServerInterceptor(challenges.Chained(
		challenges.If(authv2.ShouldEnableIncoming).Then(v2Challenge).Else(v1Challenge),
		challenges.If(session.ShouldEnableIncoming).Then(sessionAttrChallenge),
	))

	//set up update server
	updateServer := update.NewUpdateServer(lg)

	// set up plugin sync server
	binarySyncServer, err := patchserver.NewFilesystemPluginSyncServer(conf.Spec.Plugins, lg,
		patchserver.WithPluginSyncFilters(func(pm meta.PluginMeta) bool {
			if pm.ExtendedMetadata != nil {
				// only sync plugins that have the agent mode set
				return slices.Contains(pm.ExtendedMetadata.ModeList.Modes, meta.ModeAgent)
			}
			return true // default to syncing all plugins
		}),
	)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to create plugin sync server")
		panic("failed to create plugin sync server")
	}

	if err := binarySyncServer.RunGarbageCollection(ctx, storageBackend); err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to run garbage collection")
	}

	updateServer.RegisterUpdateHandler(binarySyncServer.Strategy(), binarySyncServer)

	kubernetesSyncServer, err := k8sserver.NewKubernetesSyncServer(conf.Spec.AgentUpgrades.Kubernetes, lg)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to create kubernetes agent sync server")
		panic("failed to create kubernetes agent sync server")
	}
	updateServer.RegisterUpdateHandler(kubernetesSyncServer.Strategy(), kubernetesSyncServer)

	for _, handler := range options.extraUpdateHandlers {
		updateServer.RegisterUpdateHandler(handler.Strategy(), handler)
	}

	httpServer.metricsRegisterer.MustRegister(updateServer.Collectors()...)

	// set up grpc server
	tlsConfig, pkey, err := grpcTLSConfig(&conf.Spec)
	if err != nil {
		lg.With(
			logger.Err(err),
		).Error("failed to load TLS config")
		panic("failed to load TLS config")
	}

	rateLimitOpts := []RatelimiterOption{}
	if conf.Spec.RateLimit != nil {
		rateLimitOpts = append(rateLimitOpts, WithRate(conf.Spec.RateLimit.Rate))
		rateLimitOpts = append(rateLimitOpts, WithBurst(conf.Spec.RateLimit.Burst))
	}
	grpcServer := NewGRPCServer(&conf.Spec, lg,
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainStreamInterceptor(
			NewRateLimiterInterceptor(lg, rateLimitOpts...).StreamServerInterceptor(),
			clusterAuth,
			updateServer.StreamServerInterceptor(),
			NewLastKnownDetailsApplier(storageBackend),
		),
	)

	// set up stream server
	listener := health.NewListener()
	monitor := health.NewMonitor(health.WithLogger(lg.WithGroup("monitor")))
	delegate := NewDelegateServer(storageBackend, lg)
	// set up agent connection handlers
	agentHandler := MultiConnectionHandler(listener, delegate)

	go monitor.Run(ctx, listener)
	streamSvc := NewStreamServer(agentHandler, storageBackend, httpServer.metricsRegisterer, lg)

	controlv1.RegisterHealthListenerServer(streamSvc, listener)
	streamv1.RegisterDelegateServer(streamSvc.InternalServiceRegistrar(), delegate)
	streamv1.RegisterStreamServer(grpcServer, streamSvc)
	corev1.RegisterPingerServer(streamSvc, &pinger{})
	controlv1.RegisterUpdateSyncServer(grpcServer, updateServer)

	pl.Hook(hooks.OnLoadMC(streamSvc.OnPluginLoad))

	// set up bootstrap server
	bootstrapServerV1 := bootstrap.NewServer(storageBackend, pkey, capBackendStore)
	bootstrapServerV2 := bootstrap.NewServerV2(storageBackend, pkey)
	bootstrapv1.RegisterBootstrapServer(grpcServer, bootstrapServerV1)
	bootstrapv2.RegisterBootstrapServer(grpcServer, bootstrapServerV2)

	g := &Gateway{
		GatewayOptions:  options,
		tlsConfig:       tlsConfig,
		config:          conf,
		logger:          lg,
		storageBackend:  storageBackend,
		capBackendStore: capBackendStore,
		httpServer:      httpServer,
		grpcServer:      grpcServer,
		statusQuerier:   monitor,
	}

	return g
}

type keyValueStoreServer interface {
	ServeKeyValueStore(store storage.KeyValueStore)
}

func (g *Gateway) ListenAndServe(ctx context.Context) error {
	lg := g.logger
	ctx, ca := context.WithCancel(ctx)

	// start http server
	e1 := lo.Async(func() error {
		err := g.httpServer.ListenAndServe(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				lg.Info("http server stopped")
			} else {
				lg.With(logger.Err(err)).Warn("http server exited with error")
			}
		}
		return err
	})

	// start grpc server
	e2 := lo.Async(func() error {
		err := g.grpcServer.ListenAndServe(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				lg.Info("grpc server stopped")
			} else {
				lg.With(logger.Err(err)).Warn("grpc server exited with error")
			}
		}
		return err
	})

	return util.WaitAll(ctx, ca, e1, e2)
}

// Implements management.CoreDataSource
func (g *Gateway) StorageBackend() storage.Backend {
	return g.storageBackend
}

// Implements management.CoreDataSource
func (g *Gateway) TLSConfig() *tls.Config {
	return g.httpServer.tlsConfig
}

// Implements management.CapabilitiesDataSource
func (g *Gateway) CapabilitiesStore() capabilities.BackendStore {
	return g.capBackendStore
}

// Implements management.HealthStatusDataSource
func (g *Gateway) GetClusterHealthStatus(ref *corev1.Reference) (*corev1.HealthStatus, error) {
	hs := g.statusQuerier.GetHealthStatus(ref.Id)
	if hs.Health == nil && hs.Status == nil {
		return nil, status.Error(codes.NotFound, "no health or status has been reported for this cluster yet")
	}
	return hs, nil
}

// Implements management.HealthStatusDataSource
func (g *Gateway) WatchClusterHealthStatus(ctx context.Context) <-chan *corev1.ClusterHealthStatus {
	return g.statusQuerier.WatchHealthStatus(ctx)
}

func (g *Gateway) MustRegisterCollector(collector prometheus.Collector) {
	g.httpServer.metricsRegisterer.MustRegister(collector)
}

func grpcTLSConfig(cfg *v1beta1.GatewayConfigSpec) (*tls.Config, crypto.Signer, error) {
	servingCertBundle, caPool, err := util.LoadServingCertBundle(cfg.Certs)
	if err != nil {
		return nil, nil, err
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		RootCAs:      caPool,
		Certificates: []tls.Certificate{*servingCertBundle},
		ClientAuth:   tls.NoClientCert,
	}, servingCertBundle.PrivateKey.(crypto.Signer), nil
}

func httpTLSConfig(cfg *v1beta1.GatewayConfigSpec) (*tls.Config, crypto.Signer, error) {
	servingCertBundle, caPool, err := util.LoadServingCertBundle(cfg.Certs)
	if err != nil {
		return nil, nil, err
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      caPool,
		ClientCAs:    caPool,
		Certificates: []tls.Certificate{*servingCertBundle},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, servingCertBundle.PrivateKey.(crypto.Signer), nil
}
