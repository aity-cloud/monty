package jetstream_test

import (
	"context"
	"testing"

	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/jetstream"
	"github.com/aity-cloud/monty/pkg/test"
	. "github.com/aity-cloud/monty/pkg/test/conformance/storage"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/aity-cloud/monty/pkg/util/future"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJetStream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JetStream Storage Suite")
}

var store = future.New[*jetstream.JetStreamStore]()
var lmF = future.New[*jetstream.LockManager]()

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env := test.Environment{}
		env.Start(
			test.WithEnableGateway(false),
			test.WithEnableEtcd(false),
			test.WithEnableJetstream(true),
		)

		s, err := jetstream.NewJetStreamStore(context.Background(), env.JetStreamConfig())
		Expect(err).NotTo(HaveOccurred())
		store.Set(s)

		lm, err := jetstream.NewJetstreamLockManager(context.Background(), env.JetStreamConfig())
		Expect(err).NotTo(HaveOccurred())
		lmF.Set(lm)

		DeferCleanup(env.Stop, "Test Suite Finished")
	})
})

var _ = Describe("Jetstream Token Store", Ordered, Label("integration", "slow"), TokenStoreTestSuite(store))
var _ = Describe("Jetstream Cluster Store", Ordered, Label("integration", "slow"), ClusterStoreTestSuite(store))
var _ = Describe("Jetstream RBAC Store", Ordered, Label("integration", "slow"), RBACStoreTestSuite(store))
var _ = Describe("Jetstream Keyring Store", Ordered, Label("integration", "slow"), KeyringStoreTestSuite(store))
var _ = Describe("Jetstream KV Store", Ordered, Label("integration", "slow"), KeyValueStoreTestSuite(store, NewBytes, Equal))
var _ = Describe("Jetstream Lock Manager", Ordered, Label("integration", "slow"), LockManagerTestSuite(lmF))

var _ = Context("Error Codes", func() {
	Specify("Nats KeyNotFound errors should be equal to ErrNotFound", func() {
		Expect(storage.IsNotFound(jetstream.JetstreamGrpcError(nats.ErrKeyNotFound))).To(BeTrue())
	})
	Specify("Nats KeyExists errors should be equal to ErrAlreadyExists", func() {
		Expect(storage.IsAlreadyExists(jetstream.JetstreamGrpcError(nats.ErrKeyExists))).To(BeTrue())
	})
	Specify("Nats KV sequence errors shoud be comparable with IsConflict", func() {
		Expect(storage.IsConflict(jetstream.JetstreamGrpcError(&nats.APIError{
			Code:        400,
			ErrorCode:   nats.JSErrCodeStreamWrongLastSequence,
			Description: "wrong last sequence: 1234",
		}))).To(BeTrue())
	})
})
