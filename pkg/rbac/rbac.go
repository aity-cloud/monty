package rbac

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "rbac_user_id"
)

type Provider interface {
	SubjectAccess(context.Context, *corev1.SubjectAccessRequest) (*corev1.ReferenceList, error)
}

func AuthorizedUserID(c *gin.Context) (string, bool) {
	userId, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	return userId.(string), true
}
