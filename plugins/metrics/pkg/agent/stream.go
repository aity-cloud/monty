package agent

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"
	"github.com/aity-cloud/monty/pkg/clients"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/metrics/apis/node"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remoteread"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remotewrite"
	"google.golang.org/grpc"
)

// StreamServers implements stream.StreamAPIExtension.
func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[capabilityv1.NodeServer](&capabilityv1.Node_ServiceDesc, p.node),
		util.PackService[remoteread.RemoteReadAgentServer](&remoteread.RemoteReadAgent_ServiceDesc, p.node),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	nodeClient := node.NewNodeMetricsCapabilityClient(cc)
	healthListenerClient := controlv1.NewHealthListenerClient(cc)
	identityClient := controlv1.NewIdentityClient(cc)

	p.httpServer.SetRemoteWriteClient(clients.NewLocker(cc, remotewrite.NewRemoteWriteClient))
	p.ruleStreamer.SetRemoteWriteClient(remotewrite.NewRemoteWriteClient(cc))
	p.node.SetRemoteWriter(clients.NewLocker(cc, remotewrite.NewRemoteWriteClient))

	p.node.SetClients(
		nodeClient,
		identityClient,
		healthListenerClient,
	)
}
