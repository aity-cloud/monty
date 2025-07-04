package management

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *Server) GetClusterHealthStatus(
	_ context.Context,
	ref *corev1.Reference,
) (*corev1.HealthStatus, error) {
	if m.healthStatusDataSource == nil {
		return nil, status.Error(codes.Unavailable, "health API not configured")
	}
	if err := validation.Validate(ref); err != nil {
		return nil, err
	}

	return m.healthStatusDataSource.GetClusterHealthStatus(ref)
}

func (m *Server) WatchClusterHealthStatus(
	_ *emptypb.Empty,
	stream managementv1.Management_WatchClusterHealthStatusServer,
) error {
	if m.healthStatusDataSource == nil {
		return status.Error(codes.Unavailable, "health API not configured")
	}

	healthStatusC := m.healthStatusDataSource.WatchClusterHealthStatus(stream.Context())
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case healthStatus, ok := <-healthStatusC:
			if !ok {
				return nil
			}
			if err := stream.Send(healthStatus); err != nil {
				return err
			}
		}
	}
}
