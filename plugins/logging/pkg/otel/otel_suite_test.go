package otel_test

import (
	"testing"

	"github.com/aity-cloud/monty/pkg/test"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var env *test.Environment

func TestOtel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Otel Suite")
}

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env = &test.Environment{
			TestBin: "../../../../testbin/bin",
		}
		Expect(env).NotTo(BeNil())
	})
})
