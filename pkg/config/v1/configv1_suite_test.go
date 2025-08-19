package configv1_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config V1 Suite")
}
