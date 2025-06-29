package collector

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	montyloggingv1beta1 "github.com/aity-cloud/monty/apis/logging/v1beta1"
	monitoringv1beta1 "github.com/aity-cloud/monty/apis/monitoring/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/otel"
	"github.com/aity-cloud/monty/pkg/resources"
	montymeta "github.com/aity-cloud/monty/pkg/util/meta"
	"github.com/samber/lo"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	receiversKey  = "receivers.yaml"
	mainKey       = "config.yaml"
	aggregatorKey = "aggregator.yaml"

	collectorImageRepo = "ghcr.io"
	collectorImage     = "rancher-sandbox/monty-otel-collector"
	collectorVersion   = "v0.1.4-rc1-0.85.0"
	reloaderImage      = "rancher-sandbox/config-reloader"
	reloaderVersion    = "v0.1.2"

	otelColBinaryName  = "otelcol-custom"
	otelConfigDir      = "/etc/otel"
	otelFileStorageDir = "/var/otel/filestorage"

	otlpGRPCPort    = int32(4317)
	rke2AgentLogDir = "/var/lib/rancher/rke2/agent/logs/"
	machineID       = "/etc/machine-id"
)

var directoryOrCreate = corev1.HostPathDirectoryOrCreate

func (r *Reconciler) agentConfigMapName() string {
	return fmt.Sprintf("%s-agent-config", r.collector.Name)
}

func (r *Reconciler) aggregatorConfigMapName() string {
	return fmt.Sprintf("%s-aggregator-config", r.collector.Name)
}

func (r *Reconciler) getDaemonConfig(loggingReceivers []string) otel.NodeConfig {
	return otel.NodeConfig{
		Instance: r.collector.Name,
		Logs: otel.LoggingConfig{
			Enabled:   r.collector.Spec.LoggingConfig != nil,
			Receivers: loggingReceivers,
		},
		Traces: otel.TraceConfig{
			Enabled: r.collector.Spec.TracesConfig != nil,
		},
		Metrics:       lo.FromPtr(r.getMetricsConfig()),
		Containerized: true,
		LogLevel:      r.collector.Spec.LogLevel,
		OTELConfig:    r.getDaemonOTELConfig(),
	}
}

func (r *Reconciler) getDaemonOTELConfig() otel.NodeOTELConfig {
	nodeOTELCfg := r.collector.Spec.NodeOTELConfigSpec
	if nodeOTELCfg == nil {
		r.lg.Warn("found no config for the daemon's OTEL Collector, falling back to default")
		nodeOTELCfg = montycorev1beta1.NewDefaultNodeOTELConfigSpec()
	}

	return otel.NodeOTELConfig{
		Processors: &otel.NodeOTELProcessors{
			MemoryLimiter: memorylimiterprocessor.Config{
				CheckInterval:         time.Duration(nodeOTELCfg.Processors.MemoryLimiter.CheckIntervalSeconds) * time.Second,
				MemoryLimitMiB:        nodeOTELCfg.Processors.MemoryLimiter.MemoryLimitMiB,
				MemorySpikeLimitMiB:   nodeOTELCfg.Processors.MemoryLimiter.MemorySpikeLimitMiB,
				MemoryLimitPercentage: nodeOTELCfg.Processors.MemoryLimiter.MemoryLimitPercentage,
				MemorySpikePercentage: nodeOTELCfg.Processors.MemoryLimiter.MemorySpikePercentage,
			},
		},
		Exporters: &otel.NodeOTELExporters{
			OTLP: otlpexporter.Config{
				QueueSettings: exporterhelper.QueueSettings{
					Enabled:      nodeOTELCfg.Exporters.OTLP.SendingQueue.Enabled,
					NumConsumers: nodeOTELCfg.Exporters.OTLP.SendingQueue.NumConsumers,
					QueueSize:    nodeOTELCfg.Exporters.OTLP.SendingQueue.QueueSize,
				},
			},
		},
	}
}

func (r *Reconciler) getAggregatorConfig(
	metricsCfg otel.MetricsConfig,
) otel.AggregatorConfig {
	return otel.AggregatorConfig{
		LogsEnabled:   r.collector.Spec.LoggingConfig != nil,
		TracesEnabled: r.collector.Spec.TracesConfig != nil,
		Metrics:       metricsCfg,
		AgentEndpoint: r.collector.Spec.AgentEndpoint,
		Containerized: true,
		LogLevel:      r.collector.Spec.LogLevel,
		OTELConfig:    r.getAggregatorOTELConfig(),
	}
}

func (r *Reconciler) getAggregatorOTELConfig() otel.AggregatorOTELConfig {
	aggregatorOTELCfg := r.collector.Spec.AggregatorOTELConfigSpec
	if aggregatorOTELCfg == nil {
		r.lg.Warn("found no config for the aggregator's OTEL Collector, falling back to default")
		aggregatorOTELCfg = montycorev1beta1.NewDefaultAggregatorOTELConfigSpec()
	}
	return otel.AggregatorOTELConfig{
		Processors: &otel.AggregatorOTELProcessors{
			MemoryLimiter: memorylimiterprocessor.Config{
				CheckInterval:         time.Duration(aggregatorOTELCfg.Processors.MemoryLimiter.CheckIntervalSeconds) * time.Second,
				MemoryLimitMiB:        aggregatorOTELCfg.Processors.MemoryLimiter.MemoryLimitMiB,
				MemorySpikeLimitMiB:   aggregatorOTELCfg.Processors.MemoryLimiter.MemorySpikeLimitMiB,
				MemoryLimitPercentage: aggregatorOTELCfg.Processors.MemoryLimiter.MemoryLimitPercentage,
				MemorySpikePercentage: aggregatorOTELCfg.Processors.MemoryLimiter.MemorySpikePercentage,
			},
			Batch: batchprocessor.Config{
				Timeout:          time.Duration(aggregatorOTELCfg.Processors.Batch.TimeoutSeconds) * time.Second,
				SendBatchSize:    aggregatorOTELCfg.Processors.Batch.SendBatchSize,
				SendBatchMaxSize: aggregatorOTELCfg.Processors.Batch.SendBatchMaxSize,
			},
		},
		Exporters: &otel.AggregatorOTELExporters{
			OTLPHTTP: otlphttpexporter.Config{
				QueueSettings: exporterhelper.QueueSettings{
					Enabled:      aggregatorOTELCfg.Exporters.OTLPHTTP.SendingQueue.Enabled,
					NumConsumers: aggregatorOTELCfg.Exporters.OTLPHTTP.SendingQueue.NumConsumers,
					QueueSize:    aggregatorOTELCfg.Exporters.OTLPHTTP.SendingQueue.QueueSize,
				},
			},
		},
	}

}

func (r *Reconciler) receiverConfig() (retData []byte, retReceivers []string, retErr error) {
	if r.collector.Spec.LoggingConfig != nil {
		retData = append(retData, []byte(templateLogAgentK8sReceiver)...)
		retReceivers = append(retReceivers, logReceiverK8s)

		config, err := r.fetchLoggingCollectorConfig()
		if err != nil {
			retErr = err
			return
		}

		auditLogsReceiver, data, err := r.generateKubeAuditLogsReceiver(config)
		if err != nil {
			retErr = err
			return
		}
		retData = append(retData, data...)
		if len(auditLogsReceiver) > 0 {
			retReceivers = append(retReceivers, auditLogsReceiver)
		}

		receiver, data, err := r.generateDistributionReceiver(config)
		if err != nil {
			retErr = err
			return
		}
		retData = append(retData, data...)
		if len(receiver) > 0 {
			retReceivers = append(retReceivers, receiver...)
		}
	}
	if r.collector.Spec.MetricsConfig != nil {
		cfg := r.getDaemonConfig([]string{})
		data, err := r.metricsNodeReceiverConfig(cfg)
		if err != nil {
			retErr = err
			return
		}
		retData = append(retData, data...)
	}

	return
}

func (r *Reconciler) mainConfig(loggingReceivers []string) ([]byte, error) {
	config := r.getDaemonConfig(loggingReceivers)

	var buffer bytes.Buffer
	t, err := r.tmpl.Parse(templateMainConfig)
	if err != nil {
		return buffer.Bytes(), err
	}
	err = t.Execute(&buffer, config)
	return buffer.Bytes(), err
}

func (r *Reconciler) agentConfigMap() (resources.Resource, string) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.agentConfigMapName(),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.PartOfLabel: "monty",
			},
		},
		Data: map[string]string{},
	}

	receiverData, logReceivers, err := r.receiverConfig()
	if err != nil {
		r.lg.Error("error", logger.Err(err))
		return resources.Error(cm, err), ""
	}
	cm.Data[receiversKey] = string(receiverData)

	mainData, err := r.mainConfig(logReceivers)
	if err != nil {
		r.lg.Error("error", logger.Err(err))
		return resources.Error(cm, err), ""
	}
	cm.Data[mainKey] = string(mainData)

	combinedData := append(receiverData, mainData...)
	hash := sha256.New()
	hash.Write(combinedData)
	configHash := hex.EncodeToString(hash.Sum(nil))

	if r.collector.Spec.IsEmpty() {
		return resources.Absent(cm), ""
	}

	ctrl.SetControllerReference(r.collector, cm, r.client.Scheme())
	return resources.Present(cm), configHash
}

func (r *Reconciler) aggregatorConfigMap(curCfg otel.AggregatorConfig) (resources.Resource, string) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.aggregatorConfigMapName(),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.PartOfLabel: "monty",
			},
		},
		Data: map[string]string{},
	}

	var buffer bytes.Buffer
	t, err := r.tmpl.Parse(templateAggregatorConfig)
	if err != nil {
		return resources.Error(nil, err), ""
	}
	err = t.Execute(&buffer, curCfg)
	if err != nil {
		r.lg.Error("error", logger.Err(err))
		return resources.Error(nil, err), ""
	}
	config := buffer.Bytes()
	configStr := string(config)
	cm.Data[aggregatorKey] = configStr

	hash := sha256.New()
	hash.Write(config)
	configHash := hex.EncodeToString(hash.Sum(nil))

	if r.collector.Spec.IsEmpty() {
		return resources.Absent(cm), ""
	}

	ctrl.SetControllerReference(r.collector, cm, r.client.Scheme())
	return resources.Present(cm), configHash
}

func (r *Reconciler) daemonSet() resources.Resource {
	imageSpec := r.imageSpec()
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "collector-config",
			MountPath: otelConfigDir,
		},
		{
			Name:      "filestorage-extension",
			MountPath: otelFileStorageDir,
		},
	}
	volumes := []corev1.Volume{
		{
			Name: "collector-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: r.agentConfigMapName(),
					},
				},
			},
		},
		{
			Name: "filestorage-extension",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: otelFileStorageDir,
					Type: &directoryOrCreate,
				},
			},
		},
	}

	// Short circuit and remove the daemonset if we don't need it
	if !r.shouldDeployDaemonset() {
		ds := &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-collector-agent", r.collector.Name),
				Namespace: r.collector.Spec.SystemNamespace,
				Labels: map[string]string{
					resources.AppNameLabel:  "collector-agent",
					resources.PartOfLabel:   "monty",
					resources.InstanceLabel: r.collector.Name,
				},
			},
		}
		return resources.Absent(ds)
	}

	if r.collector.Spec.LoggingConfig != nil {
		hostVolumeMounts, hostVolumes, err := r.hostLoggingVolumes()
		if err != nil {
			return resources.Error(nil, err)
		}
		volumeMounts = append(volumeMounts, hostVolumeMounts...)
		volumes = append(volumes, hostVolumes...)
	}
	if r.collector.Spec.MetricsConfig != nil {
		hostVolumeMounts, hostVolumes, err := r.hostMetricsVolumes()
		if err != nil {
			return resources.Error(nil, err)
		}
		volumeMounts = append(volumeMounts, hostVolumeMounts...)
		volumes = append(volumes, hostVolumes...)
	}
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-collector-agent", r.collector.Name),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.AppNameLabel:  "collector-agent",
				resources.PartOfLabel:   "monty",
				resources.InstanceLabel: r.collector.Name,
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					resources.AppNameLabel:  "collector-agent",
					resources.PartOfLabel:   "monty",
					resources.InstanceLabel: r.collector.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						resources.AppNameLabel:  "collector-agent",
						resources.PartOfLabel:   "monty",
						resources.InstanceLabel: r.collector.Name,
					},
				},
				Spec: corev1.PodSpec{
					ShareProcessNamespace: lo.ToPtr(true),
					Containers: []corev1.Container{
						{
							Name: "otel-collector",
							Command: []string{
								fmt.Sprintf("/%s", otelColBinaryName),
							},
							Args: []string{
								fmt.Sprintf("--config=/etc/otel/%s", mainKey),
							},
							Image:           *imageSpec.Image,
							ImagePullPolicy: imageSpec.GetImagePullPolicy(),
							SecurityContext: &corev1.SecurityContext{
								SELinuxOptions: &corev1.SELinuxOptions{
									Type: "rke_logreader_t",
								},
								RunAsUser: lo.ToPtr[int64](0),
							},
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							VolumeMounts: volumeMounts,
						},
						r.configReloaderContainer(volumeMounts, true),
					},
					ImagePullSecrets: imageSpec.ImagePullSecrets,
					Volumes:          volumes,
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/controlplane",
							Effect:   corev1.TaintEffectNoSchedule,
							Operator: corev1.TolerationOpExists,
						},
						{
							Key:      "node-role.kubernetes.io/control-plane",
							Effect:   corev1.TaintEffectNoSchedule,
							Operator: corev1.TolerationOpExists,
						},
						{
							Key:      "node-role.kubernetes.io/etcd",
							Effect:   corev1.TaintEffectNoExecute,
							Operator: corev1.TolerationOpExists,
						},
					},
					ServiceAccountName: "monty-agent",
				},
			},
		},
	}

	if r.collector.Spec.IsEmpty() {
		return resources.Absent(ds)
	}

	ctrl.SetControllerReference(r.collector, ds, r.client.Scheme())
	return resources.Present(ds)
}

func (r *Reconciler) shouldDeployDaemonset() bool {
	if r.collector.Spec.LoggingConfig != nil {
		return true
	}

	if r.collector.Spec.MetricsConfig == nil {
		return false
	}

	metrics := &monitoringv1beta1.CollectorConfig{}
	err := r.client.Get(r.ctx, client.ObjectKey{
		Name:      r.collector.Spec.MetricsConfig.Name,
		Namespace: r.collector.Namespace,
	}, metrics)
	if err != nil {
		return false
	}

	return lo.FromPtrOr(metrics.Spec.OtelSpec.HostMetrics, false)
}

func (r *Reconciler) deployment() resources.Resource {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "collector-config",
			MountPath: otelConfigDir,
		},
	}
	volumes := []corev1.Volume{
		{
			Name: "collector-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: r.aggregatorConfigMapName(),
					},
				},
			},
		},
	}

	if r.collector.Spec.MetricsConfig != nil {
		retVolumeMounts, retVolumes := r.aggregatorMetricVolumes()
		volumeMounts = append(volumeMounts, retVolumeMounts...)
		volumes = append(volumes, retVolumes...)
	}

	imageSpec := r.imageSpec()
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-collector-aggregator", r.collector.Name),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.AppNameLabel:  "collector-aggregator",
				resources.PartOfLabel:   "monty",
				resources.InstanceLabel: r.collector.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					resources.AppNameLabel:  "collector-aggregator",
					resources.PartOfLabel:   "monty",
					resources.InstanceLabel: r.collector.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						resources.AppNameLabel:  "collector-aggregator",
						resources.PartOfLabel:   "monty",
						resources.InstanceLabel: r.collector.Name,
					},
				},
				Spec: corev1.PodSpec{
					ShareProcessNamespace: lo.ToPtr(true),
					Containers: []corev1.Container{
						{
							Name: "otel-collector",
							Command: []string{
								"/otelcol-custom",
								fmt.Sprintf("--config=/etc/otel/%s", aggregatorKey),
							},
							Image:           *imageSpec.Image,
							ImagePullPolicy: imageSpec.GetImagePullPolicy(),
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							VolumeMounts: volumeMounts,
							Ports: []corev1.ContainerPort{
								{
									Name:          "otlp-grpc",
									ContainerPort: otlpGRPCPort,
								},
								{
									Name:          "pprof",
									ContainerPort: 1777,
								},
							},
						},
						r.configReloaderContainer(volumeMounts, false),
					},
					ImagePullSecrets:   imageSpec.ImagePullSecrets,
					Volumes:            volumes,
					ServiceAccountName: "monty-agent",
				},
			},
		},
	}

	if r.collector.Spec.IsEmpty() {
		return resources.Absent(deploy)
	}

	ctrl.SetControllerReference(r.collector, deploy, r.client.Scheme())
	return resources.Present(deploy)
}

func (r *Reconciler) service() resources.Resource {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-otel-aggregator", r.collector.Name),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.AppNameLabel:  "collector-aggregator",
				resources.PartOfLabel:   "monty",
				resources.InstanceLabel: r.collector.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "otlp-grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       otlpGRPCPort,
					TargetPort: intstr.FromInt(int(otlpGRPCPort)),
				},
			},
			Selector: map[string]string{
				resources.AppNameLabel:  "collector-aggregator",
				resources.PartOfLabel:   "monty",
				resources.InstanceLabel: r.collector.Name,
			},
		},
	}
	ctrl.SetControllerReference(r.collector, svc, r.client.Scheme())

	return resources.Present(svc)
}

func (r *Reconciler) configReloaderContainer(mounts []corev1.VolumeMount, runAsRoot bool) corev1.Container {
	reloaderImageSpec := r.configReloaderImageSpec()
	return corev1.Container{
		Name:            "config-reloader",
		Image:           *reloaderImageSpec.Image,
		ImagePullPolicy: reloaderImageSpec.GetImagePullPolicy(),
		Args: []string{
			"-volume-dir",
			otelConfigDir,
			"-process",
			otelColBinaryName,
		},
		SecurityContext: func() *corev1.SecurityContext {
			ctx := &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{
						"SYS_PTRACE",
					},
				},
				RunAsUser: lo.ToPtr[int64](10001),
			}
			if runAsRoot {
				ctx.RunAsUser = lo.ToPtr[int64](0)
			}
			return ctx
		}(),
		VolumeMounts: mounts,
	}
}

func (r *Reconciler) imageSpec() montymeta.ImageSpec {
	return montymeta.ImageResolver{
		Version:       collectorVersion,
		ImageName:     collectorImage,
		DefaultRepo:   collectorImageRepo,
		ImageOverride: &r.collector.Spec.ImageSpec,
	}.Resolve()
}

func (r *Reconciler) configReloaderImageSpec() montymeta.ImageSpec {
	return montymeta.ImageResolver{
		Version:     reloaderVersion,
		ImageName:   reloaderImage,
		DefaultRepo: collectorImageRepo,
		ImageOverride: func() *montymeta.ImageSpec {
			if r.collector.Spec.ConfigReloader != nil {
				return &r.collector.Spec.ConfigReloader.ImageSpec
			}
			return nil
		}(),
	}.Resolve()
}

func (r *Reconciler) fetchLoggingCollectorConfig() (*montyloggingv1beta1.CollectorConfig, error) {
	config := &montyloggingv1beta1.CollectorConfig{}
	err := r.client.Get(r.ctx, types.NamespacedName{
		Name:      r.collector.Spec.LoggingConfig.Name,
		Namespace: r.collector.Spec.SystemNamespace,
	}, config)

	return config, err
}
