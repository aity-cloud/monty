package alerting_test

import (
	"testing"
	"time"

	"github.com/aity-cloud/monty/pkg/alerting/drivers/routing"
	"github.com/aity-cloud/monty/pkg/test"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	_ "github.com/aity-cloud/monty/plugins/alerting/test"
	_ "github.com/aity-cloud/monty/plugins/metrics/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/model"
	"github.com/samber/lo"
)

func init() {
	routing.DefaultConfig = routing.Config{
		GlobalConfig: routing.GlobalConfig{
			GroupWait:      lo.ToPtr(model.Duration(1 * time.Second)),
			RepeatInterval: lo.ToPtr(model.Duration(5 * time.Hour)),
		},
		SubtreeConfig: routing.SubtreeConfig{
			GroupWait:      lo.ToPtr(model.Duration(1 * time.Second)),
			RepeatInterval: lo.ToPtr(model.Duration(5 * time.Hour)),
		},
		FinalizerConfig: routing.FinalizerConfig{
			InitialDelay:       time.Second * 1,
			ThrottlingDuration: time.Minute * 1,
			RepeatInterval:     time.Hour * 5,
		},
		NotificationConfg: routing.NotificationConfg{},
	}
}

func TestAlerting(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alerting Suite")
}

var env *test.Environment
var tmpConfigDir string

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env = &test.Environment{}
		Expect(env).NotTo(BeNil())
		Expect(env.Start()).To(Succeed())
		DeferCleanup(env.Stop, "Test Suite Finished")
		tmpConfigDir = env.GenerateNewTempDirectory("alertmanager-config")
		Expect(tmpConfigDir).NotTo(Equal(""))
	})
})
