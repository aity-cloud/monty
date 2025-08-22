package kubernetes_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKubernetes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubernetes Ident Suite")
}
