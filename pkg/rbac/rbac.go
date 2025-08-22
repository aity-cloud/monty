package rbac

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/gin-gonic/gin"
)

type RBACHeader map[string]*corev1.ReferenceList

const (
	UserIDKey = "rbac_user_id"
)

type Provider interface {
	AccessHeader(context.Context, *corev1.ReferenceList) (RBACHeader, error)
}

func AuthorizedUserID(c *gin.Context) (string, bool) {
	userId, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	return userId.(string), true
}
