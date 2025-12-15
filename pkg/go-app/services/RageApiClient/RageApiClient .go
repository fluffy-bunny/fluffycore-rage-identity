package RageApiClient

import (
	"context"
	"encoding/json"
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	models_api_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	models_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_go_app_cookies "github.com/fluffy-bunny/fluffycore/go-app/cookies"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
)

type (
	service struct {
		cachedRageConfigResponse *fluffycore_go_app_fetch.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig]
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

// getCSRFToken retrieves the CSRF token from the _csrf cookie
func getCSRFToken() string {
	csrfToken, err := fluffycore_go_app_cookies.GetCookie[string]("_csrf")
	if err != nil {
		return ""
	}
	return csrfToken
}

// buildCustomHeaders creates custom headers map with CSRF token if it exists
func buildCustomHeaders() map[string]string {
	csrfToken := getCSRFToken()
	if csrfToken != "" {
		return map[string]string{
			"X-Csrf-Token": csrfToken,
		}
	}
	return nil
}

func (s *service) fixUpRageApi(ctx context.Context, relativePath string) string {
	rageClientConfig := s.rageClientConfigAccessor.GetRageClientConfig(ctx)
	if rageClientConfig == nil {
		panic("AppConfig is not set in AppConfigAccessor")
	}
	return rageClientConfig.ApiBaseUrl + relativePath
}
func (s *service) GetOIDCFlowAppConfig(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig], error) {
	if s.cachedRageConfigResponse != nil {
		return s.cachedRageConfigResponse, nil
	}
	var err error
	s.cachedRageConfigResponse, err = fluffycore_go_app_fetch.FetchWrappedResponseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_OIDCFlowAppConfig),
			CustomHeaders: buildCustomHeaders(),
		})

	return s.cachedRageConfigResponse, err
}

func (s *service) GetCachedManifest(ctx context.Context) *models_api_manifest.Manifest {
	return s.cachedManifest
}
func (s *service) SetCachedManifest(ctx context.Context, manifest *models_api_manifest.Manifest) {
	s.cachedManifest = manifest
}

func (s *service) LoginPhaseOne(ctx context.Context, request *models_api_login_models.LoginPhaseOneRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginPhaseOneResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginPhaseOneResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_LoginPhaseOne),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyPasswordStrength(ctx context.Context, request *models_api_password.VerifyPasswordStrengthRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_password.VerifyPasswordStrengthResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_password.VerifyPasswordStrengthResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyPasswordStrength),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetManifest(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_manifest.Manifest], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_manifest.Manifest](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Manifest),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) Signup(ctx context.Context, request *models_api_login_models.SignupRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.SignupResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.SignupResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Signup),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) LoginPassword(ctx context.Context, request *models_api_login_models.LoginPasswordRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginPasswordResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginPasswordResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_LoginPassword),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}
func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetStart),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}
func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}
func (s *service) VerifyUserName(ctx context.Context, request *models_api_verify_username.VerifyUsernameRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_username.VerifyUsernameResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_verify_username.VerifyUsernameResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyUsername),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyCodeBegin(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_verify_code.VerifyCodeBeginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCodeBegin),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCode),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) StartExternalLogin(ctx context.Context, request *models_api_external_idp.StartExternalIDPLoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_external_idp.StartExternalIDPLoginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_external_idp.StartExternalIDPLoginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Start_ExternalLogin),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

// TOTP/Authenticator APIs
func (s *service) GetTOTPStatus(ctx context.Context) ([]byte, error) {
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTP),
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) EnrollTOTP(ctx context.Context) ([]byte, error) {
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTPEnroll),
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) VerifyTOTP(ctx context.Context, code string) ([]byte, error) {
	request := map[string]string{"code": code}
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTPVerify),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) DisableTOTP(ctx context.Context) ([]byte, error) {
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "DELETE",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_UserTOTP),
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

// Passkey APIs
func (s *service) GetPasskeys(ctx context.Context) ([]byte, error) {
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys),
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) DeletePasskey(ctx context.Context, credentialID string) ([]byte, error) {
	// credentialID is already base64-encoded from the GET response
	url := s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys) + "/" + credentialID
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "DELETE",
			Url:           url,
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Response)
}

func (s *service) RenamePasskey(ctx context.Context, credentialID string, friendlyName string) ([]byte, error) {
	fmt.Printf("🔵 RenamePasskey called - credentialID (already base64): %s, friendlyName: %s\n", credentialID, friendlyName)

	// credentialID is already base64-encoded from the GET response, use it directly
	request := map[string]string{"friendlyName": friendlyName}
	url := s.fixUpRageApi(ctx, wellknown_echo.API_Passkeys) + "/" + credentialID
	fmt.Printf("🔵 RenamePasskey URL: %s\n", url)
	fmt.Printf("🔵 RenamePasskey request body: %+v\n", request)
	resp, err := fluffycore_go_app_fetch.FetchWrappedResponseT[map[string]interface{}](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "PATCH",
			Url:           url,
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
	if err != nil {
		fmt.Printf("❌ RenamePasskey error: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ RenamePasskey response: %+v\n", resp)
	return json.Marshal(resp.Response)
}
