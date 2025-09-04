package services

import (
	"context"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"buf.build/go/protovalidate"
	managementext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/storage/kvutil"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/flagutil"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/cortex/configutil"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/types"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CortexOpsService struct {
	Context types.ManagementServiceContext `option:"context"`
	*driverutil.BaseConfigServer[
		*driverutil.GetRequest,
		*cortexops.SetRequest,
		*cortexops.ResetRequest,
		*driverutil.ConfigurationHistoryRequest,
		*cortexops.ConfigurationHistoryResponse,
		*cortexops.CapabilityBackendConfigSpec,
	]
	drivers.PartialCortexOpsServer
}

var _ cortexops.CortexOpsServer = (*CortexOpsService)(nil)

func (s *CortexOpsService) Activate() error {
	defer s.Context.SetServingStatus(cortexops.CortexOps_ServiceDesc.ServiceName, managementext.Serving)

	defaultStore := kvutil.WithKey(system.NewKVStoreClient[*cortexops.CapabilityBackendConfigSpec](s.Context.KeyValueStoreClient()), "/config/cluster/default")
	activeStore := s.Context.ClusterDriver().ActiveConfigStore()

	s.BaseConfigServer = s.Build(defaultStore, activeStore, flagutil.LoadDefaults)
	s.PartialCortexOpsServer = s.Context.ClusterDriver()
	return nil
}

func (s *CortexOpsService) ManagementServices() []util.ServicePackInterface {
	return []util.ServicePackInterface{
		util.PackService[cortexops.CortexOpsServer](&cortexops.CortexOps_ServiceDesc, s),
	}
}

func (s *CortexOpsService) DryRun(ctx context.Context, req *cortexops.DryRunRequest) (*cortexops.DryRunResponse, error) {
	res, err := s.ServerDryRun(ctx, req)
	if err != nil {
		return nil, err
	}

	upstreamErrs := configutil.CollectValidationErrorLogs(res.Modified.CortexConfig)
	if len(upstreamErrs) > 0 {
		if res.ValidationErrors == nil {
			res.ValidationErrors = &protovalidate.ValidationError{}
		}
		for _, err := range upstreamErrs {
			res.ValidationErrors.Violations = append(res.ValidationErrors.Violations, &protovalidate.Violation{
				Proto: &validate.Violation{
					RuleId:  lo.ToPtr("cortex"),
					Message: lo.ToPtr(err.Error()),
				},
			})
		}
	}

	return &cortexops.DryRunResponse{
		Current:          res.Current,
		Modified:         res.Modified,
		ValidationErrors: res.ValidationErrors.ToProto(),
	}, nil
}

// Overrides BaseConfigServer.Install
func (s *CortexOpsService) Install(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	out, err := s.BaseConfigServer.Install(ctx, in)
	if err != nil {
		return nil, err
	}
	if err := BroadcastNodeSync(s.Context); err != nil {
		return nil, err
	}
	return out, nil
}

// Overrides BaseConfigServer.Uninstall
func (s *CortexOpsService) Uninstall(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	out, err := s.BaseConfigServer.Uninstall(ctx, in)
	if err != nil {
		return nil, err
	}
	if err := BroadcastNodeSync(s.Context); err != nil {
		return nil, err
	}
	return out, nil
}

func init() {
	types.Services.Register("Cortex Ops Service", func(_ context.Context, opts ...driverutil.Option) (types.Service, error) {
		svc := &CortexOpsService{}
		driverutil.ApplyOptions(svc, opts...)
		return svc, nil
	})
}
