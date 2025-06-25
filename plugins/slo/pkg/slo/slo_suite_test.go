package slo_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	_ "github.com/aity-cloud/monty/plugins/alerting/test"
	_ "github.com/aity-cloud/monty/plugins/metrics/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSlo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slo Suite")
}
