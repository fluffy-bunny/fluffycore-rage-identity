package App

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/config"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	services_WellknownCookieNames "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/services/WellknownCookieNames"
	models_api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		appConfigAccessor   contracts_config.IAppConfigAccessor
		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		// Composers
		homeComposer            contracts_App.IHomeComposer
		profileComposer         contracts_App.IProfileComposer
		passwordManagerComposer contracts_App.IPasswordManagerComposer
		passkeyManagerComposer  contracts_App.IPasskeyManagerComposer
		linkedAccountsComposer  contracts_App.ILinkedAccountsComposer

		currentPage      contracts_routes.WellknownRoute
		showCookieBanner bool
		showUserMenu     bool
		showSidebar      bool
		isAuthenticated  bool
		isClaimedDomain  bool
		profile          *models_api_profile.Profile
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IApp = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,

	appConfigAccessor contracts_config.IAppConfigAccessor,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,

	homeComposer contracts_App.IHomeComposer,
	profileComposer contracts_App.IProfileComposer,
	passwordManagerComposer contracts_App.IPasswordManagerComposer,
	passkeyManagerComposer contracts_App.IPasskeyManagerComposer,
	linkedAccountsComposer contracts_App.ILinkedAccountsComposer,

) (contracts_App.IApp, error) {

	// CRITICAL: Set cookie config IMMEDIATELY in constructor before any component OnMount calls
	appConfig := appConfigAccessor.GetAppConfig(appContext)
	services_WellknownCookieNames.WellknownCookieNamesConfig = appConfig.WellknownCookieNamesConfig

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		appConfigAccessor:       appConfigAccessor,
		managementApiClient:     managementApiClient,
		homeComposer:            homeComposer,
		profileComposer:         profileComposer,
		passwordManagerComposer: passwordManagerComposer,
		passkeyManagerComposer:  passkeyManagerComposer,
		linkedAccountsComposer:  linkedAccountsComposer,
		isAuthenticated:         false, // Default to unauthenticated
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

func (s *service) IsAuthenticated() bool {
	return s.isAuthenticated
}

func (s *service) LoginWithReturnURL(returnURL string) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Str("returnURL", returnURL).Msg("LoginWithReturnURL called")

	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)

	app.Window().Call("appAsyncFunc", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		go func() {
			// Call login API with return URL
			response, err := s.managementApiClient.Login(s.AppContext,
				&models_api_login_models.LoginRequest{
					ReturnURL: returnURL,
				})
			if err != nil {
				log.Error().Err(err).Msg("login failed")
				return
			}

			if response != nil {
				s.checkUnauthorizedAndReload(response.Code)

				switch response.Code {
				case 404:
					log.Error().Msg("login returned 404")
					return
				}

				// Check if we got a redirect URL in the response
				if response.Response != nil && response.Response.RedirectURL != "" {
					log.Info().Str("redirectURL", response.Response.RedirectURL).Msg("Redirecting to login URL")
					app.Window().Get("location").Set("href", response.Response.RedirectURL)
					return
				}
			}

			// Fallback to home page if no redirect URL
			app.Window().Get("location").Set("href", appConfig.BaseHREF)
		}()
		return nil
	}))
}

func (s *service) handleAcceptCookies(ctx app.Context, e app.Event) {
	ctx.LocalStorage().Set("cookiesAccepted", true)
	s.showCookieBanner = false
	ctx.Update()
}

func (s *service) isHomePath(path string) bool {
	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)
	baseHREF := strings.Trim(appConfig.BaseHREF, "/")

	// baseHREF is a string
	return path == "/" || path == "/"+baseHREF || path == "/"+baseHREF+"/"
}

// checkUnauthorizedAndReload checks if the response code indicates unauthorized (401/403) and redirects to home
// Only redirects if we're not already on the home page to avoid infinite redirect loops
func (s *service) checkUnauthorizedAndReload(code int) {
	if code == 401 || code == 403 {
		log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()

		// Check current path - don't redirect if we're already on the home page
		currentPath := app.Window().URL().Path
		log.Debug().Str("currentPath", currentPath).Int("code", code).Msg("Checking unauthorized response")

		// If we're on home page (/, /management, /management/), don't redirect to avoid infinite loop
		if s.isHomePath(currentPath) {
			log.Debug().Msg("Already on home page, skipping redirect")
			return
		}

		log.Info().Int("code", code).Str("currentPath", currentPath).Msg("Unauthorized response detected, redirecting to home")
		// Navigate to management home page
		appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)
		baseURL := app.Window().Get("location").Get("origin").String() + "/" + appConfig.BaseHREF + "/"
		app.Window().Get("location").Call("assign", baseURL)
	}
}

func (s *service) handleSignin(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("handleEmailSubmit called")

	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)
	e.PreventDefault()

	ctx.Async(func() {
		// Call login API
		response, err := s.managementApiClient.Login(s.AppContext,
			&models_api_login_models.LoginRequest{
				ReturnURL: appConfig.ReturnURL,
			})
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Error().Err(err).Msg("login failed")
				return
			}

			if response != nil {
				s.checkUnauthorizedAndReload(response.Code)

				switch response.Code {
				case 404:
					log.Error().Msg("login returned 404")
					return
				}

				// Check if we got a redirect URL in the response
				if response.Response != nil && response.Response.RedirectURL != "" {
					log.Info().Str("redirectURL", response.Response.RedirectURL).Msg("Redirecting to login URL")
					app.Window().Get("location").Set("href", response.Response.RedirectURL)
					return
				}
			}

			// Fallback to home page if no redirect URL
			ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
		})
		ctx.Update()
	})
}

func (s *service) handleToggleSidebar(ctx app.Context, e app.Event) {
	e.PreventDefault()
	s.showSidebar = !s.showSidebar
	ctx.Update()
}

func (s *service) handleSignOut(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("handleSignOut called")

	e.PreventDefault()

	ctx.Async(func() {
		// Call logout API
		ctx.LocalStorage().Del("email")
		response, err := s.managementApiClient.Logout(s.AppContext,
			&models_api_login_models.LogoutRequest{})
		ctx.Dispatch(func(ctx app.Context) {

			if response != nil {
				s.checkUnauthorizedAndReload(response.Code)

				switch response.Code {
				case 404:
					log.Error().Msg("logout returned 404")
					return
				}
			}
			if err != nil {
				log.Error().Err(err).Msg("logout failed")
				return
			}

			// Navigate to home page with full reload to reset app state
			log.Info().Msg("Logout successful, navigating to home page")
			// Use location.assign to navigate to base path (e.g., /management/)
			// This will trigger a full page reload at the home page
			baseURL := app.Window().Get("location").Get("origin").String() + "/management/"
			app.Window().Get("location").Call("assign", baseURL)
		})
		ctx.Update()

	})

}
