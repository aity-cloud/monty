package crds_test

import (
	"context"
	"strings"
	"testing"
	"time"

	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	monitoringv1beta1 "github.com/aity-cloud/monty/apis/monitoring/v1beta1"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/crds"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	. "github.com/aity-cloud/monty/pkg/test/conformance/storage"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testk8s"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/aity-cloud/monty/pkg/test/testutil"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/future"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/pkg/util/protorand"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCrds(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(1 * time.Minute)
	RunSpecs(t, "CRDs Storage Suite")
}

var store = future.New[*crds.CRDStore]()
var kvStore = future.New[*broker]()

type methods struct{}

// ControllerReference implements crds.ValueStoreMethods.
func (methods) ControllerReference() (client.Object, bool) {
	return nil, false
}

type broker struct {
	k8sClient  client.WithWatch
	namespaces map[string]struct{}
}

func (b *broker) KeyValueStore(namespace string) storage.KeyValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
	return inmemory.NewCustomKeyValueStore(func(key string) storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
		if b.namespaces == nil {
			b.namespaces = make(map[string]struct{})
		}
		if _, ok := b.namespaces[namespace]; !ok {
			if err := b.k8sClient.Create(context.Background(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}); err != nil {
				panic(err)
			}
			b.namespaces[namespace] = struct{}{}
		}

		key = sanitizeKey(key)
		vs := crds.NewCRDValueStore(client.ObjectKey{
			Namespace: namespace,
			Name:      key,
		}, methods{}, crds.WithClient(b.k8sClient))
		return vs
	})
}

func sanitizeKey(key string) string {
	return strings.ReplaceAll(key, "/", "--")
}

// FillConfigFromObject implements crds.ValueStoreMethods.
func (methods) FillConfigFromObject(obj *montycorev1beta1.MonitoringCluster, conf *cortexops.CapabilityBackendConfigSpec) {
	conf.Enabled = obj.Spec.Cortex.Enabled
	conf.CortexConfig = obj.Spec.Cortex.CortexConfig
	conf.CortexWorkloads = obj.Spec.Cortex.CortexWorkloads
	conf.Grafana = obj.Spec.Grafana.GrafanaConfig
}

// FillObjectFromConfig implements crds.ValueStoreMethods.
func (methods) FillObjectFromConfig(obj *montycorev1beta1.MonitoringCluster, conf *cortexops.CapabilityBackendConfigSpec) {
	if conf == nil {
		obj.Spec.Cortex = montycorev1beta1.CortexSpec{}
		obj.Spec.Grafana = montycorev1beta1.GrafanaSpec{}
		return
	}
	obj.Spec.Cortex.Enabled = conf.Enabled
	obj.Spec.Cortex.CortexConfig = conf.CortexConfig
	obj.Spec.Cortex.CortexWorkloads = conf.CortexWorkloads
	obj.Spec.Grafana.GrafanaConfig = conf.Grafana

}

func newObject(seed ...int64) *cortexops.CapabilityBackendConfigSpec {
	if len(seed) == 0 {
		return nil
	}
	rand := protorand.New[*cortexops.CapabilityBackendConfigSpec]()
	rand.MaxCollectionElements = 1
	rand.Seed(seed[0])
	out := rand.MustGen()
	out.Revision = nil
	return out
}

var _ crds.ValueStoreMethods[*montycorev1beta1.MonitoringCluster, *cortexops.CapabilityBackendConfigSpec] = methods{}

var _ = BeforeSuite(func() {
	// sanity-check the valuestoremethods impl
	conf := protorand.New[*cortexops.CapabilityBackendConfigSpec]().MustGen()
	conf.Revision = nil

	var obj montycorev1beta1.MonitoringCluster
	methods{}.FillObjectFromConfig(&obj, conf.DeepCopy())

	conf2 := util.NewMessage[*cortexops.CapabilityBackendConfigSpec]()
	methods{}.FillConfigFromObject(obj.DeepCopy(), conf2)

	Expect(conf2).To(testutil.ProtoEqual(conf))

	testruntime.IfLabelFilterMatches(Label("integration", "slow"), func() {
		ctx, ca := context.WithCancel(context.Background())
		s := scheme.Scheme
		montycorev1beta1.AddToScheme(s)
		monitoringv1beta1.AddToScheme(s)
		config, _, err := testk8s.StartK8s(ctx, []string{"../../../config/crd/bases"}, s)
		Expect(err).NotTo(HaveOccurred())

		store.Set(crds.NewCRDStore(crds.WithRestConfig(config)))

		k8sClient, err := k8sutil.NewK8sClient(k8sutil.ClientOptions{
			RestConfig: config,
		})
		Expect(err).NotTo(HaveOccurred())

		kvStore.Set(&broker{k8sClient: k8sClient})

		DeferCleanup(ca)
	})
})

var _ = Describe("CRD Token Store", Ordered, Label("integration", "slow"), TokenStoreTestSuite(store))
var _ = Describe("CRD RBAC Store", Ordered, Label("integration", "slow"), RBACStoreTestSuite(store))
var _ = Describe("CRD Keyring Store", Ordered, Label("integration", "slow"), KeyringStoreTestSuite(store))
var _ = Describe("CRD Value Store", Ordered, Label("integration", "slow"), KeyValueStoreTestSuite(kvStore, newObject, func(a any) types.GomegaMatcher {
	return testutil.ProtoEqual(a.(proto.Message))
}))
