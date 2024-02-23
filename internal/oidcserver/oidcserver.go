package oidcserver

import (
	"context"
	"strconv"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	services "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services"
	services_handlers_account_about "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/about"
	services_handlers_account_callback "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/callback"
	services_handlers_account_home "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/home"
	services_handlers_account_login "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/login"
	services_handlers_account_logout "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/logout"
	services_handlers_account_profile "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/account/profile"
	services_handlers_authorization_endpoint "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/authorization_endpoint"
	services_handlers_discovery_endpoint "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/discovery_endpoint"
	services_handlers_error "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/error"
	services_handlers_externalidp "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/externalidp"
	services_handlers_forgotpassword "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/forgotpassword"
	services_handlers_healthz "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/healthz"
	services_handlers_jwks_endpoint "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/jwks_endpoint"
	services_handlers_oauth2_callback "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/oauth2/callback"
	services_handlers_oidclogin "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/oidclogin"
	services_handlers_oidcloginpassword "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/oidcloginpassword"
	services_handlers_passwordreset "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/passwordreset"
	services_handlers_signup "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/signup"
	services_handlers_swagger "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/swagger"
	services_handlers_token_endpoint "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/token_endpoint"
	services_handlers_verifycode "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/verifycode"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	contracts_startup "github.com/fluffy-bunny/fluffycore/echo/contracts/startup"
	services_startup "github.com/fluffy-bunny/fluffycore/echo/services/startup"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type (
	startup struct {
		services_startup.StartupBase
		config *contracts_config.Config
		log    zerolog.Logger
	}
)

func init() {
	var _ contracts_startup.IStartup = (*startup)(nil)
}

// GetConfigOptions ...
func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	return &fluffycore_contracts_runtime.ConfigOptions{
		RootConfig:  []byte(contracts_config.ConfigDefaultJSON),
		Destination: s.config,
	}
}
func NewStartup() contracts_startup.IStartup {
	myStartup := &startup{
		config: &contracts_config.Config{},
		log:    zlog.With().Str("runtime", "oidcserver").Caller().Logger(),
	}
	hooks := &contracts_startup.Hooks{
		PostBuildHook:   myStartup.PostBuildHook,
		PreStartHook:    myStartup.PreStartHook,
		PreShutdownHook: myStartup.PreShutdownHook,
	}

	myStartup.AddHooks(hooks)
	return myStartup
}

// ConfigureServices ...
func (s *startup) ConfigureServices(builder di.ContainerBuilder) error {
	s.SetOptions(&contracts_startup.Options{
		Port: s.config.Echo.Port,
	})
	s.addAppHandlers(builder)
	services.ConfigureServices(context.TODO(), s.config, builder)
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

	// Account Handlers
	//--------------------------------------------------------
	services_handlers_account_home.AddScopedIHandler(builder)
	services_handlers_account_about.AddScopedIHandler(builder)
	services_handlers_account_callback.AddScopedIHandler(builder)
	services_handlers_account_login.AddScopedIHandler(builder)
	services_handlers_account_logout.AddScopedIHandler(builder)
	services_handlers_account_profile.AddScopedIHandler(builder)

	// OIDC Handlers
	//--------------------------------------------------------
	services_handlers_signup.AddScopedIHandler(builder)
	services_handlers_error.AddScopedIHandler(builder)
	services_handlers_oidclogin.AddScopedIHandler(builder)
	services_handlers_oidcloginpassword.AddScopedIHandler(builder)
	services_handlers_externalidp.AddScopedIHandler(builder)
	services_handlers_swagger.AddScopedIHandler(builder)
	services_handlers_discovery_endpoint.AddScopedIHandler(builder)
	services_handlers_jwks_endpoint.AddScopedIHandler(builder)
	services_handlers_authorization_endpoint.AddScopedIHandler(builder)
	services_handlers_token_endpoint.AddScopedIHandler(builder)
	services_handlers_oauth2_callback.AddScopedIHandler(builder)
	services_handlers_forgotpassword.AddScopedIHandler(builder)
	services_handlers_verifycode.AddScopedIHandler(builder)
	services_handlers_passwordreset.AddScopedIHandler(builder)
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
	e.Use(EnsureAuth(root))
	return nil
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
				accept = "en-US"
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

			getAuthCookieResponse, err := wellknownCookies.GetAuthCookie(c)
			if err != nil ||
				getAuthCookieResponse == nil ||
				getAuthCookieResponse.AuthCookie == nil ||
				getAuthCookieResponse.AuthCookie.Identity == nil {
				return next(c)
			}
			rootIdentity := getAuthCookieResponse.AuthCookie.Identity

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
