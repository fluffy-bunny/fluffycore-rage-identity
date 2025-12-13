package Home

import (
	"regexp"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/config"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		email      string
		emailError string

		isLoading bool

		currentPage       string
		errorMessage      string
		showError         bool
		rageApiCliient    contracts_go_app_RageApiClient.IRageApiClient
		appConfigAccessor contracts_config.IAppConfigAccessor
		oidcFlowConfig    *contracts_OIDCFlowAppConfig.OIDCFlowAppConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IHomeComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	appConfigAccessor contracts_config.IAppConfigAccessor,
	rageApiCliient contracts_go_app_RageApiClient.IRageApiClient,

) (contracts_App.IHomeComposer, error) {

	return &service{
		appConfigAccessor: appConfigAccessor,
		rageApiCliient:    rageApiCliient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		email:        "",
		emailError:   "",
		currentPage:  "",
		errorMessage: "",
		showError:    false,
	}, nil
}

func AddScopedIHomeComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IHomeComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	oidcFlowConfig, err := s.appConfigAccessor.GetOIDCFlowAppConfig(s.AppContext)
	if err != nil {
		log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
		log.Error().Err(err).Msg("Failed to get OIDC Flow App Config")
		s.showErrorMessage(ctx, "Configuration error. Please try again later.")
	}
	s.oidcFlowConfig = oidcFlowConfig
}

func (s *service) renderCreateAccountUI() app.UI {
	return app.Div().Body(app.Div().Class("create-account").
		Body(
			app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyDontHaveAnAccount)+" "),
			app.A().
				Href(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_CreateAccount)).
				OnClick(s.handleCreateAccountClick).
				Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCreateOne)),
		),
		app.Div().Class("forgot-password").Body(
			app.A().
				Href(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_ForgotPassword)).
				OnClick(s.handleForgotPasswordClick).
				Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyForgotPassword)),
		),
	)
}
func (s *service) Render() app.UI {
	if s.oidcFlowConfig == nil {
		return app.Div().Text("Loading...")
	}

	return app.Div().Class("step-container").Body(
		app.If(s.showError, s.renderErrorBanner),

		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyWelcomeBack)),
		app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourEmailToContinue)),

		app.Form().OnSubmit(s.handleEmailSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("email").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
				app.Input().
					Type("email").
					ID("email").
					Value(s.email).
					OnInput(s.handleEmailInput).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourEmail)).
					Required(true),
			),
			app.If(s.emailError != "",
				func() app.UI {
					return app.Div().Class("error-message").Text(s.emailError)
				},
			),
			app.Div().Class("button-group").Body(
				app.Button().
					Type("submit").
					Class("btn-primary").
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyContinue)),
			),
		),
		app.If(
			!s.oidcFlowConfig.DisableSocialAccounts,
			s.renderCreateAccountUI,
		),
		app.If(
			!s.oidcFlowConfig.DisableSocialAccounts,
			s.renderSocialLogins,
		),
	)
}

func (s *service) renderErrorBanner() app.UI {

	return app.Div().Class("error-banner").Body(
		app.Div().Class("error-content").Body(
			app.Div().Class("error-icon").Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
				</svg>`),
			),
			app.Span().Class("error-text").Text(s.errorMessage),
			app.Button().
				Class("error-dismiss").
				OnClick(s.handleDismissError).
				Body(
					app.Raw(`<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
						<path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z"/>
					</svg>`),
				),
		),
	)
}

func (s *service) renderSocialLogins() app.UI {

	socialButtons := []app.UI{}
	for _, idp := range s.oidcFlowConfig.SocialIdps {
		switch idp.Slug {
		case "google-social":
			socialButtons = append(socialButtons, s.renderGoogleButton())
		case "microsoft-social":
			socialButtons = append(socialButtons, s.renderMicrosoftButton())
		case "github-social":
			socialButtons = append(socialButtons, s.renderGithubButton())
		}
	}
	return app.Div().Class("social-login-section").Body(
		app.Div().Class("divider").Body(
			app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyOrSignInWith)),
		),
		app.Div().Class("social-buttons").Body(
			socialButtons...,
		),
	)
}

func (s *service) renderGoogleButton() app.UI {
	return app.Button().
		Class("social-btn google-btn").
		OnClick(s.handleGoogleLogin).
		Body(
			app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
				<path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
				<path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
				<path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
			</svg>`),
			app.Span().Text("Google"),
		)
}

func (s *service) renderMicrosoftButton() app.UI {
	return app.Button().
		Class("social-btn microsoft-btn").
		OnClick(s.handleMicrosoftLogin).
		Body(
			app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path d="M11.4 11.4H2V2h9.4v9.4z" fill="#F25022"/>
				<path d="M22 11.4h-9.4V2H22v9.4z" fill="#7FBA00"/>
				<path d="M11.4 22H2v-9.4h9.4V22z" fill="#00A4EF"/>
				<path d="M22 22h-9.4v-9.4H22V22z" fill="#FFB900"/>
			</svg>`),
			app.Span().Text("Microsoft"),
		)
}

func (s *service) renderGithubButton() app.UI {
	return app.Button().
		Class("social-btn github-btn").
		OnClick(s.handleGithubLogin).
		Body(
			app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
				<path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.137 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z"/>
			</svg>`),
			app.Span().Text("GitHub"),
		)
}

func (s *service) handleEmailInput(ctx app.Context, e app.Event) {
	s.email = ctx.JSSrc().Get("value").String()
	if s.emailError != "" {
		s.validateEmail(s.email)
	}
}

// validateEmail checks if the email format is valid
func (s *service) validateEmail(email string) bool {
	if email == "" {
		s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailIsRequired)
		return false
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterValidEmail)
		return false
	}

	s.emailError = ""
	return true
}

func (s *service) handleEmailSubmit(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("handleEmailSubmit called")

	e.PreventDefault()
	if !s.validateEmail(s.email) {
		return
	}
	// Set loading state
	s.isLoading = true
	ctx.Async(func() {
		// Call email lookup API
		// Store email in LocalStorage
		ctx.LocalStorage().Set("email", s.email)
		response, err := s.rageApiCliient.LoginPhaseOne(s.AppContext,
			&login_models.LoginPhaseOneRequest{
				Email: s.email,
			})
		ctx.Dispatch(func(ctx app.Context) {

			log.Info().Interface("response", response).Msg("LoginPhaseOne response received")
			s.isLoading = false

			if response != nil {
				switch response.Code {
				case 404:
					s.emailError = s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyFailedToLookUpEmailF,
						map[string]interface{}{
							"email":  s.email,
							"reason": s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailNotFound),
						})
					return
				}

			}
			if err != nil {
				s.emailError = s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyFailedToLookUpEmailF,
					map[string]interface{}{
						"email":  s.email,
						"reason": err.Error(),
					})

				log.Error().Err(err).Msgf("Email lookup failed: %s", s.email)
				return
			}

			switch response.Response.Directive {
			case login_models.DIRECTIVE_LoginPhaseOne_DisplayPasswordPage:
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Password))

			case "startExternalLogin":
				// Start external login flow with the provided slug
				if response.Response.DirectiveStartExternalLogin != nil {
					slug := response.Response.DirectiveStartExternalLogin.Slug
					log.Info().Str("slug", slug).Msg("Starting external login flow")
					s.handleSocialLogin(ctx, slug)
				} else {
					s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred)
					log.Error().Msg("StartExternalLogin directive received but slug is missing")
				}

			case login_models.DIRECTIVE_Redirect:
				// Redirect to external OAuth provider

				log.Info().Str("redirectURI", response.Response.DirectiveRedirect.RedirectURI).Msg("Redirecting to external OAuth provider")
				if response.Response.DirectiveRedirect.RedirectURI != "" {
					// In a real app, this would redirect to the OAuth URL
					app.Window().Get("location").Set("href", response.Response.DirectiveRedirect.RedirectURI)
				} else {
					s.emailError = "Invalid redirect URL from server"
				}

			case login_models.DIRECTIVE_LoginPhaseOne_UserDoesNotExist:
				s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailNotFound)

				log.Info().Str("email", s.email).Msg("Email lookup error: user does not exist")

			default:
				s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred)
				log.Info().Msg("Unknown action from server")
			}
		})

	})

}
func (s *service) handleCreateAccountClick(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("handleCreateAccountClick called")
	e.PreventDefault()
	// Clear errors when navigating away
	s.emailError = ""
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_CreateAccount))
}

func (s *service) handleForgotPasswordClick(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("handleForgotPasswordClick called")
	e.PreventDefault()
	// Clear errors when navigating away
	s.emailError = ""
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_ForgotPassword))
}

func (s *service) handleGoogleLogin(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("Google login clicked")
	e.PreventDefault()
	s.handleSocialLogin(ctx, "google-social")
}

func (s *service) handleMicrosoftLogin(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("Microsoft login clicked")
	e.PreventDefault()
	s.handleSocialLogin(ctx, "microsoft-social")
}

func (s *service) handleGithubLogin(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("GitHub login clicked")
	e.PreventDefault()
	s.handleSocialLogin(ctx, "github-social")
}

func (s *service) handleSocialLogin(ctx app.Context, idpSlug string) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Str("idpSlug", idpSlug).Logger()
	log.Info().Msg("Starting external login")

	// Set loading state
	s.isLoading = true

	ctx.Async(func() {
		response, err := s.rageApiCliient.StartExternalLogin(s.AppContext,
			&external_idp.StartExternalIDPLoginRequest{
				Slug:      idpSlug,
				Directive: "login",
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				log.Error().Err(err).Msg("StartExternalLogin API call failed")
				s.showErrorMessage(ctx, s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred))
				return
			}

			if response == nil {
				log.Error().Msg("StartExternalLogin returned nil response")
				s.showErrorMessage(ctx, s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred))
				return
			}

			// Handle error responses
			if response.Code >= 400 {
				log.Error().Int("code", response.Code).Msg("StartExternalLogin returned error code")

				switch response.Code {
				case 404:
					s.showErrorMessage(ctx, s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyProviderNotFound))
				default:
					s.showErrorMessage(ctx, s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred))
				}
				return
			}

			// Handle successful response with redirect URI
			if response.Response == nil || response.Response.RedirectURI == "" {
				log.Error().Msg("Missing redirect URI in response")
				s.showErrorMessage(ctx, s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred))
				return
			}

			log.Info().Str("redirectUri", response.Response.RedirectURI).Msg("Redirecting to external OAuth provider")
			// Redirect to external OAuth provider
			app.Window().Get("location").Set("href", response.Response.RedirectURI)
		})
	})
}

func (s *service) handleDismissError(ctx app.Context, e app.Event) {
	e.PreventDefault()
	s.showError = false
	s.errorMessage = ""
	ctx.Update()
}

func (s *service) showErrorMessage(ctx app.Context, message string) {
	// Always reset first to ensure re-render triggers even if already showing
	s.showError = false
	ctx.Update()

	// Then show new message
	ctx.Async(func() {
		s.errorMessage = message
		s.showError = true
		ctx.Dispatch(func(ctx app.Context) {
			ctx.Update()
		})
	})
}
