package RageApiClient

import (
	"context"
	"encoding/json"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	models_api_preferences "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_preferences"
	models_api_external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	models_api_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	models_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
)

type (
	service struct {
		cachedRageConfigResponse *common.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig]
		cachedManifest           *models_api_manifest.Manifest
		rageClientConfigAccessor contracts_go_app_RageApiClient.IRageClientConfigAccessor
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_RageApiClient.IRageApiClient = stemService

func (s *service) Ctor(
	rageClientConfigAccessor contracts_go_app_RageApiClient.IRageClientConfigAccessor,
) (contracts_go_app_RageApiClient.IRageApiClient, error) {
	return &service{
		rageClientConfigAccessor: rageClientConfigAccessor,
	}, nil
}

func AddScopedIRageApiClient(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_RageApiClient.IRageApiClient](cb, stemService.Ctor)
}

func (s *service) fixUpRageApi(ctx context.Context, relativePath string) string {
	rageClientConfig := s.rageClientConfigAccessor.GetRageClientConfig(ctx)
	if rageClientConfig == nil {
		panic("AppConfig is not set in AppConfigAccessor")
	}
	return rageClientConfig.ApiBaseUrl + relativePath
}
func (s *service) GetOIDCFlowAppConfig(ctx context.Context) (*common.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig], error) {
	if s.cachedRageConfigResponse != nil {
		return s.cachedRageConfigResponse, nil
	}
	var err error
	s.cachedRageConfigResponse, err = common.HTTPFetchWrappedResponseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_OIDCFlowAppConfig),
			CustomHeaders: common.BuildCustomHeaders(),
		})

	return s.cachedRageConfigResponse, err
}

func (s *service) GetCachedManifest(ctx context.Context) *models_api_manifest.Manifest {
	return s.cachedManifest
}
func (s *service) SetCachedManifest(ctx context.Context, manifest *models_api_manifest.Manifest) {
	s.cachedManifest = manifest
}

func (s *service) LoginPhaseOne(ctx context.Context, request *models_api_login_models.LoginPhaseOneRequest) (*common.WrappedResonseT[models_api_login_models.LoginPhaseOneResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.LoginPhaseOneResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_LoginPhaseOne),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) VerifyPasswordStrength(ctx context.Context, request *models_api_password.VerifyPasswordStrengthRequest) (*common.WrappedResonseT[models_api_password.VerifyPasswordStrengthResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_password.VerifyPasswordStrengthResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyPasswordStrength),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) GetManifest(ctx context.Context) (*common.WrappedResonseT[models_api_manifest.Manifest], error) {

	return common.HTTPFetchWrappedResponseT[models_api_manifest.Manifest](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Manifest),
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) Signup(ctx context.Context, request *models_api_login_models.SignupRequest) (*common.WrappedResonseT[models_api_login_models.SignupResponse], error) {

	return common.HTTPFetchWrappedResponseT[models_api_login_models.SignupResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Signup),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) LoginPassword(ctx context.Context, request *models_api_login_models.LoginPasswordRequest) (*common.WrappedResonseT[models_api_login_models.LoginPasswordResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.LoginPasswordResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_LoginPassword),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}
func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetStart),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}
func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*common.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}
func (s *service) VerifyUserName(ctx context.Context, request *models_api_verify_username.VerifyUsernameRequest) (*common.WrappedResonseT[models_api_verify_username.VerifyUsernameResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_verify_username.VerifyUsernameResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyUsername),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) VerifyCodeBegin(ctx context.Context) (*common.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_verify_code.VerifyCodeBeginResponse](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCodeBegin),
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*common.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCode),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) KeepSignedIn(ctx context.Context, request *models_api_login_models.KeepSignedInRequest) (*common.WrappedResonseT[models_api_login_models.KeepSignedInResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_login_models.KeepSignedInResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_KeepSignedIn),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

func (s *service) StartExternalLogin(ctx context.Context, request *models_api_external_idp.StartExternalIDPLoginRequest) (*common.WrappedResonseT[models_api_external_idp.StartExternalIDPLoginResponse], error) {
	return common.HTTPFetchWrappedResponseT[models_api_external_idp.StartExternalIDPLoginResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Start_ExternalLogin),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
}

// TOTP/Authenticator APIs
func (s *service) GetTOTPStatus(ctx context.Context) ([]byte, error) {
	resp, err := common.HTTPFetchWrappedResponseT[map[string]interface{}](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTP),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) EnrollTOTP(ctx context.Context) ([]byte, error) {
	resp, err := common.HTTPFetchWrappedResponseT[map[string]interface{}](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTPEnroll),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) VerifyTOTP(ctx context.Context, code string) ([]byte, error) {
	request := map[string]string{"code": code}
	resp, err := common.HTTPFetchWrappedResponseT[map[string]interface{}](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTPVerify),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) DisableTOTP(ctx context.Context) ([]byte, error) {
	resp, err := common.HTTPFetchWrappedResponseT[map[string]interface{}](ctx,
		&common.CallInput{
			Method:        "DELETE",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTP),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

// Passkey APIs - HTTP-based implementations (avoiding JS fetch interop)
// These use Go's native http.Client which automatically includes cookies
func (s *service) GetPasskeysHTTP(ctx context.Context) (*common.WrappedResonseT[*models_api_passkey.PasskeysResponse], error) {

	resp, err := common.HTTPFetchWrappedResponseT[*models_api_passkey.PasskeysResponse](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	return resp, err

}

func (s *service) DeletePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyDeleteRequest) (*common.WrappedResonseT[*models_api_passkey.PasskeyDeleteResponse], error) {

	url := s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys) + "/" + request.CredentialID

	resp, err := common.HTTPFetchWrappedResponseT[*models_api_passkey.PasskeyDeleteResponse](ctx,
		&common.CallInput{
			Method:        "DELETE",
			Url:           url,
			CustomHeaders: common.BuildCustomHeaders(),
		})
	return resp, err

}

func (s *service) RenamePasskeyHTTP(ctx context.Context, request *models_api_passkey.PasskeyRenameRequest) (*common.WrappedResonseT[*models_api_passkey.PasskeyRenameResponse], error) {
	url := s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys) + "/" + request.CredentialID
	requestBody := map[string]string{"friendlyName": request.FriendlyName}

	resp, err := common.HTTPFetchWrappedResponseT[*models_api_passkey.PasskeyRenameResponse](ctx,
		&common.CallInput{
			Method:        "PATCH",
			Url:           url,
			CustomHeaders: common.BuildCustomHeaders(),
			Data:          requestBody,
		})
	return resp, err
}

func (s *service) ClearSSOCookie(ctx context.Context) (*common.WrappedResonseT[models_api_preferences.ClearSSOResponse], error) {
	resp, err := common.HTTPFetchWrappedResponseT[models_api_preferences.ClearSSOResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_ClearSSO),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	return resp, err
}

func (s *service) GetKeepSignedInPreference(ctx context.Context) (*common.WrappedResonseT[models_api_preferences.GetKeepSignedInPreferenceResponse], error) {
	resp, err := common.HTTPFetchWrappedResponseT[models_api_preferences.GetKeepSignedInPreferenceResponse](ctx,
		&common.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_KeepSignedInPreference),
			CustomHeaders: common.BuildCustomHeaders(),
		})
	return resp, err
}

func (s *service) UpdateKeepSignedInPreference(ctx context.Context, request *models_api_preferences.UpdateKeepSignedInPreferenceRequest) (*common.WrappedResonseT[models_api_preferences.UpdateKeepSignedInPreferenceResponse], error) {
	resp, err := common.HTTPFetchWrappedResponseT[models_api_preferences.UpdateKeepSignedInPreferenceResponse](ctx,
		&common.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_KeepSignedInPreference),
			Data:          request,
			CustomHeaders: common.BuildCustomHeaders(),
		})
	return resp, err
}
