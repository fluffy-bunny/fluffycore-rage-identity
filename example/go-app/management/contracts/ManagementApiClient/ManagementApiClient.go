package ManagementApiClient

import (
	"context"

	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_linked_identities "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_linked_identities"
	models_api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	models_api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
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
		GetUserProfile(ctx context.Context) (*common.WrappedResonseT[models_api_profile.Profile], error)
		UpdateUserProfile(ctx context.Context, profile *models_api_profile.Profile) (*common.WrappedResonseT[models_api_profile.Profile], error)
		GetUserLinkedAccounts(ctx context.Context) (*common.WrappedResonseT[models_api_linked_identities.LinkedAccountsResponse], error)
		DeleteUserLinkedAccount(ctx context.Context, identity string) (*common.WrappedResonseT[models_api_linked_identities.DeleteLinkedAccountResponse], error)

		// Passkey management - HTTP-based (avoiding JS fetch interop)
		GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[models_api_passkey.PasskeysResponse], error)
		DeletePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyDeleteRequest) (*common.WrappedResonseT[models_api_passkey.PasskeyDeleteResponse], error)
		RenamePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyRenameRequest) (*common.WrappedResonseT[models_api_passkey.PasskeyRenameResponse], error)

		// Access to core RageApiClient for TOTP and other core APIs
		GetRageApiClient() contracts_go_app_RageApiClient.IRageApiClient
	}
)
