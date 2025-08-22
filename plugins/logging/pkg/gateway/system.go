package gateway

import (
	"os"

	montycorev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/task"

	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
)

func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.mgmtApi.Set(client)
	<-p.ctx.Done()
}

func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	p.kv.Set(client)
	ctrl, err := task.NewController(p.ctx, "uninstall", system.NewKVStoreClient[*montycorev1.TaskStatus](client), &UninstallTaskRunner{
		storageNamespace:  p.storageNamespace,
		opensearchManager: p.opensearchManager,
		backendDriver:     p.clusterDriver,
		storageBackend:    p.storageBackend,
		logger:            p.logger.WithGroup("uninstaller"),
	})
	if err != nil {
		p.logger.With(
			"err", err,
		).Error("failed to create task controller")
		os.Exit(1)
	}

	p.uninstallController.Set(ctrl)
	<-p.ctx.Done()
}

func (p *Plugin) UseConfigAPI(client configv1.GatewayConfigClient) {
	config, err := client.GetConfiguration(p.ctx, &driverutil.GetRequest{})
	if err != nil {
		p.logger.With(logger.Err(err)).Error("failed to get gateway configuration")
		return
	}
	backend, err := machinery.ConfigureStorageBackendV1(p.ctx, config.Storage)
	if err != nil {
		p.logger.With(
			"err", err,
		).Error("failed to configure storage backend")
		os.Exit(1)
	}
	p.storageBackend.Set(backend)

	<-p.ctx.Done()
}
