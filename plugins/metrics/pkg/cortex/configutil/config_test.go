package configutil_test

import (
	"flag"

	"github.com/aity-cloud/monty/internal/cortex/config/querier"
	v1 "github.com/aity-cloud/monty/pkg/apis/storage/v1"
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/cortex/configutil"
	"github.com/cortexproject/cortex/pkg/cortex"
	"github.com/go-kit/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

func loadAndValidateConfig(yamlData []byte) error {
	GinkgoHelper()

	var cfg cortex.Config
	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg.RegisterFlags(flagset)
	Expect(yaml.UnmarshalStrict(yamlData, &cfg)).To(Succeed())
	Expect(flagset.Parse([]string{})).To(Succeed())

	return cfg.Validate(log.NewNopLogger())
}

var _ = Describe("Config", Label("unit"), func() {
	It("should generate a valid default config", func() {
		appconfig := &cortexops.CortexApplicationConfig{
			Storage: &v1.Config{
				Backend:    lo.ToPtr(v1.Filesystem),
				Filesystem: &v1.FilesystemConfig{},
			},
		}
		conf, _, err := configutil.CortexAPISpecToCortexConfig[*cortexops.CortexApplicationConfig](appconfig, configutil.NewTargetsOverride("all")...)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.Target).To(ConsistOf("all"))
		yamlData, err := configutil.MarshalCortexConfig(conf)
		Expect(err).NotTo(HaveOccurred())

		Expect(loadAndValidateConfig(yamlData)).To(Succeed())
	})
	It("should generate a valid config with overrides", func() {
		appconfig := &cortexops.CortexApplicationConfig{}

		conf, _, err := configutil.CortexAPISpecToCortexConfig[*cortexops.CortexApplicationConfig](appconfig,
			configutil.MergeOverrideLists(
				configutil.NewTargetsOverride("all"),
				configutil.NewHostOverrides(configutil.StandardOverridesShape{
					HttpListenAddress: "127.0.0.1",
					HttpListenNetwork: "tcp",
					HttpListenPort:    8080,
					GrpcListenNetwork: "tcp",
					GrpcListenAddress: "127.0.0.1",
					GrpcListenPort:    9095,
					StorageDir:        "/dev/null",
					RuntimeConfig:     "/dev/null",
				}),
			)...,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.Target).To(ConsistOf("all"))
		yamlData, err := configutil.MarshalCortexConfig(conf)
		Expect(err).NotTo(HaveOccurred())

		Expect(loadAndValidateConfig(yamlData)).To(Succeed())
	})
	It("should generate a valid config with extra overrides", func() {
		appconfig := &cortexops.CortexApplicationConfig{}

		conf, _, err := configutil.CortexAPISpecToCortexConfig[*cortexops.CortexApplicationConfig](appconfig,
			configutil.MergeOverrideLists(
				configutil.NewTargetsOverride(configutil.CortexTargets()...),
				configutil.NewHostOverrides(configutil.StandardOverridesShape{
					HttpListenAddress: "127.0.0.1",
					HttpListenNetwork: "tcp",
					HttpListenPort:    8080,
					GrpcListenNetwork: "tcp",
					GrpcListenAddress: "127.0.0.1",
					GrpcListenPort:    9095,
					StorageDir:        "/dev/null",
					RuntimeConfig:     "/dev/null",
				}),
				configutil.NewImplementationSpecificOverrides(configutil.ImplementationSpecificOverridesShape{
					QueryFrontendAddress: "http://localhost:9095",
					MemberlistJoinAddrs:  []string{"localhost"},
					AlertManager: configutil.AlertmanagerOverrideShape{
						AlertmanagerURL: "http://localhost:9093",
					},
				}),
				[]configutil.CortexConfigOverrider{
					configutil.NewOverrider(func(t *querier.Config) bool {
						t.MaxSubquerySteps = lo.ToPtr[int64](1000)
						return true
					}),
				},
			)...,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.Target).To(ConsistOf(configutil.CortexTargets()))
		yamlData, err := configutil.MarshalCortexConfig(conf)
		Expect(err).NotTo(HaveOccurred())

		Expect(loadAndValidateConfig(yamlData)).To(Succeed())
	})
	It("should marshal uint64 fields as numbers, not strings", func() {
		cc := cortex.Config{}
		err := configutil.LoadFromAPI(&cc.Querier, &querier.Config{
			MaxSubquerySteps: lo.ToPtr[int64](1000),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(cc.Querier.MaxSubQuerySteps).To(Equal(int64(1000)))
	})
})
