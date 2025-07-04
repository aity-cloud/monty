package alarms

import (
	"context"
	"fmt"
	"sync"

	"log/slog"

	"github.com/aity-cloud/monty/pkg/alerting/server"
	alertingSync "github.com/aity-cloud/monty/pkg/alerting/server/sync"
	"github.com/aity-cloud/monty/pkg/alerting/storage/spec"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/future"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/alerting/metrics"
	notifications "github.com/aity-cloud/monty/plugins/alerting/pkg/alerting/notifications/v1"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/rules"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexadmin"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var _ rules.RuleSyncServer = (*AlarmServerComponent)(nil)

type AlarmServerComponent struct {
	alertingv1.UnsafeAlertConditionsServer
	rules.UnsafeRuleSyncServer

	util.Initializer
	ctx context.Context

	mu sync.RWMutex
	server.Config

	logger *slog.Logger

	runner        *Runner
	notifications *notifications.NotificationServerComponent

	conditionStorage future.Future[spec.ConditionStorage]
	incidentStorage  future.Future[spec.IncidentStorage]
	stateStorage     future.Future[spec.StateStorage]

	routerStorage future.Future[spec.RouterStorage]

	js future.Future[nats.JetStreamContext]

	mgmtClient      future.Future[managementv1.ManagementClient]
	adminClient     future.Future[cortexadmin.CortexAdminClient]
	cortexOpsClient future.Future[cortexops.CortexOpsClient]

	metricExporter *AlarmMetricsExporter
}

func NewAlarmServerComponent(
	ctx context.Context,
	logger *slog.Logger,
	notifications *notifications.NotificationServerComponent,
) *AlarmServerComponent {
	comp := &AlarmServerComponent{
		ctx:              ctx,
		logger:           logger,
		runner:           NewRunner(),
		notifications:    notifications,
		conditionStorage: future.New[spec.ConditionStorage](),
		incidentStorage:  future.New[spec.IncidentStorage](),
		stateStorage:     future.New[spec.StateStorage](),
		routerStorage:    future.New[spec.RouterStorage](),
		js:               future.New[nats.JetStreamContext](),
		mgmtClient:       future.New[managementv1.ManagementClient](),
		adminClient:      future.New[cortexadmin.CortexAdminClient](),
		cortexOpsClient:  future.New[cortexops.CortexOpsClient](),
		metricExporter:   NewAlarmMetricsExporter(),
	}
	return comp
}

type AlarmServerConfiguration struct {
	spec.ConditionStorage
	spec.IncidentStorage
	spec.StateStorage
	spec.RouterStorage
	Js              nats.JetStreamContext
	MgmtClient      managementv1.ManagementClient
	AdminClient     cortexadmin.CortexAdminClient
	CortexOpsClient cortexops.CortexOpsClient
}

var _ server.ServerComponent = (*AlarmServerComponent)(nil)

func (a *AlarmServerComponent) Name() string {
	return "alarm"
}

func (a *AlarmServerComponent) Status() server.Status {
	return server.Status{
		Running: a.Initialized(),
	}
}

func (a *AlarmServerComponent) Ready() bool {
	return a.Initialized()
}

func (a *AlarmServerComponent) Healthy() bool {
	return a.Initialized()
}

func (a *AlarmServerComponent) SetConfig(conf server.Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Config = conf
}

func (a *AlarmServerComponent) Sync(ctx context.Context, syncInfo alertingSync.SyncInfo) error {
	conditionStorage, err := a.conditionStorage.GetContext(a.ctx)
	if err != nil {
		return err
	}
	groupIds, err := conditionStorage.ListGroups(a.ctx)
	if err != nil {
		return err
	}
	conds := []*alertingv1.AlertCondition{}
	for _, groupId := range groupIds {
		groupConds, err := conditionStorage.Group(groupId).List(a.ctx)
		if err != nil {
			return err
		}
		conds = append(conds, groupConds...)
	}
	eg := &util.MultiErrGroup{}
	a.logger.Debug(fmt.Sprintf("syncing (%v) %d conditions", syncInfo.ShouldSync, len(conds)))
	for _, cond := range conds {
		cond := cond
		if syncInfo.ShouldSync {
			eg.Go(func() error {
				activationAttrs := []attribute.KeyValue{
					{
						Key:   "cluster_id",
						Value: attribute.StringValue(cond.GetClusterId().Id),
					},
					{
						Key:   "name",
						Value: attribute.StringValue(cond.GetName()),
					},
					{
						Key:   "datasource",
						Value: attribute.StringValue(cond.DatasourceName()),
					},
					{
						Key:   "id",
						Value: attribute.StringValue(cond.GetId()),
					},
				}
				_, err := a.applyAlarm(ctx, cond, cond.Id, syncInfo)
				if err != nil {
					activationAttrs = append(activationAttrs, attribute.KeyValue{
						Key:   "status",
						Value: attribute.StringValue("failed"),
					})
				} else {
					activationAttrs = append(activationAttrs, attribute.KeyValue{
						Key:   "status",
						Value: attribute.StringValue("success"),
					})
				}
				a.metricExporter.activationStatusCounter.Add(ctx, 1, metric.WithAttributes(
					activationAttrs...,
				))
				return err
			})
		} else {
			eg.Go(func() error {
				return a.teardownCondition(ctx, cond, cond.Id, false)
			})
		}
	}
	eg.Wait()
	if len(eg.Errors()) > 0 {
		a.logger.Error(fmt.Sprintf("successfully synced (%d/%d) conditions", len(conds)-len(eg.Errors()), len(conds)))
	}
	if err := eg.Error(); err != nil {
		return err
	}
	return nil
}

func (a *AlarmServerComponent) Initialize(conf AlarmServerConfiguration) {
	a.InitOnce(func() {
		a.mgmtClient.Set(conf.MgmtClient)
		a.adminClient.Set(conf.AdminClient)
		a.js.Set(conf.Js)
		a.cortexOpsClient.Set(conf.CortexOpsClient)
		a.conditionStorage.Set(conf.ConditionStorage)
		a.incidentStorage.Set(conf.IncidentStorage)
		a.stateStorage.Set(conf.StateStorage)
		a.routerStorage.Set(conf.RouterStorage)
	})
}

type AlarmMetricsExporter struct {
	activationStatusCounter metric.Int64Counter
}

func NewAlarmMetricsExporter() *AlarmMetricsExporter {
	return &AlarmMetricsExporter{
		activationStatusCounter: metrics.ActivationStatusCounter,
	}
}
