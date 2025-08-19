package gateway

import (
	"github.com/aity-cloud/monty/pkg/agent"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/topology/apis/node"
	"github.com/aity-cloud/monty/plugins/topology/apis/stream"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[stream.RemoteTopologyServer](&stream.RemoteTopology_ServiceDesc, &p.topologyRemoteWrite),
		util.PackService[node.NodeTopologyCapabilityServer](&node.NodeTopologyCapability_ServiceDesc, &p.topologyBackend),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.delegate.Set(streamext.NewDelegate(cc, agent.NewClientSet))
}
