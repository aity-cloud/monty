package slo

import (
	"context"

	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexadmin"
	"github.com/aity-cloud/monty/plugins/slo/apis/slo"
	"log/slog"

	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/logger"
	managementext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/future"
)

type Plugin struct {
	slo.UnsafeSLOServer
	system.UnimplementedSystemPluginClient

	ctx    context.Context
	logger *slog.Logger

	storage             future.Future[StorageAPIs]
	mgmtClient          future.Future[managementv1.ManagementClient]
	adminClient         future.Future[cortexadmin.CortexAdminClient]
	alertEndpointClient future.Future[alertingv1.AlertEndpointsClient]
}

type StorageAPIs struct {
	SLOs     storage.KeyValueStoreT[*slo.SLOData]
	Services storage.KeyValueStoreT[*slo.Service]
	Metrics  storage.KeyValueStoreT[*slo.Metric]
}

func NewPlugin(ctx context.Context) *Plugin {
	return &Plugin{
		ctx:                 ctx,
		logger:              logger.NewPluginLogger().WithGroup("slo"),
		storage:             future.New[StorageAPIs](),
		mgmtClient:          future.New[managementv1.ManagementClient](),
		adminClient:         future.New[cortexadmin.CortexAdminClient](),
		alertEndpointClient: future.New[alertingv1.AlertEndpointsClient](),
	}
}

var _ slo.SLOServer = (*Plugin)(nil)

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme()
	p := NewPlugin(ctx)
	scheme.Add(system.SystemPluginID, system.NewPlugin(p))
	scheme.Add(managementext.ManagementAPIExtensionPluginID,
		managementext.NewPlugin(util.PackService(&slo.SLO_ServiceDesc, p)))
	return scheme
}
