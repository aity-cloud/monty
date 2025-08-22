package machinery

import (
	"context"
	"fmt"

	authnoauth "github.com/aity-cloud/monty/pkg/auth/noauth"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/noauth"
)

func NewNoauthServer(
	ctx context.Context,
	ap *v1beta1.AuthProvider,
) *noauth.Server {
	mw, err := authnoauth.New(ctx, ap.Spec)
	if err != nil {
		panic(fmt.Errorf("failed to create noauth auth provider: %w", err))
	}
	return noauth.NewServer(mw.ServerConfig())
}
