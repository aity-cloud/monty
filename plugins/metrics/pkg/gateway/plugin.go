package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/metrics/collector"
	managementext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/plugins/apis/metrics"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/types"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Plugin struct {
	system.UnimplementedSystemPluginClient
	collector.CollectorServer
	ctx context.Context

	logger *slog.Logger

	*pluginContextData

	streamServices     []util.ServicePackInterface
	managementServices []util.ServicePackInterface
}

func NewPlugin(ctx context.Context, scheme meta.Scheme) *Plugin {
	reader := sdkmetric.NewManualReader(
		sdkmetric.WithAggregationSelector(types.AggregationSelector),
	)
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	metrics := types.NewMetrics(mp)

	collector := collector.NewCollectorServer(reader)
	p := &Plugin{
		ctx:             ctx,
		CollectorServer: collector,
		logger:          logger.NewPluginLogger().WithGroup("metrics"),
	}

	var pctx types.PluginContext
	pctx, p.pluginContextData = newPluginContext(ctx, metrics, p.logger)
	go func() {
		driver, err := initClusterDriver(pctx)
		if err != nil {
			p.logger.With(logger.Err(err)).Error("failed to initialize cluster driver")
			return
		}
		p.clusterDriver.C() <- driver
	}()

	types.Services.Range(func(name string, builder driverutil.Builder[types.Service]) {
		lg := p.logger.With("service", name)
		svc, err := builder(ctx, driverutil.NewOption("context", pctx))
		if err != nil {
			lg.With(logger.Err(err)).Error("failed to initialize service")
			return
		}
		lg.Info("loading service")

		ctx, ca := context.WithTimeout(ctx, 5*time.Second)
		cancelWarn := context.AfterFunc(ctx, func() {
			lg.Warn("service activation is taking longer than expected")
		})

		if svc, ok := svc.(types.PluginService); ok {
			svc.AddToScheme(scheme)
		}
		if svc, ok := svc.(types.ManagementService); ok {
			p.managementServices = append(p.managementServices, svc.ManagementServices()...)
		}
		if svc, ok := svc.(types.StreamService); ok {
			p.streamServices = append(p.streamServices, svc.StreamServices()...)
		}
		go func() {
			defer ca()
			defer cancelWarn()
			if err := svc.Activate(); err != nil {
				lg.With(logger.Err(err)).Error("failed to activate service")
				return
			}
			lg.Info("activated service")
		}()
	})

	return p
}

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme(meta.WithMode(meta.ModeGateway))
	p := NewPlugin(ctx, scheme)
	scheme.Add(system.SystemPluginID, system.NewPlugin(p))
	streamMetricReader := sdkmetric.NewManualReader()
	p.CollectorServer.AppendReader(streamMetricReader)
	scheme.Add(streamext.StreamAPIExtensionPluginID, streamext.NewGatewayPlugin(p,
		streamext.WithMetrics(streamext.GatewayStreamMetricsConfig{
			Reader:          streamMetricReader,
			LabelsForStream: p.labelsForStreamMetrics,
		})),
	)
	scheme.Add(managementext.ManagementAPIExtensionPluginID, managementext.NewPlugin(p))
	scheme.Add(metrics.MetricsPluginID, metrics.NewPlugin(p))
	return scheme
}

func initClusterDriver(ctx types.PluginContext) (drivers.ClusterDriver, error) {
	// TODO(config)
	allDrivers := drivers.ClusterDrivers.List()
	var builder driverutil.Builder[drivers.ClusterDriver]
	if len(allDrivers) > 0 {
		builder, _ = drivers.ClusterDrivers.Get(allDrivers[0])
	} else {
		return nil, fmt.Errorf("no cluster drivers found")
	}
	driver, err := builder(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cluster driver %q: %w", allDrivers[0], err)
	}
	ctx.Logger().With(
		"driver", allDrivers[0],
	).Info("initialized cluster driver")
	return driver, nil
}
