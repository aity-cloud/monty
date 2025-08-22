package alerting_manager_test

import (
	"crypto/tls"

	"github.com/aity-cloud/monty/pkg/alerting/client"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/alerting/drivers/alerting_manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("", Label("unit"), func() {
	When("We register the alering cluster driver", func() {
		It("should apply the tls config via driver options", func() {
			tlsConfig := &tls.Config{}
			opts := []driverutil.Option{
				driverutil.NewOption("tlsConfig", tlsConfig),
			}

			options := alerting_manager.AlertingDriverOptions{
				ConfigKey:          shared.AlertManagerConfigKey,
				InternalRoutingKey: shared.InternalRoutingConfigKey,
				Logger:             logger.NewPluginLogger().WithGroup("alerting").WithGroup("alerting-manager"),
			}
			driverutil.ApplyOptions(&options, opts...)
			Expect(options.TlsConfig).NotTo(BeNil())
		})

		It("should apply cluster driver subscribers via driver options", func() {
			subscriberA := make(chan client.AlertingClient)
			subscriberB := make(chan client.AlertingClient)
			opts := []driverutil.Option{
				driverutil.NewOption("subscribers", []chan client.AlertingClient{subscriberA, subscriberB}),
			}

			options := alerting_manager.AlertingDriverOptions{
				ConfigKey:          shared.AlertManagerConfigKey,
				InternalRoutingKey: shared.InternalRoutingConfigKey,
				Logger:             logger.NewPluginLogger().WithGroup("alerting").WithGroup("alerting-manager"),
			}
			driverutil.ApplyOptions(&options, opts...)
			Expect(options.Subscribers).To(HaveLen(2))
		})
	})
})
