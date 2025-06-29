package agent

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	"github.com/aity-cloud/monty/plugins/topology/apis/node"
	"github.com/aity-cloud/monty/plugins/topology/apis/stream"

	// "github.com/aity-cloud/monty/pkg/clients"
	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"

	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []streamext.Server {
	return []streamext.Server{
		{
			Desc: &capabilityv1.Node_ServiceDesc,
			Impl: p.node,
		},
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.topologyStreamer.SetTopologyStreamClient(stream.NewRemoteTopologyClient(cc))
	p.topologyStreamer.SetIdentityClient(controlv1.NewIdentityClient(cc))
	p.node.SetClient(node.NewNodeTopologyCapabilityClient(cc))
}
