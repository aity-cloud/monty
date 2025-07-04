//go:build !cli

package commands

import (
	"fmt"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strings"

	"github.com/aity-cloud/monty/apis"
	montycorev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	loggingv1beta1 "github.com/aity-cloud/monty/apis/logging/v1beta1"
	monitoringv1beta1 "github.com/aity-cloud/monty/apis/monitoring/v1beta1"
	"github.com/aity-cloud/monty/controllers"
	"github.com/aity-cloud/monty/pkg/features"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/monty/common"
	"github.com/aity-cloud/monty/pkg/tracing"
	"github.com/aity-cloud/monty/pkg/util/k8sutil"
	"github.com/aity-cloud/monty/pkg/util/manager"
	"github.com/aity-cloud/monty/pkg/versions"
	upgraderesponder "github.com/longhorn/upgrade-responder/client"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/spf13/cobra"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	upgradeResponderAddress = "https://upgrades.monty-upgrade-responder.livestock.rancher.io/v1/checkupgrade"
)

type crdFunc func() (*crd.CRD, error)

func init() {
	apis.InitScheme(scheme)
}

func BuildClientCmd() *cobra.Command {
	var (
		metricsAddr  string
		probeAddr    string
		disableUsage bool
		echoVersion  bool
		logLevel     string
		montyCentral bool
	)

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Run the Monty Client Manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			tracing.Configure("client")

			if echoVersion {
				fmt.Println(versions.Version)
				return nil
			}

			if os.Getenv("DO_NOT_TRACK") == "1" {
				disableUsage = true
			}

			level := logger.ParseLevel(logLevel)

			ctrl.SetLogger(k8sutil.NewControllerRuntimeLogger(level))

			config := ctrl.GetConfigOrDie()

			mgr, err := ctrl.NewManager(config, ctrl.Options{
				Scheme:                 scheme,
				Metrics:                server.Options{BindAddress: metricsAddr},
				WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
				HealthProbeBindAddress: probeAddr,
				LeaderElection:         false,
			})
			if err != nil {
				setupLog.Error(err, "unable to start client manager")
				return err
			}

			crdFactory, err := crd.NewFactoryFromClient(config)
			if err != nil {
				setupLog.Error(err, "unable to create crd factory")
				return err
			}

			var upgradeChecker *upgraderesponder.UpgradeChecker
			if !(disableUsage || common.DisableUsage) {
				upgradeRequester := manager.UpgradeRequester{
					Version:     versions.Version,
					InstallType: manager.InstallTypeAgent,
				}
				upgradeRequester.SetupLoggerWithManager(mgr)
				setupLog.Info("Usage tracking enabled", "current-version", versions.Version)
				upgradeChecker = upgraderesponder.NewUpgradeChecker(upgradeResponderAddress, &upgradeRequester)
				upgradeChecker.Start()
				defer upgradeChecker.Stop()
			}

			// Apply CRDs
			crds := []crd.CRD{}
			for _, crdFunc := range []crdFunc{
				montycorev1beta1.CollectorCRD,
				loggingv1beta1.CollectorConfigCRD,
				monitoringv1beta1.CollectorConfigCRD,
				montycorev1beta1.KeyringCRD,
			} {
				crd, err := crdFunc()
				if err != nil {
					setupLog.Error(err, "failed to create crd")
					return err
				}
				crds = append(crds, *crd)
			}

			// Only create prometheus crds if they don't already exist
			for _, crdFunc := range []crdFunc{
				monitoringv1beta1.ServiceMonitorCRD,
				monitoringv1beta1.PodMonitorCRD,
				// Need to include Prometheus CRD for deletes even if we're not using prometheus
				monitoringv1beta1.PrometheusCRD,
				monitoringv1beta1.PrometheusAgentCRD,
			} {
				crd, err := crdFunc()
				if err != nil {
					setupLog.Error(err, "failed to create crd")
					return err
				}
				name := strings.ToLower(crd.PluralName + "." + crd.GVK.Group)
				_, err = crdFactory.CRDClient.ApiextensionsV1().CustomResourceDefinitions().Get(cmd.Context(), name, metav1.GetOptions{})
				if err != nil {
					if k8serrors.IsNotFound(err) {
						crds = append(crds, *crd)
					} else {
						setupLog.Error(err, "failed to get crd")
						return err
					}
				}
			}

			err = crdFactory.BatchCreateCRDs(cmd.Context(), crds...).BatchWait()
			if err != nil {
				setupLog.Error(err, "failed to apply crds")
				return err
			}

			if err = (&controllers.CoreCollectorReconciler{}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "Core Collector")
				return err
			}

			// +kubebuilder:scaffold:builder

			if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
				setupLog.Error(err, "unable to set up health check")
				return err
			}
			if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
				setupLog.Error(err, "unable to set up ready check")
				return err
			}

			setupLog.Info("starting manager")
			if err := mgr.Start(cmd.Context()); err != nil {
				setupLog.Error(err, "error running manager")
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warning, error)")
	cmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":7080", "The address the metric endpoint binds to.")
	cmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":7081", "The address the probe endpoint binds to.")
	cmd.Flags().BoolVarP(&montyCentral, "central", "c", false, "run controllers in Monty central cluster mode")
	cmd.Flags().BoolVarP(&echoVersion, "version", "v", false, "print the version and exit")
	features.DefaultMutableFeatureGate.AddFlag(cmd.Flags())

	return cmd
}

func init() {
	AddCommandsToGroup(MontyComponents, BuildClientCmd())
}
