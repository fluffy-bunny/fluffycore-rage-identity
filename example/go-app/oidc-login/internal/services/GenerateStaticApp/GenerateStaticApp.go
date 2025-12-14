package GenerateStaticApp

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	service struct {
		app.Compo
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IApp = stemService

func (s *service) Ctor() (contracts_App.IApp, error) {

	return &service{}, nil
}

func AddScopedIApp(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IApp](cb, stemService.Ctor)
}
func (s *service) Render() app.UI {
	return app.Div()
}
func (s *service) SetCurrentPage(page contracts_routes.WellknownRoute) {}
func (s *service) GetCurrentPage() contracts_routes.WellknownRoute {
	return ""
}
func (s *service) OnMount(ctx app.Context) {}
