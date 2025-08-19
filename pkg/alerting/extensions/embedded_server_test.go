package extensions_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/aity-cloud/monty/pkg/alerting/client"
	"github.com/aity-cloud/monty/pkg/alerting/drivers/config"
	"github.com/aity-cloud/monty/pkg/alerting/drivers/routing"
	"github.com/aity-cloud/monty/pkg/alerting/extensions"
	"github.com/aity-cloud/monty/pkg/alerting/message"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/pkg/test/alerting"
	"github.com/aity-cloud/monty/pkg/test/freeport"
	"github.com/google/uuid"
	"github.com/kralicky/yaml/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	amCfg "github.com/prometheus/alertmanager/config"
	"github.com/samber/lo"

	"github.com/aity-cloud/monty/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildEmbeddedServerNotificationTests(
	routerConstructor func(int) routing.MontyRouting,
	dataset *alerting.RoutableDataset,
) bool {

	return Describe("EmbeddedServer test suite", Ordered, Label("integration"), func() {
		var httpClient *http.Client
		var fingerprints []string
		var id string
		var env *test.Environment
		var tmpConfigDir string

		var webPort int
		var montyPort int
		var alertingClient client.AlertingClient
		sendMsg := func(client *http.Client, msg config.WebhookMessage, montyPort int) {
			content, err := json.Marshal(msg)
			Expect(err).NotTo(HaveOccurred())
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", montyPort, shared.AlertingDefaultHookName), bytes.NewReader(content))
			Expect(err).NotTo(HaveOccurred())
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		}
		var err error
		alertingClient, err = client.NewClient(
			client.WithAlertManagerAddress(
				fmt.Sprintf("127.0.0.1:%d", webPort),
			),
			client.WithQuerierAddress(
				fmt.Sprintf("127.0.0.1:%d", montyPort),
			),
			client.WithTLSConfig(
				env.AlertingClientTLSConfig(),
			),
		)
		Expect(err).NotTo(HaveOccurred())
		sendMsgAlertManager := func(ctx context.Context, labels, annotations map[string]string, alertManagerPort int) {
			err := alertingClient.AlertClient().PostNotification(ctx, client.AlertObject{
				Id:          labels[message.NotificationPropertyMontyUuid],
				Labels:      labels,
				Annotations: annotations,
			})
			Expect(err).NotTo(HaveOccurred())

		}

		listNotif := func(client *http.Client, listReq *alertingv1.ListNotificationRequest, montyPort int) *alertingv1.ListMessageResponse {
			listReq.Sanitize()
			err := listReq.Validate()
			Expect(err).NotTo(HaveOccurred())
			content, err := protojson.Marshal(listReq)
			Expect(err).NotTo(HaveOccurred())
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", montyPort, "/notifications/list"), bytes.NewReader(content))
			Expect(err).NotTo(HaveOccurred())
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var listResp *alertingv1.ListMessageResponse
			err = json.NewDecoder(resp.Body).Decode(&listResp)
			Expect(err).NotTo(HaveOccurred())
			return listResp
		}

		listAlarm := func(client *http.Client, listReq *alertingv1.ListAlarmMessageRequest, montyPort int) *alertingv1.ListMessageResponse {
			listReq.Sanitize()
			err := listReq.Validate()
			Expect(err).NotTo(HaveOccurred())
			content, err := protojson.Marshal(listReq)
			Expect(err).NotTo(HaveOccurred())
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", montyPort, "/alarms/list"), bytes.NewReader(content))
			Expect(err).NotTo(HaveOccurred())
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var listResp *alertingv1.ListMessageResponse
			err = json.NewDecoder(resp.Body).Decode(&listResp)
			Expect(err).NotTo(HaveOccurred())
			return listResp
		}
		BeforeAll(func(ctx SpecContext) {

			env = &test.Environment{}
			Expect(env).NotTo(BeNil())
			Expect(env.Start()).To(Succeed())
			DeferCleanup(env.Stop, "Test Suite Finished")
			tmpConfigDir = env.GenerateNewTempDirectory("alertmanager-config")
			err := os.MkdirAll(tmpConfigDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			Expect(tmpConfigDir).NotTo(Equal(""))

			// start embedded alert manager with config that points to monty embedded server

			freeport := freeport.GetFreePort()
			Expect(freeport).NotTo(BeZero())
			montyPort = freeport
			extensions.StartMontyEmbeddedServer(env.Context(), fmt.Sprintf(":%d", montyPort), false)

			router := routerConstructor(montyPort)
			Expect(tmpConfigDir).NotTo(BeEmpty())
			confFile := path.Join(tmpConfigDir, "alertmanager.yml")
			Expect(confFile).NotTo(BeEmpty())

			config, err := router.BuildConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			err = os.WriteFile(confFile, util.Must(yaml.Marshal(config)), 0644)
			Expect(err).NotTo(HaveOccurred())
			ports := env.StartEmbeddedAlertManager(env.Context(), confFile, nil)
			webPort = ports.ApiPort
			httpClient = http.DefaultClient
		})

		When("we use the embedded monty embedded server", func() {
			It("should handle webhook messages indexed by Monty", func() {
				Expect(webPort).NotTo(BeZero())
				Expect(montyPort).NotTo(BeZero())
				msg := config.WebhookMessage{
					Alerts: config.Alerts{
						{
							Status: "firing",
							Labels: map[string]string{
								message.NotificationPropertyMontyUuid: uuid.New().String(),
								message.NotificationPropertySeverity:  alertingv1.MontySeverity_Info.String(),
							},
							Annotations: map[string]string{},
						},
					},
					Version:         "4",
					Receiver:        uuid.New().String(),
					TruncatedAlerts: 0,
					Status:          "firing",
					GroupKey:        uuid.New().String(),
					ExternalURL:     fmt.Sprintf("http://localhost:%d", webPort),
				}
				sendMsg(httpClient, msg, montyPort)
			})

			It("should list notification messages indexed by Monty", func() {
				Expect(webPort).NotTo(BeZero())
				Expect(montyPort).NotTo(BeZero())

				listReq := &alertingv1.ListNotificationRequest{}
				respList := listNotif(httpClient, listReq, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(1))
			})

			Specify("it should dedupe frequency-based persistenced based on group keys and id keys based on what is available", func() {
				listRequest := &alertingv1.ListNotificationRequest{
					SeverityFilters: []alertingv1.MontySeverity{
						alertingv1.MontySeverity_Warning,
					},
				}
				groupKey := uuid.New().String()
				msgId := uuid.New().String()
				msg := config.WebhookMessage{
					Alerts: config.Alerts{
						{
							Status: "firing",
							Labels: map[string]string{
								message.NotificationPropertyMontyUuid: msgId,
								message.NotificationPropertySeverity:  alertingv1.MontySeverity_Warning.String(),
								message.NotificationPropertyDedupeKey: groupKey,
							},
							Annotations: map[string]string{},
						},
					},
					Version:         "4",
					Receiver:        uuid.New().String(),
					TruncatedAlerts: 0,
					Status:          "firing",
					GroupKey:        groupKey,
					ExternalURL:     fmt.Sprintf("http://localhost:%d", webPort),
				}
				sendMsg(httpClient, msg, montyPort)
				respList := listNotif(httpClient, listRequest, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(1))

				// send the same message again with group key but different uuid
				msg.Alerts[0].Labels[message.NotificationPropertyMontyUuid] = uuid.New().String()
				sendMsg(httpClient, msg, montyPort)
				respList = listNotif(httpClient, listRequest, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(1))

				// send the same message again with uuid but different group key but same uuid
				msg.Alerts[0].Labels[message.NotificationPropertyMontyUuid] = msgId
				msg.Alerts[0].Labels[message.NotificationPropertyDedupeKey] = uuid.New().String()

				sendMsg(httpClient, msg, montyPort)
				respList = listNotif(httpClient, listRequest, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(2))

				msg.Alerts[0].Labels[message.NotificationPropertyMontyUuid] = uuid.New().String()
				msg.Alerts[0].Labels[message.NotificationPropertyDedupeKey] = uuid.New().String()
				sendMsg(httpClient, msg, montyPort)
				respList = listNotif(httpClient, listRequest, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(3))

				sendMsg(httpClient, msg, montyPort)
				respList = listNotif(httpClient, listRequest, montyPort)
				Expect(respList.Items).NotTo(BeNil())
				Expect(respList.Items).To(HaveLen(3))

			})
		})

		When("we integrate with external AlertManager(s)", func() {
			It("should reset the embedded server state", func() {
				By("verifying there is an input dataset defined")
				Expect(len(dataset.Routables)).NotTo(BeZero())
				Expect(len(dataset.ExpectedNotifications)).NotTo(BeZero())
				Expect(len(dataset.ExpectedAlarms)).NotTo(BeZero())

				By("restarting the embedded server")
				freeport := freeport.GetFreePort()
				Expect(freeport).NotTo(BeZero())
				montyPort = freeport
				extensions.StartMontyEmbeddedServer(env.Context(), fmt.Sprintf(":%d", montyPort), false)

				router := routerConstructor(montyPort)
				By("building the required routes for the routables")
				for _, r := range dataset.Routables {
					if r.Namespace() == routing.NotificationSubTreeLabel() {
						// no need to build this one
						continue
					}
					err := router.SetNamespaceSpec(
						r.Namespace(),
						r.GetRoutingLabels()[message.NotificationPropertyMontyUuid],
						&alertingv1.FullAttachedEndpoints{
							Items: []*alertingv1.FullAttachedEndpoint{},
						},
					)
					Expect(err).NotTo(HaveOccurred())
				}

				Expect(tmpConfigDir).NotTo(BeEmpty())
				confFile := path.Join(tmpConfigDir, "alertmanager.yml")
				Expect(confFile).NotTo(BeEmpty())

				config, err := router.BuildConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config).NotTo(BeNil())
				err = os.WriteFile(confFile, util.Must(yaml.Marshal(config)), 0644)
				Expect(err).NotTo(HaveOccurred())
				ports := env.StartEmbeddedAlertManager(env.Context(), confFile, nil)
				alertingClient, err = client.NewClient(
					client.WithAlertManagerAddress(
						fmt.Sprintf("127.0.0.1:%d", ports.ApiPort),
					),
					client.WithQuerierAddress(
						fmt.Sprintf("127.0.0.1:%d", montyPort),
					),
					client.WithTLSConfig(env.AlertingClientTLSConfig()),
				)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should persist the routables", func() {
				for _, r := range dataset.Routables {
					sendMsgAlertManager(env.Context(),
						lo.Assign(
							r.GetRoutingLabels(),
							map[string]string{
								message.NotificationPropertyFingerprint: "fingerprint",
							},
						),
						lo.Assign(
							r.GetRoutingAnnotations(),
							map[string]string{
								message.NotificationPropertyFingerprint: "fingerprint",
							},
						),
						webPort)
				}
				fingerprints = []string{
					uuid.New().String(),
					uuid.New().String(),
					uuid.New().String(),
					uuid.New().String(),
				}
				id = uuid.New().String()
				r := &alertingv1.AlertCondition{
					Name:        "fingerprint test",
					Description: "fingerprint test",
					Id:          id,
					Severity:    alertingv1.MontySeverity_Critical,
					AlertType: &alertingv1.AlertTypeDetails{
						Type: &alertingv1.AlertTypeDetails_System{
							System: &alertingv1.AlertConditionSystem{
								ClusterId: &corev1.Reference{Id: uuid.New().String()},
								Timeout:   durationpb.New(10 * time.Minute),
							},
						},
					},
				}
				for i := 0; i < 50; i++ {
					fingerprint := fingerprints[i%len(fingerprints)]
					sendMsgAlertManager(
						env.Context(),
						lo.Assign(
							r.GetRoutingLabels(),
							map[string]string{
								message.NotificationPropertyFingerprint: fingerprint,
							},
						),
						lo.Assign(
							r.GetRoutingAnnotations(),
							map[string]string{
								message.NotificationPropertyFingerprint: fingerprint,
							},
						),
						webPort,
					)
				}
				Eventually(func() error {
					_ = webPort
					_ = montyPort
					_ = tmpConfigDir
					for _, pair := range dataset.ExpectedNotifications {
						listResp := listNotif(httpClient, pair.A, montyPort)
						if len(listResp.Items) != pair.B {
							return fmt.Errorf(
								"notification pair failed %s : %d vs %d",
								util.Must(json.Marshal(pair.A)),
								len(listResp.Items),
								pair.B,
							)
						}
					}

					for _, pair := range dataset.ExpectedAlarms {
						listResp := listAlarm(httpClient, pair.A, montyPort)
						if len(listResp.Items) != pair.B {
							return fmt.Errorf(
								"alarm pair failed %s : %d vs %d",
								util.Must(json.Marshal(pair.A)),
								len(listResp.Items),
								pair.B,
							)
						}
					}
					return nil
				}, time.Minute, time.Second*5,
				).Should(BeNil())
			})
		})

		It("should handle fingerprints when correlating alarm incident windows to messages", func(ctx SpecContext) {
			By("verifying the alerting cluster has received unique alerts for each unique fingerprint")
			Eventually(func() error {
				ags, err := alertingClient.AlertClient().ListAlerts(ctx)
				if err != nil {
					return err
				}
				foundFingerprints := map[string]struct{}{}
				for _, ag := range ags {
					if v, ok := ag.Labels[message.NotificationPropertyFingerprint]; ok {
						foundFingerprints[v] = struct{}{}
					}
				}
				if len(lo.Intersect(lo.Keys(foundFingerprints), fingerprints)) != len(fingerprints) {
					return fmt.Errorf("never received all fingerprints %s", fingerprints)
				}
				return nil
			}, time.Minute, time.Second*5).Should(BeNil())

			By("verifying the embedded server has persisted notifications for each fingerprint")
			Eventually(func() error {
				listResp := listAlarm(httpClient,
					&alertingv1.ListAlarmMessageRequest{
						ConditionId: &alertingv1.ConditionReference{
							Id: id,
						},
						Fingerprints: fingerprints,
						Start:        timestamppb.New(time.Now().Add(-time.Hour)),
						End:          timestamppb.New(time.Now().Add(time.Hour)),
					},
					montyPort)
				if len(listResp.GetItems()) != len(fingerprints) {
					return fmt.Errorf(
						"expected to match %d=%d persisted alarm messages & number of fingerprints",
						len(listResp.GetItems()),
						len(fingerprints),
					)
				}
				return nil
			}, time.Second*120, time.Second*5).Should(BeNil())
		})
	})
}

var _ = BuildEmbeddedServerNotificationTests(func(dynamicPort int) routing.MontyRouting {
	cfg := config.WebhookConfig{
		NotifierConfig: config.NotifierConfig{
			VSendResolved: false,
		},
		URL: &amCfg.URL{
			URL: util.Must(url.Parse(fmt.Sprintf("http://localhost:%d%s", dynamicPort, shared.AlertingDefaultHookName))),
		},
	}
	return routing.NewMontyRouterV1(cfg)
}, alerting.NewRoutableDataset())
