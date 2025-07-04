import Vue, { reactive } from 'vue';
import { uninstallCapabilityStatus } from '@pkg/monty/utils/requests/management';
import { exceptionToErrorsArray } from '@pkg/monty/utils/error';
import { Management, Capability as CapabilityProto } from '@pkg/monty/api/monty';
import { Reference } from '@pkg/monty/generated/ github.com/aity-cloud/monty/pkg/apis/core/v1/core_pb';
import { InstallRequest } from '@pkg/monty/generated/ github.com/aity-cloud/monty/pkg/apis/capability/v1/capability_pb';
import GlobalEventBus from '@pkg/monty/utils/GlobalEventBus';
import { CapabilityInstallRequest, CapabilityStatusRequest } from '@pkg/monty/generated/ github.com/aity-cloud/monty/pkg/apis/management/v1/management_pb';
import { Struct, Value } from '@bufbuild/protobuf';
import { Cluster } from './Cluster';
import { Resource } from './Resource';

export interface CapabilityLog {
  capability: string;
  state: string;
  message: string;
}
export interface CapabilityStatusLogResponse {
  msg: string;
  level: number;
  timestamp: string;
}

export interface CapabilityStatusTransitionResponse {
  state: string;
  timestamp: string;
}

export enum CapabilityStatusState {
  Unknown = 0,
  Success = 1,
  Warning = 2,
  Error = 3,
}

export interface CapabilityStatusResponse {
  state: CapabilityStatusState;
  progress: null;
  metadata: 'string';
  logs: CapabilityStatusLogResponse[];
  transitions: CapabilityStatusTransitionResponse[];
}

export enum TaskState {
  Unknown = 0,
  Pending = 1,
  Running = 2,
  Completed = 3,
  Failed = 4,
  Canceled = 6,
}

export interface NodeCapabilityStatus {
  enabled: boolean;
  lastSync?: string;
  conditions?: string[];
}

export interface CapabilityStatus {
  state: string;
  shortMessage: string;
  message: string;
}

export interface ClusterStats {
  userID: string;
  ingestionRate: number;
  numSeries: number;
  APIIngestionRate: number;
  RuleIngestionRate: number;
}
export type Type = 'alerting' | 'metrics' | 'logs';

export class Capability extends Resource {
  private type: Type;
  private cluster: Cluster;
  private capabilityStatus?: CapabilityStatus;
  private stats?: ClusterStats[];

  constructor(type: Type, cluster: Cluster, vue: any) {
    super(vue);
    this.type = type;
    this.cluster = cluster;
  }

  static create(type: Type, cluster: Cluster, vue: any): Capability {
    return reactive(new Capability(type, cluster, vue));
  }

  get nameDisplay(): string {
    return {
      metrics: 'Monitoring', logs: 'Logging', alerting: 'Alerting'
    }[this.rawType];
  }

  get rawCluster() {
    return this.cluster;
  }

  get rawType() {
    return this.type;
  }

  get id() {
    return this.cluster.id;
  }

  get status() {
    return this.capabilityStatus || {
      state:        'success',
      shortMessage: 'Loading',
    };
  }

  get isInstalled() {
    return this.cluster.capabilities?.includes(this.type);
  }

  get isLocal(): boolean {
    return this.cluster.isLocal;
  }

  get localIcon(): string {
    return this.isLocal ? 'icon-checkmark text-success' : '';
  }

  get clusterNameDisplay(): string {
    return this.cluster.nameDisplay;
  }

  get capabilities(): string[] {
    return this.cluster.capabilities;
  }

  async updateCabilityLogs(dialog: boolean = false): Promise<void> {
    function getState(state: TaskState) {
      switch (state) {
      case TaskState.Completed:
        return 'info';
      case TaskState.Running:
      case TaskState.Pending:
      case TaskState.Canceled:
        return 'warning';
      default:
        return 'error';
      }
    }

    try {
      const capMeta = this.cluster.capabilitiesRaw?.find(c => c.name === this.type);

      if (!capMeta) {
        Vue.set(this, 'capabilityStatus', {
          state:        'info',
          shortMessage: 'Not Installed',
        });
      } else if (capMeta?.deletionTimestamp) {
        const log = await uninstallCapabilityStatus(this.cluster.id, this.type, this.vue);
        const pending = log.state === TaskState.Pending || log.state === TaskState.Running || (this.capabilityStatus as any)?.pending || false;
        const state = getState(log.state);

        Vue.set(this, 'capabilityStatus', {
          state,
          shortMessage: pending ? 'Pending' : (log.state === TaskState.Completed ? 'Not Installed' : 'Uninstalling'),
          message:      (log.logs || []).reverse()[0]?.msg,
        });
      } else {
        const apiStatus = await Management.service.CapabilityStatus(new CapabilityStatusRequest({
          cluster: new Reference({ id: this.cluster.id }),
          name:    this.type,
        }));

        if (apiStatus.conditions?.length === 0) {
          Vue.set(this, 'capabilityStatus', {
            state:        'success',
            shortMessage: 'Installed',
          });
        } else if (!apiStatus.enabled) {
          Vue.set(this, 'capabilityStatus', {
            state:        'warning',
            shortMessage: 'Disabled',
          });
        } else {
          Vue.set(this, 'capabilityStatus', {
            state:        'warning',
            shortMessage: 'Degraded',
            message:      apiStatus.conditions?.join(', '),
          });
        }
      }
    } catch (ex) {}
  }

  get availableActions(): any[] {
    return [
      {
        action:   'install',
        label:    'Install',
        icon:     'icon icon-upload',
        bulkable: true,
        enabled:  !this.isCapabilityInstalled && !this.isCapabilityUninstalling,
      },
      {
        action:   'uninstallPrompt',
        label:    'Uninstall',
        icon:     'icon icon-delete',
        bulkable: true,
        enabled:  this.isCapabilityInstalled && !this.isCapabilityUninstalling,
      },
      {
        action:   'cancelUninstallPrompt',
        label:    'Cancel Uninstall',
        icon:     'icon icon-x',
        bulkable: true,
        enabled:  this.isCapabilityUninstalling,
      }
    ];
  }

  uninstallPrompt() {
    GlobalEventBus.$emit('uninstallCapabilities', [this]);
  }

  cancelUninstallPrompt() {
    GlobalEventBus.$emit('cancelUninstallCapabilities', [this]);
  }

  get isCapabilityInstalled() {
    return this.capabilities.includes(this.type);
  }

  async install() {
    try {
      await Management.service.InstallCapability(new CapabilityInstallRequest({
        name:   this.type,
        target: new InstallRequest({ cluster: new Reference({ id: this.cluster.id }) }),
      }));

      setTimeout(() => {
        this.updateCabilityLogs();
      }, 10);
    } catch (ex) {
      Vue.set(this, 'capabilityStatus', {
        state:        'error',
        shortMessage: 'Error',
        message:      exceptionToErrorsArray(ex).join('; '),
      });
    }
  }

  async uninstall(deleteData: boolean = true) {
    try {
      const options = deleteData ? new Struct({ fields: { initialDelay: new Value({ kind: { case: 'stringValue', value: '60s' } }) } }) : undefined;

      const uninstallRequest = new CapabilityProto.types.UninstallRequest({ cluster: new Reference({ id: this.cluster.id }), options });
      const capabilityUninstallRequest = new Management.types.CapabilityUninstallRequest({
        name:   this.type,
        target: uninstallRequest
      });

      await Management.service.UninstallCapability(capabilityUninstallRequest);

      setTimeout(() => {
        this.updateCabilityLogs();
      }, 10);
    } catch (ex) {
      Vue.set(this, 'capabilityStatus', {
        state:        'error',
        shortMessage: 'Error',
        message:      exceptionToErrorsArray(ex).join('; '),
      });
    }
  }

  async cancelUninstall() {
    try {
      const cancelCapabilityUninstallRequest = new Management.types.CapabilityUninstallCancelRequest({
        name:    this.type,
        cluster: new Reference({ id: this.cluster.id })
      });

      await Management.service.CancelCapabilityUninstall(cancelCapabilityUninstallRequest);

      setTimeout(() => {
        this.updateCabilityLogs();
      }, 10);
    } catch (ex) {
      Vue.set(this, 'capabilityStatus', {
        state:        'error',
        shortMessage: 'Error',
        message:      exceptionToErrorsArray(ex).join('; '),
      });
    }
  }

  get isCapabilityUninstalling(): boolean {
    const cap = this.cluster.capabilitiesRaw?.find(c => c.name === this.type);

    return !!(cap?.deletionTimestamp);
  }

  get numSeries(): number | undefined {
    return this.clusterStats?.numSeries;
  }

  get sampleRate(): number | undefined {
    return Math.floor(this.clusterStats?.ingestionRate || 0);
  }

  get rulesRate(): number | undefined {
    return this.clusterStats?.RuleIngestionRate;
  }

  get clusterStats(): ClusterStats | undefined {
    return this.isInstalled ? this.stats?.find(s => s.userID === this.cluster.id) : undefined;
  }

  updateStats(stats: ClusterStats[]) {
    Vue.set(this, 'stats', stats);
  }
}
