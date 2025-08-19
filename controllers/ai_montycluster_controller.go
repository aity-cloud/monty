//go:build !minimal

/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	opsterv1 "github.com/Opster/opensearch-k8s-operator/opensearch-operator/api/v1"
	aiv1beta1 "github.com/aity-cloud/monty/apis/ai/v1beta1"
	"github.com/aity-cloud/monty/pkg/resources"
	"github.com/aity-cloud/monty/pkg/resources/montycluster"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
)

// MontyClusterReconciler reconciles a MontyCluster object
type AIMontyClusterReconciler struct {
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	Opts     []montycluster.ReconcilerOption
}

// +kubebuilder:rbac:groups=ai.monty.io,resources=montyclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.monty.io,resources=montyclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.monty.io,resources=montyclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheuses,verbs=get;list;watch
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

// Required for insights service
// +kubebuilder:rbac:groups=core,resources=namespaces;endpoints;pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;replicasets;daemonsets;statefulsets,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs;cronjobs,verbs=get;list;watch

func (r *AIMontyClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	montyCluster := &aiv1beta1.MontyCluster{}
	err := r.Get(ctx, req.NamespacedName, montyCluster)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.Opts = append(r.Opts, montycluster.WithResourceOptions(
		reconciler.WithEnableRecreateWorkload(),
		reconciler.WithScheme(r.scheme),
	))

	montyReconciler := montycluster.NewReconciler(ctx, r, r.recorder, montyCluster, r.Opts...)

	reconcilers := []resources.ComponentReconciler{
		montyReconciler.Reconcile,
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
func (r *AIMontyClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.scheme = mgr.GetScheme()
	r.recorder = mgr.GetEventRecorderFor("monty-controller")
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1beta1.MontyCluster{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&opsterv1.OpenSearchCluster{}).
		Complete(r)
}
