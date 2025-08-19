package management_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/management"
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	"github.com/aity-cloud/monty/pkg/test/freeport"
	mock_storage "github.com/aity-cloud/monty/pkg/test/mock/storage"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testdata"
	"github.com/aity-cloud/monty/pkg/test/testlog"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/pkg/util/flagutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Management Suite")
}

type testVars struct {
	ctrl           *gomock.Controller
	client         managementv1.ManagementClient
	grpcEndpoint   string
	httpEndpoint   string
	coreDataSource management.CoreDataSource
	storageBackend storage.Backend
	ifaces         struct {
		collector prometheus.Collector
	}
}

type testCoreDataSource struct {
	storageBackend storage.Backend
	tlsConfig      *tls.Config
}

func (t *testCoreDataSource) StorageBackend() storage.Backend {
	return t.storageBackend
}

func (t *testCoreDataSource) TLSConfig() *tls.Config {
	return t.tlsConfig
}

func setupManagementServer(vars **testVars, pl plugins.LoaderInterface, opts ...management.ManagementServerOption) func() {
	return func() {
		tv := &testVars{}
		if *vars != nil && (*vars).ctrl != nil {
			tv.ctrl = (*vars).ctrl
		} else {
			tv.ctrl = gomock.NewController(GinkgoT())
		}
		ctx, ca := context.WithCancel(context.Background())
		tv.storageBackend = mock_storage.NewTestStorageBackend(ctx, tv.ctrl)
		ports := freeport.GetFreePorts(2)

		defaultStore := inmemory.NewValueStore[*configv1.GatewayConfigSpec](util.ProtoClone)
		activeStore := inmemory.NewValueStore[*configv1.GatewayConfigSpec](util.ProtoClone)
		mgr := configv1.NewGatewayConfigManager(defaultStore, activeStore, flagutil.LoadDefaults)
		_, err := mgr.SetConfiguration(ctx, &configv1.SetRequest{
			Spec: &configv1.GatewayConfigSpec{
				Management: &configv1.ManagementServerSpec{
					GrpcListenAddress: lo.ToPtr(fmt.Sprintf("tcp://127.0.0.1:%d", ports[0])),
					HttpListenAddress: lo.ToPtr(fmt.Sprintf("127.0.0.1:%d", ports[1])),
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(mgr.Start(ctx)).To(Succeed())
		cert, err := tls.X509KeyPair(testdata.TestData("localhost.crt"), testdata.TestData("localhost.key"))
		Expect(err).NotTo(HaveOccurred())
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		cds := &testCoreDataSource{
			storageBackend: tv.storageBackend,
			tlsConfig:      tlsConfig,
		}
		server := management.NewServer(ctx, cds, mgr, pl, opts...)
		tv.coreDataSource = cds
		tv.ifaces.collector = server

		go func() {
			time.Sleep(100 * time.Millisecond) // TODO: can't hook into plugin loader here yet
			if err := server.ListenAndServe(ctx); err != nil {
				testlog.Log.Error("error", logger.Err(err))
			}
		}()

		cc, err := grpc.DialContext(ctx, fmt.Sprintf("127.0.0.1:%d", ports[0]), grpc.WithCredentialsBundle(insecure.NewBundle()), grpc.WithDefaultCallOptions(grpc.WaitForReady(true)))
		Expect(err).NotTo(HaveOccurred())
		tv.client = managementv1.NewManagementClient(cc)
		tv.grpcEndpoint = fmt.Sprintf("127.0.0.1:%d", ports[0])
		tv.httpEndpoint = fmt.Sprintf("http://127.0.0.1:%d", ports[1])
		*vars = tv
		DeferCleanup(ca)
	}
}
