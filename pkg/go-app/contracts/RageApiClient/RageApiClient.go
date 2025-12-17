package RageApiClient

import (
	"context"

	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	models_api_external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	models_api_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	models_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
)

type (
	RageClientConfig struct {
		ApiBaseUrl string `json:"apiBaseUrl"`
	}
	IRageClientConfigAccessor interface {
		GetRageClientConfig(ctx context.Context) *RageClientConfig
	}

	IRageApiClient interface {
		GetOIDCFlowAppConfig(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig], error)
		GetManifest(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_manifest.Manifest], error)
		GetCachedManifest(ctx context.Context) *models_api_manifest.Manifest
		SetCachedManifest(ctx context.Context, manifest *models_api_manifest.Manifest)
		LoginPhaseOne(ctx context.Context, request *models_api_login_models.LoginPhaseOneRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginPhaseOneResponse], error)
		VerifyPasswordStrength(ctx context.Context, request *models_api_password.VerifyPasswordStrengthRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_password.VerifyPasswordStrengthResponse], error)
		VerifyUserName(ctx context.Context, request *models_api_verify_username.VerifyUsernameRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_username.VerifyUsernameResponse], error)
		Signup(ctx context.Context, request *models_api_login_models.SignupRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.SignupResponse], error)
		StartExternalLogin(ctx context.Context, request *models_api_external_idp.StartExternalIDPLoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_external_idp.StartExternalIDPLoginResponse], error)
		LoginPassword(ctx context.Context, request *models_api_login_models.LoginPasswordRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginPasswordResponse], error)
		PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error)
		PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error)
		VerifyCodeBegin(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error)
		VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error)

		// TOTP/Authenticator APIs
		GetTOTPStatus(ctx context.Context) ([]byte, error)
		EnrollTOTP(ctx context.Context) ([]byte, error)
		VerifyTOTP(ctx context.Context, code string) ([]byte, error)
		DisableTOTP(ctx context.Context) ([]byte, error)

		// Passkey APIs - HTTP-based (avoiding JS fetch interop)
		GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[*models.PasskeysResponse], error)
		DeletePasskeyHTTP(ctx context.Context, request *models.PasskeyDeleteRequest) (*common.WrappedResonseT[*models.PasskeyDeleteResponse], error)
		RenamePasskeyHTTP(ctx context.Context, request *models.PasskeyRenameRequest) (*common.WrappedResonseT[*models.PasskeyRenameResponse], error)
	}
)
