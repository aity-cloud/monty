package backend

import (
	"context"
	"os"

	montycorev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/logger"
)

func (b *LoggingBackend) updateClusterMetadata(ctx context.Context, event *managementv1.WatchEvent) error {
	incomingLabels := event.GetCluster().GetMetadata().GetLabels()
	previousLabels := event.GetPrevious().GetMetadata().GetLabels()
	var newName, oldName string
	if _, ok := incomingLabels[montycorev1.NameLabel]; ok {
		newName = incomingLabels[montycorev1.NameLabel]
	}
	if _, ok := previousLabels[montycorev1.NameLabel]; ok {
		oldName = previousLabels[montycorev1.NameLabel]
	}
	if newName == oldName {
		b.Logger.With(
			"oldName", oldName,
			"newName", newName,
		).Debug("cluster was not renamed")
		return nil
	}
	b.Logger.With(
		"oldName", oldName,
		"newName", newName,
	).Debug("cluster was renamed")

	if err := b.ClusterDriver.StoreClusterMetadata(ctx, event.Cluster.GetId(), newName); err != nil {
		b.Logger.With(
			logger.Err(err),
			"cluster", event.Cluster.Id,
		).Debug("could not update cluster metadata")
		return nil
	}

	return nil
}

func (b *LoggingBackend) watchClusterEvents(ctx context.Context) {
	clusterClient, err := b.MgmtClient.WatchClusters(ctx, &managementv1.WatchClustersRequest{})
	if err != nil {
		b.Logger.With(logger.Err(err)).Error("failed to watch clusters, existing")
		os.Exit(1)
	}

	b.Logger.Info("watching cluster events")

outer:
	for {
		select {
		case <-clusterClient.Context().Done():
			b.Logger.Info("context cancelled, stoping cluster event watcher")
			break outer
		default:
			event, err := clusterClient.Recv()
			if err != nil {
				b.Logger.With(logger.Err(err)).Error("failed to receive cluster event")
				continue
			}

			b.watcher.HandleEvent(event)
		}
	}
}

func (b *LoggingBackend) reconcileClusterMetadata(ctx context.Context, clusters []*montycorev1.Cluster) (retErr error) {
	for _, cluster := range clusters {
		err := b.ClusterDriver.StoreClusterMetadata(ctx, cluster.GetId(), cluster.Metadata.Labels[montycorev1.NameLabel])
		if err != nil {
			b.Logger.With(
				logger.Err(err),
				"cluster", cluster.Id,
			).Warn("could not update cluster metadata")
			retErr = err
		}
	}
	return
}
