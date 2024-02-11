package oauth2factory

import (
	"context"

	oauth2 "golang.org/x/oauth2"
)

type (
	GetConfigRequest struct {
		IDPSlug string
	}
	GetConfigResponse struct {
		Config *oauth2.Config
	}
	IOAuth2Factory interface {
		GetConfig(ctx context.Context, request *GetConfigRequest) (*GetConfigResponse, error)
	}
)
