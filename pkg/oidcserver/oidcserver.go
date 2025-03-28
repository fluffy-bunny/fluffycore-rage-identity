package oidcserver

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	services "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services"
	services_ScopedMemoryCache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/ScopedMemoryCache"
	services_handlers_api "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/api"
	services_handlers_authorization_endpoint "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/authorization_endpoint"
	services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"
	services_handlers_discovery_endpoint "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/discovery_endpoint"
	services_handlers_error "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/error"
	services_handlers_externalidp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/externalidp"
	services_handlers_externalidp_api "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/externalidp/api"
	services_handlers_forgotpassword "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/forgotpassword"
	services_handlers_healthz "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/healthz"
	services_handlers_jwks_endpoint "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/jwks_endpoint"
	services_handlers_oauth2_callback "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/oauth2/callback"
	services_handlers_oidclogin "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/oidclogin"
	services_handlers_oidcloginpasskey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/oidcloginpasskey"
	services_handlers_oidcloginpassword "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/oidcloginpassword"
	services_handlers_oidclogintotp "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/oidclogintotp"
	services_handlers_passwordreset "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/passwordreset"
	services_handlers_rest_api_appsettings "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_appsettings"
	services_handlers_rest_api_isauthorized "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_isauthorized"
	services_handlers_rest_api_login_password "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_login_password"
	services_handlers_rest_api_login_username_phase_one "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_login_username_phase_one"
	services_handlers_rest_api_logout "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_logout"
	services_handlers_rest_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_manifest"
	services_handlers_rest_api_password_reset_finish "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_password_reset_finish"
	services_handlers_rest_api_password_reset_start "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_password_reset_start"
	services_handlers_rest_api_signup "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_signup"
	services_handlers_rest_api_start_over "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_start_over"
	services_handlers_rest_api_user_identity_info "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_user_identity_info"
	services_handlers_rest_api_user_linked_accounts "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_user_linked_accounts"
	services_handlers_rest_api_user_remove_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_user_remove_passkey"
	services_handlers_rest_api_verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_verify_code"
	services_handlers_rest_api_verify_code_begin "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_verify_code_begin"
	services_handlers_rest_api_verify_password_strength "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_verify_password_strength"
	services_handlers_rest_api_verify_username "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/rest/api_verify_username"
	services_handlers_signup "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/signup"
	services_handlers_swagger "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/swagger"
	services_handlers_token_endpoint "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/token_endpoint"
	services_handlers_verifycode "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/verifycode"
	services_handlers_webauthn_loginbegin "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn/loginbegin"
	services_handlers_webauthn_loginfinish "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn/loginfinish"
	services_handlers_webauthn_registrationbegin "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn/registrationbegin"
	services_handlers_webauthn_registrationfinish "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn/registrationfinish"
	services_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/oidc_session"
	pkg_types "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/types"
	pkg_version "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/version"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidcuser "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
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
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

type (
	startup struct {
		services_startup.StartupBase
		config *contracts_config.Config
		log    zerolog.Logger
		ext    pkg_types.ConfigureServices
	}
)

func init() {
	var _ contracts_startup.IStartup = (*startup)(nil)
}

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
		log:    zlog.With().Str("runtime", "oidcserver").Caller().Logger(),
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
		Port: s.config.Echo.Port,
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
func (s *startup) OnLoadSeedUsers(ctx context.Context,
	container di.Container) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnLoadSeedUsers").Logger()
	config := s.config
	fileContent, err := os.ReadFile(config.ConfigFiles.SeedUsersPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read OIDCClientPath - may not be a problem if clients are comming from a DB")
		return nil
	}
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")

	rageUsers := &proto_oidc_models.RageUsers{}

	err = protojson.Unmarshal([]byte(fixedFileContent), rageUsers)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	log.Info().Interface("rageUsers", rageUsers).Msg("rageUsers")

	rageUserService := di.Get[proto_oidcuser.IFluffyCoreRageUserServiceServer](container)
	for _, rageUser := range rageUsers.Users {
		_, err := rageUserService.CreateRageUser(ctx, &proto_oidcuser.CreateRageUserRequest{
			User: rageUser,
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				continue
			}
			log.Error().Err(err).Msg("failed to CreateRageUser")
			return err
		}
	}
	return nil
}
func (s *startup) PostBuildHook(container di.Container) error {
	s.log.Info().Msg("PostBuildHook")
	s.OnLoadSeedUsers(context.Background(), container)
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
	services_handlers_signup.AddScopedIHandler(builder)
	services_handlers_error.AddScopedIHandler(builder)
	services_handlers_oidclogin.AddScopedIHandler(builder)
	services_handlers_oidcloginpassword.AddScopedIHandler(builder)
	services_handlers_rest_api_appsettings.AddScopedIHandler(builder, s.config.ApiAppSettings)
	services_handlers_rest_api_manifest.AddScopedIHandler(builder)
	services_handlers_rest_api_start_over.AddScopedIHandler(builder)
	services_handlers_rest_api_verify_code.AddScopedIHandler(builder)
	services_handlers_rest_api_verify_code_begin.AddScopedIHandler(builder)
	services_handlers_rest_api_signup.AddScopedIHandler(builder)
	services_handlers_rest_api_logout.AddScopedIHandler(builder)
	services_handlers_rest_api_user_identity_info.AddScopedIHandler(builder)
	services_handlers_rest_api_user_linked_accounts.AddScopedIHandler(builder)
	services_handlers_rest_api_verify_username.AddScopedIHandler(builder)
	services_handlers_rest_api_password_reset_start.AddScopedIHandler(builder)
	services_handlers_rest_api_password_reset_finish.AddScopedIHandler(builder)
	services_handlers_rest_api_login_password.AddScopedIHandler(builder)
	services_handlers_rest_api_isauthorized.AddScopedIHandler(builder)
	services_handlers_rest_api_login_username_phase_one.AddScopedIHandler(builder)
	services_handlers_rest_api_verify_password_strength.AddScopedIHandler(builder)
	services_handlers_oidclogintotp.AddScopedIHandler(builder)
	services_handlers_externalidp.AddScopedIHandler(builder)
	services_handlers_externalidp_api.AddScopedIHandler(builder)
	services_handlers_swagger.AddScopedIHandler(builder)
	services_handlers_discovery_endpoint.AddScopedIHandler(builder)
	services_handlers_jwks_endpoint.AddScopedIHandler(builder)
	services_handlers_authorization_endpoint.AddScopedIHandler(builder)
	services_handlers_token_endpoint.AddScopedIHandler(builder)
	services_handlers_oauth2_callback.AddScopedIHandler(builder, s.config.OIDCConfig.OAuth2CallbackPath)
	services_handlers_forgotpassword.AddScopedIHandler(builder)
	services_handlers_verifycode.AddScopedIHandler(builder)
	services_handlers_passwordreset.AddScopedIHandler(builder)
	services_handlers_api.AddScopedIHandler(builder)

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

	if s.config.WebAuthNConfig != nil && s.config.WebAuthNConfig.Enabled {
		services_handlers_oidcloginpasskey.AddScopedIHandler(builder)
		services_handlers_rest_api_user_remove_passkey.AddScopedIHandler(builder)
		// WebAuthN Handlers
		//--------------------------------------------------------

		services_handlers_webauthn_registrationbegin.AddScopedIHandler(builder)
		services_handlers_webauthn_registrationfinish.AddScopedIHandler(builder)
		services_handlers_webauthn_loginbegin.AddScopedIHandler(builder)
		services_handlers_webauthn_loginfinish.AddScopedIHandler(builder)
	}
	// sessions
	//----------------
	fluffycore_echo_services_sessions_memory_session_store.AddSingletonBackendSessionStore(builder)
	fluffycore_echo_services_sessions_cookie_session_store.AddSingletonCookieSessionStore(builder)
	fluffycore_echo_services_sessions_memory_session.AddTransientBackendSession(builder)
	fluffycore_echo_services_sessions_cookie_session.AddTransientCookieSession(builder)
	fluffycore_echo_services_sessions_session_factory.AddScopedSessionFactory(builder)
	services_oidc_session.AddScopedIOIDCSession(builder)

	services_ScopedMemoryCache.AddScopedIScopedMemoryCache(builder)

}
func (s *startup) RegisterStaticRoutes(e *echo.Echo) error {
	// i.e. e.Static("/css", "./css")
	e.Static("/static", "./static")
	return nil
}

// Configure
func (s *startup) Configure(e *echo.Echo, root di.Container) error {

	var sameSite http.SameSite = http.SameSiteStrictMode
	httpOnly := false
	if s.config.Echo.DisableSecureCookies {
		sameSite = 0
	}
	e.Use(echo_middleware.CSRFWithConfig(echo_middleware.CSRFConfig{
		TokenLookup:    "header:X-Csrf-Token,form:csrf",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: httpOnly,
		CookieSameSite: sameSite,
		CookieDomain:   s.config.CookieConfig.Domain,
		Skipper: func(c echo.Context) bool {
			if s.config.CSRFConfig.SkipApi {
				if strings.Contains(c.Request().URL.Path, "/api") {
					return true
				}
			}
			csrfSkipperPaths := CSRFSkipperPaths()
			currentPath := c.Request().URL.Path
			_, ok := csrfSkipperPaths[currentPath]
			return ok

		},
	}))

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
	wellknown_echo.ManagementPath,
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
