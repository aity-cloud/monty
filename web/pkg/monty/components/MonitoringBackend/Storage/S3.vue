<script>
import { LabeledInput } from '@components/Form/LabeledInput';
import LabeledSelect from '@shell/components/form/LabeledSelect';
import { Checkbox } from '@components/Form/Checkbox';
import UnitInput from '@shell/components/form/UnitInput';
import { Storage } from '@pkg/monty/api/monty';
import { createComputedDuration } from '@pkg/monty/utils/computed';
import { S3_REGIONS, S3_REGION_TO_ENDPOINT } from '@pkg/monty/utils/storage';

export const SECONDS_IN_DAY = 86400;

export default {
  components: {
    Checkbox, UnitInput, LabeledInput, LabeledSelect
  },

  props: {
    value: {
      type:     Object,
      required: true
    },
  },

  created() {
    if (!this.value.cortexConfig.storage?.s3?.endpoint) {
      this.updateEndpoint();
    }
  },

  data() {
    return {
      signatureVersionOptions: [
        {
          label: 'v4',
          value: 'v4'
        },
        {
          label: 'v2',
          value: 'v2'
        },
      ],
      regions:  S3_REGIONS,
      sseTypes: [
        { label: 'None', value: '' },
        { label: 'SSE-KMS', value: 'SSE-KMS' },
        { label: 'SSE-S3', value: 'SSE-S3' },
      ],
      Storage,
    };
  },

  computed: {
    s3RetentionPeriod: createComputedDuration('value.cortexConfig.limits.compactorBlocksRetentionPeriod', SECONDS_IN_DAY),

    s3IdleConnTimeout: createComputedDuration('value.cortexConfig.storage.s3.http.idleConnTimeout'),

    s3ResponseHeaderTimeout: createComputedDuration('value.cortexConfig.storage.s3.http.responseHeaderTimeout'),

    s3TlsHandshakeTimeout: createComputedDuration('value.cortexConfig.storage.s3.http.tlsHandshakeTimeout'),

    s3ExpectContinueTimeout: createComputedDuration('value.cortexConfig.storage.s3.http.expectContinueTimeout'),
  },

  methods: {
    updateEndpoint() {
      if (this.value.cortexConfig.storage?.s3?.region) {
        return this.$set(this.value.cortexConfig.storage.s3, 'endpoint', `${ S3_REGION_TO_ENDPOINT[this.value.cortexConfig.storage.s3.region] }`);
      }
    },
  },
};
</script>
<template>
  <div class="m-0">
    <h3>Target</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledSelect v-model="value.cortexConfig.storage.s3.region" :options="regions" label="Region" @input="updateEndpoint" />
      </div>
      <div class="col span-6">
        <LabeledInput v-model="value.cortexConfig.storage.s3.bucketName" label="Bucket Name" :required="true" />
      </div>
    </div>
    <div class="row mb-10 border">
      <div class="col span-6">
        <LabeledInput v-model="value.cortexConfig.storage.s3.endpoint" label="Endpoint" :required="true" />
      </div>
      <div class="col span-6 middle">
        <Checkbox v-model="value.cortexConfig.storage.s3.insecure" label="Insecure" />
      </div>
    </div>
    <h3>Access</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput v-model="value.cortexConfig.storage.s3.accessKeyId" label="Access Key ID" :required="true" />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model="value.cortexConfig.storage.s3.secretAccessKey"
          label="Secret Access Key"
          :required="true"
          type="password"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledSelect
          v-model="value.cortexConfig.storage.s3.signatureVersion"
          :options="signatureVersionOptions"
          label="Signature Version"
        />
      </div>
    </div>
    <h3>Server Side Encryption</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledSelect v-model="value.cortexConfig.storage.s3.sse.type" :options="sseTypes" label="Type" />
      </div>
    </div>
    <div v-if="value.cortexConfig.storage.s3.sse.type === 'SSE-KMS'" class="row mb-10">
      <div class="col span-6">
        <LabeledInput v-model="value.cortexConfig.storage.s3.sse.kmsKeyID" label="KMS Key Id" :required="true" />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model="value.cortexConfig.storage.s3.sse.kmsEncryptionContext"
          label="KMS Encryption Context"
          :required="true"
        />
      </div>
    </div>
    <h3>Connection</h3>
    <div class="row mb-10">
      <div class="col span-4">
        <UnitInput v-model="s3IdleConnTimeout" label="Idle Connection Timeout" placeholder="e.g. 30, 60" suffix="s" />
      </div>
      <div class="col span-4">
        <UnitInput
          v-model="s3ResponseHeaderTimeout"
          label="Response Header Timeout"
          placeholder="e.g. 30, 60"
          suffix="s"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-4">
        <UnitInput
          v-model="s3TlsHandshakeTimeout"
          label="TLS Handshake Timeout"
          placeholder="e.g. 30, 60"
          suffix="s"
        />
      </div>
      <div class="col span-4 middle">
        <Checkbox v-model="value.cortexConfig.storage.s3.http.insecureSkipVerify" label="Insecure Skip Verify" />
      </div>
      <div class="col span-4">
        <UnitInput
          v-model="s3ExpectContinueTimeout"
          label="Expect Continue Timeout"
          placeholder="e.g. 30, 60"
          suffix="s"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-4">
        <UnitInput v-model="value.cortexConfig.storage.s3.http.maxIdleConnections" label="Max Idle Connections" suffix="" :input-exponent="0" base-unit="" />
      </div>
      <div class="col span-4">
        <UnitInput
          v-model="value.cortexConfig.storage.s3.http.maxIdleConnectionsPerHost"
          label="Max Idle Connections Per Host"
          suffix=""
          :input-exponent="0"
          base-unit=""
        />
      </div>
      <div class="col span-4">
        <UnitInput v-model="value.cortexConfig.storage.s3.http.maxConnectionsPerHost" label="Max Connections Per Host" suffix="" :input-exponent="0" base-unit="" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
header {
  width: 100%;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0;
}

::v-deep {
  .not-enabled {
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
  }

  .enabled {
    width: 100%;
  }
}
</style>
