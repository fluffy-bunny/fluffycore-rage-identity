package KeepSignedIn

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/ComposerBase"
	contracts_go_app_RageApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/contracts/RageApiClient"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_ComposerBase.ComposerBase

		rageApiClient contracts_go_app_RageApiClient.IRageApiClient

		keepSignedIn bool
		isLoading    bool
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

	// TODO: Call API to set keep signed in preference
	// For now, this is just a skeleton - will be implemented when REST APIs are fixed

	ctx.Async(func() {
		// Simulate API call
		// In the future, this will call the backend to persist the keep-signed-in preference

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			// TODO: Navigate to completion or next step
			// For now, just log the selection
			log.Info().Bool("keepSignedIn", s.keepSignedIn).Msg("Keep signed in preference recorded")
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
