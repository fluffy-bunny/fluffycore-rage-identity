package RageApiClient

import (
	"context"

	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	models_api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	models_api_external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	models_api_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	models_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
)

type (
	RageClientConfig struct {
		ApiBaseUrl string `json:"apiBaseUrl"`
	}
	IRageClientConfigAccessor interface {
		GetRageClientConfig(ctx context.Context) *RageClientConfig
	}

	IRageApiClient interface {
		GetOIDCFlowAppConfig(ctx context.Context) (*common.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig], error)
		GetManifest(ctx context.Context) (*common.WrappedResonseT[models_api_manifest.Manifest], error)
		GetCachedManifest(ctx context.Context) *models_api_manifest.Manifest
		SetCachedManifest(ctx context.Context, manifest *models_api_manifest.Manifest)
		LoginPhaseOne(ctx context.Context, request *models_api_login_models.LoginPhaseOneRequest) (*common.WrappedResonseT[models_api_login_models.LoginPhaseOneResponse], error)
		VerifyPasswordStrength(ctx context.Context, request *models_api_password.VerifyPasswordStrengthRequest) (*common.WrappedResonseT[models_api_password.VerifyPasswordStrengthResponse], error)
		VerifyUserName(ctx context.Context, request *models_api_verify_username.VerifyUsernameRequest) (*common.WrappedResonseT[models_api_verify_username.VerifyUsernameResponse], error)
		Signup(ctx context.Context, request *models_api_login_models.SignupRequest) (*common.WrappedResonseT[models_api_login_models.SignupResponse], error)
		StartExternalLogin(ctx context.Context, request *models_api_external_idp.StartExternalIDPLoginRequest) (*common.WrappedResonseT[models_api_external_idp.StartExternalIDPLoginResponse], error)
		LoginPassword(ctx context.Context, request *models_api_login_models.LoginPasswordRequest) (*common.WrappedResonseT[models_api_login_models.LoginPasswordResponse], error)
		PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error)
		PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error)
		VerifyCodeBegin(ctx context.Context) (*common.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error)
		VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*common.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error)

		// TOTP/Authenticator APIs
		GetTOTPStatus(ctx context.Context) ([]byte, error)
		EnrollTOTP(ctx context.Context) ([]byte, error)
		VerifyTOTP(ctx context.Context, code string) ([]byte, error)
		DisableTOTP(ctx context.Context) ([]byte, error)

		// Passkey APIs - HTTP-based (avoiding JS fetch interop)
		GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[*models_api_passkey.PasskeysResponse], error)
		DeletePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyDeleteRequest) (*common.WrappedResonseT[*models_api_passkey.PasskeyDeleteResponse], error)
		RenamePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyRenameRequest) (*common.WrappedResonseT[*models_api_passkey.PasskeyRenameResponse], error)
	}
)
