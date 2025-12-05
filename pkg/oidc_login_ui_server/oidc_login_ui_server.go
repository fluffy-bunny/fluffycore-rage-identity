package oidc_login_ui_server

import (
	"context"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services"
	services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"
	services_handlers_healthz "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/healthz"
	services_handlers_rest_api_oidc_login_app_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_oidc_login_app_config"
	pkg_types "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/types"
	pkg_version "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/version"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	contracts_startup "github.com/fluffy-bunny/fluffycore/echo/contracts/startup"
	services_startup "github.com/fluffy-bunny/fluffycore/echo/services/startup"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type (
	startup struct {
		services_startup.StartupBase
		config *contracts_config.Config
		log    zerolog.Logger
		ext    pkg_types.ConfigureServices
	}
)

var _ contracts_startup.IStartup = (*startup)(nil)

type WithOption func(startup *startup)

func WithConfigureServices(ext pkg_types.ConfigureServices) WithOption {
	return func(startup *startup) {
		startup.ext = ext
	}
}

// GetConfigOptions ...
func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	return &fluffycore_contracts_runtime.ConfigOptions{
		RootConfig:  []byte(contracts_config.ConfigDefaultJSON),
		Destination: s.config,
		EnvPrefix:   "RAGE",
	}
}
func NewStartup(options ...WithOption) contracts_startup.IStartup {
	myStartup := &startup{
		config: &contracts_config.Config{},
		log:    zlog.With().Str("runtime", "oidc_login_ui_server").Caller().Logger(),
	}
	hooks := &contracts_startup.Hooks{
		PostBuildHook:   myStartup.PostBuildHook,
		PreStartHook:    myStartup.PreStartHook,
		PreShutdownHook: myStartup.PreShutdownHook,
	}

	myStartup.AddHooks(hooks)
	for _, option := range options {
		option(myStartup)
	}
	return myStartup
}

// ConfigureServices ...
func (s *startup) ConfigureServices(builder di.ContainerBuilder) error {

	s.SetOptions(&contracts_startup.Options{
		Port: s.config.EchoOIDCUI.Port,
	})
	s.addAppHandlers(builder)
	services.ConfigureServices(context.TODO(), s.config, builder)
	if s.ext != nil {
		s.ext(context.TODO(), s.config, builder)
	}
	return nil
}

func (s *startup) PreStartHook(echo *echo.Echo) error {
	s.log.Info().Msg("PreStartHook")

	return nil
}

func (s *startup) PostBuildHook(container di.Container) error {
	s.log.Info().Msg("PostBuildHook")
	return nil
}
func (s *startup) PreShutdownHook(echo *echo.Echo) error {
	s.log.Info().Msg("PreShutdownHook")
	return nil
}
func (s *startup) addAppHandlers(builder di.ContainerBuilder) {
	// App Handlers
	//--------------------------------------------------------
	services_handlers_healthz.AddScopedIHandler(builder)

	// OIDC Handlers
	//--------------------------------------------------------

	guid := xid.New().String()
	if pkg_version.Version() != "dev-build" {
		guid = pkg_version.Version()
	}
	s.config.OIDCUIConfig.CacheBustingConfig.ReplaceParams = []*contracts_config.KeyValuePair{
		{
			Key:   "{title}",
			Value: s.config.ApplicationName,
		},
		{
			Key:   "{version}",
			Value: guid,
		},
	}
	services_handlers_cache_busting_static_html.AddScopedIHandler(builder, s.config.OIDCUIConfig.CacheBustingConfig)

	services_handlers_rest_api_oidc_login_app_config.AddScopedIHandler(builder)
}
func (s *startup) RegisterStaticRoutes(e *echo.Echo) error {
	// i.e. e.Static("/css", "./css")
	e.Static("/static", "./static")
	return nil
}

// Configure
func (s *startup) Configure(e *echo.Echo, root di.Container) error {

	e.Use(noCacheMiddleware)

	return nil
}

var noCachePaths = []string{
	wellknown_echo.HomePath,
	wellknown_echo.OIDCLoginUIPath,
}

func noCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestUrlPath := c.Request().URL.Path
		noCacheIt := strings.HasSuffix(requestUrlPath, "index.html")
		for _, path := range noCachePaths {
			if noCacheIt || requestUrlPath == path {
				noCacheIt = true
				break
			}
		}
		if noCacheIt {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Response().Header().Set("Pragma", "no-cache")
			c.Response().Header().Set("Expires", "0")
		}
		return next(c)
	}
}
