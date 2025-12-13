package App

import (
	"strings"

	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Str("currentPage", string(s.currentPage)).Msg("OnMount called for App")

	// Check if cookies have been accepted
	var cookiesAccepted bool
	ctx.LocalStorage().Get("cookiesAccepted", &cookiesAccepted)
	s.showCookieBanner = !cookiesAccepted

	s.fetchOIDCFlowAppConfig(ctx)
	s.fetchManifest(ctx)
}

func (s *service) fetchOIDCFlowAppConfig(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Fetching OIDC Flow App Config")

	ctx.Async(func() {
		_, err := s.appConfigAccessor.GetOIDCFlowAppConfig(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch OIDC Flow App Config")
			return
		}
	})
}

// fetchManifest fetches the configuration from the server and merges it with defaults
func (s *service) fetchManifest(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	ctx.Async(func() {

		log.Info().Msg("Fetching manifest from server...")
		getManifestResponse, err := s.rageApiClient.GetManifest(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching config from server")
		}
		log.Info().Interface("manifest", getManifestResponse.Response).Msg("Manifest loaded from server")

		code := getManifestResponse.Code
		manifest := getManifestResponse.Response
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Error().Err(err).Int("code", code).Msg("Config fetch failed, using defaults")
			} else {
				if manifest == nil {
					log.Warn().Int("code", code).Msg("No config received from server, using defaults")
				} else {
					// Merge the server config with default config
					// Server config takes precedence over defaults
					s.manifest = manifest
					s.rageApiClient.SetCachedManifest(ctx, manifest)
					// Check if manifest specifies a landing page
					// Only navigate if we're still at the root path (initial load)
					currentPath := ctx.Page().URL().Path
					rootPrefix := app.Getenv("GOAPP_ROOT_PREFIX")
					isAtRoot := currentPath == "/" || currentPath == rootPrefix || currentPath == strings.TrimRight(rootPrefix, "/")

					log.Info().
						Str("currentPath", currentPath).
						Str("rootPrefix", rootPrefix).
						Bool("isAtRoot", isAtRoot).
						Bool("hasLandingPage", manifest.LandingPage != nil).
						Msg("fetchManifest")

					// Only navigate if we're at the root
					if isAtRoot {
						if manifest.LandingPage != nil {
							log.Info().Interface("landingPage", manifest.LandingPage.Page).Msg("Manifest has landing page")
							// Map landing page to wizard step
							switch manifest.LandingPage.Page {
							case models_api_manifest.PageVerifyCode:
								log.Info().Msg("Navigating to VerifyCode page from manifest")
								ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_VerifyCode))
							default:
								// Default login page
								ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Home))
							}
						}
					}
				}
			}
		})
	})
}
