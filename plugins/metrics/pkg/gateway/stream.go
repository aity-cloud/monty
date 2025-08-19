package gateway

import (
	"context"

	"github.com/aity-cloud/monty/pkg/agent"
	"github.com/aity-cloud/monty/pkg/auth/cluster"
	"github.com/aity-cloud/monty/pkg/metrics"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remoteread"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

// StreamServers implements streamext.StreamAPIExtension.
func (p *Plugin) StreamServers() []util.ServicePackInterface {
	return p.streamServices
}

// UseStreamClient implements streamext.StreamClientHandler.
func (p *Plugin) UseStreamClient(cc grpc.ClientConnInterface) {
	p.streamClient.C() <- cc
	type clientset struct {
		agent.ClientSet
		remoteread.RemoteReadAgentClient
	}
	p.delegate.C() <- streamext.NewDelegate(cc, func(cci grpc.ClientConnInterface) types.MetricsAgentClientSet {
		return &clientset{
			ClientSet:             agent.NewClientSet(cci),
			RemoteReadAgentClient: remoteread.NewRemoteReadAgentClient(cci),
		}
	})
}

func (p *Plugin) labelsForStreamMetrics(ctx context.Context) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key(metrics.LabelImpersonateAs).String(cluster.StreamAuthorizedID(ctx)),
		attribute.Key("handler").String("plugin_metrics"),
	}
}
