<script>
import { mapGetters } from 'vuex';
import { LabeledInput } from '@components/Form/LabeledInput';
import CopyCode from '@shell/components/CopyCode';
import { Checkbox } from '@components/Form/Checkbox';
import KeyValue from '@shell/components/form/KeyValue';
import Loading from '@shell/components/Loading';
import { Card } from '@components/Card';
import { Banner } from '@components/Banner';
import { createAgent } from '../utils/requests';
import {
  updateCluster, getClusterFingerprint, createToken, getDefaultImageRepository, getGatewayConfig
} from '../utils/requests/management';
import { generateName } from '../utils/nameGenerator';

const TOKEN_EXPIRATION_SECONDS = 3600;
const NEW_TOKEN_INTERVAL_SECONDS = TOKEN_EXPIRATION_SECONDS - 100;

export default {
  components: {
    Banner,
    Card,
    CopyCode,
    LabeledInput,
    Loading,
    KeyValue,
    Checkbox,
  },

  async fetch() {
    const [token, clusterFingerprint, defaultImageRepository, gatewayConfig] = await Promise.all([
      createToken(`${ TOKEN_EXPIRATION_SECONDS }s`, '', []),
      getClusterFingerprint(),
      getDefaultImageRepository(),
      getGatewayConfig()
    ]);

    const hostname = gatewayConfig?.map(g => g.json)?.find(g => g.kind === 'GatewayConfig')?.spec?.hostname;
    const port = gatewayConfig?.map(g => g.json)?.find(g => g.kind === 'GatewayConfig')?.spec?.grpcListenAddress?.split(':')[1];

    this.$set(this, 'token', token.id);
    this.$set(this, 'pin', clusterFingerprint);
    this.$set(this, 'defaultImageRepository', defaultImageRepository || '');
    this.$set(this, 'gatewayAddress', `${ hostname }:${ port }` || '');
  },

  data() {
    const placeholderText = 'Select a token to view install command';

    return {
      isManualOpen:           false,
      token:                  null,
      labels:                 {},
      name:                   '',
      newCluster:             null,
      newClusterFound:        false,
      pin:                    null,
      placeholderText,
      error:                  '',
      agentVersion:           'v2',
      namespace:              'monty-agent',
      gatewayAddress:         '',
      installPrometheus:      false,
      useOCI:                 false,
      defaultImageRepository: '',
      newTokenInervalHandle:  null,
      generatedName:          generateName()
    };
  },

  created() {
    this.newTokenInervalHandle = setInterval(async() => {
      const token = await createToken(`${ TOKEN_EXPIRATION_SECONDS }s`, '', []);

      this.$set(this, 'token', token.id);
    }, NEW_TOKEN_INTERVAL_SECONDS * 1000);
  },

  beforeDestroy() {
    if (this.newTokenInervalHandle) {
      clearInterval(this.newTokenInervalHandle);
    }
  },

  methods: {
    save() {
      updateCluster(this.newCluster.id, this.friendlyName, { ...(this.newCluster.labels || {}), ...this.labels });

      this.$router.replace({ name: 'agents' });
    },

    createAgent() {
      createAgent(this.token);
    },

    lookForNewCluster() {
      const foundCluster = this.clusters.find(c => c.name === this.friendlyName);

      if (foundCluster) {
        this.newCluster = foundCluster;
        this.newClusterFound = true;
      }
    },

    isArgEmpty(arg) {
      return !arg || arg === 'false';
    },
  },
  computed: {
    ...mapGetters({ clusters: 'monty/clusters' }),
    friendlyName() {
      return this.name || this.generatedName;
    },
    installCommand() {
      const prometheus = this.installPrometheus ? '--set kube-prometheus-stack.enabled=true' : '';
      const defaultImageRepository = this.defaultImageRepository ? `--set image.repository=${ this.defaultImageRepository }` : '';
      const imageMain = this.useOCI ? 'oci://ghcr.io/rancher/monty-agent' : 'monty/monty-agent';

      return `helm -n ${ this.namespace } install monty-agent ${ imageMain } ${ prometheus } --set address=${ this.gatewayAddress } --set pin=${ this.pin } --set token=${ this.token } --set friendlyName=${ this.friendlyName } --create-namespace ${ defaultImageRepository }`;
    },
    gatewayUrl() {
      return this.installCommand.match(/gateway-url=.+/s)?.[0]?.replace('gateway-url=', '').replace('"  ', '').replace('https://', '').replace('http://', '');
    },
    addressUrl() {
      return this.installCommand.match(/address=.+ /s)?.[0]?.replace('address=', '').replace('"  ', '').replace('https://', '').replace('http://', '');
    },
    url() {
      return this.addressUrl || this.gatewayUrl || '';
    },
    manualIcon() {
      return this.isManualOpen ? 'icon-chevron-up' : 'icon-chevron-down';
    },
    clusterCount() {
      return this.clusters.length;
    }
  },

  watch: {
    clusterCount() {
      this.lookForNewCluster();
    }
  }
};
</script>
<template>
  <Loading v-if="$fetchState.pending" />
  <div v-else>
    <div class="row">
      <div class="col span-12">
        <LabeledInput v-model="name" label="Name (optional)" />
      </div>
    </div>
    <div class="mb-20">
      <h4 class="text-default-text">
        Labels
      </h4>
      <KeyValue
        v-model="labels"
        mode="edit"
        :read-allowed="false"
        :value-multiline="false"
        label="Labels"
        add-label="Add Label"
      />
    </div>
    <Card class="m-0 mb-10" :show-highlight-border="false" :show-actions="false">
      <h4 slot="title" class="text-default-text">
        Install Command
      </h4>
      <div slot="body">
        <div v-if="token" class="row">
          <div class="options col span-4 mb-10">
            <LabeledInput
              v-model="namespace"
              label="Namespace"
            />
          </div>
          <div class="options col span-4 mb-10">
            <LabeledInput
              v-model="gatewayAddress"
              label="Gateway Address"
            />
          </div>
          <div class="checkboxes col span-4 mb-10">
            <div>
              <Checkbox
                v-model="installPrometheus"
                label="Install Prometheus Operator"
              />
            </div>
            <div>
              <Checkbox
                v-model="useOCI"
                label="Use OCI Image"
              />
            </div>
          </div>
        </div>
        <CopyCode
          class="install-command"
          :class="installCommand === placeholderText && 'placeholder-text'"
        >
          {{ installCommand }}
        </CopyCode>
        <Banner v-if="!newCluster" color="info">
          Copy and run the above command to install the Monty agent on the new cluster.
          <a
            v-if="$config.dev"
            class="btn bg-info mr-5"
            @click="createAgent"
          >
            [Developer Mode] Start Agent
          </a>
        </Banner>
        <Banner v-if="newClusterFound" color="success">
          New cluster added successfully
          <a
            class="btn bg-success mr-5"
            @click="save"
          >
            Finish
          </a>
        </Banner>
      </div>
    </Card>
    <Card class="m-0 mt-20 manual-install" :show-highlight-border="false" :show-actions="false">
      <h4 slot="title" class="text-default-text" @click="isManualOpen = !isManualOpen">
        <span>Manual Install Information</span><i class="icon" :class="manualIcon" />
      </h4>
      <div slot="body">
        <div v-if="isManualOpen">
          <div v-if="!token || !pin" class="mt-10">
            <CopyCode class="placeholder-text">
              {{ placeholderText }}
            </CopyCode>
          </div>
          <div v-else>
            <div class="mt-10">
              Bootstrap Token: <CopyCode class="ml-5">
                {{ token }}
              </CopyCode>
            </div>
            <div class="mt-10">
              Certificate Pin: <CopyCode class="ml-5">
                {{ pin }}
              </CopyCode>
            </div>
            <div class="mt-10">
              Gateway URL: <CopyCode class="ml-5">
                {{ gatewayAddress }}
              </CopyCode>
            </div>
            <div class="mt-10">
              Friendly Name: <CopyCode class="ml-5">
                {{ friendlyName }}
              </CopyCode>
            </div>
          </div>
        </div>
      </div>
    </Card>
  </div>
</template>

<style lang="scss" scoped>
.resource-footer {
  display: flex;
  flex-direction: row;

  justify-content: flex-end;
}

.install-command {
  width: 100%;
}

.placeholder-text {
  font-family: $body-font;
}

.token-select ::v-deep #vs2__combobox > div.vs__selected-options > span {
  font-family: $mono-font;
}

::v-deep .info, ::v-deep .success {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

::v-deep {
  .manual-install {
    &.card-container {
      min-height: initial;
    }

    h4 {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      width: 100%;
      cursor: pointer;
      margin-bottom: 0;

      i {
        font-size: 30px;
        line-height: 20px;
      }
    }
  }
}

.options {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.checkboxes {
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  align-items: flex-start;
}
</style>
