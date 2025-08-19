package sync

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
)

type SyncTask func(ctx context.Context, syncInfo SyncInfo) error

type SyncInfo struct {
	ShouldSync bool
	// clusterId -> cluster
	Clusters map[string]*corev1.Cluster
}
