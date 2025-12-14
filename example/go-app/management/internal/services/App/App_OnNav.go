package App

import (
	"strings"

	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

// OnNav is called when ctx.Navigate() is used within the app
func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()

	// Get the full path
	path := app.Window().URL().Path
	log.Info().Msgf("OnNav - path: %s", path)

	// Strip any route prefix (e.g., /management)
	// The path might be /management/password, we need just /password
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 1 {
		// First part might be the prefix, second part is the actual route
		path = "/" + parts[len(parts)-1]
	}

	log.Info().Msgf("OnNav - normalized path: %s", path)

	// Update currentPage based on the path
	pathRoute := contracts_routes.WellknownRoute(path)

	// Note: Authentication check removed from here to avoid race condition
	// OnMount will handle route protection after GetUserInfo completes

	switch pathRoute {
	case contracts_routes.WellknownRoute_Profile:
		s.currentPage = contracts_routes.WellknownRoute_Profile
	case contracts_routes.WellknownRoute_PasswordManager:
		s.currentPage = contracts_routes.WellknownRoute_PasswordManager
	case contracts_routes.WellknownRoute_LinkedAccounts:
		s.currentPage = contracts_routes.WellknownRoute_LinkedAccounts
	default:
		s.currentPage = contracts_routes.WellknownRoute_Home
	}

	log.Info().Msgf("OnNav - currentPage set to: %s", s.currentPage)
}
