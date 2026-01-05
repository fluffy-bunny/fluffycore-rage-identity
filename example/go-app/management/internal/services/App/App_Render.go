package App

import (
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

func (s *service) Render() app.UI {

	// Render dashboard layout for authenticated pages
	return app.Div().Class("dashboard-layout").Body(
		s.renderDashboardHeader(),
		app.Div().Class("dashboard-body").Body(
			s.renderSidebar(),
			app.Main().Class("dashboard-main").Body(
				s.renderCurrentPage(),
			),
		),
		app.If(s.showCookieBanner, s.renderCookieBanner),
	)
}

func (s *service) renderCurrentPage() app.UI {
	// Protected routes - require authentication
	switch s.currentPage {
	case contracts_routes.WellknownRoute_Profile,
		contracts_routes.WellknownRoute_PasswordManager,
		contracts_routes.WellknownRoute_PasskeyManager,
		contracts_routes.WellknownRoute_LinkedAccounts,
		contracts_routes.WellknownRoute_Preferences:
		// Check if user is authenticated before rendering protected content
		if !s.isAuthenticated {
			// Show loading or blank state while auth check is in progress
			return app.Div().Class("loading-container").Body(
				app.Div().Class("loading-spinner"),
			)
		}
	}

	// Render the appropriate page
	switch s.currentPage {
	case contracts_routes.WellknownRoute_Profile:
		return s.profileComposer
	case contracts_routes.WellknownRoute_PasswordManager:
		return s.passwordManagerComposer
	case contracts_routes.WellknownRoute_PasskeyManager:
		return s.passkeyManagerComposer
	case contracts_routes.WellknownRoute_LinkedAccounts:
		return s.linkedAccountsComposer
	case contracts_routes.WellknownRoute_Preferences:
		return s.preferencesComposer
	default:
		return s.renderHomePage()
	}
}

func (s *service) renderHomePage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Home Page")
	if s.isAuthenticated {
		// Authenticated users see the home page at /
		return s.homeComposer
	}
	// Unauthenticated users also see home page
	return s.homeComposer
}
