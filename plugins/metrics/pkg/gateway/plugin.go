package gateway

import (
	"context"
	"crypto/tls"

	"log/slog"

	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/auth"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/metrics/collector"
	httpext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/http"
	managementext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/plugins/apis/capability"
	"github.com/aity-cloud/monty/pkg/plugins/apis/metrics"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/task"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/future"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexadmin"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/aity-cloud/monty/plugins/metrics/apis/node"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remoteread"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/backend"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/cortex"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Plugin struct {
	// capabilityv1.UnsafeBackendServer
	system.UnimplementedSystemPluginClient
	collector.CollectorServer

	ctx    context.Context
	logger *slog.Logger

	cortexAdmin       cortex.CortexAdminServer
	cortexHttp        cortex.HttpApiServer
	cortexRemoteWrite cortex.RemoteWriteForwarder
	metrics           backend.MetricsBackend
	uninstallRunner   cortex.UninstallTaskRunner

	config              future.Future[*v1beta1.GatewayConfig]
	authMw              future.Future[map[string]auth.Middleware]
	mgmtClient          future.Future[managementv1.ManagementClient]
	storageBackend      future.Future[storage.Backend]
	cortexTlsConfig     future.Future[*tls.Config]
	cortexClientSet     future.Future[cortex.ClientSet]
	uninstallController future.Future[*task.Controller]
	clusterDriver       future.Future[drivers.ClusterDriver]
	delegate            future.Future[streamext.StreamDelegate[backend.MetricsAgentClientSet]]
	backendKvClients    future.Future[*backend.KVClients]
}

func NewPlugin(ctx context.Context) *Plugin {
	cortexReader := metric.NewManualReader(
		metric.WithAggregationSelector(cortex.CortexAggregationSelector),
	)
	mp := metric.NewMeterProvider(
		metric.WithReader(cortexReader),
	)
	cortex.RegisterMeterProvider(mp)

	collector := collector.NewCollectorServer(cortexReader)
	p := &Plugin{
		CollectorServer: collector,
		ctx:             ctx,
		logger:          logger.NewPluginLogger().WithGroup("metrics"),

		config:              future.New[*v1beta1.GatewayConfig](),
		authMw:              future.New[map[string]auth.Middleware](),
		mgmtClient:          future.New[managementv1.ManagementClient](),
		storageBackend:      future.New[storage.Backend](),
		cortexTlsConfig:     future.New[*tls.Config](),
		cortexClientSet:     future.New[cortex.ClientSet](),
		uninstallController: future.New[*task.Controller](),
		clusterDriver:       future.New[drivers.ClusterDriver](),
		delegate:            future.New[streamext.StreamDelegate[backend.MetricsAgentClientSet]](),
		backendKvClients:    future.New[*backend.KVClients](),
	}
	p.metrics.OpsBackend = &backend.OpsServiceBackend{MetricsBackend: &p.metrics}
	p.metrics.NodeBackend = &backend.NodeServiceBackend{MetricsBackend: &p.metrics}

	future.Wait2(p.cortexClientSet, p.config,
		func(cortexClientSet cortex.ClientSet, config *v1beta1.GatewayConfig) {
			p.cortexAdmin.Initialize(cortex.CortexAdminServerConfig{
				CortexClientSet: cortexClientSet,
				Config:          &config.Spec,
				Logger:          p.logger.WithGroup("cortex-admin"),
			})
		})

	future.Wait2(p.cortexClientSet, p.config,
		func(cortexClientSet cortex.ClientSet, config *v1beta1.GatewayConfig) {
			p.cortexRemoteWrite.Initialize(cortex.RemoteWriteForwarderConfig{
				CortexClientSet: cortexClientSet,
				Config:          &config.Spec,
				Logger:          p.logger.WithGroup("cortex-rw"),
			})
		})

	future.Wait3(p.cortexClientSet, p.config, p.storageBackend,
		func(cortexClientSet cortex.ClientSet, config *v1beta1.GatewayConfig, storageBackend storage.Backend) {
			p.uninstallRunner.Initialize(cortex.UninstallTaskRunnerConfig{
				CortexClientSet: cortexClientSet,
				Config:          &config.Spec,
				StorageBackend:  storageBackend,
			})
		})
	future.Wait2(p.config, p.backendKvClients, func(
		config *v1beta1.GatewayConfig,
		backendKvClients *backend.KVClients,
	) {
		driverName := config.Spec.Cortex.Management.ClusterDriver
		if driverName == "" {
			p.logger.Warn("no cluster driver configured")
		}
		builder, ok := drivers.ClusterDrivers.Get(driverName)
		if !ok {
			p.logger.With(
				"driver", driverName,
			).Error("unknown cluster driver, using fallback noop driver")
			builder, ok = drivers.ClusterDrivers.Get("noop")
			if !ok {
				panic("bug: noop cluster driver not found")
			}
		}
		driver, err := builder(p.ctx,
			driverutil.NewOption("defaultConfigStore", backendKvClients.DefaultClusterConfigurationSpec),
		)
		if err != nil {
			p.logger.With(
				"driver", driverName,
				logger.Err(err),
			).Error("failed to initialize cluster driver")
			panic("failed to initialize cluster driver")
			return
		}
		p.logger.With(
			"driver", driverName,
		).Info("initialized cluster driver")
		p.clusterDriver.Set(driver)
	})
	future.Wait6(p.storageBackend, p.mgmtClient, p.uninstallController, p.clusterDriver, p.delegate, p.backendKvClients,
		func(
			storageBackend storage.Backend,
			mgmtClient managementv1.ManagementClient,
			uninstallController *task.Controller,
			clusterDriver drivers.ClusterDriver,
			delegate streamext.StreamDelegate[backend.MetricsAgentClientSet],
			backendKvClients *backend.KVClients,
		) {
			p.metrics.Initialize(backend.MetricsBackendConfig{
				Logger:              p.logger.WithGroup("metrics-backend"),
				StorageBackend:      storageBackend,
				MgmtClient:          mgmtClient,
				UninstallController: uninstallController,
				ClusterDriver:       clusterDriver,
				Delegate:            delegate,
				KV:                  backendKvClients,
			})
		})

	future.Wait6(p.mgmtClient, p.cortexClientSet, p.config, p.cortexTlsConfig, p.storageBackend, p.authMw,
		func(
			mgmtApi managementv1.ManagementClient,
			cortexClientSet cortex.ClientSet,
			config *v1beta1.GatewayConfig,
			tlsConfig *tls.Config,
			storageBackend storage.Backend,
			authMiddlewares map[string]auth.Middleware,
		) {
			p.cortexHttp.Initialize(cortex.HttpApiServerConfig{
				PluginContext:    p.ctx,
				ManagementClient: mgmtApi,
				CortexClientSet:  cortexClientSet,
				Config:           &config.Spec,
				CortexTLSConfig:  tlsConfig,
				Logger:           p.logger.WithGroup("cortex-http"),
				StorageBackend:   storageBackend,
				AuthMiddlewares:  authMiddlewares,
			})
		})
	return p
}

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme(meta.WithMode(meta.ModeGateway))
	p := NewPlugin(ctx)
	scheme.Add(system.SystemPluginID, system.NewPlugin(p))
	scheme.Add(httpext.HTTPAPIExtensionPluginID, httpext.NewPlugin(&p.cortexHttp))
	streamMetricReader := metric.NewManualReader()
	p.CollectorServer.AppendReader(streamMetricReader)
	scheme.Add(streamext.StreamAPIExtensionPluginID, streamext.NewGatewayPlugin(p,
		streamext.WithMetrics(streamext.GatewayStreamMetricsConfig{
			Reader:          streamMetricReader,
			LabelsForStream: p.labelsForStreamMetrics,
		})),
	)
	scheme.Add(managementext.ManagementAPIExtensionPluginID, managementext.NewPlugin(
		util.PackService(&cortexadmin.CortexAdmin_ServiceDesc, &p.cortexAdmin),
		util.PackService(&cortexops.CortexOps_ServiceDesc, p.metrics.OpsBackend),
		util.PackService(&remoteread.RemoteReadGateway_ServiceDesc, &p.metrics),
		util.PackService(&node.NodeConfiguration_ServiceDesc, p.metrics.NodeBackend),
	))
	scheme.Add(capability.CapabilityBackendPluginID, capability.NewPlugin(&p.metrics))
	scheme.Add(metrics.MetricsPluginID, metrics.NewPlugin(p))
	return scheme
}
