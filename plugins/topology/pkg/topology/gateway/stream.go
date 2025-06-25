package gateway

import (
	"github.com/aity-cloud/monty/pkg/agent"
	"github.com/aity-cloud/monty/pkg/capabilities/wellknown"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/plugins/topology/apis/node"
	"github.com/aity-cloud/monty/plugins/topology/apis/stream"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []streamext.Server {
	return []streamext.Server{
		{
			Desc:              &stream.RemoteTopology_ServiceDesc,
			Impl:              &p.topologyRemoteWrite,
			RequireCapability: wellknown.CapabilityTopology,
		},
		{
			Desc:              &node.NodeTopologyCapability_ServiceDesc,
			Impl:              &p.topologyBackend,
			RequireCapability: wellknown.CapabilityTopology,
		},
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.delegate.Set(streamext.NewDelegate(cc, agent.NewClientSet))
}
