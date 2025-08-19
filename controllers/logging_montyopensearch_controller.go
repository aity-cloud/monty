//go:build !minimal

package controllers

import (
	"context"

	opsterv1 "github.com/Opster/opensearch-k8s-operator/opensearch-operator/api/v1"
	loggingv1beta1 "github.com/aity-cloud/monty/apis/logging/v1beta1"
	"github.com/aity-cloud/monty/pkg/resources"
	"github.com/aity-cloud/monty/pkg/resources/montyopensearch"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LoggingMontyOpensearchReconciler struct {
	client.Client
	scheme *runtime.Scheme
	Opts   []montyopensearch.ReconcilerOption
}

// +kubebuilder:rbac:groups=logging.monty.io,resources=montyopensearches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.monty.io,resources=montyopensearches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=logging.monty.io,resources=montyopensearches/finalizers,verbs=update
// +kubebuilder:rbac:groups=opensearch.opster.io,resources=opensearchclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=opensearch.opster.io,resources=opensearchclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=opensearch.opster.io,resources=opensearchclusters/finalizers,verbs=update

func (r *LoggingMontyOpensearchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	montyOpensearch := &loggingv1beta1.MontyOpensearch{}
	err := r.Get(ctx, req.NamespacedName, montyOpensearch)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	montyOpensearchReconciler := montyopensearch.NewReconciler(ctx, montyOpensearch, r.Client, r.Opts...)

	reconcilers := []resources.ComponentReconciler{
		montyOpensearchReconciler.Reconcile,
	}

	for _, rec := range reconcilers {
		op := k8sutil.LoadResult(rec())
		if op.ShouldRequeue() {
			return op.Result()
		}
	}

	return k8sutil.DoNotRequeue().Result()
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoggingMontyOpensearchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.scheme = mgr.GetScheme()
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.MontyOpensearch{}).
		Owns(&opsterv1.OpenSearchCluster{}).
		Owns(&loggingv1beta1.MulticlusterRoleBinding{}).
		Complete(r)
}
