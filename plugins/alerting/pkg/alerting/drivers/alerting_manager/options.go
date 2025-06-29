package alerting_manager

import (
	"crypto/tls"
	"log/slog"

	alertingClient "github.com/aity-cloud/monty/pkg/alerting/client"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AlertingDriverOptions struct {
	Logger             *slog.Logger                         `option:"logger"`
	K8sClient          client.Client                        `option:"k8sClient"`
	GatewayRef         types.NamespacedName                 `option:"gatewayRef"`
	ConfigKey          string                               `option:"configKey"`
	InternalRoutingKey string                               `option:"internalRoutingKey"`
	AlertingOptions    *shared.AlertingClusterOptions       `option:"alertingOptions"`
	Subscribers        []chan alertingClient.AlertingClient `option:"subscribers"`
	TlsConfig          *tls.Config                          `option:"tlsConfig"`
}
