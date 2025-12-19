package ManagementApiClient

import (
	"context"

	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
)

type (
	ManagementClientConfig struct {
		ApiBaseUrl string `json:"apiBaseUrl"`
	}
	IManagementClientConfigAccessor interface {
		GetManagementClientConfig(ctx context.Context) *ManagementClientConfig
	}
	IManagementApiClient interface {
		Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*common.WrappedResonseT[models_api_login_models.LoginResponse], error)
		Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*common.WrappedResonseT[models_api_login_models.LogoutResponse], error)
		PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error)
		PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error)
		VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*common.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error)
		GetUserProfile(ctx context.Context) (*common.WrappedResonseT[models.Profile], error)
		UpdateUserProfile(ctx context.Context, profile *models.Profile) (*common.WrappedResonseT[models.Profile], error)
		GetUserLinkedAccounts(ctx context.Context) (*common.WrappedResonseT[models.LinkedAccountsResponse], error)
		DeleteUserLinkedAccount(ctx context.Context, identity string) (*common.WrappedResonseT[models.DeleteLinkedAccountResponse], error)

		// Passkey management - HTTP-based (avoiding JS fetch interop)
		GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[models.PasskeysResponse], error)
		DeletePasskeyHTTP(ctx context.Context, request *models.PasskeyDeleteRequest) (*common.WrappedResonseT[models.PasskeyDeleteResponse], error)
		RenamePasskeyHTTP(ctx context.Context, request *models.PasskeyRenameRequest) (*common.WrappedResonseT[models.PasskeyRenameResponse], error)

		// Access to core RageApiClient for TOTP and other core APIs
		GetRageApiClient() contracts_go_app_RageApiClient.IRageApiClient
	}
)
