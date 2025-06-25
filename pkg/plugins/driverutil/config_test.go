package driverutil_test

import (
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	conformance_driverutil "github.com/aity-cloud/monty/pkg/test/conformance/driverutil"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testdata/plugins/ext"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
)

func newValueStore() storage.ValueStoreT[*ext.SampleConfiguration] {
	return inmemory.NewValueStore[*ext.SampleConfiguration](util.ProtoClone)
}

var _ = Describe("Defaulting Config Tracker", Label("unit"), conformance_driverutil.DefaultingConfigTrackerTestSuite(newValueStore, newValueStore))
