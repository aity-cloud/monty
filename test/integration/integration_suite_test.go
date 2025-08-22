package integration_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	_ "github.com/aity-cloud/monty/plugins/example/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
