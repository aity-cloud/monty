package machinery

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type subjectAccessCapableStore struct {
	client managementv1.ManagementClient
}

func SubjectAccessCapableStore(client managementv1.ManagementClient) storage.SubjectAccessCapableStore {
	return &subjectAccessCapableStore{
		client: client,
	}
}

func (s *subjectAccessCapableStore) ListClusters(
	ctx context.Context,
	matchLabels *corev1.LabelSelector,
	matchOptions corev1.MatchOptions,
) (*corev1.ClusterList, error) {
	return s.client.ListClusters(ctx, &managementv1.ListClustersRequest{
		MatchLabels:  matchLabels,
		MatchOptions: matchOptions,
	})
}

func (s *subjectAccessCapableStore) GetRole(
	ctx context.Context,
	ref *corev1.Reference,
) (*corev1.Role, error) {
	return s.client.GetRole(ctx, ref)
}

func (s *subjectAccessCapableStore) ListRoleBindings(
	ctx context.Context,
) (*corev1.RoleBindingList, error) {
	return s.client.ListRoleBindings(ctx, &emptypb.Empty{})
}
