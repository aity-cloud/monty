package drivers

import (
	"context"

	"github.com/aity-cloud/monty/pkg/alerting/drivers/config"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/plugins/alerting/apis/alertops"
)

type ClusterDriver interface {
	alertops.AlertingAdminServer
	// ShouldDisableNode is called during node sync for nodes which otherwise
	// have this capability enabled. If this function returns an error, the
	// node will be set to disabled instead, and the error will be logged.
	ShouldDisableNode(*corev1.Reference) error
	GetDefaultReceiver() *config.WebhookConfig
	GetAlertingServiceConfig() AlertingServiceConfig
}

var Drivers = driverutil.NewCache[ClusterDriver]()

type NoopClusterDriver struct {
	alertops.UnimplementedAlertingAdminServer
}

// GetAlertingServiceConfig implements ClusterDriver.
func (d *NoopClusterDriver) GetAlertingServiceConfig() AlertingServiceConfig {
	return AlertingServiceConfig{}
}

func (d *NoopClusterDriver) ShouldDisableNode(*corev1.Reference) error {
	// the noop driver will never forcefully disable a node
	return nil
}

func (d *NoopClusterDriver) GetRuntimeOptions() shared.AlertingClusterOptions {
	return shared.AlertingClusterOptions{}
}

func (d *NoopClusterDriver) GetDefaultReceiver() *config.WebhookConfig {
	return &config.WebhookConfig{}
}

type AlertingServiceConfig struct {
	Certs MTLSConfig
}

type MTLSConfig struct {
	ServerCA   string
	ClientCA   string
	ClientCert string
	ClientKey  string
}

func init() {
	Drivers.Register("noop", func(ctx context.Context, opts ...driverutil.Option) (ClusterDriver, error) {
		return &NoopClusterDriver{}, nil
	})
}
