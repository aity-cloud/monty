<script>
import { mapGetters } from 'vuex';
import SortableTable from '@shell/components/SortableTable';
import Loading from '@shell/components/Loading';
import GlobalEventBus from '@pkg/monty/utils/GlobalEventBus';
import { InstallState, getClusterStatus } from '../utils/requests/alerts';
import { getSLOs } from '../utils/requests/slo';
import CloneToClustersDialog from './dialogs/CloneToClustersDialog';

export default {
  components: {
    CloneToClustersDialog, Loading, SortableTable
  },
  async fetch() {
    await this.load();
    await this.updateStatuses();
  },

  data() {
    return {
      loading:             false,
      statsInterval:       null,
      slos:                [],
      isAlertingEnabled: false,
      headers:             [
        {
          name:          'status',
          labelKey:      'monty.tableHeaders.status',
          value:         'status',
          formatter:     'StatusBadge',
          width:     100
        },
        {
          name:          'nameDisplay',
          labelKey:      'monty.tableHeaders.name',
          value:         'nameDisplay',
          width:         undefined
        },
        {
          name:          'tags',
          labelKey:      'monty.tableHeaders.tags',
          value:         'tags',
          formatter:     'ListBubbles'
        },
        {
          name:      'period',
          labelKey:  'monty.tableHeaders.period',
          value:     'period'
        },
      ]
    };
  },

  created() {
    GlobalEventBus.$on('remove', this.onRemove);
    this.$on('clone', this.onClone);
    this.statsInterval = setInterval(this.updateStatuses, 10000);
  },

  beforeDestroy() {
    GlobalEventBus.$off('remove');
    this.$off('clone');
    if (this.statsInterval) {
      clearInterval(this.statsInterval);
    }
  },

  methods: {
    onRemove() {
      this.load();
    },

    onClone(slo) {
      this.$refs.dialog.open(slo, slo.clusterId);
    },

    async load() {
      try {
        this.loading = true;
        const status = (await getClusterStatus()).state;
        const isAlertingEnabled = status === InstallState.Installed;

        this.$set(this, 'isAlertingEnabled', isAlertingEnabled);

        if (!isAlertingEnabled) {
          return;
        }

        this.$set(this, 'slos', await getSLOs(this));
      } finally {
        this.loading = false;
      }
    },
    async updateStatuses() {
      const promises = this.slos.map(slo => slo.updateStatus());

      await Promise.all(promises);
    }
  },

  computed: {
    ...mapGetters({ clusters: 'monty/clusters' }),

    hasOneMonitoring() {
      return this.clusters.some(c => c.isCapabilityInstalled('metrics'));
    }
  }
};
</script>
<template>
  <Loading v-if="loading || $fetchState.pending" />
  <div v-else>
    <header>
      <div class="title">
        <h1>SLOs</h1>
      </div>
      <div v-if="isAlertingEnabled && hasOneMonitoring" class="actions-container">
        <n-link class="btn role-primary" :to="{name: 'slo-create'}">
          Create
        </n-link>
      </div>
    </header>
    <SortableTable
      v-if="isAlertingEnabled && hasOneMonitoring"
      :rows="slos"
      :headers="headers"
      :search="false"
      group-by="clusterNameDisplay"
      default-sort-by="expirationDate"
      key-field="id"
      :rows-per-page="15"
    >
      <template #group-by="{group: thisGroup}">
        <div v-trim-whitespace class="group-tab">
          Cluster: {{ thisGroup.ref }}
        </div>
      </template>
    </SortableTable>
    <div v-else-if="!isAlertingEnabled" class="not-enabled">
      <h4>
        Alerting must be enabled to use SLOs. <n-link :to="{name: 'alerting'}">
          Click here
        </n-link> to enable alerting.
      </h4>
    </div>
    <div v-else class="not-enabled">
      <h4>
        At least one cluster must have Monitoring installed to use SLOs. <n-link :to="{name: 'monitoring'}">
          Click here
        </n-link> to enable Monitoring.
      </h4>
    </div>
    <CloneToClustersDialog ref="dialog" :clusters="clusters" @save="load" />
  </div>
</template>

<style lang="scss" scoped>
::v-deep {
  .nowrap {
    white-space: nowrap;
  }

  .monospace {
    font-family: $mono-font;
  }
}

.not-enabled {
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>
