load('ext://min_k8s_version', 'min_k8s_version')
load('ext://helm_resource', 'helm_resource')
load('ext://namespace', 'namespace_create')

set_team('52cc75cc-c4ed-462f-8ea7-a543d398a381')

version = '0.12.1'
config.define_string_list('allowedContexts')
config.define_string_list('montyChartValues')
config.define_string('defaultRegistry')
config.define_string('valuesPath')
config.define_bool('buildCharts')

cfg = config.parse()

allow_k8s_contexts(cfg.get('allowedContexts'))

min_k8s_version('1.23')

namespace_create('monty')

update_settings (
  max_parallel_updates=1,
  k8s_upsert_timeout_secs=300,
)

ignore=[
  '**/*.pb.go',
  '**/*.pb.*.go',
  '**/*.swagger.json',
  'pkg/test/mock/*',
  'pkg/sdk/crd/*',
  '**/zz_generated.*',
  'packages/'
]

if cfg.get('buildCharts') == True:
    local_resource('build charts',
      deps='packages/**/templates',
      cmd='dagger run go run ./dagger --charts.git.export',
      ignore=ignore,
    )

k8s_yaml(helm('./charts/monty-crd/'+version,
  name='monty-crd',
  namespace='monty',
), allow_duplicates=True)

if cfg.get('valuesPath') != None:
  k8s_yaml(helm('./charts/monty/'+version,
    name='monty',
    namespace='monty',
    values=cfg.get('valuesPath')
  ), allow_duplicates=True)
else:
  k8s_yaml(helm('./charts/monty/'+version,
    name='monty',
    namespace='monty',
    set=cfg.get('montyChartValues')
  ), allow_duplicates=True)

if cfg.get('defaultRegistry') != None:
  default_registry(cfg.get('defaultRegistry'))

custom_build("registry.aity.tech/monty/monty",
  command="dagger run go run ./dagger --images.monty.push --images.monty.repo=registry.aity.tech/${EXPECTED_IMAGE} --images.monty.tag=${EXPECTED_TAG}",
  deps=['controllers', 'apis', 'pkg', 'plugins'],
  ignore=ignore,
  skips_local_docker=True,
)
