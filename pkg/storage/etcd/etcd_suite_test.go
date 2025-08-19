package etcd_test

import (
	"context"
	"testing"

	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/etcd"
	"github.com/aity-cloud/monty/pkg/test"
	. "github.com/aity-cloud/monty/pkg/test/conformance/storage"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testlog"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/aity-cloud/monty/pkg/util/future"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

func TestEtcd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcd Storage Suite")
}

var store = future.New[*etcd.EtcdStore]()

var lmF = future.New[storage.LockManager]()
var lmSet = future.New[lo.Tuple3[
	storage.LockManager, storage.LockManager, storage.LockManager,
]]()

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env := test.Environment{}
		env.Start(
			test.WithEnableGateway(false),
			test.WithEnableEtcd(true),
			test.WithEnableJetstream(false),
		)

		client, err := etcd.NewEtcdStore(context.Background(), env.EtcdConfig(),
			etcd.WithPrefix("test"),
		)
		Expect(err).To(Succeed())

		cli, err := etcd.NewEtcdClient(context.Background(), env.EtcdConfig(), testlog.Log)
		Expect(err).To(Succeed())

		lm := etcd.NewEtcdLockManager(cli, "test", logger.NewNop())
		lmF.Set(lm)

		x, err := etcd.NewEtcdClient(context.Background(), env.EtcdConfig(), testlog.Log)
		Expect(err).To(Succeed())
		lmX := etcd.NewEtcdLockManager(x, "test", logger.NewNop())

		y, err := etcd.NewEtcdClient(context.Background(), env.EtcdConfig(), testlog.Log)
		Expect(err).To(Succeed())
		lmY := etcd.NewEtcdLockManager(y, "test", logger.NewNop())
		Expect(err).To(Succeed())

		z, err := etcd.NewEtcdClient(context.Background(), env.EtcdConfig(), testlog.Log)
		Expect(err).To(Succeed())
		lmZ := etcd.NewEtcdLockManager(z, "test", logger.NewNop())

		lmSet.Set(lo.Tuple3[storage.LockManager, storage.LockManager, storage.LockManager]{
			A: lmX, B: lmY, C: lmZ,
		})

		Expect(err).NotTo(HaveOccurred())
		store.Set(client)
		Expect(err).To(Succeed())
		DeferCleanup(env.Stop, "Test Suite Finished")
	})
})

var _ = Describe("Etcd Token Store", Ordered, Label("integration", "slow"), TokenStoreTestSuite(store))
var _ = Describe("Etcd Cluster Store", Ordered, Label("integration", "slow"), ClusterStoreTestSuite(store))
var _ = Describe("Etcd RBAC Store", Ordered, Label("integration", "slow"), RBACStoreTestSuite(store))
var _ = Describe("Etcd Keyring Store", Ordered, Label("integration", "slow"), KeyringStoreTestSuite(store))
var _ = Describe("Etcd KV Store", Ordered, Label("integration", "slow"), KeyValueStoreTestSuite(store, NewBytes, Equal))
var _ = Describe("Etcd Lock Manager", Ordered, Label("integration", "slow"), LockManagerTestSuite(lmF, lmSet))
