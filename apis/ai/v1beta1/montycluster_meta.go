package v1beta1

import (
	montymeta "github.com/aity-cloud/monty/pkg/util/meta"
	corev1 "k8s.io/api/core/v1"
)

type ServiceKind int

const (
	InferenceService ServiceKind = iota
	DrainService
	PreprocessingService
	PayloadReceiverService
	GPUControllerService
	MetricsService
	TrainingControllerService
	OpensearchUpdateService
)

type OpensearchRole string

const (
	OpensearchDataRole       OpensearchRole = "data"
	OpensearchClientRole     OpensearchRole = "client"
	OpensearchMasterRole     OpensearchRole = "master"
	OpensearchDashboardsRole OpensearchRole = "kibana"
)

func (s ServiceKind) String() string {
	switch s {
	case InferenceService:
		return "inference"
	case DrainService:
		return "drain"
	case PreprocessingService:
		return "preprocessing"
	case PayloadReceiverService:
		return "payload-receiver"
	case GPUControllerService:
		return "gpu-controller"
	case MetricsService:
		return "metrics"
	case OpensearchUpdateService:
		return "opensearch-update"
	case TrainingControllerService:
		return "training-controller"
	default:
		return ""
	}
}

func (s ServiceKind) ServiceName() string {
	return "monty-svc-" + s.String()
}

func (s ServiceKind) ImageName() string {
	switch s {
	case GPUControllerService:
		return "monty-gpu-service-controller"
	default:
		return "monty-" + s.String() + "-service"
	}
}

func (s ServiceKind) GetImageSpec(montyCluster *MontyCluster) *montymeta.ImageSpec {
	switch s {
	case InferenceService:
		return &montyCluster.Spec.Services.Inference.ImageSpec
	case DrainService:
		return &montyCluster.Spec.Services.Drain.ImageSpec
	case PreprocessingService:
		return &montyCluster.Spec.Services.Preprocessing.ImageSpec
	case PayloadReceiverService:
		return &montyCluster.Spec.Services.PayloadReceiver.ImageSpec
	case GPUControllerService:
		return &montyCluster.Spec.Services.GPUController.ImageSpec
	case MetricsService:
		return &montyCluster.Spec.Services.Metrics.ImageSpec
	case OpensearchUpdateService:
		return &montyCluster.Spec.Services.OpensearchUpdate.ImageSpec
	case TrainingControllerService:
		return &montyCluster.Spec.Services.TrainingController.ImageSpec
	default:
		return nil
	}
}

func (s ServiceKind) GetNodeSelector(montyCluster *MontyCluster) map[string]string {
	switch s {
	case InferenceService:
		return montyCluster.Spec.Services.Inference.NodeSelector
	case DrainService:
		return montyCluster.Spec.Services.Drain.NodeSelector
	case PreprocessingService:
		return montyCluster.Spec.Services.Preprocessing.NodeSelector
	case PayloadReceiverService:
		return montyCluster.Spec.Services.PayloadReceiver.NodeSelector
	case GPUControllerService:
		return montyCluster.Spec.Services.GPUController.NodeSelector
	case MetricsService:
		return montyCluster.Spec.Services.Metrics.NodeSelector
	case OpensearchUpdateService:
		return montyCluster.Spec.Services.OpensearchUpdate.NodeSelector
	case TrainingControllerService:
		return montyCluster.Spec.Services.TrainingController.NodeSelector
	default:
		return map[string]string{}
	}
}

func (s ServiceKind) GetTolerations(montyCluster *MontyCluster) []corev1.Toleration {
	switch s {
	case InferenceService:
		return montyCluster.Spec.Services.Inference.Tolerations
	case DrainService:
		return montyCluster.Spec.Services.Drain.Tolerations
	case PreprocessingService:
		return montyCluster.Spec.Services.Preprocessing.Tolerations
	case PayloadReceiverService:
		return montyCluster.Spec.Services.PayloadReceiver.Tolerations
	case GPUControllerService:
		return montyCluster.Spec.Services.GPUController.Tolerations
	case MetricsService:
		return montyCluster.Spec.Services.Metrics.Tolerations
	case OpensearchUpdateService:
		return montyCluster.Spec.Services.OpensearchUpdate.Tolerations
	case TrainingControllerService:
		return montyCluster.Spec.Services.TrainingController.Tolerations
	default:
		return []corev1.Toleration{}
	}
}

func (c *MontyCluster) GetState() string {
	return string(c.Status.State)
}

func (c *MontyCluster) GetConditions() []string {
	return c.Status.Conditions
}
