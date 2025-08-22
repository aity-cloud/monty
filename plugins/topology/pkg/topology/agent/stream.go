package agent

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/topology/apis/node"
	"github.com/aity-cloud/monty/plugins/topology/apis/stream"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[capabilityv1.NodeServer](&capabilityv1.Node_ServiceDesc, p.node),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.topologyStreamer.SetTopologyStreamClient(stream.NewRemoteTopologyClient(cc))
	p.topologyStreamer.SetIdentityClient(controlv1.NewIdentityClient(cc))
	p.node.SetClient(node.NewNodeTopologyCapabilityClient(cc))
}
