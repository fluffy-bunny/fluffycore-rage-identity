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
	switch s.currentPage {

	case contracts_routes.WellknownRoute_Profile:
		return s.profileComposer
	case contracts_routes.WellknownRoute_PasswordManager:
		return s.passwordManagerComposer
	case contracts_routes.WellknownRoute_LinkedAccounts:
		return s.linkedAccountsComposer
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
