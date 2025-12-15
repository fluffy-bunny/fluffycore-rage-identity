package ManagementApiClient

import (
	"context"

	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
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
		Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error)
		Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error)
		PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error)
		PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error)
		VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error)
		GetUserProfile(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error)
		UpdateUserProfile(ctx context.Context, profile *models.Profile) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error)
		GetUserLinkedAccounts(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.LinkedAccountsResponse], error)
		DeleteUserLinkedAccount(ctx context.Context, identity string) (*fluffycore_go_app_fetch.WrappedResonseT[models.DeleteLinkedAccountResponse], error)

		// Passkey management (delegated to core RageApiClient)
		GetPasskeys(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeysResponse], error)
		DeletePasskey(ctx context.Context, credentialID string) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyDeleteResponse], error)
		RenamePasskey(ctx context.Context, credentialID string, body *models.PasskeyRenameRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models.PasskeyRenameResponse], error)

		// Access to core RageApiClient for TOTP and other core APIs
		GetRageApiClient() contracts_go_app_RageApiClient.IRageApiClient
	}
)
