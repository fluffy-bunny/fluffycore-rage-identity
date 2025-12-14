package ManagementApiClient

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_go_app_cookies "github.com/fluffy-bunny/fluffycore/go-app/cookies"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
)

type (
	service struct {
		managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_ManagementApiClient.IManagementApiClient = stemService

func (s *service) Ctor(
	managementClientConfigAccessor contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor,
) (contracts_go_app_ManagementApiClient.IManagementApiClient, error) {
	return &service{
		managementClientConfigAccessor: managementClientConfigAccessor,
	}, nil
}

func AddScopedIManagementApiClient(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_ManagementApiClient.IManagementApiClient](cb, stemService.Ctor)
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

func (s *service) fixupApiPath(ctx context.Context, relativePath string) string {
	appConfig := s.managementClientConfigAccessor.GetManagementClientConfig(ctx)
	if appConfig == nil {
		panic("AppConfig is not set in AppConfigAccessor")
	}
	return appConfig.ApiBaseUrl + relativePath
}

func (s *service) Login(ctx context.Context, request *models_api_login_models.LoginRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Login),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LogoutResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_Logout),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyCodeBegin(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_verify_code.VerifyCodeBeginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_verify_code.VerifyCodeBeginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_VerifyCodeBegin),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) VerifyCode(ctx context.Context, request *models_api_login_models.VerifyCodeRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.VerifyCodeResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.VerifyCodeResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_VerifyCode),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) PasswordResetStart(ctx context.Context, request *models_api_login_models.PasswordResetStartRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetStartResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetStartResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetStart),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) PasswordResetFinish(ctx context.Context, request *models_api_login_models.PasswordResetFinishRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.PasswordResetFinishResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.PasswordResetFinishResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_PasswordResetFinish),
			Data:          request,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetUserProfile(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.Profile](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) UpdateUserProfile(ctx context.Context, profile *models.Profile) (*fluffycore_go_app_fetch.WrappedResonseT[models.Profile], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.Profile](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "POST",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserProfilePath),
			Data:          profile,
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetUserLinkedAccounts(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models.LinkedAccountsResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.LinkedAccountsResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts"),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) DeleteUserLinkedAccount(ctx context.Context, identity string) (*fluffycore_go_app_fetch.WrappedResonseT[models.DeleteLinkedAccountResponse], error) {

	return fluffycore_go_app_fetch.FetchWrappedResponseT[models.DeleteLinkedAccountResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "DELETE",
			Url:           s.fixupApiPath(ctx, "/api/linked-accounts") + "?identity=" + identity,
			CustomHeaders: buildCustomHeaders(),
		})
}
