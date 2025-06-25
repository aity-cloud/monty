package noauth

import (
	"context"

	"github.com/aity-cloud/monty/pkg/auth"
	"github.com/aity-cloud/monty/pkg/auth/openid"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/noauth"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type NoauthMiddleware struct {
	openidMiddleware auth.HTTPMiddleware
	noauthConfig     *noauth.ServerConfig
	logger           *slog.Logger
}

var _ auth.Middleware = (*NoauthMiddleware)(nil)

func New(ctx context.Context, config v1beta1.AuthProviderSpec) (*NoauthMiddleware, error) {
	lg := logger.New().WithGroup("noauth")
	conf, err := util.DecodeStruct[noauth.ServerConfig](config.Options)
	if err != nil {
		return nil, err
	}
	openidConf, err := util.DecodeStruct[map[string]any](conf.OpenID)
	if err != nil {
		return nil, err
	}
	openidMw, err := openid.New(ctx, v1beta1.AuthProviderSpec{
		Type:    "openid",
		Options: *openidConf,
	})
	if err != nil {
		return nil, err
	}
	m := &NoauthMiddleware{
		openidMiddleware: openidMw,
		noauthConfig:     conf,
		logger:           lg,
	}
	m.noauthConfig.Logger = lg

	return m, nil
}

func (m *NoauthMiddleware) Handle(c *gin.Context) {
	m.openidMiddleware.Handle(c)
}

func (m *NoauthMiddleware) ServerConfig() *noauth.ServerConfig {
	return m.noauthConfig
}
