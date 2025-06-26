//go:build !noinfra

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aity-cloud/monty/infra/pkg/aws"
	"github.com/aity-cloud/monty/infra/pkg/resources"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	. "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/mod/semver"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
)

func main() {
	Run(run)
}

func run(ctx *Context) (runErr error) {
	defer func() {
		if sr, ok := runErr.(interface {
			StackTrace() errors.StackTrace
		}); ok {
			st := sr.StackTrace()
			ctx.Log.Error(fmt.Sprintf("%+v", st), &LogArgs{})
		}
	}()

	conf := LoadConfig(ctx)

	var provisioner resources.Provisioner

	switch conf.Cloud {
	case "aws":
		provisioner = aws.NewProvisioner()
	default:
		return errors.Errorf("unsupported cloud: %s", conf.Cloud)
	}

	var id StringOutput
	rand, err := random.NewRandomId(ctx, "id", &random.RandomIdArgs{
		ByteLength: Int(4),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	id = rand.Hex

	tags := map[string]string{}
	for k, v := range conf.Tags {
		tags[k] = v
	}
	tags["PulumiProjectName"] = ctx.Project()
	tags["PulumiStackName"] = ctx.Stack()

	mainCluster, err := provisioner.ProvisionMainCluster(ctx, resources.MainClusterConfig{
		ID:                   id,
		NamePrefix:           conf.NamePrefix,
		NodeInstanceType:     conf.Cluster.NodeInstanceType,
		NodeGroupMinSize:     conf.Cluster.NodeGroupMinSize,
		NodeGroupMaxSize:     conf.Cluster.NodeGroupMaxSize,
		NodeGroupDesiredSize: conf.Cluster.NodeGroupDesiredSize,
		ZoneID:               conf.ZoneID,
		Tags:                 tags,
		UseIdInDnsNames:      conf.UseIdInDnsNames,
	})
	if err != nil {
		return err
	}
	var montyCrdChart, montyPrometheusCrdChart, montyChart string
	var chartRepoOpts *helm.RepositoryOptsArgs
	if conf.UseLocalCharts {
		var ok bool
		if montyCrdChart, ok = findLocalChartDir("monty-crd"); !ok {
			return errors.New("could not find local monty-crd chart")
		}
		if montyPrometheusCrdChart, ok = findLocalChartDir("monty-prometheus-crd"); !ok {
			return errors.New("could not find local monty-prometheus-crd chart")
		}
		if montyChart, ok = findLocalChartDir("monty"); !ok {
			return errors.New("could not find local monty chart")
		}
	} else {
		chartRepoOpts = &helm.RepositoryOptsArgs{
			Repo: StringPtr(conf.ChartsRepo),
		}
		montyCrdChart = "monty-crd"
		montyPrometheusCrdChart = "monty-prometheus-crd"
		montyChart = "monty"
	}

	montyServiceLB := mainCluster.Provider.ApplyT(func(k *kubernetes.Provider) (StringOutput, error) {
		montyCrd, err := helm.NewRelease(ctx, "monty-crd", &helm.ReleaseArgs{
			Chart:          String(montyCrdChart),
			RepositoryOpts: chartRepoOpts,
			Version:        StringPtr(conf.ChartVersion),
			Namespace:      String("monty"),
			Timeout:        Int(60),
		}, Provider(k), RetainOnDelete(true))
		if err != nil {
			return StringOutput{}, errors.WithStack(err)
		}

		var montyChartExtraDeps []Resource

		if conf.PrometheusCrdChartMode == "separate" {
			montyPrometheusCrd, err := helm.NewRelease(ctx, "monty-prometheus-crd", &helm.ReleaseArgs{
				Chart:          String(montyPrometheusCrdChart),
				RepositoryOpts: chartRepoOpts,
				Namespace:      String("monty"),
				Timeout:        Int(60),
			}, Provider(k), RetainOnDelete(true))
			if err != nil {
				return StringOutput{}, errors.WithStack(err)
			}
			montyChartExtraDeps = append(montyChartExtraDeps, montyPrometheusCrd)
		}

		monty, err := helm.NewRelease(ctx, "monty", &helm.ReleaseArgs{
			Name:           String("monty"),
			Chart:          String(montyChart),
			RepositoryOpts: chartRepoOpts,
			Version:        StringPtr(conf.ChartVersion),
			SkipCrds:       Bool(true),
			Namespace:      String("monty"),
			Values: Map{
				"image": Map{
					"repository": String(conf.ImageRepo),
					"tag":        String(conf.ImageTag),
				},
				"monty-prometheus-crd": Map{
					"enabled": Bool(conf.PrometheusCrdChartMode == "embedded"),
				},
				"gateway": Map{
					"enabled": Bool(true),
					"alerting": Map{
						"enabled": Bool(true),
					},
					"hostname":    mainCluster.GatewayHostname,
					"serviceType": String("LoadBalancer"),
					"auth": Map{
						"provider": String("openid"),
						"openid": Map{
							"discovery": Map{
								"issuer": mainCluster.OAuth.Issuer,
							},
							"identifyingClaim":  String("email"),
							"clientID":          mainCluster.OAuth.ClientID,
							"clientSecret":      mainCluster.OAuth.ClientSecret,
							"scopes":            ToStringArray([]string{"openid", "profile", "email"}),
							"roleAttributePath": String(`'"custom:grafana_role"'`),
						},
					},
				},
				"monty-agent": Map{
					"image": Map{
						"repository": String(conf.ImageRepo),
						"tag":        String(conf.MinimalImageTag),
					},
					"enabled":          Bool(true),
					"address":          String("monty"),
					"fullnameOverride": String("monty-agent"),
					"bootstrapInCluster": Map{
						"enabled": Bool(true),
					},
					"agent": Map{
						"version": String("v2"),
					},
					"persistence": Map{
						"mode": String("pvc"),
					},
					"kube-prometheus-stack": Map{
						"enabled": Bool(!conf.DisableKubePrometheusStack),
					},
				},
			},
			Timeout:     Int(300),
			WaitForJobs: Bool(true),
		}, Provider(k), DependsOn(append([]Resource{montyCrd}, montyChartExtraDeps...)), RetainOnDelete(true))
		if err != nil {
			return StringOutput{}, errors.WithStack(err)
		}

		montyServiceLB := All(monty.Status.Namespace(), monty.Status.Name()).
			ApplyT(func(args []any) (StringOutput, error) {
				namespace := args[0].(*string)
				montyLBSvc, err := corev1.GetService(ctx, "monty", ID(
					fmt.Sprintf("%s/monty", *namespace),
				), nil, Provider(k), Parent(monty))
				if err != nil {
					return StringOutput{}, err
				}
				return montyLBSvc.Status.LoadBalancer().Ingress().Index(Int(0)).Hostname().Elem(), nil
			}).(StringOutput)
		return montyServiceLB, nil
	}).(StringOutput)

	_, err = provisioner.ProvisionDNSRecord(ctx, "gateway", resources.DNSRecordConfig{
		Name:    mainCluster.GatewayHostname,
		Type:    "CNAME",
		ZoneID:  conf.ZoneID,
		Records: StringArray{montyServiceLB},
		TTL:     60,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = provisioner.ProvisionDNSRecord(ctx, "grafana", resources.DNSRecordConfig{
		Name:    mainCluster.GrafanaHostname,
		Type:    "CNAME",
		ZoneID:  conf.ZoneID,
		Records: StringArray{mainCluster.LoadBalancerHostname},
		TTL:     60,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = provisioner.ProvisionDNSRecord(ctx, "opensearch", resources.DNSRecordConfig{
		Name:    mainCluster.OpensearchHostname,
		Type:    "CNAME",
		ZoneID:  conf.ZoneID,
		Records: StringArray{mainCluster.LoadBalancerHostname},
		TTL:     60,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	ctx.Export("kubeconfig", mainCluster.Kubeconfig.ApplyT(func(kubeconfig any) (string, error) {
		jsonData, err := json.Marshal(kubeconfig)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return string(jsonData), nil
	}).(StringOutput))
	ctx.Export("gateway_url", mainCluster.GatewayHostname)
	ctx.Export("grafana_url", mainCluster.GrafanaHostname.ApplyT(func(hostname string) string {
		return fmt.Sprintf("https://%s", hostname)
	}).(StringOutput))
	ctx.Export("opensearch_url", mainCluster.OpensearchHostname.ApplyT(func(hostname string) string {
		return fmt.Sprintf("https://%s", hostname)
	}).(StringOutput))
	ctx.Export("s3_bucket", mainCluster.S3.Bucket)
	ctx.Export("s3_endpoint", mainCluster.S3.Endpoint)
	ctx.Export("s3_region", mainCluster.S3.Region)
	ctx.Export("s3_access_key_id", mainCluster.S3.AccessKeyID)
	ctx.Export("s3_secret_access_key", mainCluster.S3.SecretAccessKey)
	ctx.Export("oauth_client_id", mainCluster.OAuth.ClientID)
	ctx.Export("oauth_client_secret", mainCluster.OAuth.ClientSecret)
	ctx.Export("oauth_issuer_url", mainCluster.OAuth.Issuer)
	return nil
}

func findLocalChartDir(chartName string) (string, bool) {
	// find charts from ../charts/<chartName> and return the latest version
	dir := fmt.Sprintf("../charts/%s", chartName)
	if _, err := os.Stat(dir); err != nil {
		return "", false
	}
	versions, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}
	if len(versions) == 0 {
		return "", false
	}
	names := make([]string, 0, len(versions))
	for _, version := range versions {
		if version.IsDir() {
			names = append(names, version.Name())
		}
	}
	semver.Sort(names)
	return filepath.Join(dir, names[len(names)-1]), true
}
