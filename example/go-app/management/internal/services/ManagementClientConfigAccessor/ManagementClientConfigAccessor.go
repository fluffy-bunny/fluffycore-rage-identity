package ManagementClientConfigAccessor

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/config"
)

type (
	service struct {
		appConfigAccessor contracts_config.IAppConfigAccessor
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor = stemService

func (s *service) Ctor(
	appConfigAccessor contracts_config.IAppConfigAccessor,
) (contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor, error) {
	return &service{
		appConfigAccessor: appConfigAccessor,
	}, nil
}
func AddScopedIManagementClientConfigAccessor(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_ManagementApiClient.IManagementClientConfigAccessor](cb, stemService.Ctor)
}
func (s *service) GetManagementClientConfig(ctx context.Context) *contracts_go_app_ManagementApiClient.ManagementClientConfig {
	cc := s.appConfigAccessor.GetAppConfig(ctx)
	return &contracts_go_app_ManagementApiClient.ManagementClientConfig{
		ApiBaseUrl: cc.AccountManagementBaseURL,
	}
}
