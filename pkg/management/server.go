package management

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"log/slog"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/caching"
	"github.com/aity-cloud/monty/pkg/capabilities"
	"github.com/aity-cloud/monty/pkg/config"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/pkp"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions"
	"github.com/aity-cloud/monty/pkg/plugins/hooks"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/plugins/types"
	"github.com/aity-cloud/monty/pkg/rbac"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jhump/protoreflect/desc"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	channelzservice "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CoreDataSource provides a way to obtain data which the management
// server needs to serve its core API
type CoreDataSource interface {
	StorageBackend() storage.Backend
	TLSConfig() *tls.Config
}

// CapabilitiesDataSource provides a way to obtain data which the management
// server needs to serve capabilities-related endpoints
type CapabilitiesDataSource interface {
	CapabilitiesStore() capabilities.BackendStore
}

type HealthStatusDataSource interface {
	GetClusterHealthStatus(ref *corev1.Reference) (*corev1.HealthStatus, error)
	WatchClusterHealthStatus(ctx context.Context) <-chan *corev1.ClusterHealthStatus
}

type apiExtension struct {
	client      apiextensions.ManagementAPIExtensionClient
	clientConn  *grpc.ClientConn
	serviceDesc *desc.ServiceDescriptor
	httpRules   []*managementv1.HTTPRuleDescriptor
}

type Server struct {
	managementv1.UnsafeManagementServer
	managementServerOptions
	config            *v1beta1.ManagementSpec
	logger            *slog.Logger
	rbacProvider      rbac.Provider
	coreDataSource    CoreDataSource
	grpcServer        *grpc.Server
	dashboardSettings *DashboardSettingsManager

	apiExtMu      sync.RWMutex
	apiExtensions []apiExtension
}

var _ managementv1.ManagementServer = (*Server)(nil)

type managementServerOptions struct {
	lifecycler             config.Lifecycler
	capabilitiesDataSource CapabilitiesDataSource
	healthStatusDataSource HealthStatusDataSource
}

type ManagementServerOption func(*managementServerOptions)

func (o *managementServerOptions) apply(opts ...ManagementServerOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithLifecycler(lc config.Lifecycler) ManagementServerOption {
	return func(o *managementServerOptions) {
		o.lifecycler = lc
	}
}

func WithCapabilitiesDataSource(src CapabilitiesDataSource) ManagementServerOption {
	return func(o *managementServerOptions) {
		o.capabilitiesDataSource = src
	}
}

func WithHealthStatusDataSource(src HealthStatusDataSource) ManagementServerOption {
	return func(o *managementServerOptions) {
		o.healthStatusDataSource = src
	}
}

func NewServer(
	ctx context.Context,
	conf *v1beta1.ManagementSpec,
	cds CoreDataSource,
	pluginLoader plugins.LoaderInterface,
	opts ...ManagementServerOption,
) *Server {
	lg := logger.New().WithGroup("mgmt")
	options := managementServerOptions{}
	options.apply(opts...)

	m := &Server{
		managementServerOptions: options,
		config:                  conf,
		logger:                  lg,
		coreDataSource:          cds,
		rbacProvider:            storage.NewRBACProvider(cds.StorageBackend()),
		dashboardSettings: &DashboardSettingsManager{
			kv:     cds.StorageBackend().KeyValueStore("dashboard"),
			logger: lg,
		},
	}

	director := m.configureApiExtensionDirector(ctx, pluginLoader)
	m.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnknownServiceHandler(unknownServiceHandler(director)),
		grpc.ChainStreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.ChainUnaryInterceptor(
			caching.NewClientGrpcTtlCacher().UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor()),
	)
	managementv1.RegisterManagementServer(m.grpcServer, m)
	channelzservice.RegisterChannelzServiceToServer(m.grpcServer)

	pluginLoader.Hook(hooks.OnLoadM(func(sp types.SystemPlugin, md meta.PluginMeta) {
		go sp.ServeManagementAPI(m)
		go func() {
			if err := sp.ServeAPIExtensions(m.config.GRPCListenAddress); err != nil {
				lg.With(
					"plugin", md.Module,
					logger.Err(err),
				).Error("failed to serve plugin API extensions")
			}
		}()
	}))

	return m
}

type managementApiServer interface {
	ServeManagementAPI(managementv1.ManagementServer)
}

func (m *Server) ListenAndServe(ctx context.Context) error {
	ctx, ca := context.WithCancel(ctx)

	e1 := lo.Async(func() error {
		err := m.listenAndServeGrpc(ctx)
		if err != nil {
			return fmt.Errorf("management grpc server exited with error: %w", err)
		}
		return nil
	})

	e2 := lo.Async(func() error {
		err := m.listenAndServeHttp(ctx)
		if err != nil {
			return fmt.Errorf("management http server exited with error: %w", err)
		}
		return nil
	})

	return util.WaitAll(ctx, ca, e1, e2)
}

func (m *Server) listenAndServeGrpc(ctx context.Context) error {
	if m.config.GRPCListenAddress == "" {
		return errors.New("GRPCListenAddress not configured")
	}
	lg := m.logger
	listener, err := util.NewProtocolListener(m.config.GRPCListenAddress)
	if err != nil {
		return err
	}
	lg.With(
		"address", listener.Addr().String(),
	).Info("management gRPC server starting")

	errC := lo.Async(func() error {
		return m.grpcServer.Serve(listener)
	})
	select {
	case <-ctx.Done():
		m.grpcServer.Stop()
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (m *Server) listenAndServeHttp(ctx context.Context) error {
	if m.config.GRPCListenAddress == "" {
		return errors.New("GRPCListenAddress not configured")
	}
	lg := m.logger
	lg.With(
		"address", m.config.HTTPListenAddress,
	).Info("management HTTP server starting")
	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/json", &LegacyJsonMarshaler{}),
		runtime.WithMarshalerOption("application/octet-stream", &DynamicV1Marshaler{}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &DynamicV1Marshaler{}),
	)

	m.configureManagementHttpApi(ctx, gwmux)
	m.configureHttpApiExtensions(gwmux)
	mux.Handle("/", gwmux)
	server := &http.Server{
		Addr:    m.config.HTTPListenAddress,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	errC := lo.Async(func() error {
		return server.ListenAndServe()
	})
	select {
	case <-ctx.Done():
		server.Close()
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (m *Server) CertsInfo(_ context.Context, _ *emptypb.Empty) (*managementv1.CertsInfoResponse, error) {
	resp := &managementv1.CertsInfoResponse{
		Chain: []*corev1.CertInfo{},
	}
	for _, tlsCert := range m.coreDataSource.TLSConfig().Certificates[:1] {
		for _, der := range tlsCert.Certificate {
			cert, err := x509.ParseCertificate(der)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			resp.Chain = append(resp.Chain, &corev1.CertInfo{
				Issuer:      cert.Issuer.String(),
				Subject:     cert.Subject.String(),
				IsCA:        cert.IsCA,
				NotBefore:   cert.NotBefore.Format(time.RFC3339),
				NotAfter:    cert.NotAfter.Format(time.RFC3339),
				Fingerprint: pkp.NewSha256(cert).Encode(),
				Raw:         cert.Raw,
			})
		}
	}
	return resp, nil
}

func (m *Server) ListCapabilities(ctx context.Context, in *emptypb.Empty) (*managementv1.CapabilityList, error) {
	if m.capabilitiesDataSource == nil {
		return nil, status.Error(codes.Unavailable, "capability backend store not configured")
	}

	clusters, err := m.ListClusters(ctx, &managementv1.ListClustersRequest{})
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int32)
	for _, cluster := range clusters.Items {
		for _, cap := range cluster.GetCapabilities() {
			counts[cap.Name]++
		}
	}

	names := m.capabilitiesDataSource.CapabilitiesStore().List()
	var items []*managementv1.CapabilityInfo
	for _, name := range names {
		capability, err := m.capabilitiesDataSource.CapabilitiesStore().Get(name)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		details, err := capability.Info(ctx, in)
		if err != nil {
			m.logger.With(
				logger.Err(err),
				"capability", name,
			).Error("failed to fetch capability details")
			continue
		}
		items = append(items, &managementv1.CapabilityInfo{
			Details:   details,
			NodeCount: counts[name],
		})
	}

	return &managementv1.CapabilityList{
		Items: items,
	}, nil
}

func (m *Server) CapabilityInstaller(
	_ context.Context,
	req *managementv1.CapabilityInstallerRequest,
) (*managementv1.CapabilityInstallerResponse, error) {
	if m.capabilitiesDataSource == nil {
		return nil, status.Error(codes.Unavailable, "capability backend store not configured")
	}

	cmd, err := m.capabilitiesDataSource.
		CapabilitiesStore().
		RenderInstaller(req.Name, capabilities.UserInstallerTemplateSpec{
			Token: req.Token,
			Pin:   req.Pin,
		})
	if err != nil {
		return nil, err
	}
	return &managementv1.CapabilityInstallerResponse{
		Command: cmd,
	}, nil
}

func (m *Server) Server() *grpc.Server {
	return m.grpcServer
}
