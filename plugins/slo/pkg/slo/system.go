package slo

import (
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/slo/shared"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexadmin"
	sloapi "github.com/aity-cloud/monty/plugins/slo/apis/slo"
)

func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.mgmtClient.Set(client)
	<-p.ctx.Done()
}

func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	p.storage.Set(StorageAPIs{
		SLOs:     system.NewKVStoreClient[*sloapi.SLOData](client),
		Services: system.NewKVStoreClient[*sloapi.Service](client),
		Metrics:  system.NewKVStoreClient[*sloapi.Metric](client),
	})
	<-p.ctx.Done()
}

func (p *Plugin) UseAPIExtensions(intf system.ExtensionClientInterface) {
	cc, err := intf.GetClientConn(p.ctx, "CortexAdmin", "AlertEndpoints")
	if err != nil {
		if p.ctx.Err() != nil {
			// Plugin is shutting down, don't exit
			return
		}
	}
	adminClient := cortexadmin.NewCortexAdminClient(cc)
	alertingEndpointClient := alertingv1.NewAlertEndpointsClient(cc)

	p.adminClient.Set(adminClient)
	p.alertEndpointClient.Set(alertingEndpointClient)
	RegisterDatasource(
		shared.MonitoringDatasource,
		NewSLOMonitoringStore(p, p.logger),
		NewMonitoringServiceBackend(p, p.logger),
	)
	<-p.ctx.Done()
}
