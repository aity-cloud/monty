package prometheusrule_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrometheusrule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prometheusrule Suite")
}
