package Home

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/config"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		appConfigAccessor   contracts_config.IAppConfigAccessor
		isAuthenticated     bool
		isClaimedDomain     bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IHomeComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	appConfigAccessor contracts_config.IAppConfigAccessor,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,

) (contracts_App.IHomeComposer, error) {

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		appConfigAccessor:   appConfigAccessor,
		managementApiClient: managementApiClient,
		isAuthenticated:     false,
	}, nil
}

func AddScopedIHomeComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IHomeComposer](cb, stemService.Ctor)
}

func (s *service) Render() app.UI {
	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)

	return app.Div().Class("home-container").Body(
		// Hero Section
		app.Div().Class("home-hero").Body(
			app.Div().Class("home-hero-content").Body(
				app.H1().Class("home-title").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyWelcomeToRageAccounts)),
				app.P().Class("home-subtitle").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySecurelyManageYourIdentity)),
			),
		),

		// Feature Cards
		app.Div().Class("home-features").Body(
			s.renderFeatureCard(
				"security-icon",
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySecureAuthentication),
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyProtectYourAccountWithAdvancedSecurity),
				`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
					</svg>`,
				contracts_routes.WellknownRoute_PasswordManager,
			),
			app.If(appConfig.EnabledWebAuthN, func() app.UI {
				return s.renderFeatureCard(
					"passkey-icon",
					s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasskeys),
					s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyManageYourPasskeys),
					`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<circle cx="7" cy="7" r="2"></circle>
						<path d="M7 9v4a2 2 0 0 0 2 2h4"></path>
						<circle cx="19" cy="15" r="4"></circle>
						<path d="M19 11v-1"></path>
						<path d="M22 15h-1"></path>
					</svg>`,
					contracts_routes.WellknownRoute_PasskeyManager,
				)
			}),
			s.renderFeatureCard(
				"accounts-icon",
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLinkedAccounts),
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConnectMultipleAccountsInOnePlace),
				`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
					<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
				</svg>`,
				contracts_routes.WellknownRoute_LinkedAccounts,
			),
			s.renderFeatureCard(
				"privacy-icon",
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyYourPrivacyMatters),
				s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyWeKeepYourDataSafeAndPrivate),
				`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
					<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
				</svg>`,
				contracts_routes.WellknownRoute_Profile,
			),
		),

		// CTA Section
		app.Div().Class("home-cta").Body(
			app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyReadyToGetStarted)),
			app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignInToAccessYourAccount)),
		),
	)
}

func (s *service) renderFeatureCard(iconClass, title, description, iconSVG string, route contracts_routes.WellknownRoute) app.UI {
	return app.Div().
		Class("home-feature-card").
		OnClick(func(ctx app.Context, e app.Event) {
			e.PreventDefault()

			// Check if user is authenticated
			if s.isAuthenticated {
				// If authenticated, navigate normally
				ctx.Navigate(contracts_routes.GetFixedRoute(route))
			} else {
				// If not authenticated, initiate login with return URL
				s.handleLoginWithReturnURL(ctx, route)
			}
		}).
		Body(
			app.Div().Class("home-feature-icon "+iconClass).Body(
				app.Raw(iconSVG),
			),
			app.H3().Class("home-feature-title").Text(title),
			app.P().Class("home-feature-description").Text(description),
		)
}

func (s *service) handleLoginWithReturnURL(ctx app.Context, route contracts_routes.WellknownRoute) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	returnURL := contracts_routes.GetFixedRoute(route)
	log.Info().Str("returnURL", returnURL).Msg("Initiating login with return URL")

	ctx.Async(func() {
		// Call login API with return URL
		response, err := s.managementApiClient.Login(s.AppContext,
			&models_api_login_models.LoginRequest{
				ReturnURL: returnURL,
			})
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Error().Err(err).Msg("login failed")
				return
			}

			if response != nil {
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
		})
	})
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Home page mounted")

	// Check authentication status on mount
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserProfile(s.AppContext)
		ctx.Dispatch(func(ctx app.Context) {
			if err == nil && response != nil && response.Code == 200 {
				s.isAuthenticated = true
				if response.Response != nil {
					s.isClaimedDomain = response.Response.IsClaimedDomain
				}
				log.Info().Bool("isClaimedDomain", s.isClaimedDomain).Msg("User is authenticated")
			} else {
				s.isAuthenticated = false
				log.Info().Msg("User is not authenticated")
			}
			ctx.Update()
		})
	})
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Home page navigated")
}

func (s *service) OnDismount() {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Home page dismounted")
}
