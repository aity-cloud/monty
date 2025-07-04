package storage_test

import (
	"testing"

	"github.com/aity-cloud/monty/pkg/test"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAlertstorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alertstorage Suite")
}

var embeddedJetstream nats.JetStreamContext
var testKv nats.KeyValue
var testKv2 nats.KeyValue
var testKv3 nats.KeyValue
var testObj nats.ObjectStore
var env *test.Environment

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env = &test.Environment{}
		Expect(env.Start()).To(Succeed())
		nc, err := env.StartEmbeddedJetstream()
		Expect(err).NotTo(HaveOccurred())
		embeddedJetstream, err = nc.JetStream()
		Expect(err).NotTo(HaveOccurred())
		Expect(embeddedJetstream).NotTo(BeNil())
		kv, err := embeddedJetstream.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      "test-kv",
			Description: "bucket for testing",
		})
		Expect(err).To(Succeed())
		Expect(kv).NotTo(BeNil())
		kv2, err := embeddedJetstream.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      "test-kv2",
			Description: "bucket for testing",
		})
		Expect(err).To(Succeed())
		Expect(kv2).NotTo(BeNil())
		kv3, err := embeddedJetstream.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      "test-kv3",
			Description: "bucket for testing",
		})
		Expect(err).To(Succeed())
		Expect(kv3).NotTo(BeNil())

		obj, err := embeddedJetstream.CreateObjectStore(&nats.ObjectStoreConfig{
			Bucket:      "test-obj",
			Description: "bucket for testing",
		})
		Expect(err).To(Succeed())
		Expect(obj).NotTo(BeNil())

		testKv = kv
		testKv2 = kv2
		testKv3 = kv3
		testObj = obj

		DeferCleanup(func() {
			embeddedJetstream.DeleteKeyValue("test-kv")
			embeddedJetstream.DeleteKeyValue("test-kv2")
			embeddedJetstream.DeleteKeyValue("test-kv3")
			embeddedJetstream.DeleteObjectStore("test-obj")
		})
	})
})
