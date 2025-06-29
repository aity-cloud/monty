package metrics

import (
	"context"

	"github.com/aity-cloud/monty/pkg/metrics/collector"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	MetricsPluginID = "monty.Metrics"
	ServiceID       = "collector.RemoteCollector"
)

type metricsPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	rcServer collector.RemoteCollectorServer
}

func NewPlugin(srv collector.RemoteCollectorServer) plugin.Plugin {
	return &metricsPlugin{
		rcServer: srv,
	}
}

func (p *metricsPlugin) GRPCServer(
	_ *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	collector.RegisterRemoteCollectorServer(s, p.rcServer)
	return nil
}

func (p *metricsPlugin) GRPCClient(
	ctx context.Context,
	_ *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	if err := plugins.CheckAvailability(ctx, c, ServiceID); err != nil {
		return nil, err
	}
	client := collector.NewRemoteCollectorClient(c)
	return collector.NewRemoteProducer(client), nil
}

func init() {
	plugins.GatewayScheme.Add(MetricsPluginID, NewPlugin(nil))
}
