package Password

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		rageApiClient contracts_go_app_RageApiClient.IRageApiClient

		email         string
		password      string
		passwordError string
		isLoading     bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IPasswordComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_App.IPasswordComposer, error) {

	return &service{
		rageApiClient: rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
	}, nil
}

func AddScopedIPasswordComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IPasswordComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordComposer").Logger()
	log.Info().Msg("OnMount called for PasswordComposer")
	// Retrieve email from LocalStorage
	ctx.LocalStorage().Get("email", &s.email)

	// If no email, redirect back to home
	if s.email == "" {
		ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
		return
	}
}

func (s *service) Render() app.UI {
	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourPassword)),
		app.P().Text(s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyEnterYourPasswordF, map[string]interface{}{"email": s.email})),

		app.Form().OnSubmit(s.handlePasswordSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("password").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPassword)),
				app.Input().
					Type("password").
					ID("password").
					Value(s.password).
					OnInput(func(ctx app.Context, e app.Event) {
						s.password = ctx.JSSrc().Get("value").String()
					}).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourPassword)).
					Required(true),
				app.If(s.passwordError != "", func() app.UI {
					return app.Div().Class("error-message").Text(s.passwordError)
				}),
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
					Body(
						app.If(s.isLoading, func() app.UI {
							return app.Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySigningInDotDot))
						}).Else(func() app.UI {
							return app.Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignIn))
						}),
					),
			),
		),
	)
}

func (s *service) validatePassword() bool {
	s.password = strings.TrimSpace(s.password)
	if s.password == "" {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordRequired)
		return false
	}
	s.passwordError = ""
	return true
}

func (s *service) handlePasswordSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordComposer").Logger()
	log.Info().Msg("handlePasswordSubmit called")

	// Clear previous errors
	s.passwordError = ""

	// Validate password
	if !s.validatePassword() {
		log.Info().Str("email", s.email).Str("passwordError", s.passwordError).Msg("Password validation failed")
		return
	}

	// Set loading state
	s.isLoading = true

	ctx.Async(func() {
		// Call LoginPassword API
		response, err := s.rageApiClient.LoginPassword(s.AppContext,
			&login_models.LoginPasswordRequest{
				Email:    s.email,
				Password: s.password,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				// Network or parsing error
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msgf("Password login failed for email: %s", s.email)
				return
			}

			if response.Response == nil {
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msgf("Password login failed for email: %s, invalid response from server", s.email)
				return
			}

			// Check HTTP status code for errors
			if response.Code >= 400 {
				// Server returned an error
				switch response.Code {
				case 401:
					// Invalid password
					s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidPassword)
				case 404:
					// User not found (shouldn't happen at this stage, but handle gracefully)
					s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidPassword)
				default:
					s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				}
				log.Error().Msgf("Password login failed for email: %s, server error code: %d", s.email, response.Code)
				return
			}

			// Success - check directive for next action
			if response.Response.Directive == login_models.DIRECTIVE_Redirect {
				// Final redirect - authentication complete
				if response.Response.DirectiveRedirect != nil {
					redirectURI := response.Response.DirectiveRedirect.RedirectURI
					log.Info().Msgf("Login successful, redirecting to: %s", redirectURI)
					// Use window.location.href for external redirects (full page navigation)
					app.Window().Get("location").Set("href", redirectURI)
				}
			} else if response.Response.Directive == login_models.DIRECTIVE_VerifyCode_DisplayVerifyCodePage {
				// Need to verify email code (MFA)
				log.Info().Msg("Email verification required, navigating to verify-code")
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_VerifyCode))
			} else {
				// Unexpected directive
				log.Warn().Msgf("Unexpected directive: %s", response.Response.Directive)
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
			}
		})
	})
}

func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}
