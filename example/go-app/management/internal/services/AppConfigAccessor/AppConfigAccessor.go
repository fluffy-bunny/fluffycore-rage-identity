package AppConfigAccessor

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	management_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/config"
	fluffycore_go_app_js_loader "github.com/fluffy-bunny/fluffycore/go-app/js_loader"
)

type (
	service struct {
		appConfig *management_contracts_config.AppConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_config.IAppConfigAccessor = stemService

func (s *service) Ctor() (contracts_config.IAppConfigAccessor, error) {
	config, err := fluffycore_go_app_js_loader.LoadConfigFromJS[management_contracts_config.AppConfig](
		&fluffycore_go_app_js_loader.LoadConfigOptions{
			IsReadyFuncName:   "isAppConfigReady",
			GetConfigFuncName: "getAppConfig",
		},
	)
	if err != nil {
		return nil, err
	}
	return &service{
		appConfig: config,
	}, nil
}

func AddScopedIAppConfigAccessor(cb di.ContainerBuilder) {
	di.AddScoped[contracts_config.IAppConfigAccessor](cb, stemService.Ctor)
}

func (s *service) GetAppConfig(ctx context.Context) *management_contracts_config.AppConfig {
	return s.appConfig
}
