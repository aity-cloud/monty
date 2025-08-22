package dryrun_test

import (
	"testing"

	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testdata/plugins/ext"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDryrun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DryRun Suite")
}

func newValueStore() storage.ValueStoreT[*ext.SampleConfiguration] {
	return inmemory.NewValueStore[*ext.SampleConfiguration](util.ProtoClone)
}

func newKeyValueStore() storage.KeyValueStoreT[*ext.SampleConfiguration] {
	return inmemory.NewKeyValueStore[*ext.SampleConfiguration](util.ProtoClone)
}
