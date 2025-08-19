package alerting

import (
	"context"

	"log/slog"

	corev1 "github.com/aity-cloud/monty/apis/core/v1"
	corev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/resources"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	reconciler.ResourceReconciler
	ctx    context.Context
	client client.Client
	ac     *corev1beta1.AlertingCluster
	gw     *corev1.Gateway
	lg     *slog.Logger
}

func NewReconciler(
	ctx context.Context,
	client client.Client,
	instance *corev1beta1.AlertingCluster,
) *Reconciler {
	logger := logger.New().WithGroup("controller").WithGroup("alerting")
	return &Reconciler{
		ResourceReconciler: reconciler.NewReconcilerWith(
			client,
			reconciler.WithEnableRecreateWorkload(),
			reconciler.WithRecreateErrorMessageCondition(reconciler.MatchImmutableErrorMessages),
			reconciler.WithLog(log.FromContext(ctx)),
			reconciler.WithScheme(client.Scheme()),
		),
		ctx:    ctx,
		client: client,
		lg:     logger,
		ac:     instance,
	}
}

func (r *Reconciler) Reconcile() (reconcile.Result, error) {
	gw := &corev1.Gateway{}
	err := r.client.Get(r.ctx, client.ObjectKey{
		Name:      r.ac.Spec.Gateway.Name,
		Namespace: r.ac.Namespace,
	}, gw)
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	r.gw = gw

	if gw.DeletionTimestamp != nil {
		return k8sutil.DoNotRequeue().Result()
	}
	allResources := []resources.Resource{}

	alertingRsc, err := r.alerting()
	if err != nil {
		r.lg.Error(err.Error())
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, alertingRsc...)

	if op := resources.ReconcileAll(r, allResources); op.ShouldRequeue() {
		return op.Result()
	}
	return k8sutil.DoNotRequeue().Result()
}
