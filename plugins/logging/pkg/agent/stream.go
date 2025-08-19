package agent

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/logging/apis/node"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[capabilityv1.NodeServer](&capabilityv1.Node_ServiceDesc, p.node),
		util.PackService[collogspb.LogsServiceServer](&collogspb.LogsService_ServiceDesc, p.otelForwarder.LogsForwarder),
		util.PackService[coltracepb.TraceServiceServer](&coltracepb.TraceService_ServiceDesc, p.otelForwarder.TraceForwarder),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	nodeClient := node.NewNodeLoggingCapabilityClient(cc)
	p.node.SetClient(nodeClient)
	p.otelForwarder.TraceForwarder.SetClient(cc)
	p.otelForwarder.LogsForwarder.SetClient(cc)
}
