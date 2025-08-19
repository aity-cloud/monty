package driverutil_test

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	conformance_driverutil "github.com/aity-cloud/monty/pkg/test/conformance/driverutil"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testdata/plugins/ext"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func newValueStore() storage.ValueStoreT[*ext.SampleConfiguration] {
	return inmemory.NewValueStore[*ext.SampleConfiguration](util.ProtoClone)
}

func newKeyValueStore() storage.KeyValueStoreT[*ext.SampleConfiguration] {
	return inmemory.NewKeyValueStore[*ext.SampleConfiguration](util.ProtoClone)
}

var _ = Describe("Defaulting Config Tracker", Label("unit"), conformance_driverutil.DefaultingConfigTrackerTestSuite(newValueStore, newValueStore))

type testContextKey struct {
	*corev1.ClusterStatus
}

func (t testContextKey) ContextKey() protoreflect.FieldDescriptor {
	return t.ProtoReflect().Descriptor().Fields().ByName("cluster")
}

var _ = Describe("Context Keys", func() {
	It("should correctly obtain context key values from ContextKeyable messages", func() {
		getReq := &ext.SampleGetRequest{
			Key: lo.ToPtr("foo"),
		}
		ctx := context.Background()
		ctx = driverutil.ContextWithKey(ctx, getReq)
		key := driverutil.KeyFromContext(ctx)

		Expect(key).To(Equal("foo"))
	})
	It("should correctly obtain context key values if the key field is a core.Reference", func() {
		testKeyable := &testContextKey{
			ClusterStatus: &corev1.ClusterStatus{
				Cluster: &corev1.Reference{
					Id: "foo",
				},
			},
		}
		ctx := context.Background()
		ctx = driverutil.ContextWithKey(ctx, testKeyable)
		key := driverutil.KeyFromContext(ctx)

		Expect(key).To(Equal("foo"))
	})
})
