package Home

import (
	"regexp"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	"github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/common"
	oidc_login_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/contracts/config"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/config"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/utils"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	external_idp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	models_errors "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/errors"
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
		appConfig         *oidc_login_contracts_config.AppConfig
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
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("HomeComposer OnMount called")

	oidcFlowConfig, err := s.appConfigAccessor.GetOIDCFlowAppConfig(s.AppContext)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get OIDC Flow App Config")
		s.showErrorMessage(ctx, "Configuration error. Please try again later.")
		return
	}
	s.oidcFlowConfig = oidcFlowConfig

	s.checkAndDisplayErrorCookie(ctx)
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("HomeComposer OnNav called")
	s.checkAndDisplayErrorCookie(ctx)
}
func stringKeyMapToInterface[T any](input map[string]T) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		result[k] = v
	}
	return result
}
func (s *service) checkAndDisplayErrorCookie(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()

	// Check for error cookie (e.g., from OAuth2 callback redirects)

	errorCookie, err := utils.GetCookie[contracts_cookies.ErrorCookie]("_error")

	log.Info().
		Err(err).
		Interface("errorCookie", errorCookie).
		Interface("params", errorCookie.Params).
		Int("paramsLen", len(errorCookie.Params)).
		Msg("Checking for error cookie")
	if err == nil && errorCookie.Error != "" {
		msg := ""
		errorCookieParams := stringKeyMapToInterface(errorCookie.Params)
		log.Info().
			Str("error", errorCookie.Error).
			Str("code", errorCookie.Code).
			Interface("params", errorCookie.Params).
			Interface("convertedParams", errorCookieParams).
			Msg("Found error cookie from backend, displaying to user")
		switch errorCookie.Code {
		case string(models_errors.FlowError_UserNotFound):
			msg = s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyUserNotFoundF, errorCookieParams)
		default:
			msg = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnexpectedErrorOccurred)
		}
		s.showErrorMessage(ctx, msg)

		// Delete the error cookie after reading it
		// Use the cookie domain from the OIDC config if available
		cookieDomain := ""
		if s.appConfig != nil && s.appConfig.CookieDomain != "" {
			cookieDomain = s.appConfig.CookieDomain
			log.Info().Str("cookieDomain", cookieDomain).Msg("Using cookie domain from config")
		}

		utils.DeleteCookie("_error", utils.CookieOptions{
			Path:   "/",
			Domain: cookieDomain,
		})
	}
}

func (s *service) renderPasskeyLogin() app.UI {
	return app.Div().Class("passkey-login-section").Body(
		app.Div().Class("divider").Body(
			app.Span().Text("or"),
		),
		app.Button().
			Class("passkey-btn").
			OnClick(s.handlePasskeyLogin).
			Body(
				app.Raw(common.PasskeyIconSmallSVG),
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySigninWithPasskey)),
			),
	)
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
		return app.Div().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLoadingDotDot))
	}

	return app.Div().Class("step-container").Body(
		app.If(s.showError, s.renderErrorBanner),

		//app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyWelcomeBack)),
		//app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourEmailToContinue)),

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
					Class("passkey-btn").
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyContinue)),
			),
		),
		app.If(
			s.oidcFlowConfig.EnabledWebAuthN,
			s.renderPasskeyLogin,
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
			app.Raw(common.GoogleIconSmallSVG),
			app.Span().Text("Google"),
		)
}

func (s *service) renderMicrosoftButton() app.UI {
	return app.Button().
		Class("social-btn microsoft-btn").
		OnClick(s.handleMicrosoftLogin).
		Body(
			app.Raw(common.MicrosoftIconSmallSVG),
			app.Span().Text("Microsoft"),
		)
}

func (s *service) renderGithubButton() app.UI {
	return app.Button().
		Class("social-btn github-btn").
		OnClick(s.handleGithubLogin).
		Body(
			app.Raw(common.GitHubIconSmallSVG),
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

func (s *service) handlePasskeyLogin(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "HomeComposer").Logger()
	log.Info().Msg("Passkey login clicked")
	e.PreventDefault()

	// Clear any existing email field to ensure clean state
	s.email = ""
	s.emailError = ""

	// Always use discoverable credentials flow for passkey login
	// This ensures the user can select any passkey for any account
	log.Info().Msg("Using discoverable credentials flow for passkey login")

	// Create error callback that will be called from JavaScript
	errorCallback := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		if len(args) > 0 {
			errorMsg := args[0].String()
			log.Error().Str("error", errorMsg).Msg("Passkey login failed")
			ctx.Dispatch(func(ctx app.Context) {
				s.showErrorMessage(ctx, errorMsg)
			})
		}
		return nil
	})

	// Call the WebAuthn JavaScript function directly
	// This will trigger the browser's passkey selection UI
	result := app.Window().Call("LoginUser", "", false, errorCallback)

	if !result.Truthy() {
		log.Error().Msg("Failed to initiate passkey login")
		s.showErrorMessage(ctx, "Failed to start passkey authentication")
		return
	}

	log.Info().Msg("Passkey authentication initiated")
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
