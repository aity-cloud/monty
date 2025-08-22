package metrics_test

import (
	"testing"
	"time"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	_ "github.com/aity-cloud/monty/plugins/metrics/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	SetDefaultEventuallyTimeout(5 * time.Second)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Suite")
}
