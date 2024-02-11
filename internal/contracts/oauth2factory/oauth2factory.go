package oauth2factory

import (
	"context"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	oauth2 "golang.org/x/oauth2"
)

type (
	GetConfigRequest struct {
		IDPSlug string
	}
	GetConfigResponse struct {
		Config *oauth2.Config
	}
	GetOIDCProviderRequest struct {
		IDPSlug string
	}
	GetOIDCProviderResponse struct {
		OIDCProvider *oidc.Provider
	}
	IOIDCProviderFactory interface {
		GetOIDCProvider(ctx context.Context, request *GetOIDCProviderRequest) (*GetOIDCProviderResponse, error)
	}
	IOAuth2Factory interface {
		GetConfig(ctx context.Context, request *GetConfigRequest) (*GetConfigResponse, error)
	}
)
