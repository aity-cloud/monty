package metrics_test

import (
	"context"
	"net"

	"github.com/aity-cloud/monty/apis"
	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/machinery"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/crds"
	"github.com/aity-cloud/monty/pkg/storage/kvutil"
	"github.com/aity-cloud/monty/pkg/test"
	conformance_driverutil "github.com/aity-cloud/monty/pkg/test/conformance/driverutil"
	"github.com/aity-cloud/monty/pkg/test/testk8s"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Config Storage Tests", Ordered, Label("integration"), func() {
	Context("Defaulting Config Tracker", func() {
		Context("CRD active store + Etcd default store", func() {
			var k8sClient client.WithWatch
			var kvstoreClient storage.KeyValueStoreT[*cortexops.CapabilityBackendConfigSpec]
			BeforeAll(func() {
				environment := &test.Environment{}
				Expect(environment.Start(test.WithEnableEtcd(true), test.WithEnableJetstream(false), test.WithStorageBackend(v1beta1.StorageTypeEtcd), test.WithEnableGateway(false))).To(Succeed())
				DeferCleanup(environment.Stop)
				restConfig, scheme, err := testk8s.StartK8s(environment.Context(), []string{"../../../config/crd/bases"}, apis.NewScheme())
				Expect(err).NotTo(HaveOccurred())
				k8sClient, err = k8sutil.NewK8sClient(k8sutil.ClientOptions{
					RestConfig: restConfig,
					Scheme:     scheme,
				})
				Expect(err).NotTo(HaveOccurred())

				backend, err := machinery.ConfigureStorageBackend(environment.Context(), &v1beta1.StorageSpec{
					Type: v1beta1.StorageTypeEtcd,
					Etcd: environment.EtcdConfig(),
				})
				Expect(err).NotTo(HaveOccurred())

				server := system.NewKVStoreServer(backend.KeyValueStore("test"))

				listener := bufconn.Listen(1024)
				srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
				system.RegisterKeyValueStoreServer(srv, server)
				go srv.Serve(listener)
				DeferCleanup(srv.Stop)

				conn, err := grpc.DialContext(environment.Context(), "bufconn", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return listener.Dial()
				}), grpc.WithTransportCredentials(insecure.NewCredentials()))
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(conn.Close)

				kvstoreClient = system.NewKVStoreClient[*cortexops.CapabilityBackendConfigSpec](system.NewKeyValueStoreClient(conn))
			})
			Context("Test Suite", conformance_driverutil.DefaultingConfigTrackerTestSuite[*cortexops.CapabilityBackendConfigSpec](
				func() storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
					return kvutil.WithKey(kvstoreClient, uuid.New().String())
				},
				func() storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
					ns := uuid.New().String()
					k8sClient.Create(context.Background(), &corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: ns,
						},
					})

					return crds.NewCRDValueStore(types.NamespacedName{
						Namespace: ns,
						Name:      "test",
					}, methods{}, crds.WithClient(k8sClient))
				},
			))
		})
		Context("Etcd active store + Etcd default store", func() {
			var kvstoreClient storage.KeyValueStoreT[*cortexops.CapabilityBackendConfigSpec]
			BeforeAll(func() {
				environment := &test.Environment{}
				Expect(environment.Start(test.WithEnableEtcd(true), test.WithEnableJetstream(false), test.WithStorageBackend(v1beta1.StorageTypeEtcd), test.WithEnableGateway(false))).To(Succeed())
				DeferCleanup(environment.Stop)

				backend, err := machinery.ConfigureStorageBackend(environment.Context(), &v1beta1.StorageSpec{
					Type: v1beta1.StorageTypeEtcd,
					Etcd: environment.EtcdConfig(),
				})
				Expect(err).NotTo(HaveOccurred())

				server := system.NewKVStoreServer(backend.KeyValueStore("test"))

				listener := bufconn.Listen(1024)
				srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
				system.RegisterKeyValueStoreServer(srv, server)
				go srv.Serve(listener)
				DeferCleanup(srv.Stop)

				conn, err := grpc.DialContext(environment.Context(), "bufconn", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return listener.Dial()
				}), grpc.WithTransportCredentials(insecure.NewCredentials()))
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(conn.Close)

				kvstoreClient = system.NewKVStoreClient[*cortexops.CapabilityBackendConfigSpec](system.NewKeyValueStoreClient(conn))
			})
			Context("Test Suite", conformance_driverutil.DefaultingConfigTrackerTestSuite[*cortexops.CapabilityBackendConfigSpec](
				func() storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
					return kvutil.WithKey(kvstoreClient, uuid.New().String())
				},
				func() storage.ValueStoreT[*cortexops.CapabilityBackendConfigSpec] {
					return kvutil.WithKey(kvstoreClient, uuid.New().String())
				},
			))
		})
	})
})

type methods struct{}

// ControllerReference implements crds.ValueStoreMethods.
func (m methods) ControllerReference() (client.Object, bool) {
	return nil, true
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
	obj.Spec.Cortex.Enabled = conf.Enabled
	obj.Spec.Cortex.CortexConfig = conf.CortexConfig
	obj.Spec.Cortex.CortexWorkloads = conf.CortexWorkloads
	obj.Spec.Grafana.GrafanaConfig = conf.Grafana
}
