package agent

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/node"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/rules"
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
	nodeClient := node.NewNodeAlertingCapabilityClient(cc)
	healthListenerClient := controlv1.NewHealthListenerClient(cc)
	identityClient := controlv1.NewIdentityClient(cc)
	ruleSyncClient := rules.NewRuleSyncClient(cc)

	p.node.SetClients(
		healthListenerClient,
		nodeClient,
		identityClient,
	)

	p.ruleStreamer.SetClients(ruleSyncClient)
}
