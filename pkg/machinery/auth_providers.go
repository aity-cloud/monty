package machinery

import (
	"context"
	"fmt"

	"github.com/aity-cloud/monty/pkg/auth"
	"github.com/aity-cloud/monty/pkg/auth/noauth"
	"github.com/aity-cloud/monty/pkg/auth/openid"
	"github.com/aity-cloud/monty/pkg/auth/test"
	"github.com/aity-cloud/monty/pkg/config/meta"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
)

func LoadAuthProviders(ctx context.Context, objects meta.ObjectList) map[string]auth.Middleware {
	authProviders := make(map[string]auth.Middleware)
	objects.Visit(
		func(ap *v1beta1.AuthProvider) {
			switch ap.Spec.Type {
			case v1beta1.AuthProviderOpenID:
				mw, err := openid.New(ctx, ap.Spec)
				if err != nil {
					panic(fmt.Errorf("failed to create OpenID auth provider: %w", err))
				}
				authProviders[ap.Name] = mw
			case v1beta1.AuthProviderNoAuth:
				mw, err := noauth.New(ctx, ap.Spec)
				if err != nil {
					panic(fmt.Errorf("failed to create noauth auth provider: %w", err))
				}
				authProviders[ap.Name] = mw
			case "test":
				authProviders["test"] = &test.TestAuthMiddleware{
					Strategy: test.AuthStrategyUserIDInAuthHeader,
				}
			default:
				panic(fmt.Errorf("unsupported auth provider type: %s", ap.Spec.Type))
			}
		},
	)
	return authProviders
}
