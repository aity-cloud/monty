package capabilities_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
)

var ctrl *gomock.Controller

func TestCapabilities(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Capabilities Suite")
}

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
})

type testCapability string

func (tc testCapability) Equal(other testCapability) bool {
	return tc == other
}

func (tc testCapability) String() string {
	return "test"
}

type testResourceWithMetadata struct {
	capabilities    []testCapability
	labels          map[string]string
	resourceVersion string
}

var _ corev1.MetadataAccessor[testCapability] = (*testResourceWithMetadata)(nil)

func (t *testResourceWithMetadata) GetCapabilities() []testCapability {
	return t.capabilities
}

func (t *testResourceWithMetadata) SetCapabilities(capabilities []testCapability) {
	t.capabilities = capabilities
}

func (t *testResourceWithMetadata) GetLabels() map[string]string {
	return t.labels
}

func (t *testResourceWithMetadata) SetLabels(labels map[string]string) {
	t.labels = labels
}

func (t *testResourceWithMetadata) GetResourceVersion() string {
	return t.resourceVersion
}

func (t *testResourceWithMetadata) SetResourceVersion(resourceVersion string) {
	t.resourceVersion = resourceVersion
}
