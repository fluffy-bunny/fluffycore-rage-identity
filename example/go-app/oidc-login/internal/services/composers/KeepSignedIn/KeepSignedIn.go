package KeepSignedIn

import (
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

		keepSignedIn bool
		isLoading    bool
		errorMessage string
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IKeepSignedInComposer = stemService

func (s *service) Ctor(
	container di.Container,
	localizer contracts_Localizer.ILocalizer,
	appContext contracts_App.AppContext,
	rageApiClient contracts_go_app_RageApiClient.IRageApiClient,
) (contracts_App.IKeepSignedInComposer, error) {

	return &service{
		rageApiClient: rageApiClient,
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		keepSignedIn: false, // Default to not keeping signed in
	}, nil
}

func AddScopedIKeepSignedInComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IKeepSignedInComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "KeepSignedInComposer").Logger()
	log.Info().Msg("OnMount called for KeepSignedInComposer")

	// TODO: Validate session state
	// This page should only be shown if there's a valid auth session
}

func (s *service) Render() app.UI {
	return app.Div().Class("step-container").Body(
		app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyKeepMeSignedIn)),
		app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyKeepMeSignedInDescription)),

		app.If(s.errorMessage != "",
			func() app.UI {
				return app.Div().Class("error-message").Text(s.errorMessage)
			},
		),

		app.Form().OnSubmit(s.handleSubmit).Body(
			app.Div().Class("form-group").Body(
				app.Label().Class("checkbox-label").Body(
					app.Input().
						Type("checkbox").
						ID("keep-signed-in").
						Checked(s.keepSignedIn).
						OnChange(func(ctx app.Context, e app.Event) {
							s.keepSignedIn = ctx.JSSrc().Get("checked").Bool()
						}),
					app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyKeepMeSignedInCheckbox)),
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
					Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyContinue)),
			),
		),
	)
}

func (s *service) handleSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	log := zerolog.Ctx(s.AppContext).With().Str("component", "KeepSignedInComposer").Logger()
	log.Info().Bool("keepSignedIn", s.keepSignedIn).Msg("Keep signed in form submitted")

	s.isLoading = true
	s.errorMessage = ""

	ctx.Async(func() {
		// Call the keep-signed-in API
		response, err := s.rageApiClient.KeepSignedIn(s.AppContext,
			&login_models.KeepSignedInRequest{
				KeepSignedIn: s.keepSignedIn,
			})

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				// Network or parsing error
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Err(err).Msg("KeepSignedIn API call failed")
				return
			}

			if response.Response == nil {
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyInvalidResponseFromServer)
				log.Error().Msg("KeepSignedIn API returned invalid response")
				return
			}

			// Check HTTP status code for errors
			if response.Code >= 400 {
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
				log.Error().Msgf("KeepSignedIn API failed with code %d", response.Code)
				return
			}

			// Success - check directive for next action
			if response.Response.Directive == login_models.DIRECTIVE_Redirect {
				// Final redirect - authentication complete
				if response.Response.DirectiveRedirect != nil {
					redirectURI := response.Response.DirectiveRedirect.RedirectURI
					log.Info().Msgf("Keep signed in successful, redirecting to: %s", redirectURI)
					// Use window.location.href for external redirects (full page navigation)
					app.Window().Get("location").Set("href", redirectURI)
				}
			} else {
				// Unexpected directive
				log.Warn().Msgf("Unexpected directive: %s", response.Response.Directive)
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)
			}
		})
	})
}

func (s *service) handleBackClick(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "KeepSignedInComposer").Logger()
	log.Info().Msg("Back button clicked")

	// TODO: Navigate to previous page in login flow
	// This will depend on the flow state
	ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
}
