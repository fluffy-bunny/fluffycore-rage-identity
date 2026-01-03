package WASMRuntime

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	services_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/App"
	services_AppConfigAccessor "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/AppConfigAccessor"
	services_GenerateStaticApp "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/GenerateStaticApp"
	services_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ManagementApiClient"
	services_ManagementClientConfigAccessor "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ManagementClientConfigAccessor"
	services_RageClientConfigAccessor "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/RageClientConfigAccessor"
	services_composers_Home "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/composers/Home"
	services_composers_LinkedAccounts "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/composers/LinkedAccounts"
	services_composers_PasskeyManager "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/composers/PasskeyManager"
	services_composers_PasswordManager "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/composers/PasswordManager"
	services_composers_Profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/composers/Profile"
	servies_i18n_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/i18n/Localizer"
	services_i18n_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/i18n/LocalizerBundle"
	services_logging "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/logging"
	services_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/services/RageApiClient"
	services_WellknownCookieNames "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/services/WellknownCookieNames"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

func RegisterGenerateStaticServices(ctx context.Context, cb di.ContainerBuilder) {
	var appContext contracts_App.AppContext = ctx
	di.AddInstance[contracts_App.AppContext](cb, appContext)
	services_GenerateStaticApp.AddScopedIApp(cb)
}
func RegisterServices(ctx context.Context, cb di.ContainerBuilder) {
	// Register services here
	services_App.AddScopedIApp(cb)
	servies_i18n_Localizer.AddScopedILocalizer(cb)
	services_i18n_LocalizerBundle.AddScopedILocalizerBundle(cb)
	services_composers_Home.AddScopedIHomeComposer(cb)
	services_composers_Profile.AddScopedIProfileComposer(cb)
	services_composers_PasswordManager.AddScopedIPasswordManagerComposer(cb)
	services_composers_PasskeyManager.AddScopedIPasskeyManagerComposer(cb)
	services_composers_LinkedAccounts.AddScopedILinkedAccountsComposer(cb)
	services_AppConfigAccessor.AddScopedIAppConfigAccessor(cb)
	services_RageClientConfigAccessor.AddScopedIRageClientConfigAccessor(cb)
	services_ManagementClientConfigAccessor.AddScopedIManagementClientConfigAccessor(cb)
	// Register RageApiClient first, as ManagementApiClient depends on it
	services_go_app_RageApiClient.AddScopedIRageApiClient(cb)
	services_go_app_ManagementApiClient.AddScopedIManagementApiClient(cb)
	var appContext contracts_App.AppContext = ctx
	di.AddInstance[contracts_App.AppContext](cb, appContext)
	// when we load the wasm make sure we set the WellknownCookieNamesConfig as a global
	services_WellknownCookieNames.AddSingletonIWellknownCookieNames(cb)
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
		// register a empty app because for static generation we only want the js files and the index.html file.
		RegisterGenerateStaticServices(ctx, cb)
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
		// Register Profile route (must be before home route)
		routeProfile := fixupRoute(contracts_routes.WellknownRoute_Profile)
		app.Route(routeProfile, func() app.Composer {
			log.Info().Msg("Routing to " + routeProfile)
			return newScopedWizardApp(contracts_routes.WellknownRoute_Profile)
		})

		// Register PasswordManager route
		routePasswordManager := fixupRoute(contracts_routes.WellknownRoute_PasswordManager)
		app.Route(routePasswordManager, func() app.Composer {
			log.Info().Msg("Routing to " + routePasswordManager)
			return newScopedWizardApp(contracts_routes.WellknownRoute_PasswordManager)
		})

		// Register PasskeyManager route
		routePasskeyManager := fixupRoute(contracts_routes.WellknownRoute_PasskeyManager)
		app.Route(routePasskeyManager, func() app.Composer {
			log.Info().Msg("Routing to " + routePasskeyManager)
			return newScopedWizardApp(contracts_routes.WellknownRoute_PasskeyManager)
		})

		// Register LinkedAccounts route
		routeLinkedAccounts := fixupRoute(contracts_routes.WellknownRoute_LinkedAccounts)
		app.Route(routeLinkedAccounts, func() app.Composer {
			log.Info().Msg("Routing to " + routeLinkedAccounts)
			return newScopedWizardApp(contracts_routes.WellknownRoute_LinkedAccounts)
		})
	}

}
