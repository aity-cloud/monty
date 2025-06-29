package cortex_test

import (
	"encoding/json"

	"github.com/aity-cloud/monty/pkg/alerting/drivers/cortex"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	"github.com/aity-cloud/monty/pkg/metrics/compat"
	"github.com/aity-cloud/monty/pkg/test/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func convertToMatrix(filepath string) []*alertingv1.ActiveWindow {
	mat := testdata.TestData(filepath)
	qr, err := compat.UnmarshalPrometheusResponse(mat)
	Expect(err).To(Succeed())
	Expect(qr).ToNot(BeNil())
	matrix, err := qr.GetMatrix()
	Expect(err).To(Succeed())
	return cortex.ReducePrometheusMatrix(matrix)

}

var _ = Describe("Alerting cortex suite", Label("unit"), func() {
	When("Reducing alert matrix results", func() {
		It("should reduce deduplicate incoming metric samples", func() {
			windows := convertToMatrix("alerting/matrix/matrix.json")
			Expect(windows).NotTo(HaveLen(0))
			Expect(windows).To(HaveLen(1))
			Expect(windows[0].Fingerprints).To(HaveLen(1))

			windows = convertToMatrix("alerting/matrix/overlapping_causes.json")
			Expect(windows).NotTo(HaveLen(0))
			Expect(windows).To(HaveLen(1))
			Expect(windows[0].Fingerprints).To(HaveLen(2))
		})

		It("should discern discrete incidents", func() {
			windows := convertToMatrix("alerting/matrix/discrete_incidents.json")
			Expect(windows).NotTo(HaveLen(0))
			Expect(windows).To(HaveLen(2))
			Expect(windows[0].Fingerprints).To(HaveLen(1))
			Expect(windows[1].Fingerprints).To(HaveLen(1))
		})

		It("should discern when overlapping intervals should be treated as discrete incidents", func() {
			windows := convertToMatrix("alerting/matrix/overlapping_but_discrete.json")
			Expect(windows).NotTo(HaveLen(0))
			Expect(windows).To(HaveLen(2))
			Expect(windows[0].Fingerprints).To(HaveLen(2))
			Expect(windows[1].Fingerprints).To(HaveLen(2))
		})
	})

	When("We parse cortex webhook payloads", func() {
		It("should parse valid payloads to appropriate monty alerting payloads", func() {
			someId := shared.NewAlertingRefId()
			somename := "some-alert-name"
			exampleWorkingPayload := NewSimpleMockAlertManagerPayloadFromAnnotations(map[string]string{
				"alertname":   somename,
				"conditionId": someId,
			})
			mockBody, err := json.Marshal(&exampleWorkingPayload)
			Expect(err).To(Succeed())
			annotations, err := ParseCortexPayloadBytes(mockBody)
			Expect(err).To(Succeed())
			Expect(annotations).To(HaveLen(1))

			montyiResponses, errors := ParseAlertManagerWebhookPayload(annotations)
			Expect(errors).To(HaveLen(1))
			Expect(len(montyiResponses)).To(Equal(len(errors)))
			for _, e := range errors {
				Expect(e).To(Succeed())
			}
			Expect(montyiResponses[0].ConditionId.GetId()).To(Equal(someId))
			Expect(montyiResponses[0].Annotations["alertname"]).To(Equal(somename))
		})

		It("Should errror on invalid cortex webhook payloads", func() {
			somename := "some-alert-name"
			exampleInvalidPayload := NewSimpleMockAlertManagerPayloadFromAnnotations(map[string]string{
				"alertname": somename,
			})
			mockBody, err := json.Marshal(&exampleInvalidPayload)
			Expect(err).To(Succeed())
			annotations, err := ParseCortexPayloadBytes(mockBody)
			Expect(err).To(Succeed())
			montyRequests, errors := ParseAlertManagerWebhookPayload(annotations)
			Expect(errors).To(HaveLen(1))
			Expect(len(montyRequests)).To(Equal(len(errors)))
			Expect(errors[0]).To(HaveOccurred())
			Expect(montyRequests[0]).To(BeNil())
		})
	})
})
