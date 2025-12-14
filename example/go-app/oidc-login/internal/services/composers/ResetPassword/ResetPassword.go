package ResetPassword

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
	password_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		rageApiClient contracts_go_app_RageApiClient.IRageApiClient

		password             string
		confirmPassword      string
		passwordError        string
		confirmPasswordError string
		isLoading            bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IResetPasswordComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_App.IResetPasswordComposer, error) {

	return &service{
		rageApiClient: rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
	}, nil
}

func AddScopedIResetPasswordComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IResetPasswordComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "ResetPasswordComposer").Logger()
	log.Info().Msg("OnMount called for ResetPasswordComposer")
}

func (s *service) Render() app.UI {
	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetPassword)),
		app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourNewPasswordBelow)),

		app.Form().OnSubmit(s.handleResetPasswordSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().For("password").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyNewPassword)),
				app.Input().
					Type("password").
					ID("password").
					Value(s.password).
					OnInput(func(ctx app.Context, e app.Event) {
						s.password = ctx.JSSrc().Get("value").String()
						if s.passwordError != "" {
							s.validatePassword()
						}
					}).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourNewPassword)).
					Required(true),
				app.If(s.passwordError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.passwordError)
					},
				),
			),

			app.Div().Class("form-group").Body(
				app.Label().For("confirmPassword").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmPassword)),
				app.Input().
					Type("password").
					ID("confirmPassword").
					Value(s.confirmPassword).
					OnInput(func(ctx app.Context, e app.Event) {
						s.confirmPassword = ctx.JSSrc().Get("value").String()
						if s.confirmPasswordError != "" {
							s.validateConfirmPassword()
						}
					}).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmPassword)).
					Required(true),
				app.If(s.confirmPasswordError != "",
					func() app.UI {
						return app.Div().Class("error-message").Text(s.confirmPasswordError)
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

func (s *service) validatePassword() {
	s.password = strings.TrimSpace(s.password)
	if s.password == "" {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordIsRequired)
		return
	}
	if len(s.password) < 8 {
		s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordMinimum8Characters)
		return
	}
	s.passwordError = ""
}

func (s *service) validateConfirmPassword() {
	s.confirmPassword = strings.TrimSpace(s.confirmPassword)
	if s.confirmPassword == "" {
		s.confirmPasswordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordConfirmIsRequired)
		return
	}
	if s.password != s.confirmPassword {
		s.confirmPasswordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordsDoNotMatch)
		return
	}
	s.confirmPasswordError = ""
}

func (s *service) validatePasswords() bool {
	s.validatePassword()
	s.validateConfirmPassword()
	return s.passwordError == "" && s.confirmPasswordError == ""
}

func (s *service) handleResetPasswordSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	log := zerolog.Ctx(s.AppContext).With().Str("component", "ResetPasswordComposer").Logger()
	log.Info().Msg("handleResetPasswordSubmit called")

	// Clear previous errors
	s.passwordError = ""
	s.confirmPasswordError = ""

	// Validate passwords
	if !s.validatePasswords() {
		log.Info().Str("passwordError", s.passwordError).Str("confirmPasswordError", s.confirmPasswordError).Msg("Password validation failed")
		return
	}

	// Set loading state
	s.isLoading = true

	ctx.Async(func() {
		// First verify password strength
		strengthReq := &password_models.VerifyPasswordStrengthRequest{
			Password: s.password,
		}

		strengthResp, err := s.rageApiClient.VerifyPasswordStrength(s.AppContext, strengthReq)
		if err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				s.isLoading = false
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msg("Password strength verification failed")
			})
			return
		}

		if strengthResp == nil || strengthResp.Response == nil || !strengthResp.Response.Valid {
			ctx.Dispatch(func(ctx app.Context) {
				s.isLoading = false
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordTooWeak)
				log.Info().Msg("Password strength check failed")
			})
			return
		}

		// Password is strong enough, proceed with reset
		req := &login_models.PasswordResetFinishRequest{
			Password:        s.password,
			PasswordConfirm: s.confirmPassword,
		}

		resp, err := s.rageApiClient.PasswordResetFinish(s.AppContext, req)
		if err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				s.isLoading = false
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msg("Password reset failed")
			})
			return
		}

		if resp == nil || resp.Response == nil {
			ctx.Dispatch(func(ctx app.Context) {
				s.isLoading = false
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msg("Password reset failed, invalid response from server")
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			switch resp.Response.ErrorReason {
			case login_models.PasswordResetErrorReason_PasswordsDoNotMatch:
				s.confirmPasswordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordsDoNotMatch)
				log.Info().Msg("Password reset failed: passwords do not match")
				return
			case login_models.PasswordResetErrorReason_InvalidPassword:
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordTooWeak)
				log.Info().Msg("Password reset failed: invalid password")
				return
			}

			directive := resp.Response.Directive
			switch directive {
			case login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage:
				log.Info().Msg("Password reset successful, navigating to home")
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
			default:
				s.passwordError = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Str("directive", directive).Msg("Unknown directive")
			}
		})
	})
}

func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}
