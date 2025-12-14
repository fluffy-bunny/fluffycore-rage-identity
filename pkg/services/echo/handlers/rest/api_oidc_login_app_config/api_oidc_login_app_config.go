package api_oidc_login_app_config

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config
	}

	BannerBranding struct {
		Title             string `json:"title,omitempty"`
		LogoURL           string `json:"logoUrl,omitempty"`
		ShowBannerVersion bool   `json:"showBannerVersion,omitempty"`
	}

	OIDCLoginAppConfig struct {
		RageBaseURL    string         `json:"rageBaseUrl"`
		BannerBranding BannerBranding `json:"bannerBranding,omitempty"`
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	// Register for /api/oidc-login-app-config
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_AppConfig,
	)
	// Register for /config/app.json (alternate path for WASM app)
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.Config_AppJSON,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API OIDC Login App Config godoc
// @Summary get the OIDC login app configuration.
// @Description This is the configuration for the WASM-based OIDC login app.
// @Tags root
// @Produce json
// @Success 200 {object} api_oidc_login_app_config.OIDCLoginAppConfig
// @Router /api/oidc-login-app-config [get]
func (s *service) Do(c echo.Context) error {
	// Use the configured OIDC base URL (the OIDC server on port 9044)
	baseURL := s.config.OIDCConfig.BaseUrl

	appConfig := &OIDCLoginAppConfig{
		RageBaseURL: baseURL,
		BannerBranding: BannerBranding{
			Title:             s.config.ApplicationName,
			LogoURL:           "web/apple-touch-icon-192x192.png",
			ShowBannerVersion: s.config.SystemConfig.DeveloperMode,
		},
	}
	return c.JSON(http.StatusOK, appConfig)
}
