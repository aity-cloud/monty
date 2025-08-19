package util

import (
	"github.com/aity-cloud/monty/plugins/metrics/pkg/constants"
	"github.com/gin-gonic/gin"
)

func AuthorizedClusterIDs(c *gin.Context) []string {
	value, ok := c.Get(constants.AuthorizedClusterIDsKey)
	if !ok {
		return nil
	}
	return value.([]string)
}
