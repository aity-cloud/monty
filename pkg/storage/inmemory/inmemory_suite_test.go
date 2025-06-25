package inmemory_test

import (
	"bytes"
	"testing"

	"github.com/aity-cloud/monty/pkg/storage"
	. "github.com/aity-cloud/monty/pkg/test/conformance/storage"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/util/future"

	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInmemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inmemory Suite")
}

type testBroker struct{}

func (t testBroker) KeyValueStore(string) storage.KeyValueStore {
	return inmemory.NewKeyValueStore(bytes.Clone)
}

var _ = Describe("In-memory KV Store", Ordered, Label("integration"), KeyValueStoreTestSuite(future.Instant(testBroker{}), NewBytes, Equal))
