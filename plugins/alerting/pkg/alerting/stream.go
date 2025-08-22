package alerting

import (
	"github.com/aity-cloud/monty/pkg/agent"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/node"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/rules"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[node.NodeAlertingCapabilityServer](&node.NodeAlertingCapability_ServiceDesc, &p.node),
		util.PackService[rules.RuleSyncServer](&rules.RuleSync_ServiceDesc, p.AlarmServerComponent),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.delegate.Set(streamext.NewDelegate(cc, agent.NewClientSet))
}
