package WASMRuntime

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/common"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/App"
	service_AppConfigAccessor "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/AppConfigAccessor"
	services_GenerateStaticApp "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/GenerateStaticApp"
	services_RageClientConfigAccessor "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/RageClientConfigAccessor"
	services_composers_CreateAccount "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/CreateAccount"
	services_composers_ForgotPassword "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/ForgotPassword"
	services_composers_Home "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/Home"
	services_composers_KeepSignedIn "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/KeepSignedIn"
	services_composers_Password "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/Password"
	services_composers_ResetPassword "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/ResetPassword"
	services_composers_VerifyCode "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/composers/VerifyCode"
	servies_i18n_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/i18n/Localizer"
	services_i18n_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/i18n/LocalizerBundle"
	services_logging "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/logging"
	services_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/services/RageApiClient"
	services_WellknownCookieNames "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/services/WellknownCookieNames"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

func RegisterGenerateStaticAppService(ctx context.Context, cb di.ContainerBuilder) {
	services_GenerateStaticApp.AddScopedIApp(cb)
}
func RegisterServices(ctx context.Context, cb di.ContainerBuilder) {
	// Register services here
	service_AppConfigAccessor.AddScopedIAppConfigAccessor(cb)
	services_RageClientConfigAccessor.AddScopedIRageClientConfigAccessor(cb)
	services_App.AddScopedIApp(cb)
	servies_i18n_Localizer.AddScopedILocalizer(cb)
	services_i18n_LocalizerBundle.AddScopedILocalizerBundle(cb)
	services_composers_Home.AddScopedIHomeComposer(cb)
	services_composers_Password.AddScopedIPasswordComposer(cb)
	services_composers_CreateAccount.AddScopedICreateAccountComposer(cb)
	services_composers_ForgotPassword.AddScopedIForgotPasswordComposer(cb)
	services_composers_ResetPassword.AddScopedIResetPasswordComposer(cb)
	services_composers_VerifyCode.AddScopedIVerifyCodeComposer(cb)
	services_composers_KeepSignedIn.AddScopedIKeepSignedInComposer(cb)
	services_go_app_RageApiClient.AddScopedIRageApiClient(cb)
	services_WellknownCookieNames.AddSingletonIWellknownCookieNames(cb)

	var appContext contracts_App.AppContext = ctx
	di.AddInstance[contracts_App.AppContext](cb, appContext)

	log := zerolog.Ctx(ctx).With().Logger()
	di.AddInstance[*zerolog.Logger](cb, &log)
}

var diContainer di.Container

func NewWASMApp(ctx context.Context, generateStaticMode bool) {
	// Initialize logging service early to check localStorage
	if !generateStaticMode {
		services_logging.GetInstance()
	}
	log := zerolog.Ctx(ctx).With().Logger()

	// Log build version info
	log.Info().
		Str("version", common.AppVersion).
		Str("buildTime", common.BuildTime).
		Str("gitCommit", common.GitCommit).
		Str("gitBranch", common.GitBranch).
		Msg("Starting WASM App")

	cb := di.Builder()
	if generateStaticMode {
		RegisterGenerateStaticAppService(ctx, cb)
	} else {
		RegisterServices(ctx, cb)

	}
	diContainer = cb.Build()

	scopeFactory := di.Get[di.ScopeFactory](diContainer)

	// Create a single scope and wizard app instance (singleton)

	type NewScopedWizardAppFunc func(page contracts_routes.WellknownRoute) app.Composer
	var newScopedWizardApp NewScopedWizardAppFunc
	if generateStaticMode {
		// we have to create a new instance per route.
		newScopedWizardApp = func(page contracts_routes.WellknownRoute) app.Composer {
			scope := scopeFactory.CreateScope()
			c := scope.Container()
			myApp := di.Get[contracts_App.IApp](c)
			myApp.SetCurrentPage(page)
			return myApp.(app.Composer)
		}
	} else {
		scope := scopeFactory.CreateScope()
		c := scope.Container()
		myApp := di.Get[contracts_App.IApp](c)
		newScopedWizardApp = func(page contracts_routes.WellknownRoute) app.Composer {
			myApp.SetCurrentPage(page)
			return myApp.(app.Composer)
		}
	}
	fixupRoute := func(route contracts_routes.WellknownRoute) string {
		return string(route)
	}
	// Register home route LAST (it's a catch-all)
	route := fixupRoute(contracts_routes.WellknownRoute_Home)
	app.Route(route, func() app.Composer {
		log.Info().Msg("Routing to " + route)
		return newScopedWizardApp(contracts_routes.WellknownRoute_Home)
	})
	if !generateStaticMode {

		route = fixupRoute(contracts_routes.WellknownRoute_ResetPassword)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_ResetPassword)
		})
		// Register more specific routes FIRST (before "/" catch-all)
		route = fixupRoute(contracts_routes.WellknownRoute_CreateAccount)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_CreateAccount)
		})

		route = fixupRoute(contracts_routes.WellknownRoute_Password)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_Password)
		})

		route = fixupRoute(contracts_routes.WellknownRoute_ForgotPassword)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_ForgotPassword)
		})

		route = fixupRoute(contracts_routes.WellknownRoute_VerifyCode)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_VerifyCode)
		})

		route = fixupRoute(contracts_routes.WellknownRoute_KeepSignedIn)
		app.Route(route, func() app.Composer {
			log.Info().Msg("Routing to " + route)
			return newScopedWizardApp(contracts_routes.WellknownRoute_KeepSignedIn)
		})

	}

}
