package AppConfigAccessor

import (
	"context"
	"encoding/base64"
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	oidc_login_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/contracts/config"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/common"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/config"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/utils"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_go_app_fetch "github.com/fluffy-bunny/fluffycore/go-app/fetch"
	fluffycore_go_app_js_loader "github.com/fluffy-bunny/fluffycore/go-app/js_loader"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		appConfig            *oidc_login_contracts_config.AppConfig
		oidcFlowAppConfig    *contracts_OIDCFlowAppConfig.OIDCFlowAppConfig
		wellknownCookieNames contracts_cookies.IWellknownCookieNames
	}
)

var stemService = (*service)(nil)

var _ contracts_config.IAppConfigAccessor = stemService

func (s *service) Ctor(wellknownCookieNames contracts_cookies.IWellknownCookieNames) (contracts_config.IAppConfigAccessor, error) {
	config, err := fluffycore_go_app_js_loader.LoadConfigFromJS[oidc_login_contracts_config.AppConfig](
		&fluffycore_go_app_js_loader.LoadConfigOptions{
			IsReadyFuncName:   "isAppConfigReady",
			GetConfigFuncName: "getAppConfig",
		},
	)
	if err != nil {
		return nil, err
	}
	return &service{
		appConfig:            config,
		wellknownCookieNames: wellknownCookieNames,
	}, nil
}

func AddScopedIAppConfigAccessor(cb di.ContainerBuilder) {
	di.AddScoped[contracts_config.IAppConfigAccessor](cb, stemService.Ctor)
}

func (s *service) GetAppConfig(ctx context.Context) *oidc_login_contracts_config.AppConfig {
	return s.appConfig
}

func OIDCFlowAppConfigKeyFromState(state string) string {
	return fmt.Sprintf("oidc_config_%s", state)
}
func ClearAllOIDCFlowAppConfigCache() {
	// Clear all keys starting with "oidc_config_"
	keys := utils.GetLocalStorageKeys()
	for _, key := range keys {
		if hasPrefix(key, "oidc_config_") {
			utils.RemoveLocalStorage(key)
		}
	}
}
func (s *service) GetOIDCFlowAppConfig(ctx context.Context) (*contracts_OIDCFlowAppConfig.OIDCFlowAppConfig, error) {

	log := zerolog.Ctx(ctx).With().Str("component", "AppConfigAccessor").Logger()
	if s.oidcFlowAppConfig != nil {
		return s.oidcFlowAppConfig, nil
	}
	var oidcFlowAppConfigKey string
	// Try to load from localStorage cache first based on authorization state
	authorizationStateCookie, err := s.GetOrCreateAuthorizationStateCookie(ctx)
	if err == nil && fluffycore_utils.IsNotEmptyOrNil(authorizationStateCookie.State) {
		oidcFlowAppConfigKey = OIDCFlowAppConfigKeyFromState(authorizationStateCookie.State)

		cached := s.loadOIDCFlowAppConfigFromCache(authorizationStateCookie.State)
		if cached != nil {
			s.oidcFlowAppConfig = cached
			log.Info().Str("state", authorizationStateCookie.State).Msg("✅ Loaded OIDC config from cache")
			return s.oidcFlowAppConfig, nil
		}
	}
	if authorizationStateCookie == nil || fluffycore_utils.IsEmptyOrNil(authorizationStateCookie.State) {
		// we will fake it here as we just use it
		return nil, status.Error(codes.FailedPrecondition, "authorization state cookie is missing or invalid")
	}

	// Lazy load the OIDC config from the Rage server
	if s.appConfig == nil || s.appConfig.RageBaseURL == "" {
		return nil, status.Error(codes.FailedPrecondition, "app config or RageBaseURL is missing")
	}

	// Fetch the config with cache busting
	url := s.appConfig.RageBaseURL + wellknown_echo.API_OIDCFlowAppConfig
	if common.AppVersion != "" {
		url = url + "?v=" + authorizationStateCookie.State
	}

	wr, err := fluffycore_go_app_fetch.FetchWrappedResponseT[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig](
		ctx,
		&fluffycore_go_app_fetch.CallInput{
			Method: "GET",
			Url:    url,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch OIDC Flow App Config from server")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch OIDC Flow App Config: %v", err))
	}
	s.oidcFlowAppConfig = wr.Response

	// Cache it with the authorization state if available

	s.saveOIDCFlowAppConfigToCache(ctx, oidcFlowAppConfigKey, s.oidcFlowAppConfig)

	return s.oidcFlowAppConfig, nil
}

// GetAuthorizationStateCookie reads the _authorization_state cookie and extracts the state value
func (s *service) GetAuthorizationStateCookie(ctx context.Context) (*contracts_config.AuthorizationStateCookie, error) {
	ccName := s.wellknownCookieNames.GetCookieName(contracts_cookies.CookieName_AuthorizationState)
	cookie, err := utils.GetCookie[contracts_config.AuthorizationStateCookie](ccName)
	if err != nil {
		return nil, status.Error(codes.NotFound, "authorization state cookie not found")
	}
	return &cookie, nil
}

// SetAuthorizationStateCookie writes the _authorization_state cookie
func (s *service) SetAuthorizationStateCookie(ctx context.Context, authStateCookie *contracts_config.AuthorizationStateCookie) error {
	log := zerolog.Ctx(ctx).With().Str("component", "AppConfigAccessor").Logger()
	opts := utils.CookieOptions{
		Path:     "/",
		MaxAge:   3600,
		SameSite: "Lax",
	}
	ccName := s.wellknownCookieNames.GetCookieName(contracts_cookies.CookieName_AuthorizationState)
	err := utils.SetCookie(ccName, *authStateCookie, opts)
	if err != nil {
		return fmt.Errorf("failed to set authorization state cookie: %w", err)
	}
	log.Info().Str("state", authStateCookie.State).Msg("✅ Authorization state cookie set")
	return nil
}

// GetOrCreateAuthorizationStateCookie gets existing cookie or creates a new one
func (s *service) GetOrCreateAuthorizationStateCookie(ctx context.Context) (*contracts_config.AuthorizationStateCookie, error) {
	// Try to get existing cookie
	authStateCookie, err := s.GetAuthorizationStateCookie(ctx)
	if err == nil && authStateCookie != nil && fluffycore_utils.IsNotEmptyOrNil(authStateCookie.State) {
		return authStateCookie, nil
	}

	// Create new cookie with generated state
	newState := generateRandomState()
	newCookie := &contracts_config.AuthorizationStateCookie{
		State: newState,
	}

	// Write the cookie
	if err := s.SetAuthorizationStateCookie(ctx, newCookie); err != nil {
		return nil, err
	}

	return newCookie, nil
}

// generateRandomState creates a random state value for OAuth
func generateRandomState() string {
	// Use crypto random if available, otherwise use timestamp-based
	randomBytes := make([]byte, 32)

	// Generate random bytes using JavaScript crypto.getRandomValues
	crypto := app.Window().Get("crypto")
	if !crypto.IsUndefined() && !crypto.IsNull() {
		// Create a Uint8Array
		uint8Array := app.Window().Get("Uint8Array").New(32)
		crypto.Call("getRandomValues", uint8Array)

		// Convert to Go bytes
		for i := 0; i < 32; i++ {
			randomBytes[i] = byte(uint8Array.Index(i).Int())
		}
	} else {
		// Fallback: use timestamp + random (less secure but works)
		timestamp := app.Window().Get("Date").Call("now").Int()
		for i := 0; i < 32; i++ {
			randomBytes[i] = byte((timestamp + i) % 256)
		}
	}

	// Return base64url encoded string (URL-safe)
	return base64.RawURLEncoding.EncodeToString(randomBytes)
}

// loadOIDCFlowAppConfigFromCache retrieves cached OIDC config from localStorage
func (s *service) loadOIDCFlowAppConfigFromCache(state string) *contracts_OIDCFlowAppConfig.OIDCFlowAppConfig {
	cacheKey := fmt.Sprintf("oidc_config_%s", state)

	config, err := utils.GetLocalStorage[contracts_OIDCFlowAppConfig.OIDCFlowAppConfig](cacheKey)
	if err != nil {
		return nil
	}

	return &config
}

// saveOIDCFlowAppConfigToCache stores OIDC config in localStorage
func (s *service) saveOIDCFlowAppConfigToCache(ctx context.Context, cacheKey string, config *contracts_OIDCFlowAppConfig.OIDCFlowAppConfig) {
	log := zerolog.Ctx(ctx).With().Str("component", "AppConfigAccessor").Logger()
	ClearAllOIDCFlowAppConfigCache()

	err := utils.SetLocalStorage(cacheKey, *config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cache config")
		return
	}
	log.Info().Str("cacheKey", cacheKey).Msg("✅ OIDC config cached in localStorage")
}

// hasPrefix helper for cache key checking
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
