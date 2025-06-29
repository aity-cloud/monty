<script>
import ArrayListSelect from '@shell/components/form/ArrayListSelect';
import UnitInput from '@shell/components/form/UnitInput';
import Loading from '@shell/components/Loading';
import { createComputedTime } from '@pkg/monty/utils/computed';
import { AlertType } from '@pkg/monty/models/alerting/Condition';
import { mapClusterOptions, loadChoices } from './shared';

const TYPE = 'monitoringBackend';

export const CONSTS = {
  TYPE,
  ENUM:        AlertType.MONITORING_BACKEND,
  TYPE_OPTION: {
    label: 'Monitoring Backend',
    value: TYPE,
    enum:  AlertType.MONITORING_BACKEND
  },
  DEFAULT_CONFIG: { [TYPE]: { clusterId: { id: '' }, for: '30s' } },
};

export default {
  ...CONSTS,

  components: {
    ArrayListSelect,
    Loading,
    UnitInput,
  },

  props: {
    value: {
      type:     Object,
      required: true
    }
  },

  async fetch() {
    await this.loadChoices();
  },

  data() {
    return {
      ...CONSTS,
      clusterOptions: [],
      choices:        { clusters: [] },
      error:          '',
    };
  },

  methods: {
    async loadChoices() {
      await loadChoices(this, this.TYPE, this.ENUM);
    }
  },

  computed: {
    ...mapClusterOptions(),

    monitoringBackendBackendComponentOptions() {
      const options = this.choices.backendComponents || [];

      if ((!this.value.backendComponents || this.value.backendComponents.length === 0) && !options.find(o => o === this.value.backendComponents)) {
        this.$set(this.value, 'backendComponents', options[0] ? [options[0]] : []);
      }

      return options;
    },

    monitoringBackendFor: createComputedTime('value.for'),
  },
};
</script>
<template>
  <Loading v-if="$fetchState.pending" />
  <div v-else>
    <h4 class="mt-20">
      State
    </h4>
    <div class="row mt-10">
      <div class="col span-6">
        <UnitInput v-model="monitoringBackendFor" label="Duration" suffix="s" :required="true" />
      </div>
    </div>
    <div class="row mt-10">
      <div class="col span-12">
        <ArrayListSelect
          v-model="value.backendComponents"
          label="Backend Components"
          :disabled="monitoringBackendBackendComponentOptions.length === 0"
          :options="(monitoringBackendBackendComponentOptions)"
          :required="true"
          add-label="Add Backend Component"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
</style>
