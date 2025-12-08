package ManagementApiClient

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/account_management_server/contracts/go-app/ManagementApiClient"
	models_api_user_linked_accounts "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_user_linked_accounts"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
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

func (s *service) GetUserInfo(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.UserInfoResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.UserInfoResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_IsAuthorized),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) GetUserLinkedAccounts(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_user_linked_accounts.UserLinkedAccounts], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_user_linked_accounts.UserLinkedAccounts](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "GET",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserLinkedAccounts),
			CustomHeaders: buildCustomHeaders(),
		})
}

func (s *service) DeleteUserLinkedAccount(ctx context.Context, identity string) (*fluffycore_go_app_fetch.WrappedResonseT[string], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[string](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method:        "DELETE",
			Url:           s.fixupApiPath(ctx, wellknown_echo.API_UserLinkedAccounts+"/"+identity),
			CustomHeaders: buildCustomHeaders(),
		})
}
