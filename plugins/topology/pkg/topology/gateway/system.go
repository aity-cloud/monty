package gateway

import (
	"os"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/config/adapt"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/task"
	natsutil "github.com/aity-cloud/monty/pkg/util/nats"
	"google.golang.org/protobuf/proto"

	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
)

func (p *Plugin) UseManagementAPI(client managementv1.ManagementClient) {
	p.mgmtClient.Set(client)

	<-p.ctx.Done()
}

func (p *Plugin) UseConfigAPI(client configv1.GatewayConfigClient) {
	config, err := client.GetConfiguration(p.ctx, &driverutil.GetRequest{})
	if err != nil {
		p.logger.With(logger.Err(err)).Error("failed to get gateway configuration")
		return
	}
	p.gatewayConfig.Set(&v1beta1.GatewayConfig{
		Spec: *adapt.V1BetaConfigOf[*v1beta1.GatewayConfigSpec](config),
	})
	backend, err := machinery.ConfigureStorageBackendV1(p.ctx, config.Storage)
	if err != nil {
		p.logger.With(
			"err", err,
		).Error("failed to configure storage backend")
		os.Exit(1)
	}
	p.storageBackend.Set(backend)
	p.configureTopologyManagement()
	<-p.ctx.Done()
}

func (p *Plugin) UseKeyValueStore(client system.KeyValueStoreClient) {
	// set other futures before trying to acquire NATS connection
	ctrl, err := task.NewController(
		p.ctx,
		"topology.uninstall",
		system.NewKVStoreClient[*corev1.TaskStatus](client),
		&p.uninstallRunner)

	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("failed to create uninstall task controller")
	}
	p.uninstallController.Set(ctrl)

	p.storage.Set(ConfigStorageAPIs{
		Placeholder: system.NewKVStoreClient[proto.Message](client),
	})
	cfg := p.gatewayConfig.Get().Spec.Storage.JetStream
	natsURL := os.Getenv("NATS_SERVER_URL")
	natsSeedPath := os.Getenv("NKEY_SEED_FILENAME")
	if cfg == nil {
		cfg = &v1beta1.JetStreamStorageSpec{}
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = natsURL
	}
	if cfg.NkeySeedPath == "" {
		cfg.NkeySeedPath = natsSeedPath
	}
	nc, err := natsutil.AcquireNATSConnection(p.ctx, cfg)
	if err != nil {
		p.logger.With(
			logger.Err(err),
		).Error("fatal :  failed to acquire NATS connection")
		os.Exit(1)
	}
	p.nc.Set(nc)
	<-p.ctx.Done()
}
