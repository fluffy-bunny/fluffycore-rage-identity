package RageClientConfigAccessor

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/config"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
)

type (
	service struct {
		appConfigAccessor contracts_config.IAppConfigAccessor
	}
)

var stemService = (*service)(nil)

var _ contracts_go_app_RageApiClient.IRageClientConfigAccessor = stemService

func (s *service) Ctor(
	appConfigAccessor contracts_config.IAppConfigAccessor,
) (contracts_go_app_RageApiClient.IRageClientConfigAccessor, error) {
	return &service{
		appConfigAccessor: appConfigAccessor,
	}, nil
}
func AddScopedIRageClientConfigAccessor(cb di.ContainerBuilder) {
	di.AddScoped[contracts_go_app_RageApiClient.IRageClientConfigAccessor](cb, stemService.Ctor)
}
func (s *service) GetRageClientConfig(ctx context.Context) *contracts_go_app_RageApiClient.RageClientConfig {
	cc := s.appConfigAccessor.GetAppConfig(ctx)
	return &contracts_go_app_RageApiClient.RageClientConfig{
		ApiBaseUrl: cc.RageBaseURL,
	}
}
