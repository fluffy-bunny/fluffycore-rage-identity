package RageApiClient

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	models_api_external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	models_api_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	models_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/oidc_login_ui_server/contracts/go-app/RageApiClient"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
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
			Method: "GET",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_OIDCFlowAppConfig),
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
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_LoginPhaseOne),
			Data:   request,
		})
}

func (s *service) VerifyPasswordStrength(ctx context.Context, request *models_api_password.VerifyPasswordStrengthRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_password.VerifyPasswordStrengthResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_password.VerifyPasswordStrengthResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_VerifyPasswordStrength),
			Data:   request,
		})
}

func (s *service) GetManifest(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_manifest.Manifest], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_manifest.Manifest](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "GET",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_Manifest),
		})
}

func (s *service) Signup(ctx context.Context, request *models_api_login_models.SignupRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.SignupResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.SignupResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_Signup),
			Data:   request,
		})
}

func (s *service) LoginPassword(ctx context.Context, request *models_api_login_models.LoginPasswordRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginPasswordResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginPasswordResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_LoginPassword),
			Data:   request,
		})
}
func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetStart),
			Data:   request,
		})
}
func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:   request,
		})
}
func (s *service) VerifyUserName(ctx context.Context, request *models_api_verify_username.VerifyUsernameRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_username.VerifyUsernameResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_verify_username.VerifyUsernameResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_VerifyUsername),
			Data:   request,
		})
}

func (s *service) VerifyCodeBegin(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_verify_code.VerifyCodeBeginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "GET",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCodeBegin),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_VerifyCode),
			Data:   request,
		})
}

func (s *service) StartExternalLogin(ctx context.Context, request *models_api_external_idp.StartExternalIDPLoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_external_idp.StartExternalIDPLoginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_external_idp.StartExternalIDPLoginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixUpRageApi(ctx, wellknown_echo.API_Start_ExternalLogin),
			Data:   request,
		})
}
