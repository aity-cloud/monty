package kubernetes_test

import (
	"context"
	"fmt"
	"os"

	corev1beta1 "github.com/aity-cloud/monty/apis/core/v1beta1"
	"github.com/aity-cloud/monty/pkg/oci"
	"github.com/aity-cloud/monty/pkg/oci/kubernetes"
	"github.com/aity-cloud/monty/pkg/versions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	imageDigest  = "sha256:15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225"
	imageDigest2 = "sha256:576d76d88778a1d23c411e92701b89fbf9807cf3c8ca5f832677843ce9db1ccb"
)

var _ = Describe("Kubernetes OCI handler", Ordered, Label("unit", "slow"), func() {
	var (
		gateway *corev1beta1.Gateway
		k8sOCI  oci.Fetcher
	)
	BeforeAll(func() {
		gateway = &corev1beta1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "monty-gateway",
				Namespace: namespace,
			},
			Spec: corev1beta1.GatewaySpec{},
		}
		Expect(k8sClient.Create(context.Background(), gateway)).To(Succeed())

		var err error
		k8sOCI, err = kubernetes.NewKubernetesResolveImageDriver(
			namespace,
			kubernetes.WithRestConfig(restConfig),
		)
		Expect(err).NotTo(HaveOccurred())
	})

	When("an unknown image type is requested", func() {
		It("should return unsported image type error", func() {
			_, err := k8sOCI.GetImage(context.Background(), "unknown")
			Expect(err).To(MatchError(kubernetes.ErrUnsupportedImageType))
		})
	})

	When("gateway status is not set", func() {
		It("should not return monty image", func() {
			_, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMonty)
			Expect(err).To(HaveOccurred())
		})
		It("should not return the minimal image", func() {
			_, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMonty)
			Expect(err).To(HaveOccurred())
		})
		When("version is set", func() {
			BeforeEach(func() {
				versions.Version = "v1.0.0"
			})
			It("should not return the minimal image", func() {
				_, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMonty)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	When("gateway status is set", Ordered, func() {
		BeforeEach(func() {
			gateway.Status = corev1beta1.GatewayStatus{
				Image: fmt.Sprintf("rancher/monty-test@%s", imageDigest),
			}
			Expect(k8sClient.Status().Update(context.Background(), gateway)).To(Succeed())
		})
		When("the minimal image is available from the environment", func() {
			It("should return the minimal image", func() {
				minimalRef := fmt.Sprintf("rancher/monty-test@%s", imageDigest2)
				os.Setenv("MONTY_MINIMAL_IMAGE_REF", minimalRef)
				DeferCleanup(func() {
					os.Unsetenv("MONTY_MINIMAL_IMAGE_REF")
				})
				image, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMinimal)
				Expect(err).NotTo(HaveOccurred())
				Expect(image.String()).To(Equal(minimalRef))
			})
		})
		When("version is unset", func() {
			BeforeEach(func() {
				versions.Version = "unversioned"
			})
			It("should return the monty image", func() {
				image, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMonty)
				Expect(err).NotTo(HaveOccurred())
				Expect(image.String()).To(Equal(fmt.Sprintf("rancher/monty-test@%s", imageDigest)))
			})
			It("should return the monty image as the minimal image", func() {
				image, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMinimal)
				Expect(err).NotTo(HaveOccurred())
				Expect(image.String()).To(Equal(fmt.Sprintf("rancher/monty-test@%s", imageDigest)))
			})
		})
		When("version is set", func() {
			BeforeEach(func() {
				versions.Version = "v1.0.0"
			})
			It("should return the minimal image", func() {
				image, err := k8sOCI.GetImage(context.Background(), oci.ImageTypeMinimal)
				Expect(err).NotTo(HaveOccurred())
				Expect(image.String()).To(Equal("rancher/monty-test:v1.0.0-minimal"))
			})
		})
	})
})
