package montycluster

import (
	aiv1beta1 "github.com/aity-cloud/monty/apis/ai/v1beta1"
	"github.com/aity-cloud/monty/pkg/resources"
	montymeta "github.com/aity-cloud/monty/pkg/util/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (r *Reconciler) serviceLabels(service aiv1beta1.ServiceKind) map[string]string {
	return map[string]string{
		resources.AppNameLabel: service.ServiceName(),
		resources.ServiceLabel: service.String(),
		resources.PartOfLabel:  "monty",
	}
}

func (r *Reconciler) natsLabels() map[string]string {
	return map[string]string{
		resources.AppNameLabel:     "nats",
		resources.PartOfLabel:      "monty",
		resources.MontyClusterName: r.montyCluster.Name,
	}
}

func (r *Reconciler) pretrainedModelLabels(modelName string) map[string]string {
	return map[string]string{
		resources.PretrainedModelLabel: modelName,
	}
}

func (r *Reconciler) serviceImageSpec(service aiv1beta1.ServiceKind) montymeta.ImageSpec {
	return montymeta.ImageResolver{
		Version:             r.montyCluster.Spec.Version,
		ImageName:           service.ImageName(),
		DefaultRepo:         "registry.aity.tech/monty",
		DefaultRepoOverride: r.montyCluster.Spec.DefaultRepo,
		ImageOverride:       service.GetImageSpec(r.montyCluster),
	}.Resolve()
}

func (r *Reconciler) serviceNodeSelector(service aiv1beta1.ServiceKind) map[string]string {
	if s := service.GetNodeSelector(r.montyCluster); len(s) > 0 {
		return s
	}
	return r.montyCluster.Spec.GlobalNodeSelector
}

func (r *Reconciler) serviceTolerations(service aiv1beta1.ServiceKind) []corev1.Toleration {
	return append(r.montyCluster.Spec.GlobalTolerations, service.GetTolerations(r.montyCluster)...)
}

func addCPUInferenceLabel(deployment *appsv1.Deployment) {
	deployment.Labels[resources.MontyInferenceType] = "cpu"
	deployment.Spec.Template.Labels[resources.MontyInferenceType] = "cpu"
	deployment.Spec.Selector.MatchLabels[resources.MontyInferenceType] = "cpu"
}
