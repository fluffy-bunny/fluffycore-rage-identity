package ManagementApiClient

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/account_management_server/contracts/go-app/ManagementApiClient"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
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

func (s *service) fixupApiPath(ctx context.Context, relativePath string) string {
	appConfig := s.managementClientConfigAccessor.GetManagementClientConfig(ctx)
	if appConfig == nil {
		panic("AppConfig is not set in AppConfigAccessor")
	}
	return appConfig.ApiBaseUrl + relativePath
}

func (s *service) Login(ctx context.Context) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LoginResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LoginResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "GET",
			Url:    s.fixupApiPath(ctx, wellknown_echo.API_Login),
		})
}

func (s *service) Logout(ctx context.Context, request *models_api_login_models.LogoutRequest) (*fluffycore_go_app_fetch.WrappedResonseT[models_api_login_models.LogoutResponse], error) {
	return fluffycore_go_app_fetch.FetchWrappedResponseT[models_api_login_models.LogoutResponse](ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "POST",
			Url:    s.fixupApiPath(ctx, wellknown_echo.API_Logout),
			Data:   request,
		})
}
