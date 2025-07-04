package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/aity-cloud/monty/apis"
	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	montyloggingv1beta1 "github.com/aity-cloud/monty/apis/logging/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/otel"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/pkg/util/k8sutil/provider"
	montymeta "github.com/aity-cloud/monty/pkg/util/meta"
	"github.com/aity-cloud/monty/plugins/logging/apis/node"
	"github.com/aity-cloud/monty/plugins/logging/pkg/agent/drivers"
	"github.com/cisco-open/k8s-objectmatcher/patch"
	"github.com/lestrrat-go/backoff/v2"
	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	secretName         = "monty-opensearch-auth"
	secretKey          = "password"
	dataPrepperName    = "monty-shipper"
	dataPrepperVersion = "2.0.1"
	clusterOutputName  = "monty-output"
	clusterFlowName    = "monty-flow"
	loggingConfigName  = "monty-logging"
	traceConfigName    = "monty-tracing"
)

type KubernetesManagerDriver struct {
	KubernetesManagerDriverOptions
	provider  string
	k8sClient client.Client

	state reconcilerState
}

type reconcilerState struct {
	sync.Mutex
	running       bool
	backoffCtx    context.Context
	backoffCancel context.CancelFunc
}

type KubernetesManagerDriverOptions struct {
	Namespace  string       `option:"namespace"`
	RestConfig *rest.Config `option:"restConfig"`
	Logger     *slog.Logger `option:"logger"`
}

func NewKubernetesManagerDriver(options KubernetesManagerDriverOptions) (*KubernetesManagerDriver, error) {
	scheme := apis.NewScheme()
	if options.RestConfig == nil {
		var err error
		options.RestConfig, err = k8sutil.NewRestConfig(k8sutil.ClientOptions{
			Scheme: scheme,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(options.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	provider, err := provider.DetectProvider(context.TODO(), clientset)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	k8sClient, err := client.New(options.RestConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &KubernetesManagerDriver{
		KubernetesManagerDriverOptions: options,
		k8sClient:                      k8sClient,
		provider:                       provider,
	}, nil
}

var _ drivers.LoggingNodeDriver = (*KubernetesManagerDriver)(nil)

func (m *KubernetesManagerDriver) ConfigureNode(config *node.LoggingCapabilityConfig) {
	m.state.Lock()
	if m.state.running {
		m.state.backoffCancel()
	}
	m.state.running = true
	ctx, ca := context.WithCancel(context.TODO())
	m.state.backoffCtx = ctx
	m.state.backoffCancel = ca
	m.state.Unlock()

	p := backoff.Exponential()
	b := p.Start(ctx)
	var success bool
BACKOFF:
	for backoff.Continue(b) {
		logCollectorConf := m.buildLoggingCollectorConfig()
		if err := m.reconcileObject(logCollectorConf, config.Enabled); err != nil {
			m.Logger.With(
				"object", client.ObjectKeyFromObject(logCollectorConf).String(),
				logger.Err(err),
			).Error("error reconciling object")
			continue BACKOFF
		}

		if err := m.reconcileCollector(config.Enabled); err != nil {
			m.Logger.With(
				"object", "monty collector",
				logger.Err(err),
			).Error("error reconciling object")
		}

		success = true
		break
	}

	if !success {
		m.Logger.Error("timed out reconciling objects")
	} else {
		m.Logger.Info("objects reconciled successfully")
	}
}

func (m *KubernetesManagerDriver) buildLoggingCollectorConfig() *montyloggingv1beta1.CollectorConfig {
	collectorConfig := &montyloggingv1beta1.CollectorConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "monty-logging",
		},
		Spec: montyloggingv1beta1.CollectorConfigSpec{
			Provider: montyloggingv1beta1.LogProvider(m.provider),
			KubeAuditLogs: &montyloggingv1beta1.KubeAuditLogsSpec{
				Enabled: false,
			},
		},
	}
	return collectorConfig
}

// TODO: make this generic
func (m *KubernetesManagerDriver) reconcileObject(desired client.Object, shouldExist bool) error {
	// get the object
	key := client.ObjectKeyFromObject(desired)
	lg := m.Logger.With("object", key)
	lg.Info("reconciling object")

	// get the agent statefulset
	agentStatefulSet, err := m.getAgentSet()
	if err != nil {
		return err
	}

	current := desired.DeepCopyObject().(client.Object)
	err = m.k8sClient.Get(context.TODO(), key, current)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	// this can error if the object is cluster-scoped, but that's ok
	controllerutil.SetOwnerReference(agentStatefulSet, desired, m.k8sClient.Scheme())

	if k8serrors.IsNotFound(err) {
		if !shouldExist {
			lg.Info("object does not exist and should not exist, skipping")
			return nil
		}
		lg.Info("object does not exist, creating")
		// create the object
		return m.k8sClient.Create(context.TODO(), desired)
	} else if !shouldExist {
		// delete the object
		lg.Info("object exists and should not exist, deleting")
		return m.k8sClient.Delete(context.TODO(), current)
	}

	// update the object
	return m.patchObject(current, desired)
}

func (m *KubernetesManagerDriver) reconcileCollector(shouldExist bool) error {
	coll := &montycorev1beta1.Collector{
		ObjectMeta: metav1.ObjectMeta{
			Name: otel.CollectorName,
		},
	}
	err := m.k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(coll), coll)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	if k8serrors.IsNotFound(err) && shouldExist {
		coll = m.buildEmptyCollector()
		coll.Spec.LoggingConfig = &corev1.LocalObjectReference{
			Name: loggingConfigName,
		}
		coll.Spec.TracesConfig = &corev1.LocalObjectReference{
			Name: traceConfigName,
		}

		return m.k8sClient.Create(context.TODO(), coll)
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := m.k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(coll), coll)
		if err != nil {
			return err
		}
		if shouldExist {
			coll.Spec.LoggingConfig = &corev1.LocalObjectReference{
				Name: loggingConfigName,
			}
			coll.Spec.TracesConfig = &corev1.LocalObjectReference{
				Name: traceConfigName,
			}
		} else {
			coll.Spec.LoggingConfig = nil
			coll.Spec.TracesConfig = nil
		}

		return m.k8sClient.Update(context.TODO(), coll)
	})
	return err
}

func (m *KubernetesManagerDriver) getAgentSet() (*appsv1.Deployment, error) {
	list := &appsv1.DeploymentList{}
	if err := m.k8sClient.List(context.TODO(), list,
		client.InNamespace(m.Namespace),
		client.MatchingLabels{
			"monty.io/app": "agent",
		},
	); err != nil {
		return nil, err
	}

	if len(list.Items) != 1 {
		return nil, errors.New("deployments found not exactly 1")
	}
	return &list.Items[0], nil
}

func (m *KubernetesManagerDriver) getAgentService() (*corev1.Service, error) {
	list := &corev1.ServiceList{}
	if err := m.k8sClient.List(context.TODO(), list,
		client.InNamespace(m.Namespace),
		client.MatchingLabels{
			"monty.io/app": "agent",
		},
	); err != nil {
		return nil, err
	}

	if len(list.Items) != 1 {
		return nil, errors.New("statefulsets found not exactly 1")
	}
	return &list.Items[0], nil
}

func (m *KubernetesManagerDriver) patchObject(current client.Object, desired client.Object) error {
	// update the object
	patchResult, err := patch.DefaultPatchMaker.Calculate(current, desired, patch.IgnoreStatusFields())
	if err != nil {
		m.Logger.With(
			logger.Err(err),
		).Warn("could not match objects")
		return err
	}
	if patchResult.IsEmpty() {
		m.Logger.Info("resource is in sync")
		return nil
	}
	m.Logger.Info("resource diff")

	if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(desired); err != nil {
		m.Logger.With(
			logger.Err(err),
		).Error("failed to set last applied annotation")
	}

	metaAccessor := meta.NewAccessor()

	currentResourceVersion, err := metaAccessor.ResourceVersion(current)
	if err != nil {
		return err
	}
	if err := metaAccessor.SetResourceVersion(desired, currentResourceVersion); err != nil {
		return err
	}

	m.Logger.Info("updating resource")

	return m.k8sClient.Update(context.TODO(), desired)
}

func (m *KubernetesManagerDriver) buildEmptyCollector() *montycorev1beta1.Collector {
	var serviceName string
	service, err := m.getAgentService()
	if err == nil {
		serviceName = service.Name
	}
	return &montycorev1beta1.Collector{
		ObjectMeta: metav1.ObjectMeta{
			Name: otel.CollectorName,
		},
		Spec: montycorev1beta1.CollectorSpec{
			ImageSpec: montymeta.ImageSpec{
				ImagePullPolicy: lo.ToPtr(corev1.PullAlways),
			},
			SystemNamespace:          m.Namespace,
			AgentEndpoint:            otel.AgentEndpoint(serviceName),
			AggregatorOTELConfigSpec: montycorev1beta1.NewDefaultAggregatorOTELConfigSpec(),
			NodeOTELConfigSpec:       montycorev1beta1.NewDefaultNodeOTELConfigSpec(),
		},
	}
}

func init() {
	drivers.NodeDrivers.Register("kubernetes-manager", func(_ context.Context, opts ...driverutil.Option) (drivers.LoggingNodeDriver, error) {
		options := KubernetesManagerDriverOptions{
			Namespace: os.Getenv("POD_NAMESPACE"),
			Logger:    logger.NewPluginLogger().WithGroup("logging").WithGroup("kubernetes-manager"),
		}
		driverutil.ApplyOptions(&options, opts...)
		return NewKubernetesManagerDriver(options)
	})
}
