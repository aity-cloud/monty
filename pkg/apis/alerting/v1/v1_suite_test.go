package v1_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alerting API Suite")
}
