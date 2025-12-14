package PasswordManager

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	passwordStage string

	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		currentStage        passwordStage
		verificationCode    string
		newPassword         string
		confirmPassword     string
		errorMessage        string
		showError           bool
		isSendingCode       bool
		isVerifying         bool
		isResetting         bool
		email               string
		isClaimedDomain     bool
		isLoading           bool
	}
)

const (
	stageInitial       passwordStage = "initial"
	stageVerifyCode    passwordStage = "verify-code"
	stageResetPassword passwordStage = "reset-password"
	stageSuccess       passwordStage = "success"
)

var stemService = (*service)(nil)

var _ contracts_App.IPasswordManagerComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,
) (contracts_App.IPasswordManagerComposer, error) {

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		managementApiClient: managementApiClient,
		currentStage:        stageInitial,
		email:               "",
		isLoading:           true,
	}, nil
}

func AddScopedIPasswordManagerComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IPasswordManagerComposer](cb, stemService.Ctor)
}

func (s *service) Render() app.UI {
	// Show unavailable message for claimed domain users
	if s.isClaimedDomain {
		return app.Div().Class("profile-container").Body(
			app.Div().Class("profile-header").Body(
				app.H1().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordManager)),
				app.P().Class("profile-subtitle").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetPasswordDescription)),
			),
			app.Div().Class("profile-card").Body(
				app.Div().Class("card-header").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordResetNotAvailable)),
				),
				app.Div().Class("card-body").Body(
					app.P().Class("card-description").Style("color", "var(--text-secondary)").Text(
						s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordResetNotAvailableDescription),
					),
				),
			),
		)
	}

	return app.Div().Class("profile-container").Body(
		app.If(s.showError, s.renderErrorBanner),

		// Page Header
		app.Div().Class("profile-header").Body(
			app.H1().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordManager)),
			app.P().Class("profile-subtitle").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetPasswordDescription)),
		),

		// Stage-based content
		app.Div().Class("profile-cards").Body(
			s.renderCurrentStage(),
		),
	)
}

func (s *service) renderCurrentStage() app.UI {
	switch s.currentStage {
	case stageInitial:
		return s.renderInitialStage()
	case stageVerifyCode:
		return s.renderVerifyCodeStage()
	case stageResetPassword:
		return s.renderResetPasswordStage()
	case stageSuccess:
		return s.renderSuccessStage()
	default:
		return s.renderInitialStage()
	}
}

func (s *service) renderInitialStage() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon password-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
						<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetYourPassword)),
					app.P().Class("card-description").Text("Click the button below to receive a verification code"),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("info-rows").Body(
				app.Div().Class("info-row").Body(
					app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
					app.Span().Class("info-value").Text(s.email),
				),
			),
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-primary").
					Disabled(s.isSendingCode).
					OnClick(s.handleSendCode).
					Text(func() string {
						if s.isSendingCode {
							return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySendingCodeDotDot)
						}
						return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySendVerificationCode)
					}()),
			),
		),
	)
}

func (s *service) renderVerifyCodeStage() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon password-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
						<polyline points="22 4 12 14.01 9 11.01"></polyline>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerifyCode)),
					app.P().Class("card-description").Text(strings.Replace(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterTheVerificationCodeSentToF), "{{.Email}}", s.email, 1)),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("form-group").Body(
				app.Label().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourVerficationCode)),
				app.Input().
					Type("text").
					Value(s.verificationCode).
					Placeholder("000000").
					MaxLength(6).
					OnInput(s.handleVerificationCodeInput),
			),
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-secondary").
					OnClick(s.handleBackToInitial).
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCancel)),
				app.Button().
					Class("btn-primary").
					Disabled(s.isVerifying || len(s.verificationCode) < 6).
					OnClick(s.handleVerifyCode).
					Text(func() string {
						if s.isVerifying {
							return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerifyingDotDot)
						}
						return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyVerifyCode)
					}()),
			),
		),
	)
}

func (s *service) renderResetPasswordStage() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon password-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetPassword)),
					app.P().Class("card-description").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourNewPasswordBelow)),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("form-group").Body(
				app.Label().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyNewPassword)),
				app.Input().
					Type("password").
					Value(s.newPassword).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEnterYourNewPassword)).
					OnInput(s.handleNewPasswordInput),
			),
			app.Div().Class("form-group").Body(
				app.Label().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmPassword)),
				app.Input().
					Type("password").
					Value(s.confirmPassword).
					Placeholder(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyConfirmYourNewPassword)).
					OnInput(s.handleConfirmPasswordInput),
			),
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-secondary").
					OnClick(s.handleBackToInitial).
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCancel)),
				app.Button().
					Class("btn-primary").
					Disabled(s.isResetting || len(s.newPassword) == 0 || len(s.confirmPassword) == 0).
					OnClick(s.handleResetPassword).
					Text(func() string {
						if s.isResetting {
							return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResettingPasswordDotDot)
						}
						return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyResetPassword)
					}()),
			),
		),
	)
}

func (s *service) renderSuccessStage() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon success-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
						<polyline points="22 4 12 14.01 9 11.01"></polyline>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySuccess)),
					app.P().Class("card-description").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordResetSuccessfully)),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordResetSuccessfully)),
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-primary").
					OnClick(s.handleBackToInitial).
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyDone)),
			),
		),
	)
}

func (s *service) renderErrorBanner() app.UI {
	return app.Div().Class("error-banner").Body(
		app.Div().Class("error-content").Body(
			app.Div().Class("error-icon").Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
				</svg>`),
			),
			app.Div().Class("error-text").Body(
				app.Span().Class("error-title").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)),
				app.Span().Class("error-message").Text(s.errorMessage),
			),
		),
	)
}

// Event handlers
func (s *service) handleSendCode(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("Sending verification code")

	s.isSendingCode = true
	s.showError = false
	ctx.Update()

	ctx.Async(func() {
		response, err := s.managementApiClient.PasswordResetStart(s.AppContext,
			&models_api_login_models.PasswordResetStartRequest{
				Email: s.email,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isSendingCode = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to start password reset")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySomethingWentWrong)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Msg("Verification code sent successfully")
				s.currentStage = stageVerifyCode
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySomethingWentWrong)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) handleVerificationCodeInput(ctx app.Context, e app.Event) {
	s.verificationCode = ctx.JSSrc().Get("value").String()
	s.showError = false
	ctx.Update()
}

func (s *service) handleVerifyCode(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Str("code", s.verificationCode).Msg("Verifying code")

	s.isVerifying = true
	s.showError = false
	ctx.Update()

	ctx.Async(func() {
		response, err := s.managementApiClient.VerifyCode(s.AppContext,
			&models_api_login_models.VerifyCodeRequest{
				Code: s.verificationCode,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isVerifying = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to verify code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidCode)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 && response.Response != nil {
				// Check directive to see if verification succeeded
				if response.Response.Directive == models_api_login_models.DIRECTIVE_PasswordReset_DisplayPasswordResetPage {
					log.Info().Msg("Code verified successfully")
					s.currentStage = stageResetPassword
				} else {
					log.Warn().Str("directive", response.Response.Directive).Msg("Unexpected directive")
					s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidCode)
					s.showError = true
				}
			} else {
				log.Error().Msg("Invalid verification code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidCode)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) handleNewPasswordInput(ctx app.Context, e app.Event) {
	s.newPassword = ctx.JSSrc().Get("value").String()
	s.showError = false
	ctx.Update()
}

func (s *service) handleConfirmPasswordInput(ctx app.Context, e app.Event) {
	s.confirmPassword = ctx.JSSrc().Get("value").String()
	s.showError = false
	ctx.Update()
}

func (s *service) handleResetPassword(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("Resetting password")

	// Validate passwords match
	if s.newPassword != s.confirmPassword {
		s.showError = true
		s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordsDoNotMatch)
		ctx.Update()
		return
	}

	// Validate password length
	if len(s.newPassword) < 8 {
		s.showError = true
		s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordMustBeAtLeast8Characters)
		ctx.Update()
		return
	}

	s.isResetting = true
	s.showError = false
	ctx.Update()

	ctx.Async(func() {
		response, err := s.managementApiClient.PasswordResetFinish(s.AppContext,
			&models_api_login_models.PasswordResetFinishRequest{
				Password:        s.newPassword,
				PasswordConfirm: s.confirmPassword,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isResetting = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to reset password")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySomethingWentWrong)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 && response.Response != nil {
				// Check for errors
				if response.Response.ErrorReason == models_api_login_models.PasswordResetErrorReason_InvalidPassword {
					log.Error().Msg("Invalid password")
					s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordMustBeAtLeast8Characters)
					s.showError = true
				} else if response.Response.ErrorReason == models_api_login_models.PasswordResetErrorReason_PasswordsDoNotMatch {
					log.Error().Msg("Passwords do not match")
					s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasswordsDoNotMatch)
					s.showError = true
				} else if response.Response.ErrorReason == models_api_login_models.PasswordResetErrorReason_NoError {
					log.Info().Msg("Password reset successfully")
					s.currentStage = stageSuccess
				} else {
					log.Error().Int("errorReason", int(response.Response.ErrorReason)).Msg("Unexpected error reason")
					s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySomethingWentWrong)
					s.showError = true
				}
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySomethingWentWrong)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) handleBackToInitial(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("Returning to initial stage")

	s.currentStage = stageInitial
	s.verificationCode = ""
	s.newPassword = ""
	s.confirmPassword = ""
	s.showError = false
	s.errorMessage = ""
	ctx.Update()
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("PasswordManager page mounted")

	// Fetch user profile to check if claimed domain and get email
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserProfile(s.AppContext)
		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false
			if err == nil && response != nil && response.Code == 200 && response.Response != nil {
				s.isClaimedDomain = response.Response.IsClaimedDomain
				s.email = response.Response.Email
				log.Info().Bool("isClaimedDomain", s.isClaimedDomain).Str("email", s.email).Msg("Profile loaded")
			} else {
				log.Error().Err(err).Msg("Failed to load profile")
			}
			ctx.Update()
		})
	})
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("PasswordManager page navigated")
}

func (s *service) OnDismount() {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasswordManager").Logger()
	log.Info().Msg("PasswordManager page dismounted")
}
