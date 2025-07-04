package benchmark_storage

import (
	"context"
	"testing"

	"github.com/aity-cloud/monty/pkg/storage/etcd"
	"github.com/aity-cloud/monty/pkg/storage/jetstream"
	"github.com/aity-cloud/monty/pkg/test"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/aity-cloud/monty/pkg/util/future"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var lmsEtcdF = future.New[[]*etcd.EtcdLockManager]()
var lmsJetstreamF = future.New[[]*jetstream.LockManager]()

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env := test.Environment{}
		env.Start(
			test.WithEnableGateway(false),
			test.WithEnableEtcd(true),
			test.WithEnableJetstream(false),
		)

		lmsE := make([]*etcd.EtcdLockManager, 7)
		for i := 0; i < 7; i++ {
			l, err := etcd.NewEtcdLockManager(context.Background(), env.EtcdConfig(),
				etcd.WithPrefix("test"),
			)
			Expect(err).NotTo(HaveOccurred())
			lmsE[i] = l
		}
		lmsJ := make([]*jetstream.LockManager, 7)
		for i := 0; i < 7; i++ {
			j, err := jetstream.NewJetstreamLockManager(context.Background(), env.JetStreamConfig())
			Expect(err).NotTo(HaveOccurred())
			lmsJ[i] = j
		}

		lmsEtcdF.Set(lmsE)
		lmsJetstreamF.Set(lmsJ)
		DeferCleanup(env.Stop, "Test Suite Finished")
	})
})

// Manually enable benchmarks
var _ = XDescribe("Etcd lock manager", Ordered, Serial, Label("integration", "slow"), LockManagerBenchmark("etcd", lmsEtcdF))
var _ = XDescribe("Jetstream lock manager", Ordered, Serial, Label("integration", "slow"), LockManagerBenchmark("jetstream", lmsJetstreamF))
