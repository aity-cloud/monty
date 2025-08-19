package gateway

import (
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"

	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
)

// UseManagementAPI implements system.SystemPluginServer.
func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.managementClient.C() <- client
	<-p.ctx.Done()
}

// UseManagementAPI implements system.SystemPluginServer.
func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	p.keyValueStoreClient.C() <- client
	<-p.ctx.Done()
}

// UseManagementAPI implements system.SystemPluginServer.
func (p *Plugin) UseAPIExtensions(intf system.ExtensionClientInterface) {
	p.extensionClient.C() <- intf
	<-p.ctx.Done()
}

// UseConfigAPI implements system.SystemPluginServer.
func (p *Plugin) UseConfigAPI(client configv1.GatewayConfigClient) {
	p.gatewayConfigClient.C() <- client
	config, err := client.GetConfiguration(p.ctx, &driverutil.GetRequest{})
	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("failed to get gateway configuration")
		return
	}
	backend, err := machinery.ConfigureStorageBackendV1(p.ctx, config.Storage)
	if err != nil {
		p.logger.With(
			"err", err,
		).Error("failed to configure storage backend")
		return
	}
	p.storageBackend.C() <- backend
	<-p.ctx.Done()
}
