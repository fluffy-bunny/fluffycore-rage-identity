package ManagementApiClient

import (
	"context"

	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
)

type (
	ManagementClientConfig struct {
		ApiBaseUrl string `json:"apiBaseUrl"`
	}
	IManagementClientConfigAccessor interface {
		GetManagementClientConfig(ctx context.Context) *ManagementClientConfig
	}

	IManagementApiClient interface {
		Login(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error)
		Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error)
	}
)
