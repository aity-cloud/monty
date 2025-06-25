package stream_test

import (
	"fmt"
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStream(t *testing.T) {
	RegisterFailHandler(func(message string, callerSkip ...int) {
		fmt.Println(message)
		Fail(message, callerSkip...)
	})
	RunSpecs(t, "Stream API Extensions Suite")
}
