package gateway

import (
	"github.com/aity-cloud/monty/pkg/agent"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/logging/apis/node"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[node.NodeLoggingCapabilityServer](&node.NodeLoggingCapability_ServiceDesc, &p.logging),
		util.PackService[collogspb.LogsServiceServer](&collogspb.LogsService_ServiceDesc, p.otelForwarder.LogsForwarder),
		util.PackService[coltracepb.TraceServiceServer](&coltracepb.TraceService_ServiceDesc, p.otelForwarder.TraceForwarder),
	}
}

func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.delegate.Set(streamext.NewDelegate(cc, agent.NewClientSet))
}
