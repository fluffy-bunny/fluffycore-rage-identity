package CreateAccount

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/config"
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

		rageApiCliient    contracts_go_app_RageApiClient.IRageApiClient
		appConfigAccessor contracts_config.IAppConfigAccessor

		email             string
		confirmEmail      string
		password          string
		emailError        string
		confirmEmailError string
		passwordError     string
		isLoading         bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.ICreateAccountComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiCliient contracts_go_app_RageApiClient.IRageApiClient,
	appConfigAccessor contracts_config.IAppConfigAccessor,
) (contracts_App.ICreateAccountComposer, error) {

	return &service{
		rageApiCliient:    rageApiCliient,
		appConfigAccessor: appConfigAccessor,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
	}, nil
}

func AddScopedICreateAccountComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.ICreateAccountComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "CreateAccountComposer").Logger()
	log.Info().Msg("OnMount called for CreateAccountComposer")
}

func (s *service) Render() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "CreateAccountComposer").Logger()
	log.Info().Msg("Render called for CreateAccountComposer")

	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCreateAccount)),
		app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFillInYourDetailsToCreateANewAccount)),

		app.Form().OnSubmit(s.handleCreateAccountSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("email").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
				app.Input().
					Type("email").
					ID("email").
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourEmail)).
					Value(s.email).
					OnInput(func(ctx app.Context, e app.Event) {
						s.email = ctx.JSSrc().Get("value").String()
						if s.emailError != "" {
							s.validateEmail()
						}
					}).
					Required(true),
				app.If(s.emailError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.emailError)
					},
				),
			),

			app.Div().Class("form-group").Body(
				app.Label().For("confirm-email").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmEmail)),
				app.Input().
					Type("email").
					ID("confirm-email").
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmEmail)).
					Value(s.confirmEmail).
					OnInput(func(ctx app.Context, e app.Event) {
						s.confirmEmail = ctx.JSSrc().Get("value").String()
						if s.confirmEmailError != "" {
							s.validateEmail()
						}
					}).
					Required(true),
				app.If(s.confirmEmailError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.confirmEmailError)
					},
				),
			), app.Div().Class("form-group").Body(
				app.Label().For("password").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPassword)),
				app.Input().
					Type("password").
					ID("password").
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourPassword)).
					Value(s.password).
					OnInput(func(ctx app.Context, e app.Event) {
						s.password = ctx.JSSrc().Get("value").String()
						if s.passwordError != "" {
							s.validatePassword()
						}
					}).
					Required(true),
				app.If(s.passwordError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.passwordError)
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
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCreateAccount)),
			),
		),

		app.Div().Class("create-account").Body(
			app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyAlreadyHaveAnAccount)+" "),
			app.A().
				Href(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home)).
				OnClick(s.handleBackToLoginClick).
				Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignIn)),
		),
	)
}

func (s *service) validateEmail() bool {
	s.email = strings.TrimSpace(s.email)
	s.confirmEmail = strings.TrimSpace(s.confirmEmail)
	if s.email == "" {
		s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailIsRequired)
		return false
	}
	if s.confirmEmail == "" {
		s.confirmEmailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmEmailRequired)
		return false
	}
	if s.email != s.confirmEmail {
		s.confirmEmailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailsDoNotMatch)
		return false
	}
	s.emailError = ""
	// ensure that all fields aren't prepended or appended with spaces

	return true
}

func (s *service) validatePassword() bool {
	s.password = strings.TrimSpace(s.password)
	if s.password == "" {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordRequired)
		return false
	}
	if len(s.password) < 8 {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidPassword)
		return false
	}
	s.passwordError = ""
	return true
}

func (s *service) handleCreateAccountSubmit(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "CreateAccountComposer").Logger()
	log.Info().Msg("handleCreateAccountSubmit called")
	e.PreventDefault()

	// Clear previous errors
	s.emailError = ""
	s.confirmEmailError = ""
	s.passwordError = ""

	// Validate all fields
	emailValid := s.validateEmail()
	if !emailValid {
		log.Info().
			Str("email", s.email).Str("confirmEmail", s.confirmEmail).Str("password", s.password).
			Str("emailError", s.emailError).Str("confirmEmailError", s.confirmEmailError).Msg("Email validation failed")
		return
	}

	passwordValid := s.validatePassword()

	if !passwordValid {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordRequired)
		return
	}

	// Set loading state
	s.isLoading = true
	ctx.Async(func() {
		// Call Signup API - server will validate password strength and check if user exists
		response, err := s.rageApiCliient.Signup(s.AppContext, &login_models.SignupRequest{
			Email:    s.email,
			Password: s.password,
		})
		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false
			if err != nil {
				// Network or parsing error
				s.emailError = s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyFailedToCreateAcctountF,
					map[string]interface{}{
						"email": s.email,
					})
				log.Error().Err(err).Msgf("Account creation failed for email: %s", s.email)
				return
			}
			if response.Response == nil {
				s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msgf("Account creation failed for email: %s, invalid response from server", s.email)
				return
			}
			// Check HTTP status code for errors
			if response.Code >= 400 {
				// Server returned an error (user exists, weak password, etc.)
				// The error details should be in the response
				errorMsg := s.Localizer.GetLocalizedStringF(contracts_LocalizerBundle.LocaleKeyFailedToCreateAcctountF,
					map[string]interface{}{
						"email": s.email,
					})

				// Display error on appropriate field based on error code
				switch response.Code {
				case 409:
					// Conflict - user already exists
					s.emailError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAlreadyInUse)
				case 400:
					// Bad request - likely weak password
					s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidPassword)
				default:
					s.emailError = errorMsg
				}
				log.Error().Msgf("Account creation failed for email: %s, server error code: %d", s.email, response.Code)

				return
			}
			ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_VerifyCode))

		})
	})
}
func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}

func (s *service) handleBackToLoginClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}
