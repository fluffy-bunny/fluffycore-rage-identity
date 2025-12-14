package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	example_auth "github.com/fluffy-bunny/fluffycore-rage-identity/example/auth"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/contracts/config"
	service_AuthorizationCodeClaimsAugmentor "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/AuthorizationCodeClaimsAugmentor"
	services_AuthorizationCodeClaimsAugmentor "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/AuthorizationCodeClaimsAugmentor"
	services_EmailTemplateData "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/EmailTemplateData"
	services_EventSink "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/EventSink"
	services_handlers_account_about "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/about"
	services_handlers_account_api_api_linked_accounts "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/api/api_linked_accounts"
	services_handlers_account_api_login "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/api/api_login"
	services_handlers_account_api_api_user_profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/api/api_user_profile"
	services_handlers_account_callback "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/callback"
	services_handlers_account_home "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/home"
	services_handlers_account_logout "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/logout"
	services_handlers_account_passkey_management "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/passkey_management"
	services_handlers_account_personal_information "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/personal_information"
	services_handlers_account_profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/profile"
	services_handlers_account_totp_management "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/totp_management"
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/oidcflowstore"
	services_user_id_generator "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/user_id_generator"
	services_oidcuser_inmemory "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/userstore/inmemory"
	example_version "github.com/fluffy-bunny/fluffycore-rage-identity/example/version"
	rage_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
	rage_runtime "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/runtime"
	services_ScopedMemoryCache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/ScopedMemoryCache"
	services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"
	services_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/session_with_options"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_async "github.com/fluffy-bunny/fluffycore/async"
	fluffycore_cobracore_cmd "github.com/fluffy-bunny/fluffycore/cobracore/cmd"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
	fluffycore_contracts_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/contracts/middleware/auth/jwt"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	fluffycore_echo_services_sessions_cookie_session "github.com/fluffy-bunny/fluffycore/echo/services/sessions/cookie_session"
	fluffycore_echo_services_sessions_cookie_session_store "github.com/fluffy-bunny/fluffycore/echo/services/sessions/cookie_session_store"
	fluffycore_echo_services_sessions_memory_session "github.com/fluffy-bunny/fluffycore/echo/services/sessions/memory_session"
	fluffycore_echo_services_sessions_memory_session_store "github.com/fluffy-bunny/fluffycore/echo/services/sessions/memory_session_store"
	fluffycore_echo_services_sessions_session_factory "github.com/fluffy-bunny/fluffycore/echo/services/sessions/session_factory"
	fluffycore_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/middleware/auth/jwt"
	fluffycore_middleware_claimsprincipal "github.com/fluffy-bunny/fluffycore/middleware/claimsprincipal"
	fluffycore_middleware_correlation "github.com/fluffy-bunny/fluffycore/middleware/correlation"
	fluffycore_middleware_dicontext "github.com/fluffy-bunny/fluffycore/middleware/dicontext"
	fluffycore_middleware_logging "github.com/fluffy-bunny/fluffycore/middleware/logging"
	core_runtime "github.com/fluffy-bunny/fluffycore/runtime"
	fluffycore_services_ddprofiler "github.com/fluffy-bunny/fluffycore/services/ddprofiler"
	services_health "github.com/fluffy-bunny/fluffycore/services/health"
	fluffycore_utils_redact "github.com/fluffy-bunny/fluffycore/utils/redact"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	echo "github.com/labstack/echo/v4"
	async "github.com/reugn/async"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type (
	startup struct {
		fluffycore_contracts_runtime.UnimplementedStartup
		RootContainer di.Container

		configOptions *fluffycore_contracts_runtime.ConfigOptions
		config        *contracts_config.Config

		ddProfiler  fluffycore_contracts_ddprofiler.IDataDogProfiler
		rageStartup fluffycore_contracts_runtime.IStartup
		rageFuture  async.Future[*fluffycore_async.AsyncResponse]
	}
)

func NewStartup() fluffycore_contracts_runtime.IStartup {

	appStartup := &startup{}
	appStartup.rageStartup = rage_runtime.NewStartup(
		rage_runtime.WithConfigureServices(appStartup.MyConfigServices),
	)
	return appStartup
}
func (s *startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	dst, err := fluffycore_utils_redact.CloneAndRedact(s.configOptions.Destination)
	if err != nil {
		panic(err)
	}
	log.Info().Interface("config", dst).Msg("config")
	config := s.configOptions.Destination.(*contracts_config.Config)
	config.DDProfilerConfig.ApplicationEnvironment = config.ApplicationEnvironment
	config.DDProfilerConfig.ServiceName = config.ApplicationName
	config.DDProfilerConfig.Version = example_version.Version()
	di.AddInstance[*fluffycore_contracts_ddprofiler.Config](builder, config.DDProfilerConfig)
	di.AddInstance[*contracts_config.Config](builder, config)

	fluffycore_services_ddprofiler.AddSingletonIProfiler(builder, config.DDProfilerConfig)
	services_health.AddHealthService(builder)
	issuerConfigs := &fluffycore_contracts_middleware_auth_jwt.IssuerConfigs{}
	for idx := range s.config.JWTValidators.Issuers {
		issuerConfigs.IssuerConfigs = append(issuerConfigs.IssuerConfigs,
			&fluffycore_contracts_middleware_auth_jwt.IssuerConfig{
				OAuth2Config: &fluffycore_contracts_middleware_auth_jwt.OAuth2Config{
					Issuer:  s.config.JWTValidators.Issuers[idx],
					JWKSUrl: s.config.JWTValidators.JWKSURLS[idx],
				},
			})
	}
	fluffycore_middleware_auth_jwt.AddValidators(builder, issuerConfigs)
	service_AuthorizationCodeClaimsAugmentor.AddSingletonIAuthorizationCodeClaimsAugmentor(builder)
	services_oidcuser_inmemory.AddSingletonIFluffyCoreUserServiceServer(builder)
	services_user_id_generator.AddSingletonIUserIdGenerator(builder)

}
func (s *startup) MyConfigServices(ctx context.Context, config *rage_contracts_config.Config, builder di.ContainerBuilder) {
	log := zerolog.Ctx(ctx).With().Logger()
	// this extension point is called by the runtime at the end of the startup process
	// it allows you to swap out external services like the user store
	rootContainer := s.RootContainer
	log.Info().Msg("MyConfigServices")

	di.AddInstance[*rage_contracts_config.Config](builder, config)

	// these objects are registered in our main container.
	// we are registering them in the rage container by pulling them first from ours.
	di.AddSingletonFromContainer[proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer](builder, rootContainer)
	di.AddSingletonFromContainer[proto_oidc_user.IFluffyCoreRageUserServiceServer](builder, rootContainer)
	di.AddSingletonFromContainer[proto_external_user.IFluffyCoreUserServiceServer](builder, rootContainer)
	//di.AddSingletonFromContainer[contracts_userservice.ISingletonUserService](builder, rootContainer)
	di.AddSingletonFromContainer[contracts_identity.IUserIdGenerator](builder, rootContainer)
	di.AddSingletonFromContainer[contracts_tokenservice.IAuthorizationCodeClaimsAugmentor](builder, rootContainer)
	di.AddSingletonFromContainer[contracts_events.IEventSink](builder, rootContainer)
	di.AddSingletonFromContainer[contracts_email.IEmailTemplateData](builder, rootContainer)

	services_user_id_generator.AddSingletonIUserIdGenerator(builder)
	services_oidcflowstore.AddSingletonAuthorizationRequestStateStoreServer(builder)
	services_AuthorizationCodeClaimsAugmentor.AddSingletonIAuthorizationCodeClaimsAugmentor(builder)
	services_EventSink.AddSingletonIEventSink(builder)
	services_EmailTemplateData.AddSingletonIEmailTemplateData(builder)
	// Account Handlers
	//--------------------------------------------------------
	services_handlers_account_about.AddScopedIHandler(builder)
	services_handlers_account_callback.AddScopedIHandler(builder)
	services_handlers_account_api_api_user_profile.AddScopedIHandler(builder)
	services_handlers_account_api_api_linked_accounts.AddScopedIHandler(builder)
	services_handlers_account_home.AddScopedIHandler(builder)
	services_handlers_account_api_login.AddScopedIHandler(builder)
	services_handlers_account_logout.AddScopedIHandler(builder)
	services_handlers_account_personal_information.AddScopedIHandler(builder)
	services_handlers_account_passkey_management.AddScopedIHandler(builder)
	services_handlers_account_profile.AddScopedIHandler(builder)
	services_handlers_account_totp_management.AddScopedIHandler(builder)

	guid := xid.New().String()
	if example_version.Version() != "dev-build" {
		guid = example_version.Version()
	}
	managementCacheBustingHTMLConfig := &rage_contracts_config.CacheBustingHTMLConfig{
		Version:    guid,
		FilePath:   "./static/go-app/management/static_output/index_template.html",
		StaticPath: "./static/go-app/management/static_output/",
		EchoPath:   "/management/*",
		RootPath:   "/management/",
		ReplaceParams: []*rage_contracts_config.KeyValuePair{
			{
				Key:   "{basehref}",
				Value: s.config.ManagementAppConfig.BaseHREF,
			},
			{
				Key:   "{title}",
				Value: s.config.ManagementAppConfig.BannerBranding.Title,
			},
			{
				Key:   "{version}",
				Value: guid,
			},
		},

		RoutePatterns: []*rage_contracts_config.RoutePattern{
			{
				Pattern: "/web/app.wasm",
				Handler: func(c echo.Context, filePath string) (bool, error) {
					// Get file info to set Content-Length
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						return false, err
					}
					// Set correct MIME type and Content-Length for WASM files
					c.Response().Header().Set("Content-Type", "application/wasm")
					c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
					return true, c.File(filePath)
				},
			},
			{
				Pattern: "web/app.json",
				Handler: func(c echo.Context, filePath string) (bool, error) {
					jsonB, err := json.Marshal(s.config.ManagementAppConfig)
					if err != nil {
						return false, err
					}

					// Get version from query param
					version := c.QueryParam("v")

					// Replace {version} placeholder
					modifiedContent := strings.ReplaceAll(string(jsonB), "{version}", version)

					// Serve with appropriate content type
					return true, c.JSONBlob(http.StatusOK, []byte(modifiedContent))
				},
			},
		},
	}
	services_handlers_cache_busting_static_html.AddScopedIHandler(builder, managementCacheBustingHTMLConfig)

	oidcloginCacheBustingHTMLConfig := &rage_contracts_config.CacheBustingHTMLConfig{
		Version:    guid,
		FilePath:   "./static/go-app/oidc-login/static_output/index_template.html",
		StaticPath: "./static/go-app/oidc-login/static_output/",
		EchoPath:   "/oidc-login/*",
		RootPath:   "/oidc-login/",
		ReplaceParams: []*rage_contracts_config.KeyValuePair{
			{
				Key:   "{basehref}",
				Value: s.config.OIDCLoginAppConfig.BaseHREF,
			},
			{
				Key:   "{title}",
				Value: s.config.OIDCLoginAppConfig.BannerBranding.Title,
			},
			{
				Key:   "{version}",
				Value: guid,
			},
		},

		RoutePatterns: []*rage_contracts_config.RoutePattern{
			{
				Pattern: "/web/app.wasm",
				Handler: func(c echo.Context, filePath string) (bool, error) {
					// Get file info to set Content-Length
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						return false, err
					}
					// Set correct MIME type and Content-Length for WASM files
					c.Response().Header().Set("Content-Type", "application/wasm")
					c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
					return true, c.File(filePath)
				},
			},
			{
				Pattern: "web/app.json",
				Handler: func(c echo.Context, filePath string) (bool, error) {
					jsonB, err := json.Marshal(s.config.OIDCLoginAppConfig)
					if err != nil {
						return false, err
					}

					// Get version from query param
					version := c.QueryParam("v")

					// Replace {version} placeholder
					modifiedContent := strings.ReplaceAll(string(jsonB), "{version}", version)

					// Serve with appropriate content type
					return true, c.JSONBlob(http.StatusOK, []byte(modifiedContent))
				},
			},
		},
	}
	services_handlers_cache_busting_static_html.AddScopedIHandler(builder, oidcloginCacheBustingHTMLConfig)

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

func (s *startup) SetRootContainer(container di.Container) {
	s.RootContainer = container
}
func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	log := log.With().Caller().Str("method", "GetConfigOptions").Logger()
	// here we load a config file and merge it over the default.
	initialConfigOptions := &fluffycore_contracts_runtime.ConfigOptions{
		Destination: &contracts_config.InitialConfig{},
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	err := core_runtime.LoadConfig(initialConfigOptions)
	if err != nil {
		panic(err)
	}
	err = onLoadMyAppConfig(context.Background(),
		initialConfigOptions.Destination.(*contracts_config.InitialConfig).ConfigFiles.MyAppPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to onLoadMyAppConfig")
		panic(err)
	}
	defaultConfig := &contracts_config.Config{}
	err = json.Unmarshal([]byte(contracts_config.ConfigDefaultJSON), defaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
	}
	log.Info().Interface("defaultConfig", defaultConfig).Msg("config after merge")
	s.config = &contracts_config.Config{}
	s.configOptions = &fluffycore_contracts_runtime.ConfigOptions{
		Destination: s.config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	return s.configOptions
}

func (s *startup) ConfigureServerOpts(ctx context.Context) []grpc.ServerOption {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	var serverOpts []grpc.ServerOption

	// dicontext is responsible of create a scoped context for each request.
	log.Info().Msg("adding ChainUnaryInterceptor: fluffycore_middleware_dicontext.ScopedContextUnaryServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(fluffycore_middleware_dicontext.ScopedContextUnaryServerInterceptor(s.RootContainer)))
	log.Info().Msg("adding ChainStreamInterceptor: fluffycore_middleware_dicontext.ScopedContextStreamServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainStreamInterceptor(fluffycore_middleware_dicontext.ScopedContextStreamServerInterceptor(s.RootContainer)))

	log.Info().Msg("adding ChainUnaryInterceptor: fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor()))
	log.Info().Msg("adding ChainStreamInterceptor: fluffycore_middleware_logging.EnsureContextLoggingStreamServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainStreamInterceptor(fluffycore_middleware_logging.EnsureContextLoggingStreamServerInterceptor()))

	// log correlation and spans
	log.Info().Msg("adding ChainUnaryInterceptor: fluffycore_middleware_correlation.EnsureCorrelationIDUnaryServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(fluffycore_middleware_correlation.EnsureCorrelationIDUnaryServerInterceptor()))

	// auth
	log.Info().Msg("adding ChainUnaryInterceptor: fluffycore_middleware_auth_jwt.UnaryServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(fluffycore_middleware_auth_jwt.UnaryServerInterceptor(s.RootContainer)))

	// Here the gating happens
	grpcEntrypointClaimsMap := example_auth.BuildGrpcEntrypointPermissionsClaimsMap()
	// claims principal
	log.Info().Msg("adding ChainUnaryInterceptor: fluffycore_middleware_claimsprincipal.UnaryServerInterceptor")
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(fluffycore_middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMapWithZeroTrustV2(grpcEntrypointClaimsMap)))

	// last is the recovery middleware
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(opts...)))

	return serverOpts
}

// OnPreServerStartup ...
func (s *startup) OnPreServerStartup(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerStartup").Logger()

	s.rageFuture = fluffycore_async.ExecuteWithPromiseAsync(func(promise async.Promise[*fluffycore_async.AsyncResponse]) {
		var err error
		defer func() {
			promise.Success(&fluffycore_async.AsyncResponse{
				Message: "End Serve - rageStartup",
				Error:   err,
			})
		}()
		log.Info().Msg("rageStartup starting up")
		fluffycore_cobracore_cmd.Execute(s.rageStartup)
	})

	s.ddProfiler = di.Get[fluffycore_contracts_ddprofiler.IDataDogProfiler](s.RootContainer)
	s.ddProfiler.Start(ctx)
	return nil
}

// OnPreServerShutdown ...
func (s *startup) OnPreServerShutdown(ctx context.Context) {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerShutdown").Logger()

	log.Info().Msg("Stopping Datadog Tracer and Profiler")
	s.ddProfiler.Stop(ctx)
	log.Info().Msg("Datadog Tracer and Profiler stopped")

}
