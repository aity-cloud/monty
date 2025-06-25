package setup

import (
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
)

func init() {
	gin.SetMode(gin.TestMode)
	logger.DefaultWriter = ginkgo.GinkgoWriter
}
