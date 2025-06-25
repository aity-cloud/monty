package gateway

import (
	"context"
	"os"

	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/aity-cloud/monty/plugins/metrics/apis/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/task"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/backend"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/cortex"

	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
	"github.com/aity-cloud/monty/pkg/storage/kvutil"
)

func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.mgmtClient.Set(client)
	cfg, err := client.GetConfig(context.Background(), &emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("failed to get config")
		os.Exit(1)
	}
	objectList, err := machinery.LoadDocuments(cfg.Documents)
	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("failed to load config")
		os.Exit(1)
	}
	machinery.LoadAuthProviders(p.ctx, objectList)
	objectList.Visit(func(config *v1beta1.GatewayConfig) {
		backend, err := machinery.ConfigureStorageBackend(p.ctx, &config.Spec.Storage)
		if err != nil {
			p.logger.With(
				logger.Err(err),
			).Error("failed to configure storage backend")
			os.Exit(1)
		}
		p.storageBackend.Set(backend)
		p.config.Set(config)
		tlsConfig := p.loadCortexCerts()
		p.cortexTlsConfig.Set(tlsConfig)
		clientset, err := cortex.NewClientSet(p.ctx, &config.Spec.Cortex, tlsConfig)
		if err != nil {
			p.logger.With(
				logger.Err(err),
			).Error("failed to configure cortex clientset")
			os.Exit(1)
		}
		p.cortexClientSet.Set(clientset)
	})

	p.authMw.Set(machinery.LoadAuthProviders(p.ctx, objectList))
	<-p.ctx.Done()
}

func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	ctrl, err := task.NewController(p.ctx, "uninstall", system.NewKVStoreClient[*corev1.TaskStatus](client), &p.uninstallRunner)
	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("failed to create task controller")
		os.Exit(1)
	}
	p.uninstallController.Set(ctrl)

	p.backendKvClients.Set(&backend.KVClients{
		DefaultClusterConfigurationSpec: kvutil.WithKey(system.NewKVStoreClient[*cortexops.CapabilityBackendConfigSpec](client), "/config/cluster/default"),
		DefaultCapabilitySpec:           kvutil.WithKey(system.NewKVStoreClient[*node.MetricsCapabilitySpec](client), "/config/capability/default"),
		NodeCapabilitySpecs:             kvutil.WithPrefix(system.NewKVStoreClient[*node.MetricsCapabilitySpec](client), "/config/capability/nodes/"),
	})
	<-p.ctx.Done()
}
