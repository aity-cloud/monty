package drivers

import (
	"context"
	"fmt"
	"os"

	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/crds"
	"github.com/aity-cloud/monty/pkg/util/flagutil"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/cortex/configutil"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MontyManagerClusterDriverOptions struct {
	K8sClient          client.WithWatch                                            `option:"k8sClient"`
	MonitoringCluster  types.NamespacedName                                        `option:"monitoringCluster"`
	GatewayRef         types.NamespacedName                                        `option:"gatewayRef"`
	DefaultConfigStore storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] `option:"defaultConfigStore"`
}

func (k MontyManagerClusterDriverOptions) newMonitoringCluster() *montycorev1beta1.MonitoringCluster {
	return &montycorev1beta1.MonitoringCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.MonitoringCluster.Name,
			Namespace: k.MonitoringCluster.Namespace,
		},
	}
}

func (k MontyManagerClusterDriverOptions) newGateway() *montycorev1beta1.Gateway {
	return &montycorev1beta1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.GatewayRef.Name,
			Namespace: k.GatewayRef.Namespace,
		},
	}
}

type MontyManager struct {
	cortexops.UnsafeCortexOpsServer
	MontyManagerClusterDriverOptions
	*driverutil.BaseConfigServer[*cortexops.ResetRequest, *cortexops.ConfigurationHistoryResponse, *cortexops.CapabilityBackendConfigSpec]
	configTracker *driverutil.DefaultingConfigTracker[*cortexops.CapabilityBackendConfigSpec]
}

type methods struct {
	controllerRef client.Object
}

// ControllerReference implements crds.ValueStoreMethods.
func (m methods) ControllerReference() (client.Object, bool) {
	return m.controllerRef, true
}

// FillConfigFromObject implements crds.ValueStoreMethods.
func (methods) FillConfigFromObject(obj *montycorev1beta1.MonitoringCluster, conf *cortexops.CapabilityBackendConfigSpec) {
	conf.Enabled = obj.Spec.Cortex.Enabled
	conf.CortexConfig = obj.Spec.Cortex.CortexConfig
	conf.CortexWorkloads = obj.Spec.Cortex.CortexWorkloads
	conf.Grafana = obj.Spec.Grafana.GrafanaConfig
}

// FillObjectFromConfig implements crds.ValueStoreMethods.
func (m methods) FillObjectFromConfig(obj *montycorev1beta1.MonitoringCluster, conf *cortexops.CapabilityBackendConfigSpec) {
	obj.Spec.Cortex.Enabled = conf.Enabled
	obj.Spec.Cortex.CortexConfig = conf.CortexConfig
	obj.Spec.Cortex.CortexWorkloads = conf.CortexWorkloads
	obj.Spec.Grafana.GrafanaConfig = conf.Grafana
	obj.Spec.Gateway.Name = m.controllerRef.GetName()
}

func NewMontyManagerClusterDriver(options MontyManagerClusterDriverOptions) (*MontyManager, error) {
	if options.K8sClient == nil {
		s := scheme.Scheme
		montycorev1beta1.AddToScheme(s)
		c, err := k8sutil.NewK8sClient(k8sutil.ClientOptions{
			Scheme: s,
			QPS:    50,
			Burst:  100,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
		}
		options.K8sClient = c
	}
	if options.DefaultConfigStore == nil {
		return nil, fmt.Errorf("missing required option: DefaultConfigStore")
	}

	gateway := options.newGateway()
	err := options.K8sClient.Get(context.TODO(), options.GatewayRef, gateway)
	if err != nil {
		return nil, err
	}
	activeStore := crds.NewCRDValueStore(options.MonitoringCluster, methods{
		controllerRef: gateway,
	}, crds.WithClient(options.K8sClient))

	configSrv := driverutil.NewBaseConfigServer[
		*cortexops.ResetRequest,
		*cortexops.ConfigurationHistoryResponse,
	](options.DefaultConfigStore, activeStore, flagutil.LoadDefaults)

	return &MontyManager{
		BaseConfigServer:                 configSrv,
		MontyManagerClusterDriverOptions: options,
		configTracker:                    configSrv.Tracker(),
	}, nil
}

// ListPresets implements cortexops.CortexOpsServer.
func (k *MontyManager) ListPresets(context.Context, *emptypb.Empty) (*cortexops.PresetList, error) {
	return &cortexops.PresetList{
		Items: []*cortexops.Preset{
			{
				Id: &corev1.Reference{Id: "all-in-one"},
				Metadata: &driverutil.PresetMetadata{
					DisplayName: "All In One",
					Description: "Minimal Cortex deployment with all components running in a single process",
					Notes: []string{
						"Warning: this configuration is not recommended for production use.",
					},
				},
				Spec: &cortexops.CapabilityBackendConfigSpec{
					CortexWorkloads: &cortexops.CortexWorkloadsConfig{
						Targets: map[string]*cortexops.CortexWorkloadSpec{
							"all": {Replicas: lo.ToPtr[int32](1)},
						},
					},
					CortexConfig: &cortexops.CortexApplicationConfig{
						LogLevel: lo.ToPtr("debug"),
					},
				},
			},
			{
				Id: &corev1.Reference{Id: "highly-available"},
				Metadata: &driverutil.PresetMetadata{
					DisplayName: "Highly Available",
					Description: "Basic HA Cortex deployment with all components running in separate processes",
					Notes: []string{
						"Additional storage configuration is required. Note that filesystem storage cannot be used in HA mode.",
						"Not all components are scaled to multiple replicas by default. The replica count for each component can be modified at any time.",
					},
				},
				Spec: &cortexops.CapabilityBackendConfigSpec{
					CortexWorkloads: &cortexops.CortexWorkloadsConfig{
						Targets: map[string]*cortexops.CortexWorkloadSpec{
							"distributor":    {Replicas: lo.ToPtr[int32](1)},
							"query-frontend": {Replicas: lo.ToPtr[int32](1)},
							"purger":         {Replicas: lo.ToPtr[int32](1)},
							"ruler":          {Replicas: lo.ToPtr[int32](3)},
							"compactor":      {Replicas: lo.ToPtr[int32](3)},
							"store-gateway":  {Replicas: lo.ToPtr[int32](3)},
							"ingester":       {Replicas: lo.ToPtr[int32](3)},
							"querier":        {Replicas: lo.ToPtr[int32](3)},
						},
					},
					CortexConfig: &cortexops.CortexApplicationConfig{
						LogLevel: lo.ToPtr("debug"),
					},
				},
			},
		},
	}, nil
}

// Status implements cortexops.CortexOpsServer.
func (k *MontyManager) Status(ctx context.Context, _ *emptypb.Empty) (*driverutil.InstallStatus, error) {
	status := &driverutil.InstallStatus{
		ConfigState:  driverutil.ConfigurationState_NotConfigured,
		InstallState: driverutil.InstallState_NotInstalled,
		AppState:     driverutil.ApplicationState_NotRunning,
		Metadata: map[string]string{
			"driver": "monty-manager",
		},
	}

	cluster := k.newMonitoringCluster()
	err := k.K8sClient.Get(ctx, k.MonitoringCluster, cluster)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get monitoring cluster: %w", err)
		}
	} else {
		status.ConfigState = driverutil.ConfigurationState_Configured
		if cluster.Spec.Cortex.Enabled != nil && *cluster.Spec.Cortex.Enabled {
			status.InstallState = driverutil.InstallState_Installed
		}
		mcStatus := cluster.Status.Cortex
		if err != nil {
			return nil, err
		}
		status.Version = mcStatus.Version
		if cluster.GetDeletionTimestamp() != nil {
			status.InstallState = driverutil.InstallState_Uninstalling
			status.AppState = driverutil.ApplicationState_Running
		} else {
			if mcStatus.WorkloadsReady {
				status.AppState = driverutil.ApplicationState_Running
			} else {
				status.AppState = driverutil.ApplicationState_Pending
				status.Warnings = append(status.Warnings, mcStatus.Conditions...)
			}
		}
	}

	return status, nil
}

func (k *MontyManager) ShouldDisableNode(_ *corev1.Reference) error {
	stat, err := k.Status(context.TODO(), &emptypb.Empty{})
	if err != nil {
		// can't determine cluster status, so don't disable the node
		return nil
	}
	switch stat.InstallState {
	case driverutil.InstallState_NotInstalled, driverutil.InstallState_Uninstalling:
		return status.Error(codes.Unavailable, fmt.Sprintf("Cortex cluster is not installed"))
	case driverutil.InstallState_Installed:
		return nil
	default:
		// can't determine cluster status, so don't disable the node
		return nil
	}
}

func (k *MontyManager) DryRun(ctx context.Context, req *cortexops.DryRunRequest) (*cortexops.DryRunResponse, error) {
	res, err := k.configTracker.DryRun(ctx, req)
	if err != nil {
		return nil, err
	}
	return &cortexops.DryRunResponse{
		Current:          res.Current,
		Modified:         res.Modified,
		ValidationErrors: configutil.ValidateConfiguration(res.Modified),
	}, nil
}

func init() {
	drivers.ClusterDrivers.Register("monty-manager", func(_ context.Context, opts ...driverutil.Option) (drivers.ClusterDriver, error) {
		options := MontyManagerClusterDriverOptions{
			MonitoringCluster: types.NamespacedName{
				Namespace: os.Getenv("POD_NAMESPACE"),
				Name:      "monty",
			},
			GatewayRef: types.NamespacedName{
				Namespace: os.Getenv("POD_NAMESPACE"),
				Name:      os.Getenv("GATEWAY_NAME"),
			},
		}
		if err := driverutil.ApplyOptions(&options, opts...); err != nil {
			return nil, err
		}

		return NewMontyManagerClusterDriver(options)
	})
}
