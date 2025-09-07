//go:build !minimal && !cli

package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aity-cloud/monty/apis"
	montycorev1 "github.com/aity-cloud/monty/apis/core/v1"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/pkg/config/reactive"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/dashboard"
	"github.com/aity-cloud/monty/pkg/gateway"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/management"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/tracing"
	"github.com/aity-cloud/monty/pkg/update/noop"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/fieldmask"
	"github.com/aity-cloud/monty/pkg/util/flagutil"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/pkg/validation"
	"github.com/nsf/jsondiff"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/reflect/protopath"
	"k8s.io/apimachinery/pkg/types"

	_ "github.com/aity-cloud/monty/pkg/oci/kubernetes"
	_ "github.com/aity-cloud/monty/pkg/oci/noop"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/plugins/hooks"
	"github.com/aity-cloud/monty/pkg/storage/crds"
	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
	"github.com/aity-cloud/monty/pkg/storage/kvutil"
)

func BuildGatewayCmd() *cobra.Command {
	var inCluster bool
	host, hostOk := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	port, portOk := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	if hostOk && portOk && host != "" && port != "" {
		inCluster = true
	}

	var applyDefaultFlags bool
	var ignoreValidationErrors bool
	storageConfig := &configv1.StorageSpec{}
	config := &configv1.GatewayConfigSpec{}
	configDefaults := &configv1.GatewayConfigSpec{}
	flagutil.LoadDefaults(configDefaults)

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Run the Monty Gateway",
		Long: `
Flags and Configuration
===========================
Flags for this command are split into two categories: 'storage' flags and
'default' flags. Though the storage spec is contained within the complete
gateway configuration, but is separated for startup and initialization purposes,


Storage flags (--storage.*)
===========================
These flags configure the storage backend used to persist and retrieve
the active configuration. The values set by these flags will take precedence
over the corresponding values in the active configuration, if it exists.
* If running in Kubernetes, these flags are ignored.


Default flags (--default.*)
===========================
These flags allow adjusting the default values used when starting up the
gateway for the first time. They are all optional, and are ignored once
an active configuration has been created. However, if --apply-default-flags
is set, if there is an existing active configuration, any --default.* flags
listed on the command line will be applied to the active configuration before
starting the gateway.

* If running in Kubernetes, these flags are ignored, and --apply-default-flags
  has no effect.


Startup Logic
===========================
If the gateway is running inside a Kubernetes cluster, see the section below
for startup logic specific to Kubernetes. Regardless of runtime environment,
the runtime config APIs work the same way.

The gateway startup logic is as follows:

When the gateway starts up, it uses its storage flags (--storage.*) to connect
to a KV store and look for an active configuration.
- If there is no existing active config in the KV store, it will create one
  with the default values from its flags (--default.*).
- If there is an existing active config, it will use that. Additionally, before
  starting the gateway, if --apply-default-flags is set, it will apply any
  --default.* flags listed on the command line to the active configuration,
  overwriting any existing values, and persisting the changes to the KV store.


Startup Logic (Kubernetes)
===========================
When the gateway is running inside Kubernetes, it will look for the Gateway
custom resource that controls the running pod. Because this custom resource
controls the deployment of the gateway pod itself, it assumes it will always
exist. If it does not exist, the gateway will exit with an error.

When running in Kubernetes, the 'config' field in the Gateway custom resource
is the corresponding "active" configuration. Using the dashboard or CLI to
change these settings will update the custom resource; if possible, avoid
editing it directly.

In Kubernetes, all flags are ignored. Because the active configuration is
always present in the Gateway custom resource, it does not need to supply
its own default values or storage settings.


Runtime Config APIs
===========================
Once the gateway has started, the active configuration can be modified at
runtime using the dashboard UI or the 'monty config' CLI command.

Changes to the active configuration will be persisted to the KV store, and
relevant components will be notified of changed fields and will reload their
configuration accordingly (restarting servers, etc).

Changes to the default configuration, while possible, will not be persisted
across restarts. The gateway only stores its default configuration in-memory.
Note that similar configuration APIs for other Monty components generally do
persist their default configurations in the KV store.
`[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			lg := logger.New()
			ctx := cmd.Context()
			tracing.Configure("gateway")

			var storageBackend storage.Backend

			defaultStore := inmemory.NewValueStore[*configv1.GatewayConfigSpec](util.ProtoClone)

			if !inCluster {
				v, err := validation.NewValidator()
				if err != nil {
					return err
				}
				if err := v.Validate(config); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
					return errors.New("exiting due to validation errors")
				}
				defaultConfig := util.ProtoClone(config)
				defaultConfig.Storage = storageConfig
				if err := defaultStore.Put(ctx, defaultConfig); err != nil {
					return fmt.Errorf("failed to set defaults from flags: %w", err)
				}
			}
			var activeStore storage.ValueStoreT[*configv1.GatewayConfigSpec]

			if inCluster {
				lg.Info("loading config (in-cluster)")
				k8sClient, err := k8sutil.NewK8sClient(k8sutil.ClientOptions{
					Scheme: apis.NewScheme(),
				})
				if err != nil {
					return err
				}
				activeStore = crds.NewCRDValueStore[*montycorev1.Gateway, *configv1.GatewayConfigSpec](types.NamespacedName{
					Namespace: os.Getenv("POD_NAMESPACE"),
					Name:      os.Getenv("GATEWAY_NAME"),
				}, montycorev1.ValueStoreMethods{}, crds.WithClient(k8sClient))
				active, err := activeStore.Get(ctx)
				if err != nil {
					return err
				}
				lg := lg.With("backend",
					active.GetStorage().GetBackend(),
				)
				switch active.GetStorage().GetBackend() {
				case configv1.StorageBackend_Etcd:
					lg = lg.With("endpoints", active.GetStorage().GetEtcd().GetEndpoints())
				case configv1.StorageBackend_JetStream:
					lg = lg.With(
						"endpoint", active.GetStorage().GetJetStream().GetEndpoint(),
						"nkeySeedPath", active.GetStorage().GetJetStream().GetNkeySeedPath(),
					)
				}
				lg.Info("configuring storage backend")
				storageBackend, err = machinery.ConfigureStorageBackendV1(ctx, active.GetStorage())
				if err != nil {
					return err
				}
			} else {
				lg.Info("loading config")
				var err error
				storageBackend, err = machinery.ConfigureStorageBackendV1(ctx, storageConfig)
				if err != nil {
					return err
				}
				activeStore = kvutil.WithMessageCodec[*configv1.GatewayConfigSpec](
					kvutil.WithKey(storageBackend.KeyValueStore("gateway"), "config"))
			}
			lg.Info("storage configured", "backend", storageConfig.GetBackend().String())
			defer storageBackend.Close()

			mgr := configv1.NewGatewayConfigManager(
				defaultStore, activeStore,
				flagutil.LoadDefaults,
				configv1.WithControllerOptions(
					reactive.WithLogger(lg.WithGroup("config")),
					reactive.WithDiffMode(reactive.DiffFull),
				),
			)

			if !inCluster {
				var rev *corev1.Revision
				if ac, err := mgr.Tracker().ActiveStore().Get(context.Background()); err != nil {
					if storage.IsNotFound(err) {
						lg.Info("no previous configuration found, creating from defaults")
						_, err := mgr.SetConfiguration(context.Background(), &configv1.SetRequest{})
						if err != nil {
							return fmt.Errorf("failed to set configuration: %w", err)
						}
					}
				} else {
					rev = ac.GetRevision()
					lg.Info("loaded existing configuration", "rev", ac.GetRevision().GetRevision())
				}

				if applyDefaultFlags {
					updateMask := fieldmask.Leaves(fieldmask.Diff(configDefaults.ProtoReflect(), config.ProtoReflect()), config.ProtoReflect().Descriptor())
					resp, err := mgr.DryRun(ctx, &configv1.DryRunRequest{
						Target:   driverutil.Target_ActiveConfiguration,
						Action:   driverutil.Action_Reset,
						Revision: rev,
						Patch:    config,
						Mask:     updateMask,
					})
					if err != nil {
						return err
					}
					opts := jsondiff.DefaultConsoleOptions()
					opts.SkipMatches = true
					diff, anyChanges := driverutil.RenderJsonDiff(resp.Current, resp.Modified, opts)
					stat := driverutil.DiffStat(diff, opts)
					if anyChanges {
						lg.Info("applying default flags to active configuration", "diff", stat)
						lg.Info("⤷ diff:\n" + diff)
					} else {
						lg.Warn("--apply-default-flags was set, but no changes would be made to the active configuration")
					}
					if resp.GetValidationErrors() != nil {
						lg.Error("refusing to apply default flags due to validation errors (re-run with --ignore-validation-errors to skip this check)")
						return validation.ErrorsFromProto(resp.ValidationErrors)
					}
					if anyChanges {
						_, err := mgr.ResetConfiguration(ctx, &configv1.ResetRequest{
							Revision: rev,
							Mask:     updateMask,
							Patch:    config,
						})
						if err != nil {
							return err
						}
					}
				}
			}

			lg.Debug("starting config manager")
			if err := mgr.Start(ctx); err != nil {
				return fmt.Errorf("failed to start config manager: %w", err)
			}
			lg.Debug("config manager started")

			pluginLoader := plugins.NewPluginLoader(plugins.WithLogger(lg.WithGroup("gateway")))

			g := gateway.NewGateway(ctx, mgr, storageBackend, pluginLoader,
				gateway.WithExtraUpdateHandlers(noop.NewSyncServer()),
			)

			m := management.NewServer(ctx, g, mgr, pluginLoader,
				management.WithCapabilitiesDataSource(g.CapabilitiesDataSource()),
				management.WithHealthStatusDataSource(g),
			)

			g.MustRegisterCollector(m)

			d, err := dashboard.NewServer(ctx, mgr, pluginLoader, g, dashboard.WithLogger(lg.WithGroup("dashboard")))
			if err != nil {
				return err
			}

			go func() {
				ctx, ca := context.WithCancel(ctx)
				defer ca()
				w := reactive.Message[*configv1.PluginsSpec](mgr.Reactive(protopath.Path(config.ProtoPath().Plugins()))).Watch(ctx)
				pluginLoader.Hook(hooks.OnLoadingCompleted(func(numLoaded int) {
					lg.Info(fmt.Sprintf("loaded %d plugins", numLoaded))
				}))
				conf := <-w
				if conf == nil {
					lg.Error("no plugin configuration found")
					return
				}
				lg.Info("loaded plugin configuration", "dir", conf.GetDir())
				pluginLoader.LoadPlugins(ctx, conf.GetDir(), plugins.GatewayScheme, plugins.WithFilters(conf.GetFilters()))
			}()

			var eg errgroup.Group
			eg.Go(func() error {
				lg.Debug("starting gateway server")
				err := g.ListenAndServe(ctx)
				if errors.Is(err, context.Canceled) {
					lg.Info("gateway server stopped")
				} else if err != nil {
					lg.With(logger.Err(err)).Warn("gateway server exited with error")
				}
				return err
			})
			eg.Go(func() error {
				lg.Debug("starting management server")
				err := m.ListenAndServe(ctx)
				if errors.Is(err, context.Canceled) {
					lg.Info("management server stopped")
				} else if err != nil {
					lg.With(logger.Err(err)).Warn("management server exited with error")
				}
				return err
			})
			eg.Go(func() error {
				lg.Debug("starting dashboard server")
				err := d.ListenAndServe(ctx)
				if errors.Is(err, context.Canceled) {
					lg.Info("dashboard server stopped")
				} else if err != nil {
					lg.With(logger.Err(err)).Warn("dashboard server exited with error")
				}
				return err
			})

			lg.Debug("gateway startup complete")
			return eg.Wait()
		},
	}
	if !inCluster {
		cmd.Flags().AddFlagSet(storageConfig.FlagSet("storage"))
		cmd.Flags().AddFlagSet(config.FlagSet("defaults"))
		cmd.Flags().BoolVar(&applyDefaultFlags, "apply-default-flags", false,
			"Apply default flags listed on the command-line to the active configuration on startup")
		cmd.Flags().BoolVar(&ignoreValidationErrors, "ignore-validation-errors", false, "Ignore validation errors when applying default flags")
		cmd.Flags().MarkHidden("ignore-validation-errors")
	}
	return cmd
}

func init() {
	AddCommandsToGroup(MontyComponents, BuildGatewayCmd())
}
