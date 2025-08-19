package types

import (
	"context"
	"log/slog"

	"github.com/aity-cloud/monty/pkg/agent"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	streamext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remoteread"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers"
	"golang.org/x/tools/pkg/memoize"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type PluginContext interface {
	context.Context
	Logger() *slog.Logger
	Metrics() *Metrics
	Memoize(key any, fn memoize.Function) *memoize.Promise

	ManagementClient() managementv1.ManagementClient
	KeyValueStoreClient() system.KeyValueStoreClient
	StreamClient() grpc.ClientConnInterface
	GatewayConfigClient() configv1.GatewayConfigClient
	ClusterDriver() drivers.ClusterDriver
	// AuthMiddlewares() map[string]auth.Middleware
	ExtensionClient() system.ExtensionClientInterface
}

type MetricsAgentClientSet interface {
	agent.ClientSet
	remoteread.RemoteReadAgentClient
}

type ServiceContext interface {
	PluginContext
	StorageBackend() storage.Backend
	Delegate() streamext.StreamDelegate[MetricsAgentClientSet]
}

type ManagementServiceContext interface {
	ServiceContext
	SetServingStatus(serviceName string, status healthpb.HealthCheckResponse_ServingStatus)
}

type StreamServiceContext interface {
	ServiceContext
}
