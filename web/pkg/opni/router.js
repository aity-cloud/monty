import { isStandalone } from './utils/standalone';
import MonitoringBackend from './pages/MonitoringBackend';
import Clusters from './pages/Clusters';
import Cluster from './pages/Cluster';
import Roles from './pages/Roles';
import Role from './pages/Role';
import RoleBindings from './pages/RoleBindings';
import RoleBinding from './pages/RoleBinding';
import Configuration from './pages/Configuration';
import LoggingBackend from './pages/LoggingBackend';
import SLOs from './pages/SLOs';
import SLO from './pages/SLO';
import WorkloadModelConfig from './pages/WorkloadModelConfig';

import Endpoints from './pages/Endpoints';
import Endpoint from './pages/Endpoint';
import Alarms from './pages/Alarms';
import Alarm from './pages/Alarm';
import AlertingOverview from './pages/AlertingOverview';
import AlertingBackend from './pages/AlertingBackend';
import LogAnomalyBackend from './pages/LogAnomalyBackend';
import PretrainedModels from './pages/PretrainedModels';
import Snapshots from './pages/Snapshots';
import Snapshot from './pages/Snapshot';
import Root from './pages/Root';

export const NAVIGATION = {
  routes: [
    {
      name:      isStandalone() ? 'index' : 'monty',
      path:      '/',
      component: Root,
      display:   false
    },
    {
      name:      'agents',
      path:      '/agents',
      labelKey:  'monty.nav.agents',
      component: Clusters,
      routes:    [
        {
          name:      'agent-create',
          path:      '/create',
          labelKey:  'monty.nav.clusters',
          component: Cluster,
          display:   false
        },
      ]
    },
    {
      name:      'logging-config',
      path:      '/logging',
      labelKey:  'monty.nav.loggingConfig',
      component: LoggingBackend,
      display:   true,
      routes:    [
        {
          name:      'snapshots',
          path:      '/snapshots',
          labelKey:  'monty.nav.snapshots',
          component: Snapshots,
          display:   true,
          routes:     [
            {
              name:      'snapshot-create',
              path:      '/create',
              component: Snapshot,
              display:   false,
            },
            {
              name:      'snapshot',
              path:      '/:id',
              component: Snapshot,
              display:   false
            },
          ]
        }
      ]
    },
    {
      name:      'monitoring',
      path:      '/monitoring',
      labelKey:  'monty.nav.monitoring',
      component: MonitoringBackend,
      routes:    [
        {
          name:     'rbac',
          path:     '/rbac',
          labelKey: 'monty.nav.rbac',
          display:   true,
          redirect: { name: 'roles' },
          routes:   [
            {
              name:      'roles',
              path:      '/roles',
              labelKey:  'monty.nav.roles',
              component: Roles,
              routes:    [
                {
                  name:      'role-create',
                  path:      '/create',
                  labelKey:  'monty.nav.roles',
                  component: Role,
                  display:   false
                },
              ]
            },
            {
              name:      'role-bindings',
              path:      '/role-bindings',
              labelKey:  'monty.nav.roleBindings',
              component: RoleBindings,
              routes:    [
                {
                  name:      'role-binding-create',
                  path:      '/create',
                  labelKey:  'monty.nav.roleBindings',
                  component: RoleBinding,
                  display:   false
                },
              ]
            },
          ]
        },
      ]
    },
    {
      name:      'alerting',
      path:      '/alerting',
      labelKey:  'monty.nav.alerting',
      display:   true,
      component: AlertingBackend,
      // redirect:  { name: 'slos' },
      routes:    [
        {
          name:      'alerting-overview',
          path:      '/overview',
          labelKey:  'monty.nav.alertingOverview',
          component: AlertingOverview,
          display:   true
        },
        {
          name:      'endpoints',
          path:      '/endpoints',
          labelKey:  'monty.nav.endpoints',
          component: Endpoints,
          display:   true,
          routes:    [
            {
              name:      'endpoint-create',
              path:      '/create',
              labelKey:  'monty.nav.endpoints',
              component: Endpoint,
              display:   false
            },
            {
              name:      'endpoint',
              path:      '/:id',
              labelKey:  'monty.nav.endpoints',
              component: Endpoint,
              display:   false
            },
          ]
        },
        {
          name:      'alarms',
          path:      '/alarms',
          labelKey:  'monty.nav.alarms',
          component: Alarms,
          display:   true,
          routes:    [
            {
              name:      'alarm-create',
              path:      '/create',
              labelKey:  'monty.nav.alarms',
              component: Alarm,
              display:   false
            },
            {
              name:      'alarm',
              path:      '/:id',
              labelKey:  'monty.nav.alarms',
              component: Alarm,
              display:   false
            },
          ]
        },
        {
          name:      'slos',
          path:      '/slos',
          labelKey:  'monty.nav.slos',
          component: SLOs,
          display:   true,
          routes:    [
            {
              name:      'slo',
              path:      '/:id',
              labelKey:  'monty.nav.slos',
              component: SLO,
              display:   false
            },
            {
              name:      'slo-create',
              path:      '/create',
              labelKey:  'monty.nav.slos',
              component: SLO,
              display:   false
            },
          ]
        },
      ]
    },
    {
      name:      'ai-ops',
      path:      '/ai-ops',
      labelKey:  'monty.nav.aiOps',
      display:   true,
      redirect:  { name: 'log-anomaly' },
      routes:    [
        {
          name:      'log-anomaly',
          path:      '/log-anomaly',
          labelKey:  'monty.nav.logAnomaly',
          display:   true,
          component: LogAnomalyBackend,
          routes:    [
            {
              name:      'pretrained-models',
              path:      '/pretrained-models',
              labelKey:  'monty.nav.pretrainedModels',
              display:   true,
              component: PretrainedModels,
            },
            {
              name:      'workload-model-config',
              path:      '/workload-insights',
              labelKey:  'monty.nav.workloadModel',
              component: WorkloadModelConfig,
              display:   true
            },
          ]
        },

      ]
    },
    {
      name:      'configuration',
      path:      '/configuration',
      labelKey:  'monty.nav.configuration',
      component: Configuration,
      display:   true
    },
  ]
};
