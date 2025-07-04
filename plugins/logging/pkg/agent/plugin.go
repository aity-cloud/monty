package agent

import (
	"context"
	"log/slog"

	healthpkg "github.com/aity-cloud/monty/pkg/health"
	"github.com/aity-cloud/monty/pkg/logger"
	httpext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/http"
	"github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/plugins/apis/capability"
	"github.com/aity-cloud/monty/pkg/plugins/apis/health"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/logging/pkg/agent/drivers"
	"github.com/aity-cloud/monty/plugins/logging/pkg/otel"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type Plugin struct {
	ctx           context.Context
	logger        *slog.Logger
	node          *LoggingNode
	otelForwarder *otel.Forwarder
}

func NewPlugin(ctx context.Context) *Plugin {
	lg := logger.NewPluginLogger().WithGroup("logging")

	ct := healthpkg.NewDefaultConditionTracker(lg)

	p := &Plugin{
		ctx:    ctx,
		logger: lg,
		node:   NewLoggingNode(ct, lg),
		otelForwarder: otel.NewForwarder(
			otel.NewLogsForwarder(
				otel.WithLogger(lg.WithGroup("otel-logs-forwarder")),
				otel.WithPrivileged(true)),
			otel.NewTraceForwarder(
				otel.WithLogger(lg.WithGroup("otel-trace-forwarder")),
				otel.WithPrivileged(true)),
		),
	}

	for _, d := range drivers.NodeDrivers.List() {
		driverBuilder, _ := drivers.NodeDrivers.Get(d)
		driver, err := driverBuilder(ctx,
			driverutil.NewOption("logger", lg),
		)
		if err != nil {
			lg.With(
				"driver", d,
				logger.Err(err),
			).Error("failed to initialize logging node driver")
			continue
		}

		p.node.AddConfigListener(drivers.NewListenerFunc(ctx, driver.ConfigureNode))
	}
	return p
}

var (
	_ collogspb.LogsServiceServer   = (*otel.LogsForwarder)(nil)
	_ coltracepb.TraceServiceServer = (*otel.TraceForwarder)(nil)
)

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme(meta.WithMode(meta.ModeAgent))
	p := NewPlugin(ctx)
	scheme.Add(capability.CapabilityBackendPluginID, capability.NewAgentPlugin(p.node))
	scheme.Add(health.HealthPluginID, health.NewPlugin(p.node))
	scheme.Add(stream.StreamAPIExtensionPluginID, stream.NewAgentPlugin(p))
	scheme.Add(httpext.HTTPAPIExtensionPluginID, httpext.NewPlugin(p.otelForwarder))
	return scheme
}
