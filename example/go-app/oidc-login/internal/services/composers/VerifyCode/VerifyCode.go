package VerifyCode

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		rageApiClient contracts_go_app_RageApiClient.IRageApiClient

		email     string
		code      string
		codeError string
		isLoading bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IVerifyCodeComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_App.IVerifyCodeComposer, error) {

	return &service{
		rageApiClient: rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
	}, nil
}

func AddScopedIVerifyCodeComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IVerifyCodeComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "VerifyCodeComposer").Logger()
	log.Info().Msg("OnMount called for VerifyCodeComposer")

	// Call VerifyCodeBegin to validate session and get verification details
	// This handles cases where redirect from external IDP requires verification
	s.callVerifyCodeBegin(ctx)
}

// callVerifyCodeBegin calls the VerifyCodeBegin API
func (s *service) callVerifyCodeBegin(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "VerifyCodeComposer").Logger()
	log.Info().Msg("callVerifyCodeBegin")
	ctx.Async(func() {
		response, err := s.rageApiClient.VerifyCodeBegin(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Error().Err(err).Msg("VerifyCodeBegin error - auth session invalid, redirecting to home")
				// Auth session is bad, redirect to home (error page)
				// Only way out is for initiating website to start new OIDC authorization session
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
				return
			}

			if response.Code >= 400 {
				log.Error().Msgf("VerifyCodeBegin failed with code %d - auth session invalid, redirecting to home", response.Code)
				// Auth session is bad, redirect to home (error page)
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
				return
			}

			log.Info().Msg("VerifyCodeBegin successful")

			// Store email from response (required)
			if response.Response != nil && response.Response.Email != "" {
				s.email = response.Response.Email
				log.Info().Msgf("Retrieved email from VerifyCodeBegin: %s", s.email)
			}

			// If response contains a code (dev mode), prefill it
			if response.Response != nil && response.Response.Code != "" {
				s.code = response.Response.Code
				log.Info().Msgf("Dev mode: pre-filled verification code from API: %s", s.code)
			}

			ctx.Update()
		})
	})
}

func (s *service) Render() app.UI {
	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerifyCode)),
		app.P().Text(s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyEnterTheVerificationCodeSentToF, map[string]interface{}{"email": s.email})),

		app.Form().OnSubmit(s.handleVerifyCodeSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("code").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerificationCode)),
				app.Input().
					Type("text").
					ID("code").
					Value(s.code).
					OnInput(func(ctx app.Context, e app.Event) {
						s.code = ctx.JSSrc().Get("value").String()
						if s.codeError != "" {
							s.validateCode()
						}
					}).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourVerficationCode)).
					Required(true),
				app.If(s.codeError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.codeError)
					},
				),
			),

			app.Div().Class("button-group").Body(
				app.Button().
					Type("button").
					Class("btn-secondary").
					OnClick(s.handleBackClick).
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyBack)),
				app.Button().
					Type("submit").
					Class("btn-primary").
					Disabled(s.isLoading).
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySubmit)),
			),
		),
	)
}

func (s *service) validateCode() bool {
	s.code = strings.TrimSpace(s.code)
	if s.code == "" {
		s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerificationCodeRequired)
		return false
	}
	s.codeError = ""
	return true
}

func (s *service) handleVerifyCodeSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	log := zerolog.Ctx(s.AppContext).With().Str("component", "VerifyCodeComposer").Logger()
	log.Info().Msg("handleVerifyCodeSubmit called")

	// Clear previous errors
	s.codeError = ""

	// Validate code
	if !s.validateCode() {
		log.Info().Str("code", s.code).Str("codeError", s.codeError).Msg("Code validation failed")
		return
	}

	// Set loading state
	s.isLoading = true

	ctx.Async(func() {
		// Call VerifyCode API
		response, err := s.rageApiClient.VerifyCode(s.AppContext,
			&login_models.VerifyCodeRequest{
				Code: s.code,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				// Network or parsing error
				s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msgf("Code verification failed for email: %s", s.email)
				return
			}

			if response.Response == nil {
				s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msgf("Code verification failed for email: %s, invalid response from server", s.email)
				return
			}

			// Check HTTP status code for errors
			if response.Code >= 400 {
				// Server returned an error
				switch response.Code {
				case 400:
					// Invalid code
					s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidVerificationCode)
				case 404:
					// Email or code not found
					s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidVerificationCode)
				default:
					s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				}
				log.Error().Msgf("Code verification failed for email: %s, server error code: %d", s.email, response.Code)
				return
			}

			// Success - check directive for next action
			switch response.Response.Directive {
			case login_models.DIRECTIVE_Redirect:
				// Final redirect - authentication complete
				if response.Response.DirectiveRedirect != nil {
					redirectURI := response.Response.DirectiveRedirect.RedirectURI
					log.Info().Msgf("Code verification successful, redirecting to: %s", redirectURI)
					// Use window.location.href for external redirects (full page navigation)
					app.Window().Get("location").Set("href", redirectURI)
				}
			case login_models.DIRECTIVE_PasswordReset_DisplayPasswordResetPage:
				// Need to set/reset password
				log.Info().Msg("Password reset required, navigating to password reset page")
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_ResetPassword))
			case login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage:
				// Go back to home/login page
				log.Info().Msg("Navigating back to home page")
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
			default:
				// Unexpected directive
				log.Warn().Msgf("Unexpected directive: %s", response.Response.Directive)
				s.codeError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)

			}

		})
	})
}

func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_ForgotPassword))
}
