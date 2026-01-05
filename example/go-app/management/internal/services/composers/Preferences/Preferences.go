package Preferences

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	models_api_preferences "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_preferences"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		keepSignedIn        bool
		doNotShowAgain      bool
		errorMessage        string
		showError           string
		successMessage      string
		showSuccess         bool
		isLoading           bool
		isSavingPreference  bool
		isClearingSSO       bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IPreferencesComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,
) (contracts_App.IPreferencesComposer, error) {

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		managementApiClient: managementApiClient,
		isLoading:           true,
	}, nil
}

func AddScopedIPreferencesComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IPreferencesComposer](cb, stemService.Ctor)
}

func (s *service) Render() app.UI {
	return app.Div().Class("profile-container").Body(
		app.If(s.showError != "", s.renderErrorBanner),
		app.If(s.showSuccess, s.renderSuccessBanner),

		// Page Header
		app.Div().Class("profile-header").Body(
			app.H1().Text("Preferences"),
			app.P().Class("profile-subtitle").Text("Manage your sign-in preferences and security settings"),
		),

		// Show loading or preference cards
		app.If(s.isLoading, func() app.UI {
			return app.Div().Class("loading-container").Body(
				app.Div().Class("loading-spinner"),
				app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLoadingProfileDotDot)),
			)
		}).Else(func() app.UI {
			// Preference Cards Container
			return app.Div().Class("profile-cards").Body(
				s.renderKeepSignedInCard(),
				s.renderSSOCookieCard(),
			)
		}),
	)
}

func (s *service) renderKeepSignedInCard() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon personal-info-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text("Keep Me Signed In Preferences"),
					app.P().Class("card-description").Text("Control whether you see the 'Keep me signed in' page after authentication"),
				),
			),
		),

		app.Div().Class("card-body").Body(
			app.Div().Class("info-rows").Body(
				// First toggle: Keep me signed in
				app.Div().Class("preference-row").Body(
					app.Div().Class("preference-info").Body(
						app.Span().Class("info-label").Text("Keep me signed in"),
						app.Span().Class("info-description").Text("Stay signed in across sessions"),
					),
					app.Div().Class("preference-control").Body(
						app.Label().Class("toggle-switch").Body(
							app.Input().
								Type("checkbox").
								Checked(s.keepSignedIn).
								OnChange(s.handleKeepSignedInToggle).
								Disabled(s.isSavingPreference),
							app.Span().Class("toggle-slider"),
						),
					),
				),
				// Second toggle: Do not show dialog
				app.Div().Class("preference-row").Body(
					app.Div().Class("preference-info").Body(
						app.Span().Class("info-label").Text("Do not show 'Keep me signed in' dialog"),
						app.Span().Class("info-description").Text("Skip the prompt and use your preference automatically"),
					),
					app.Div().Class("preference-control").Body(
						app.Label().Class("toggle-switch").Body(
							app.Input().
								Type("checkbox").
								Checked(s.doNotShowAgain).
								OnChange(s.handleDoNotShowAgainToggle).
								Disabled(s.isSavingPreference),
							app.Span().Class("toggle-slider"),
						),
					),
				),
			),
		),
	)
}

func (s *service) renderSSOCookieCard() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon personal-info-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
						<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text("Session Management"),
					app.P().Class("card-description").Text("Manage your active sessions and sign-in state"),
				),
			),
		),

		app.Div().Class("card-body").Body(
			app.Div().Class("info-rows").Body(
				app.Div().Class("preference-row").Body(
					app.Div().Class("preference-info").Body(
						app.Span().Class("info-label").Text("Clear SSO Cookie"),
						app.Span().Class("info-description").Text("Remove your single sign-on session. You'll need to sign in again on your next visit."),
					),
					app.Div().Class("preference-control").Body(
						app.Button().
							Class("btn-secondary").
							OnClick(s.handleClearSSO).
							Disabled(s.isClearingSSO).
							Body(
								app.If(s.isClearingSSO, func() app.UI {
									return app.Text("Clearing...")
								}).Else(func() app.UI {
									return app.Text("Clear SSO")
								}),
							),
					),
				),
			),
		),
	)
}

func (s *service) renderErrorBanner() app.UI {
	return app.Div().Class("error-banner").Body(
		app.Div().Class("error-content").Body(
			app.Span().Class("error-icon").Text("⚠️"),
			app.Span().Class("error-message").Text(s.showError),
		),
		app.Button().
			Class("error-close").
			OnClick(s.handleCloseError).
			Text("×"),
	)
}

func (s *service) renderSuccessBanner() app.UI {
	return app.Div().Class("success-notification").Body(
		app.Div().Class("success-notification-icon").Body(
			app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
				<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
			</svg>`),
		),
		app.Div().Class("success-notification-text").Text(s.successMessage),
		app.Button().
			Class("success-notification-close").
			OnClick(s.handleCloseSuccess).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
				</svg>`),
			),
	)
}

func (s *service) handleKeepSignedInToggle(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PreferencesComposer").Logger()

	// Get the new value from the checkbox
	checked := e.JSValue().Get("target").Get("checked").Bool()

	log.Info().
		Bool("keepSignedIn", checked).
		Msg("Toggling keep signed in preference")

	s.isSavingPreference = true
	s.showError = ""
	s.showSuccess = false
	ctx.Update()

	ctx.Async(func() {
		// Call API to update preference, preserve doNotShowAgain value
		rageClient := s.managementApiClient.GetRageApiClient()
		response, err := rageClient.UpdateKeepSignedInPreference(s.AppContext,
			&models_api_preferences.UpdateKeepSignedInPreferenceRequest{
				DoNotShowAgain: s.doNotShowAgain,
				KeepSignedIn:   checked,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isSavingPreference = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to update keep signed in preference")
				s.showError = "Failed to save preference. Please try again."
				// Revert the toggle
				s.keepSignedIn = !checked
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Msg("Keep signed in preference updated successfully")
				s.keepSignedIn = checked
				s.successMessage = "Preference saved successfully"
				s.showSuccess = true
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.showError = "Failed to save preference. Please try again."
				// Revert the toggle
				s.keepSignedIn = !checked
			}

			ctx.Update()
		})
	})
}

func (s *service) handleDoNotShowAgainToggle(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PreferencesComposer").Logger()

	// Get the new value from the checkbox
	checked := e.JSValue().Get("target").Get("checked").Bool()

	log.Info().
		Bool("doNotShowAgain", checked).
		Msg("Toggling do not show again preference")

	s.isSavingPreference = true
	s.showError = ""
	s.showSuccess = false
	ctx.Update()

	ctx.Async(func() {
		// Call API to update preference, preserve keepSignedIn value
		rageClient := s.managementApiClient.GetRageApiClient()
		response, err := rageClient.UpdateKeepSignedInPreference(s.AppContext,
			&models_api_preferences.UpdateKeepSignedInPreferenceRequest{
				DoNotShowAgain: checked,
				KeepSignedIn:   s.keepSignedIn,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isSavingPreference = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to update do not show again preference")
				s.showError = "Failed to save preference. Please try again."
				// Revert the toggle
				s.doNotShowAgain = !checked
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Msg("Do not show again preference updated successfully")
				s.doNotShowAgain = checked
				s.successMessage = "Preference saved successfully"
				s.showSuccess = true
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.showError = "Failed to save preference. Please try again."
				// Revert the toggle
				s.doNotShowAgain = !checked
			}

			ctx.Update()
		})
	})
}

func (s *service) handleClearSSO(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PreferencesComposer").Logger()

	log.Info().Msg("Clearing SSO cookie")

	s.isClearingSSO = true
	s.showError = ""
	s.showSuccess = false
	ctx.Update()

	ctx.Async(func() {
		rageClient := s.managementApiClient.GetRageApiClient()
		response, err := rageClient.ClearSSOCookie(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			s.isClearingSSO = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to clear SSO cookie")
				s.showError = "Failed to clear SSO cookie. Please try again."
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Msg("SSO cookie cleared successfully")
				s.successMessage = "SSO session cleared successfully"
				s.showSuccess = true
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.showError = "Failed to clear SSO cookie. Please try again."
			}

			ctx.Update()
		})
	})
}

func (s *service) handleCloseError(ctx app.Context, e app.Event) {
	s.showError = ""
	ctx.Update()
}

func (s *service) handleCloseSuccess(ctx app.Context, e app.Event) {
	s.showSuccess = false
	ctx.Update()
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PreferencesComposer").Logger()
	log.Info().Msg("Preferences page mounted")

	// Load preference data
	ctx.Async(func() {
		rageClient := s.managementApiClient.GetRageApiClient()
		response, err := rageClient.GetKeepSignedInPreference(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err == nil && response != nil && response.Code == 200 {
				log.Info().Msg("Preferences loaded successfully")
				s.keepSignedIn = response.Response.KeepSignedIn
				s.doNotShowAgain = response.Response.DoNotShowAgain
			} else {
				// Show error instead of redirecting - App's OnMount already handled auth
				log.Error().
					Err(err).
					Interface("response", response).
					Msg("Failed to load preferences")
				s.showError = "Failed to load preferences. Please try refreshing the page."
			}

			ctx.Update()
		})
	})
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Preferences page navigated")

	// Reset state and reload data on navigation (e.g., manual refresh)
	s.isLoading = true
	s.showError = ""
	s.showSuccess = false
	ctx.Update()

	// Load preference data
	ctx.Async(func() {
		rageClient := s.managementApiClient.GetRageApiClient()
		response, err := rageClient.GetKeepSignedInPreference(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err == nil && response != nil && response.Code == 200 {
				log.Info().Msg("Preferences loaded successfully")
				s.keepSignedIn = response.Response.KeepSignedIn
				s.doNotShowAgain = response.Response.DoNotShowAgain
			} else {
				// Show error instead of redirecting - App's OnMount already handled auth
				log.Error().
					Err(err).
					Interface("response", response).
					Msg("Failed to load preferences")
				s.showError = "Failed to load preferences. Please try refreshing the page."
			}

			ctx.Update()
		})
	})
}

func (s *service) OnDismount() {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Preferences page dismounted")
}
