package rules_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var ctrl *gomock.Controller

func TestRules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rules Suite")
}

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
})
