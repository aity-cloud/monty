package proxy

import (
	"context"

	proxyv1 "github.com/aity-cloud/monty/pkg/apis/proxy/v1"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	ProxyPluginID = "monty.Proxy"
	ServiceID     = "proxy.RegisterProxy"
)

type proxyBackendPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	proxySrv proxyv1.RegisterProxyServer
}

var _ plugin.GRPCPlugin = (*proxyBackendPlugin)(nil)
var _ plugin.Plugin = (*proxyBackendPlugin)(nil)

func (p *proxyBackendPlugin) GRPCServer(
	_ *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	proxyv1.RegisterRegisterProxyServer(s, p.proxySrv)
	return nil
}

func (p *proxyBackendPlugin) GRPCClient(
	ctx context.Context,
	_ *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	if err := plugins.CheckAvailability(ctx, c, ServiceID); err != nil {
		return nil, err
	}
	return proxyv1.NewRegisterProxyClient(c), nil
}

func NewPlugin(proxy proxyv1.RegisterProxyServer) plugin.Plugin {
	return &proxyBackendPlugin{
		proxySrv: proxy,
	}
}

func init() {
	plugins.GatewayScheme.Add(ProxyPluginID, NewPlugin(nil))
}
