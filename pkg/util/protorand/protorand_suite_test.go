package protorand_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProtorand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protorand Suite")
}
