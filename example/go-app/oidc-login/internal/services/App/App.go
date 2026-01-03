package App

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/config"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		appConfigAccessor contracts_config.IAppConfigAccessor
		rageApiClient     contracts_go_app_RageApiClient.IRageApiClient
		// Composers
		homeComposer           contracts_App.IHomeComposer
		passwordComposer       contracts_App.IPasswordComposer
		createAccountComposer  contracts_App.ICreateAccountComposer
		forgotPasswordComposer contracts_App.IForgotPasswordComposer
		resetPasswordComposer  contracts_App.IResetPasswordComposer
		verifyCodeComposer     contracts_App.IVerifyCodeComposer
		keepSignedInComposer   contracts_App.IKeepSignedInComposer

		currentPage      contracts_routes.WellknownRoute
		showCookieBanner bool

		manifest *models_api_manifest.Manifest
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IApp = stemService

func (s *service) Ctor(
	container di.Container,
	appConfigAccessor contracts_config.IAppConfigAccessor,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	homeComposer contracts_App.IHomeComposer,
	passwordComposer contracts_App.IPasswordComposer,
	createAccountComposer contracts_App.ICreateAccountComposer,
	forgotPasswordComposer contracts_App.IForgotPasswordComposer,
	resetPasswordComposer contracts_App.IResetPasswordComposer,
	verifyCodeComposer contracts_App.IVerifyCodeComposer,
	keepSignedInComposer contracts_App.IKeepSignedInComposer,
) (contracts_App.IApp, error) {

	return &service{
		appConfigAccessor: appConfigAccessor,
		rageApiClient:     rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			AppContext: appContext,
			Container:  container,
			Localizer:  localizer,
		},
		homeComposer:           homeComposer,
		passwordComposer:       passwordComposer,
		createAccountComposer:  createAccountComposer,
		forgotPasswordComposer: forgotPasswordComposer,
		verifyCodeComposer:     verifyCodeComposer,
		resetPasswordComposer:  resetPasswordComposer,
		keepSignedInComposer:   keepSignedInComposer,
	}, nil
}

func AddScopedIApp(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IApp](cb, stemService.Ctor)
}

func (s *service) SetCurrentPage(page contracts_routes.WellknownRoute) {
	s.currentPage = page
}

func (s *service) GetCurrentPage() contracts_routes.WellknownRoute {
	return s.currentPage
}

func (s *service) handleAcceptCookies(ctx app.Context, e app.Event) {
	ctx.LocalStorage().Set("cookiesAccepted", true)
	s.showCookieBanner = false
}
