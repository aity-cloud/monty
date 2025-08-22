package routing_test

import (
	"testing"

	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/test"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	_ "github.com/aity-cloud/monty/plugins/alerting/test"
	_ "github.com/aity-cloud/monty/plugins/metrics/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRouting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Routing Suite")
}

var (
	env          *test.Environment
	tmpConfigDir string
)

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env = &test.Environment{}
		Expect(env).NotTo(BeNil())
		Expect(env.Start(test.WithStorageBackend(v1beta1.StorageTypeJetStream), test.WithEnableNodeExporter(true))).To(Succeed())
		DeferCleanup(env.Stop, "Test Suite Finished")
		tmpConfigDir = env.GenerateNewTempDirectory("alertmanager-config")
		Expect(tmpConfigDir).NotTo(Equal(""))
	})
})
