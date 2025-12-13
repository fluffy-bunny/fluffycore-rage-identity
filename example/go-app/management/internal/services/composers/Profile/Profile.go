package Profile

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		subject             string
		firstName           string
		lastName            string
		phoneNumber         string
		email               string
		isClaimedDomain     bool
		errorMessage        string
		showError           bool
		isEditing           bool
		isSaving            bool
		isLoading           bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IProfileComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,
) (contracts_App.IProfileComposer, error) {

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

func AddScopedIProfileComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IProfileComposer](cb, stemService.Ctor)
}

func (s *service) Render() app.UI {
	return app.Div().Class("profile-container").Body(
		app.If(s.showError, s.renderErrorBanner),

		// Page Header
		app.Div().Class("profile-header").Body(
			app.H1().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyYourProfile)),
			app.P().Class("profile-subtitle").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyManageYourPersonalInfo)),
		),

		// Show loading or profile cards
		app.If(s.isLoading, func() app.UI {
			return app.Div().Class("loading-container").Body(
				app.Div().Class("loading-spinner"),
				app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLoadingProfileDotDot)),
			)
		}).Else(func() app.UI {
			// Profile Cards Container
			return app.Div().Class("profile-cards").Body(
				s.renderPersonalInfoCard(),
			)
		}),
	)
}

func (s *service) renderPersonalInfoCard() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon personal-info-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPersonalInformation)),
					app.P().Class("card-description").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyYourNameAndBasicInfo)),
				),
			),
			app.If(!s.isEditing && !s.isSaving, func() app.UI {
				return app.Button().
					Class("btn-edit").
					OnClick(s.handleEditClick).
					Body(
						app.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
							<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
						</svg>`),
						app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEdit)),
					)
			}),
		),

		app.Div().Class("card-body").Body(
			app.If(!s.isEditing, func() app.UI {
				return app.Div().Class("info-rows").Body(
					app.Div().Class("info-row").Body(
						app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
						app.Span().Class("info-value").Text(s.email),
					),
					app.Div().Class("info-row").Body(
						app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFirstName)),
						app.Span().Class("info-value").Text(s.firstName),
					),
					app.Div().Class("info-row").Body(
						app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLastName)),
						app.Span().Class("info-value").Text(s.lastName),
					),
					app.Div().Class("info-row").Body(
						app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPhoneNumber)),
						app.Span().Class("info-value").Text(s.phoneNumber),
					),
				)
			}).Else(func() app.UI {
				saveButtonText := s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySave)
				if s.isSaving {
					saveButtonText = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySavingDotDot)
				}
				return app.Form().OnSubmit(s.handleSavePersonalInfo).Body(
					app.Div().Class("form-group").Body(
						app.Label().For("email").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyEmailAddress)),
						app.Input().
							Type("email").
							ID("email").
							Value(s.email).
							Disabled(true).
							Class("disabled-input"),
					),
					app.Div().Class("form-group").Body(
						app.Label().For("firstName").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFirstName)),
						app.Input().
							Type("text").
							ID("firstName").
							Value(s.firstName).
							OnInput(s.handleFirstNameInput).
							Required(true).
							Disabled(s.isSaving),
					),
					app.Div().Class("form-group").Body(
						app.Label().For("lastName").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLastName)),
						app.Input().
							Type("text").
							ID("lastName").
							Value(s.lastName).
							OnInput(s.handleLastNameInput).
							Required(true).
							Disabled(s.isSaving),
					),
					app.Div().Class("form-group").Body(
						app.Label().For("phoneNumber").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPhoneNumber)),
						app.Input().
							Type("tel").
							ID("phoneNumber").
							Value(s.phoneNumber).
							OnInput(s.handlePhoneNumberInput).
							Disabled(s.isSaving).
							Placeholder("+1 (555) 123-4567"),
					),
					app.Div().Class("button-group").Body(
						app.Button().
							Type("button").
							Class("btn-secondary").
							OnClick(s.handleCancelEdit).
							Disabled(s.isSaving).
							Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCancel)),
						app.Button().
							Type("submit").
							Class("btn-primary").
							Disabled(s.isSaving).
							Text(saveButtonText),
					),
				)
			}),
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
			app.Div().Class("error-text").Body(
				app.Span().Class("error-title").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)),
				app.Span().Class("error-message").Text(s.errorMessage),
			),
		),
		app.Button().
			Class("error-close").
			OnClick(s.handleCloseError).
			Body(
				app.Raw(`<svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
				</svg>`),
			),
	)
}

func (s *service) handleEditClick(ctx app.Context, e app.Event) {
	s.isEditing = true
	ctx.Update()
}

func (s *service) handleCancelEdit(ctx app.Context, e app.Event) {
	s.isEditing = false
	ctx.Update()
}

func (s *service) handleFirstNameInput(ctx app.Context, e app.Event) {
	s.firstName = ctx.JSSrc().Get("value").String()
}

func (s *service) handleLastNameInput(ctx app.Context, e app.Event) {
	s.lastName = ctx.JSSrc().Get("value").String()
}

func (s *service) handlePhoneNumberInput(ctx app.Context, e app.Event) {
	s.phoneNumber = ctx.JSSrc().Get("value").String()
}

func (s *service) handleSavePersonalInfo(ctx app.Context, e app.Event) {
	e.PreventDefault()

	log := zerolog.Ctx(s.AppContext).With().Str("component", "ProfileComposer").Logger()
	log.Info().
		Str("firstName", s.firstName).
		Str("lastName", s.lastName).
		Str("phoneNumber", s.phoneNumber).
		Msg("Saving profile")

	s.isSaving = true
	s.showError = false
	ctx.Update()

	ctx.Async(func() {
		response, err := s.managementApiClient.UpdateUserProfile(s.AppContext, &models.Profile{
			Subject:     s.subject,
			GivenName:   s.firstName,
			FamilyName:  s.lastName,
			PhoneNumber: s.phoneNumber,
		})

		ctx.Dispatch(func(ctx app.Context) {
			s.isSaving = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to update profile")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToSaveProfile)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Msg("Profile saved successfully")
				s.isEditing = false
				if response.Response != nil {
					s.firstName = response.Response.GivenName
					s.lastName = response.Response.FamilyName
					s.phoneNumber = response.Response.PhoneNumber
				}
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToSaveProfile)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) handleCloseError(ctx app.Context, e app.Event) {
	s.showError = false
	ctx.Update()
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "ProfileComposer").Logger()
	log.Info().Msg("Profile page mounted")

	// Load profile data
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserProfile(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to load profile")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToLoadProfile)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 && response.Response != nil {
				log.Info().Msg("Profile loaded successfully")
				profile := response.Response
				s.subject = profile.Subject
				s.email = profile.Email
				s.firstName = profile.GivenName
				s.lastName = profile.FamilyName
				s.phoneNumber = profile.PhoneNumber
				s.isClaimedDomain = profile.IsClaimedDomain
			} else {
				log.Error().Int("code", response.Code).Msg("Unexpected response code")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToLoadProfile)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Profile page navigated")
}

func (s *service) OnDismount() {
	log := zerolog.Ctx(s.AppContext).With().Logger()
	log.Info().Msg("Profile page dismounted")
}
