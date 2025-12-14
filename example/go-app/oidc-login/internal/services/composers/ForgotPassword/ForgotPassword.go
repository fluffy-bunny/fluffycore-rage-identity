package ForgotPassword

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

		email      string
		emailError string
		isLoading  bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IForgotPasswordComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_App.IForgotPasswordComposer, error) {

	return &service{
		rageApiClient: rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
	}, nil
}

func AddScopedIForgotPasswordComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IForgotPasswordComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "ForgotPasswordComposer").Logger()
	log.Info().Msg("OnMount called for ForgotPasswordComposer")
}

func (s *service) Render() app.UI {
	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyForgotPassword)),
		app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterEmailToContinue)),

		app.Form().OnSubmit(s.handleForgotPasswordSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("email").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
				app.Input().
					Type("email").
					ID("email").
					Value(s.email).
					OnInput(func(ctx app.Context, e app.Event) {
						s.email = ctx.JSSrc().Get("value").String()
						if s.emailError != "" {
							s.validateEmail()
						}
					}).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourEmail)).
					Required(true),
				app.If(s.emailError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.emailError)
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

func (s *service) validateEmail() bool {
	s.email = strings.TrimSpace(s.email)
	if s.email == "" {
		s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailIsRequired)
		return false
	}
	s.emailError = ""
	return true
}

func (s *service) handleForgotPasswordSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	log := zerolog.Ctx(s.AppContext).With().Str("component", "ForgotPasswordComposer").Logger()
	log.Info().Msg("handleForgotPasswordSubmit called")

	// Clear previous errors
	s.emailError = ""

	// Validate email
	if !s.validateEmail() {
		log.Info().Str("email", s.email).Str("emailError", s.emailError).Msg("Email validation failed")
		return
	}

	// Set loading state
	s.isLoading = true

	ctx.Async(func() {
		// Call PasswordResetStart API
		response, err := s.rageApiClient.PasswordResetStart(s.AppContext,
			&login_models.PasswordResetStartRequest{
				Email: s.email,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				// Network or parsing error
				s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msgf("Password reset failed for email: %s", s.email)
				return
			}

			if response.Response == nil {
				s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msgf("Password reset failed for email: %s, invalid response from server", s.email)
				return
			}

			// Check HTTP status code for errors
			if response.Code >= 400 {
				// Server returned an error
				switch response.Code {
				case 404:
					// Email not found - but we don't want to reveal this for security
					// Still navigate to verify code page
					break
				default:
					s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
					log.Error().Msgf("Password reset failed for email: %s, server error code: %d", s.email, response.Code)
					return
				}
			}

			// Store email in LocalStorage and navigate to verify code page
			// Even if email doesn't exist, we navigate to avoid revealing account existence
			ctx.LocalStorage().Set("email", s.email)
			ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_VerifyCode))
		})
	})
}

func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}
