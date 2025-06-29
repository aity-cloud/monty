import { CoreStoreSpecifics, CoreStoreConfig } from '@shell/core/types';
import { Cluster as ClusterModel } from '@pkg/monty/models/Cluster';
import * as core from '@pkg/monty/generated/ github.com/aity-cloud/monty/pkg/apis/core/v1/core_pb';
import * as management from '@pkg/monty/generated/ github.com/aity-cloud/monty/pkg/apis/management/v1/management_pb';

const config: CoreStoreConfig = { namespace: 'monty' };

export interface ClustersState {
  clusters: ClusterModel[];
}

// Getters
export const Clusters = 'monty/clusters';

// Mutations
export const UpdateCluster = 'monty/updateCluster';
export const UpdateClusterHealthStatus = 'monty/updateClusterHealthStatus';
export const DeleteCluster = 'monty/deleteCluster';

// Actions
export const HandleClusterWatchEvent = 'monty/handleClusterWatchEvent';
export const HandleClusterHealthStatusEvent = 'monty/handleClusterHealthStatusEvent';

const createStore = (): CoreStoreSpecifics => {
  const initialHealthStatusCache = new Map<string, core.HealthStatus>();

  return {
    state(): ClustersState {
      return { clusters: [] };
    },
    getters:   {
      clusters: (state: ClustersState) => {
        return state.clusters;
      }
    },
    mutations: {
      updateCluster(state: ClustersState, cluster: core.Cluster) {
        const clusterId = cluster.id;

        const clusterIndex = state.clusters.findIndex(c => c.id === clusterId);

        let c: ClusterModel;

        if (clusterIndex === -1) {
          c = ClusterModel.create(this._vm);
          state.clusters.push(c);
        } else {
          c = state.clusters[clusterIndex];
        }

        c.onClusterUpdated(cluster);
        // check if the cluster has an initial health status that came in out of order
        const hs = initialHealthStatusCache.get(clusterId);

        if (hs) {
          c.onHealthStatusUpdated(hs);
          initialHealthStatusCache.delete(clusterId);
        }
      },
      updateClusterHealthStatus(state: ClustersState, payload: { id: string; healthStatus: core.HealthStatus }) {
        const clusterIndex = state.clusters.findIndex(c => c.id === payload.id);

        if (clusterIndex !== -1) {
          const c = state.clusters[clusterIndex];

          c.onHealthStatusUpdated(payload.healthStatus);
        } else {
          initialHealthStatusCache.set(payload.id, payload.healthStatus);
        }
      },
      deleteCluster(state: ClustersState, clusterId: string) {
        const clusterIndex = state.clusters.findIndex(c => c.id === clusterId);

        if (clusterIndex !== -1) {
          state.clusters.splice(clusterIndex, 1);
        }
      }
    },
    actions: {
      handleClusterWatchEvent(ctx: any, event: management.WatchEvent) {
        if (!event.cluster?.id) {
          return;
        }

        const clusterId = event.cluster.id;

        if (event.type === management.WatchEventType.Put) {
          ctx.commit('updateCluster', event.cluster);
        } else if (event.type === management.WatchEventType.Delete) {
          ctx.commit('deleteCluster', clusterId);
        }
      },
      handleClusterHealthStatusEvent(ctx: any, hs: core.ClusterHealthStatus) {
        if (!hs.cluster || !hs.healthStatus) {
          return;
        }
        const id = hs.cluster.id;
        const healthStatus = hs.healthStatus;

        ctx.commit('updateClusterHealthStatus', { id, healthStatus });
      },
    },
  };
};

export default {
  specifics: createStore(),
  config,
};
