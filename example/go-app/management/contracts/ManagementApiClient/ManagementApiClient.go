package ManagementApiClient

import (
	"context"

	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
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
		Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error)
		Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error)
		VerifyCodeBegin(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error)
		VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error)
		PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error)
		PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error)
		GetUserProfile(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error)
		UpdateUserProfile(ctx context.Context, profile *models.Profile) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error)
		GetUserLinkedAccounts(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.LinkedAccountsResponse], error)
		DeleteUserLinkedAccount(ctx context.Context, identity string) (*fluffycore_go_app_fetch.WrappedResonseT[models.DeleteLinkedAccountResponse], error)
	}
)
