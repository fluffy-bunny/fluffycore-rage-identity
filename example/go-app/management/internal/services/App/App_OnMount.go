package App

import (
	"strings"

	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rs/zerolog"
)

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()

	// Get the actual current path from the browser
	currentPath := app.Window().URL().Path
	log.Info().
		Str("currentPage", string(s.currentPage)).
		Str("actualPath", currentPath).
		Msg("App OnMount")

	// Check if cookies have been accepted
	var cookiesAccepted bool
	ctx.LocalStorage().Get("cookiesAccepted", &cookiesAccepted)
	s.showCookieBanner = !cookiesAccepted

	// Check authentication status by calling GetUserInfo
	// Note: GetUserInfo method needs to be added to IManagementApiClient
	// Expected signature: GetUserInfo(ctx context.Context) (*WrappedResonseT[UserInfoResponse], error)
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserProfile(s.AppContext)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Debug().Err(err).Msg("GetUserInfo failed - user not authenticated")
				s.isAuthenticated = false
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 && response.Response != nil {
				// User is authenticated
				log.Info().
					Interface("response", response.Response).
					Msg("User authenticated")

				s.isAuthenticated = true
				s.isClaimedDomain = response.Response.IsClaimedDomain
				s.profile = response.Response

				// If authenticated, stay on current route - don't auto-navigate
				currentPath := app.Window().URL().Path
				log.Info().
					Str("currentPath", currentPath).
					Bool("isClaimedDomain", s.isClaimedDomain).
					Msg("User authenticated, staying on current route")
			} else {
				// Not authorized
				log.Debug().Int("code", response.Code).Msg("GetUserInfo returned non-200 - user not authenticated")
				s.isAuthenticated = false

				// Get current path to determine action
				currentPath := app.Window().URL().Path

				// If we're on the home page, just show it without auth (it's public)
				if s.isHomePath(currentPath) {
					log.Debug().Str("currentPath", currentPath).Msg("User not authenticated on home page - this is OK")
					ctx.Update()
					return
				}

				// Check if current path is a protected route - redirect regardless of response code
				parts := strings.Split(strings.TrimPrefix(currentPath, "/"), "/")
				var normalizedPath string
				if len(parts) > 1 {
					normalizedPath = "/" + parts[len(parts)-1]
				} else {
					normalizedPath = "/" + parts[0]
				}

				pathRoute := contracts_routes.WellknownRoute(normalizedPath)
				switch pathRoute {
				case contracts_routes.WellknownRoute_Profile,
					contracts_routes.WellknownRoute_PasswordManager,
					contracts_routes.WellknownRoute_PasskeyManager,
					contracts_routes.WellknownRoute_LinkedAccounts:
					log.Info().Str("currentPath", currentPath).Msg("User not authenticated on protected route, redirecting to home")
					appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)
					baseURL := app.Window().Get("location").Get("origin").String() + "/" + appConfig.BaseHREF + "/"
					app.Window().Get("location").Call("assign", baseURL)
					return
				}
			}

			ctx.Update()
		})
	})

	// Page is already set via SetCurrentPage() from the route handler
}
