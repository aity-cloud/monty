package client_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/aity-cloud/monty/pkg/alerting/client"
	"github.com/aity-cloud/monty/pkg/alerting/drivers/config"
	"github.com/aity-cloud/monty/pkg/alerting/drivers/routing"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/pkg/test/freeport"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	amCfg "github.com/prometheus/alertmanager/config"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

var (
	env  *test.Environment
	cl   client.AlertingClient
	clHA client.AlertingClient
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = BeforeSuite(func() {
	testruntime.IfIntegration(func() {
		env = &test.Environment{
			TestBin: "../../../testbin/bin",
		}
		Expect(env.Start()).To(Succeed())

		montyPort := freeport.GetFreePort()

		singleEmitterCfg := config.WebhookConfig{
			NotifierConfig: config.NotifierConfig{
				VSendResolved: false,
			},
			URL: &amCfg.URL{
				URL: util.Must(url.Parse(fmt.Sprintf("http://localhost:%d%s", montyPort, shared.AlertingDefaultHookName))),
			},
		}

		// set up a config file
		router := routing.NewMontyRouterV1(singleEmitterCfg)
		cfg, err := router.BuildConfig()
		Expect(err).To(Succeed())
		dir := env.GenerateNewTempDirectory("alertmanager_client")
		Expect(os.MkdirAll(dir, 0755)).To(Succeed())
		file, err := os.Create(fmt.Sprintf("%s/alertmanager.yaml", dir))
		Expect(err).To(Succeed())
		err = yaml.NewEncoder(file).Encode(cfg)
		Expect(err).To(Succeed())

		// start alertmanager
		ports := env.StartEmbeddedAlertManager(env.Context(), file.Name(), lo.ToPtr(montyPort))
		clA, err := client.NewClient(
			client.WithAlertManagerAddress(
				fmt.Sprintf("127.0.0.1:%d", ports.ApiPort),
			),
			client.WithQuerierAddress(
				fmt.Sprintf("127.0.0.1:%d", ports.EmbeddedPort),
			),
			client.WithTLSConfig(env.AlertingClientTLSConfig()),
		)
		Expect(err).To(Succeed())
		cl = clA

		msgPort := freeport.GetFreePort()
		haEmitterCfg := config.WebhookConfig{
			NotifierConfig: config.NotifierConfig{
				VSendResolved: false,
			},
			URL: &amCfg.URL{
				URL: util.Must(url.Parse(fmt.Sprintf("http://localhost:%d%s", msgPort, shared.AlertingDefaultHookName))),
			},
		}
		haRouter := routing.NewMontyRouterV1(haEmitterCfg)
		haCfg, err := haRouter.BuildConfig()
		Expect(err).To(Succeed())
		haFile, err := os.Create(fmt.Sprintf("%s/ha_alertmanager.yaml", dir))
		Expect(err).To(Succeed())
		err = yaml.NewEncoder(haFile).Encode(haCfg)
		Expect(err).To(Succeed())
		haPorts := env.StartEmbeddedAlertManager(env.Context(), haFile.Name(), lo.ToPtr(msgPort))

		replica1 := env.StartEmbeddedAlertManager(
			env.Context(),
			haFile.Name(),
			lo.ToPtr(freeport.GetFreePort()),
			fmt.Sprintf("127.0.0.1:%d", haPorts.ClusterPort),
		)
		replica2 := env.StartEmbeddedAlertManager(
			env.Context(),
			haFile.Name(),
			lo.ToPtr(freeport.GetFreePort()),
			fmt.Sprintf("127.0.0.1:%d", haPorts.ClusterPort),
			fmt.Sprintf("127.0.0.1:%d", replica1.ClusterPort),
		)

		clHA, err = client.NewClient(
			client.WithAlertManagerAddress(
				fmt.Sprintf("127.0.0.1:%d", ports.ApiPort),
			),
			client.WithQuerierAddress(
				fmt.Sprintf("127.0.0.1:%d", ports.EmbeddedPort),
			),
			client.WithTLSConfig(env.AlertingClientTLSConfig()),
		)
		Expect(err).To(Succeed())

		clHA.MemberlistClient().SetKnownPeers([]client.AlertingPeer{
			{
				ApiAddress:      fmt.Sprintf("127.0.0.1:%d", haPorts.ApiPort),
				EmbeddedAddress: fmt.Sprintf("127.0.0.1:%d", haPorts.EmbeddedPort),
			},
			{
				ApiAddress:      fmt.Sprintf("127.0.0.1:%d", replica2.ApiPort),
				EmbeddedAddress: fmt.Sprintf("127.0.0.1:%d", replica2.EmbeddedPort),
			},
			{
				ApiAddress:      fmt.Sprintf("127.0.0.1:%d", replica1.ApiPort),
				EmbeddedAddress: fmt.Sprintf("127.0.0.1:%d", replica1.EmbeddedPort),
			},
		})

		DeferCleanup(func() {
			Expect(env.Stop()).To(Succeed())
		})
	})

})
