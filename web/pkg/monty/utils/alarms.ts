import { CONSTS as AgentDisconnectConsts } from '@pkg/monty/components/Alarm/AgentDisconnect.vue';
import { CONSTS as KubeStateConsts } from '@pkg/monty/components/Alarm/KubeState.vue';
import { CONSTS as DownstreamCapabilityConsts } from '@pkg/monty/components/Alarm/DownstreamCapability.vue';
import { CONSTS as MonitoringBackendConsts } from '@pkg/monty/components/Alarm/MonitoringBackend.vue';
import { CONSTS as PrometheusConsts } from '@pkg/monty/components/Alarm/Prometheus.vue';

export const CONDITION_TYPES = [
  AgentDisconnectConsts.TYPE_OPTION,
  KubeStateConsts.TYPE_OPTION,
  DownstreamCapabilityConsts.TYPE_OPTION,
  PrometheusConsts.TYPE_OPTION,
  MonitoringBackendConsts.TYPE_OPTION,
];

export const CONDITION_TYPES_WITHOUT_MONITORING = [
  AgentDisconnectConsts.TYPE_OPTION,
  DownstreamCapabilityConsts.TYPE_OPTION,
  PrometheusConsts.TYPE_OPTION,
  MonitoringBackendConsts.TYPE_OPTION,
];
