package account_management_server

import (
	"context"
	"strconv"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	services_handlers_account_callback "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/callback"
	account_management_server_api_app_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/account_management_server/api_app_config"
	account_management_server_api_login "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/account_management_server/api_login"
	account_management_server_api_logout "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/account_management_server/api_logout"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	services "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services"
	services_ScopedMemoryCache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/ScopedMemoryCache"
	services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"
	services_handlers_healthz "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/healthz"
	services_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/session_with_options"
	pkg_types "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/types"
	pkg_version "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/version"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	contracts_startup "github.com/fluffy-bunny/fluffycore/echo/contracts/startup"
	fluffycore_echo_services_sessions_cookie_session "github.com/fluffy-bunny/fluffycore/echo/services/sessions/cookie_session"
	fluffycore_echo_services_sessions_cookie_session_store "github.com/fluffy-bunny/fluffycore/echo/services/sessions/cookie_session_store"
	fluffycore_echo_services_sessions_memory_session "github.com/fluffy-bunny/fluffycore/echo/services/sessions/memory_session"
	fluffycore_echo_services_sessions_memory_session_store "github.com/fluffy-bunny/fluffycore/echo/services/sessions/memory_session_store"
	fluffycore_echo_services_sessions_session_factory "github.com/fluffy-bunny/fluffycore/echo/services/sessions/session_factory"
	services_startup "github.com/fluffy-bunny/fluffycore/echo/services/startup"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
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
		log:    zlog.With().Str("runtime", "account_management_server").Caller().Logger(),
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
		Port: s.config.EchoAccount.Port,
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
	account_management_server_api_login.AddScopedIHandler(builder)
	account_management_server_api_logout.AddScopedIHandler(builder)
	account_management_server_api_app_config.AddScopedIHandler(builder)

	services_handlers_cache_busting_static_html.AddScopedIHandler(builder, s.config.AccountUIConfig.CacheBustingConfig)
	services_handlers_account_callback.AddScopedIHandler(builder)

	// sessions
	//----------------
	fluffycore_echo_services_sessions_memory_session_store.AddSingletonBackendSessionStore(builder)
	fluffycore_echo_services_sessions_cookie_session_store.AddSingletonCookieSessionStore(builder)
	fluffycore_echo_services_sessions_memory_session.AddTransientBackendSession(builder)
	fluffycore_echo_services_sessions_cookie_session.AddTransientCookieSession(builder)
	fluffycore_echo_services_sessions_session_factory.AddScopedSessionFactory(builder)
	services_session_with_options.AddScopedISessionWithOptions(builder,
		&contracts_session_with_options.SessionWithOptions{
			Name: "_rage_account_management_session",
		})

	services_ScopedMemoryCache.AddScopedIScopedMemoryCache(builder)

}
func (s *startup) RegisterStaticRoutes(e *echo.Echo) error {
	// i.e. e.Static("/css", "./css")
	e.Static("/static", "./static")
	return nil
}

// Configure
func (s *startup) Configure(e *echo.Echo, root di.Container) error {

	e.Use(EnsureCookieClaimsPrincipal(root))
	e.Use(EnsureLocalizer(root))
	if s.config.CORSConfig.Enabled {
		e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
			AllowOrigins:                             s.config.CORSConfig.AllowedOrigins,
			AllowMethods:                             s.config.CORSConfig.AllowedMethods,
			AllowHeaders:                             s.config.CORSConfig.AllowedHeaders,
			AllowCredentials:                         s.config.CORSConfig.AllowCredentials,
			UnsafeWildcardOriginWithAllowCredentials: s.config.CORSConfig.UnsafeWildcardOriginWithAllowCredentials,
			ExposeHeaders:                            s.config.CORSConfig.ExposeHeaders,
			MaxAge:                                   s.config.CORSConfig.MaxAge,
		}))
	}
	e.Use(EnsureAuth(root))
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

func EnsureLocalizer(_ di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//ctx := c.Request().Context()
			subContainer, ok := c.Get(fluffycore_echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			if !ok {
				return next(c)
			}
			// pull the SCOPED ILocalizer and initialize it.
			localizer := di.Get[contracts_localizer.ILocalizer](subContainer)
			accept := strings.TrimSpace(c.Request().Header.Get("Accept-Language"))
			if accept == "" {
				accept = "en"
			}
			localizer.Initialize(c)
			return next(c)
		}
	}
}

// EnsureContextLogger ...
func EnsureCookieClaimsPrincipal(_ di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//ctx := c.Request().Context()
			subContainer, ok := c.Get(fluffycore_echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			if !ok {
				return next(c)
			}
			wellknownCookies := di.Get[contracts_cookies.IWellknownCookies](subContainer)
			claimsPrincipal := di.Get[fluffycore_contracts_common.IClaimsPrincipal](subContainer)
			scopedMemoryCache := di.Get[contracts_cache.IScopedMemoryCache](subContainer)
			getAuthCookieResponse, err := wellknownCookies.GetAuthCookie(c)
			if err != nil ||
				getAuthCookieResponse == nil ||
				getAuthCookieResponse.AuthCookie == nil ||
				getAuthCookieResponse.AuthCookie.Identity == nil {
				return next(c)
			}
			rootIdentity := getAuthCookieResponse.AuthCookie.Identity
			scopedMemoryCache.Set("rootIdentity", rootIdentity)
			claimsPrincipal.AddClaim(
				fluffycore_contracts_common.Claim{
					Type:  fluffycore_echo_wellknown.ClaimTypeAuthenticated,
					Value: "true",
				},
				fluffycore_contracts_common.Claim{
					Type:  fluffycore_echo_wellknown.ClaimTypeSubject,
					Value: rootIdentity.Subject,
				}, fluffycore_contracts_common.Claim{
					Type:  "email",
					Value: rootIdentity.Email,
				}, fluffycore_contracts_common.Claim{
					Type:  "email_verified",
					Value: strconv.FormatBool(rootIdentity.EmailVerified),
				})

			return next(c)
		}
	}
}
