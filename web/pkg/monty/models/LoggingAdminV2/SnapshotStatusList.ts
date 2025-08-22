import { LoggingAdmin } from '@pkg/monty/api/monty';
import { SnapshotStatus } from './SnapshotStatus';

export class SnapshotStatusList {
    base: LoggingAdmin.Types.SnapshotStatusList;

    constructor(base: LoggingAdmin.Types.SnapshotStatusList) {
      this.base = base;
    }

    get statuses(): SnapshotStatus[] {
      return this.base.statuses.map(s => new SnapshotStatus(s));
    }
}
