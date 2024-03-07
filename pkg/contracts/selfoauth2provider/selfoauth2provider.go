package selfoauth2provider

import (
	"context"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	oauth2 "golang.org/x/oauth2"
)

type (
	GetConfigResponse struct {
		Config   *oauth2.Config
		Verifier *oidc.IDTokenVerifier
	}
	ISelfOAuth2Provider interface {
		GetConfig(ctx context.Context) (*GetConfigResponse, error)
	}
)
