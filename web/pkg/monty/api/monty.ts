import * as CortexOpsService from '@pkg/monty/generated/github.com/aity-cloud/monty/plugins/metrics/apis/cortexops/cortexops_svc';
import * as CortexOpsTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/plugins/metrics/apis/cortexops/cortexops_pb';
import * as StorageTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/internal/cortex/config/storage/storage_pb';
import * as ManagementService from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/apis/management/v1/management_svc';
import * as ManagementTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/apis/management/v1/management_pb';
import * as ConfigService from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/config/v1/config_server_svc';
import * as ConfigTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/config/v1/config_server_pb';
import * as CoreTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/apis/core/v1/core_pb';
import * as DriverUtilTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/plugins/driverutil/types_pb';
import * as CapabilityTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/pkg/apis/capability/v1/capability_pb';
import * as LoggingAdminService from '@pkg/monty/generated/github.com/aity-cloud/monty/plugins/logging/apis/loggingadmin/loggingadmin_svc';
import * as LoggingAdminTypes from '@pkg/monty/generated/github.com/aity-cloud/monty/plugins/logging/apis/loggingadmin/loggingadmin_pb';

export const CortexOps = {
  service: CortexOpsService,
  types:   CortexOpsTypes,
};

export const Management = {
  service: ManagementService,
  types:   ManagementTypes,
};

export const Config = {
  service: ConfigService,
  types:   ConfigTypes,
};

export const DriverUtil = { types: DriverUtilTypes };

export const Storage = { types: StorageTypes };

export const Core = { types: CoreTypes };

export const Capability = { types: CapabilityTypes };

export namespace LoggingAdmin {
  export import Service = LoggingAdminService;
  export import Types = LoggingAdminTypes;
}
