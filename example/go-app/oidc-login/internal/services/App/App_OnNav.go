package App

import (
	"strings"

	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

// OnNav is called when ctx.Navigate() is used within the app
func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()

	// Get the full path
	path := app.Window().URL().Path
	log.Info().Msgf("OnNav - path: %s", path)

	// Strip any route prefix (e.g., /oidc-login)
	// The path might be /oidc-login/password, we need just /password
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 1 {
		// First part might be the prefix, second part is the actual route
		path = "/" + parts[len(parts)-1]
	}

	log.Info().Msgf("OnNav - normalized path: %s", path)

	// Update currentPage based on the path
	switch path {
	case string(contracts_routes.WellknownRoute_Password):
		s.currentPage = contracts_routes.WellknownRoute_Password
	case string(contracts_routes.WellknownRoute_CreateAccount):
		s.currentPage = contracts_routes.WellknownRoute_CreateAccount
	case string(contracts_routes.WellknownRoute_ForgotPassword):
		s.currentPage = contracts_routes.WellknownRoute_ForgotPassword
	case string(contracts_routes.WellknownRoute_ResetPassword):
		s.currentPage = contracts_routes.WellknownRoute_ResetPassword
	case string(contracts_routes.WellknownRoute_VerifyCode):
		s.currentPage = contracts_routes.WellknownRoute_VerifyCode
	default:
		s.currentPage = contracts_routes.WellknownRoute_Home
	}

	log.Info().Msgf("OnNav - currentPage set to: %s", s.currentPage)
}
