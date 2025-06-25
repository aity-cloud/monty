package main

import (
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func LoadConfig(ctx *pulumi.Context) *Config {
	var clusterConfig ClusterConfig
	namePrefix := config.Require(ctx, "monty:namePrefix")
	zoneID := config.Require(ctx, "monty:zoneID")
	useLocalCharts := config.GetBool(ctx, "monty:useLocalCharts")
	chartsRepo := config.Get(ctx, "monty:chartsRepo")
	chartVersion := config.Get(ctx, "monty:chartVersion")
	config.GetObject(ctx, "monty:cluster", &clusterConfig)
	tags := map[string]string{}
	config.GetObject(ctx, "monty:tags", &tags)
	useIdInDnsNames := config.GetBool(ctx, "monty:useIdInDnsNames")
	prometheusCrdChartMode := config.Get(ctx, "monty:prometheusCrdChartMode")
	disableKubePrometheusStack := config.GetBool(ctx, "monty:disableKubePrometheusStack")
	clusterConfig.LoadDefaults()

	var cloud, imageRepo, imageTag, minimalImageTag string

	if value, ok := os.LookupEnv("CLOUD"); ok {
		cloud = value
	} else {
		cloud = config.Get(ctx, "monty:cloud")
	}
	if value, ok := os.LookupEnv("IMAGE_REPO"); ok {
		imageRepo = value
	} else {
		imageRepo = config.Get(ctx, "monty:imageRepo")
	}
	if value, ok := os.LookupEnv("IMAGE_TAG"); ok {
		imageTag = value
	} else {
		imageTag = config.Get(ctx, "monty:imageTag")
	}
	if value, ok := os.LookupEnv("MINIMAL_IMAGE_TAG"); ok {
		minimalImageTag = value
	} else {
		minimalImageTag = config.Get(ctx, "monty:minimalImageTag")
	}

	conf := &Config{
		NamePrefix:                 namePrefix,
		ZoneID:                     zoneID,
		Cloud:                      cloud,
		ImageRepo:                  imageRepo,
		ImageTag:                   imageTag,
		MinimalImageTag:            minimalImageTag,
		UseLocalCharts:             useLocalCharts,
		ChartsRepo:                 chartsRepo,
		ChartVersion:               chartVersion,
		Cluster:                    clusterConfig,
		Tags:                       tags,
		UseIdInDnsNames:            useIdInDnsNames,
		DisableKubePrometheusStack: disableKubePrometheusStack,
		PrometheusCrdChartMode:     prometheusCrdChartMode,
	}
	conf.LoadDefaults()
	return conf
}

type Config struct {
	NamePrefix                 string            `json:"namePrefix"`
	ZoneID                     string            `json:"zoneID"`
	Cloud                      string            `json:"cloud"`
	ImageRepo                  string            `json:"imageRepo"`
	ImageTag                   string            `json:"imageTag"`
	MinimalImageTag            string            `json:"minimalImageTag"`
	UseLocalCharts             bool              `json:"useLocalCharts"`
	ChartsRepo                 string            `json:"chartsRepo"`
	ChartVersion               string            `json:"chartVersion"`
	Cluster                    ClusterConfig     `json:"cluster"`
	Tags                       map[string]string `json:"tags"`
	UseIdInDnsNames            bool              `json:"useIdInDnsNames"`
	DisableKubePrometheusStack bool              `json:"disableKubePrometheusStack"`

	// "separate" to deploy the monty-prometheus-crd chart separately
	// "embedded" to deploy the monty-prometheus-crd chart as a subchart of monty
	// "skip" to skip deploying the monty-prometheus-crd chart
	PrometheusCrdChartMode string `json:"prometheusCrdChartMode"`
}

type ClusterConfig struct {
	NodeInstanceType     string `json:"nodeInstanceType"`
	NodeGroupMinSize     int    `json:"nodeGroupMinSize"`
	NodeGroupMaxSize     int    `json:"nodeGroupMaxSize"`
	NodeGroupDesiredSize int    `json:"nodeGroupDesiredSize"`
}

func (c *Config) LoadDefaults() {
	if c.Cloud == "" {
		c.Cloud = "aws"
	}
	if c.ImageRepo == "" {
		c.ImageRepo = "rancher/monty"
	}
	if c.ImageTag == "" {
		c.ImageTag = "latest"
	}
	if c.MinimalImageTag == "" {
		c.MinimalImageTag = "latest-minimal"
	}
	if c.ChartsRepo == "" {
		c.ChartsRepo = "https://raw.githubusercontent.com/rancher/monty/charts-repo/"
	}
	if c.PrometheusCrdChartMode == "" {
		c.PrometheusCrdChartMode = "separate"
	}
	c.Cluster.LoadDefaults()
}

func (c *ClusterConfig) LoadDefaults() {
	if c.NodeInstanceType == "" {
		c.NodeInstanceType = "r6a.xlarge"
	}
	if c.NodeGroupMinSize == 0 {
		c.NodeGroupMinSize = 3
	}
	if c.NodeGroupMaxSize == 0 {
		c.NodeGroupMaxSize = 3
	}
	if c.NodeGroupDesiredSize == 0 {
		c.NodeGroupDesiredSize = 3
	}
}
