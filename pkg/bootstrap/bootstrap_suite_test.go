package bootstrap_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	_ "github.com/aity-cloud/monty/plugins/example/test" // Required for incluster_test.go
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestBootstrap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bootstrap Suite")
}

var ctrl *gomock.Controller

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
})
