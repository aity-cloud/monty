package gateway

import (
	"context"

	"log/slog"

	corev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/resources"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	reconciler.ResourceReconciler
	ctx    context.Context
	client client.Client
	gw     *corev1beta1.Gateway
	lg     *slog.Logger
}

func NewReconciler(
	ctx context.Context,
	client client.Client,
	instance *corev1beta1.Gateway,
) *Reconciler {
	return &Reconciler{
		ResourceReconciler: reconciler.NewReconcilerWith(client,
			reconciler.WithEnableRecreateWorkload(),
			reconciler.WithRecreateErrorMessageCondition(reconciler.MatchImmutableErrorMessages),
			reconciler.WithLog(log.FromContext(ctx)),
			reconciler.WithScheme(client.Scheme()),
		),
		gw:     instance,
		ctx:    ctx,
		client: client,
		lg:     logger.New().WithGroup("controller").WithGroup("gateway"),
	}
}

func (r *Reconciler) Reconcile() (retResult reconcile.Result, retErr error) {
	updated, err := r.updateImageStatus()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	if updated {
		return k8sutil.Requeue().Result()
	}

	allResources := []resources.Resource{}
	etcdResources, err := r.etcd()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, etcdResources...)
	configMap, configDigest, err := r.configMap()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, configMap)
	certs, err := r.certs()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, certs...)
	keys, err := r.ephemeralKeys()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, keys...)
	deployment, err := r.deployment(map[string]string{
		resources.MontyConfigHash: configDigest,
	})
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, deployment...)
	services, err := r.services()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, services...)
	rbac, err := r.rbac()
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	allResources = append(allResources, r.amtoolConfigMap())
	allResources = append(allResources, rbac...)
	allResources = append(allResources, r.serviceMonitor())

	// allResources = append(allResources, r.alerting()...)

	if op := resources.ReconcileAll(r, allResources); op.ShouldRequeue() {
		return op.Result()
	}

	// Post initial reconcile we need to build the gateway secret for ingresses
	object, op := r.gatewayIngressSecret()
	if op != nil {
		return op.Result()
	}

	result, err := r.ReconcileResource(object, reconciler.StatePresent)
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}
	if result != nil {
		return *result, err
	}

	// Post-reconcile, wait for the public service's load balancer to be ready
	if op := r.waitForServiceEndpoints(); op.ShouldRequeue() {
		return op.Result()
	}
	if r.gw.Spec.ServiceType == corev1.ServiceTypeLoadBalancer {
		if op := r.waitForLoadBalancer(); op.ShouldRequeue() {
			return op.Result()
		}
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := r.client.Get(r.ctx, client.ObjectKeyFromObject(r.gw), r.gw)
		if err != nil {
			return err
		}
		r.gw.Status.Ready = true

		return r.client.Status().Update(r.ctx, r.gw)
	})
	if err != nil {
		return k8sutil.RequeueErr(err).Result()
	}

	return k8sutil.DoNotRequeue().Result()
}
