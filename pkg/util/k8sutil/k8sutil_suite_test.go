package k8sutil_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestK8sutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sutil Suite")
}
