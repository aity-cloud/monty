package montycluster

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	aiv1beta1 "github.com/aity-cloud/monty/apis/ai/v1beta1"
	"github.com/aity-cloud/monty/pkg/opensearch/certs"
	"github.com/aity-cloud/monty/pkg/resources"
	"github.com/aity-cloud/monty/pkg/resources/montycluster/elastic/indices"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	opensearchv1 "opensearch.opster.io/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	prometheusRuleFinalizer = "monty.io/prometheusRule"
)

var (
	ErrOpensearchUpgradeFailed = errors.New("opensearch upgrade failed")
)

type Reconciler struct {
	reconciler.ResourceReconciler
	ReconcilerOptions
	ctx               context.Context
	client            client.Client
	recorder          record.EventRecorder
	montyCluster      *aiv1beta1.MontyCluster
	opensearchCluster *opensearchv1.OpenSearchCluster
}

type ReconcilerOptions struct {
	continueOnIndexError bool
	certMgr              certs.OpensearchCertReader
	resourceOptions      []reconciler.ResourceReconcilerOption
}

type ReconcilerOption func(*ReconcilerOptions)

func (o *ReconcilerOptions) apply(opts ...ReconcilerOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithContinueOnIndexError() ReconcilerOption {
	return func(o *ReconcilerOptions) {
		o.continueOnIndexError = true
	}
}

func WithCertManager(certMgr certs.OpensearchCertReader) ReconcilerOption {
	return func(o *ReconcilerOptions) {
		o.certMgr = certMgr
	}
}

func WithResourceOptions(opts ...reconciler.ResourceReconcilerOption) ReconcilerOption {
	return func(o *ReconcilerOptions) {
		o.resourceOptions = opts
	}
}

func NewReconciler(
	ctx context.Context,
	client client.Client,
	recorder record.EventRecorder,
	instance *aiv1beta1.MontyCluster,
	opts ...ReconcilerOption,
) *Reconciler {
	options := ReconcilerOptions{}
	options.apply(opts...)

	return &Reconciler{
		ReconcilerOptions: options,
		ResourceReconciler: reconciler.NewReconcilerWith(client,
			append(options.resourceOptions, reconciler.WithLog(log.FromContext(ctx)))...),
		ctx:          ctx,
		client:       client,
		recorder:     recorder,
		montyCluster: instance,
	}
}

func (r *Reconciler) Reconcile() (retResult *reconcile.Result, retErr error) {
	allResults := reconciler.CombinedResult{}
	result := reconciler.CombinedResult{}
	lg := log.FromContext(r.ctx)
	conditions := []string{}

	defer func() {
		// When the reconciler is done, figure out what the state of the montycluster
		// is and set it in the state field accordingly.
		if r.continueOnIndexError {
			retErr = result.Err
		}
		op := k8sutil.LoadResult(retResult, retErr)
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			if err := r.client.Get(r.ctx, client.ObjectKeyFromObject(r.montyCluster), r.montyCluster); err != nil {
				return err
			}
			r.montyCluster.Status.Conditions = conditions
			if op.ShouldRequeue() {
				if retErr != nil {
					// If an error occurred, the state should be set to error
					r.montyCluster.Status.State = aiv1beta1.MontyClusterStateError
				} else {
					// If no error occurred, but we need to requeue, the state should be
					// set to working
					r.montyCluster.Status.State = aiv1beta1.MontyClusterStateWorking
				}
			} else if len(r.montyCluster.Status.Conditions) == 0 {
				// If we are not requeueing and there are no conditions, the state should
				// be set to ready
				r.montyCluster.Status.State = aiv1beta1.MontyClusterStateReady
			}
			return r.client.Status().Update(r.ctx, r.montyCluster)
		})

		if err != nil {
			lg.Error(err, "failed to update status")
		}
	}()

	// Handle finalizer
	if r.montyCluster.DeletionTimestamp != nil && controllerutil.ContainsFinalizer(r.montyCluster, prometheusRuleFinalizer) {
		retResult, retErr = r.cleanupPrometheusRule()
		if retErr != nil || retResult != nil {
			return
		}
	}

	if r.montyCluster.Spec.Opensearch == nil {
		retErr = errors.New("opensearch cluster reference not specified")
		return
	}

	r.opensearchCluster = &opensearchv1.OpenSearchCluster{}
	err := r.client.Get(r.ctx, r.montyCluster.Spec.Opensearch.ObjectKeyFromRef(), r.opensearchCluster)
	if err != nil {
		retErr = errors.Wrap(err, "could not fetch opensearch object")
		return
	}

	allResults.Combine(r.reconcileIndices())
	if !r.continueOnIndexError {
		if allResults.Err != nil {
			return nil, allResults.Err
		}
	}

	allResources := []resources.Resource{}
	montyServices, err := r.montyServices()
	if err != nil {
		result.CombineErr(err)
		allResults.CombineErr(err)
		retErr = allResults.Err
		conditions = append(conditions, err.Error())
		return nil, err
	}
	pretrained, err := r.pretrainedModels()
	if err != nil {
		result.CombineErr(err)
		allResults.CombineErr(err)
		conditions = append(conditions, err.Error())
		lg.Error(err, "Error when reconciling pretrained models, will retry.")
		// Keep going, we can reconcile the rest of the deployments and come back
		// to this later.
	}

	workloadDrain, err := r.workloadDrain()
	if err != nil {
		retErr = errors.Combine(retErr, err)
		conditions = append(conditions, err.Error())
		lg.Error(err, "Error when reconciling workload drain, will retry.")
	}

	var es []resources.Resource
	es, err = r.externalOpensearchConfig()
	if err != nil {
		result.CombineErr(err)
		allResults.CombineErr(err)
		retErr = allResults.Err
		conditions = append(conditions, err.Error())
		lg.Error(err, "Error when setting external opensearch config")
		return
	}

	var s3 []resources.Resource
	intS3, err := r.internalS3()
	if err != nil {
		result.CombineErr(err)
		allResults.CombineErr(err)
		retErr = allResults.Err
		conditions = append(conditions, err.Error())
		lg.Error(err, "Error when reconciling s3, cannot continue.")
		return
	}
	extS3, err := r.externalS3()
	if err != nil {
		result.CombineErr(err)
		allResults.CombineErr(err)
		retErr = allResults.Err
		conditions = append(conditions, err.Error())
		lg.Error(err, "Error when reconciling s3, cannot continue.")
		return
	}
	s3 = append(s3, intS3...)
	s3 = append(s3, extS3...)

	// Order is important here
	// nats, s3, and elasticsearch reconcilers will add fields to the montyCluster status object
	// which are used by other reconcilers.
	allResources = append(allResources, s3...)
	allResources = append(allResources, es...)
	allResources = append(allResources, montyServices...)
	allResources = append(allResources, workloadDrain...)
	allResources = append(allResources, pretrained...)

	for _, factory := range allResources {
		o, state, err := factory()
		if err != nil {
			err = errors.WrapIf(err, "failed to create object")
			result.CombineErr(err)
			allResults.CombineErr(err)
			retErr = allResults.Err
			return
		}
		if o == nil {
			panic(fmt.Sprintf("reconciler %#v created a nil object", factory))
		}
		recResult, err := r.ReconcileResource(o, state)
		if err != nil {
			err = errors.WrapWithDetails(err, "failed to reconcile resource",
				"resource", o.GetObjectKind().GroupVersionKind())
			result.CombineErr(err)
			allResults.CombineErr(err)
			retErr = allResults.Err
		}
		if recResult != nil {
			allResults.Combine(recResult, err)
			retResult = &allResults.Result
		}
	}

	if allResults.Err != nil {
		return nil, result.Err
	}

	return
}

func (r *Reconciler) cleanupPrometheusRule() (retResult *reconcile.Result, retErr error) {
	namespace := r.montyCluster.Status.PrometheusRuleNamespace
	if namespace == "" {
		retErr = errors.New("prometheusRule namespace is unknown")
		return
	}
	prometheusRule := &monitoringv1.PrometheusRule{}
	err := r.client.Get(r.ctx, types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", aiv1beta1.MetricsService.ServiceName(), r.generateSHAID()),
		Namespace: namespace,
	}, prometheusRule)
	if k8serrors.IsNotFound(err) {
		retErr = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			if err := r.client.Get(r.ctx, client.ObjectKeyFromObject(r.montyCluster), r.montyCluster); err != nil {
				return err
			}
			controllerutil.RemoveFinalizer(r.montyCluster, prometheusRuleFinalizer)
			return r.client.Update(r.ctx, r.montyCluster)
		})
		return
	} else if err != nil {
		retErr = err
		return
	}
	retErr = r.client.Delete(r.ctx, prometheusRule)
	if retErr != nil {
		return
	}
	return &reconcile.Result{RequeueAfter: time.Second}, nil
}

func (r *Reconciler) reconcileIndices() (*reconcile.Result, error) {
	indicesReconciler, err := indices.NewReconciler(
		r.ctx,
		r.montyCluster,
		r.opensearchCluster,
		r.client,
		indices.WithCertManager(r.certMgr),
	)
	if err != nil {
		return nil, err
	}
	return indicesReconciler.Reconcile()

}

func RegisterWatches(builder *builder.Builder) *builder.Builder {
	return builder
}
