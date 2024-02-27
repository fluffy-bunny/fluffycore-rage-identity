package runtime

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	internal_auth "github.com/fluffy-bunny/fluffycore-rage-identity/internal/auth"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/config"
	oidcserver "github.com/fluffy-bunny/fluffycore-rage-identity/internal/oidcserver"
	services "github.com/fluffy-bunny/fluffycore-rage-identity/internal/services"
	services_health "github.com/fluffy-bunny/fluffycore-rage-identity/internal/services/health"
	internal_types "github.com/fluffy-bunny/fluffycore-rage-identity/internal/types"
	"github.com/fluffy-bunny/fluffycore-rage-identity/internal/utils"
	internal_version "github.com/fluffy-bunny/fluffycore-rage-identity/internal/version"
	fluffycore_async "github.com/fluffy-bunny/fluffycore/async"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
	fluffycore_contracts_middleware "github.com/fluffy-bunny/fluffycore/contracts/middleware"
	fluffycore_contracts_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/contracts/middleware/auth/jwt"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	core_echo_runtime "github.com/fluffy-bunny/fluffycore/echo/runtime"
	fluffycore_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/middleware/auth/jwt"
	fluffycore_middleware_claimsprincipal "github.com/fluffy-bunny/fluffycore/middleware/claimsprincipal"
	fluffycore_middleware_correlation "github.com/fluffy-bunny/fluffycore/middleware/correlation"
	fluffycore_middleware_dicontext "github.com/fluffy-bunny/fluffycore/middleware/dicontext"
	fluffycore_middleware_logging "github.com/fluffy-bunny/fluffycore/middleware/logging"
	core_runtime "github.com/fluffy-bunny/fluffycore/runtime"
	fluffycore_services_ddprofiler "github.com/fluffy-bunny/fluffycore/services/ddprofiler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	fluffycore_utils_redact "github.com/fluffy-bunny/fluffycore/utils/redact"
	status "github.com/gogo/status"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	async "github.com/reugn/async"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	startup struct {
		fluffycore_contracts_runtime.UnimplementedStartup
		RootContainer di.Container

		configOptions *fluffycore_contracts_runtime.ConfigOptions
		config        *contracts_config.Config

		ddProfiler        fluffycore_contracts_ddprofiler.IDataDogProfiler
		oidcserverFuture  async.Future[fluffycore_async.AsyncResponse]
		oidcserverRuntime *core_echo_runtime.Runtime
		ext               internal_types.ConfigureServices
	}
)
type WithOption func(startup *startup)

func WithConfigureServices(ext internal_types.ConfigureServices) WithOption {
	return func(startup *startup) {
		startup.ext = ext
	}
}
func NewStartup(options ...WithOption) fluffycore_contracts_runtime.IStartup {
	var s = &startup{}
	for _, option := range options {
		option(s)
	}
	return s
}
func (s *startup) SetRootContainer(container di.Container) {
	s.RootContainer = container

}
func onLoadRageConfig(ctx context.Context, ragePath string) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadIDPs").Logger()
	fileContent, err := os.ReadFile(ragePath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read IDPsPath - may not be a problem if idps are comming from a DB")
		return nil
	}
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")
	overlay := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(fixedFileContent)).Decode(&overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ragePath")
		return err
	}
	src := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(string(contracts_config.ConfigDefaultJSON))).Decode(&src)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
		return err
	}
	err = utils.ReplaceMergeMap(overlay, src)
	if err != nil {
		log.Error().Err(err).Msg("failed to ReplaceMergeMap")
		return err
	}
	bb, err := json.Marshal(overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal overlay")
		return err
	}
	contracts_config.ConfigDefaultJSON = bb

	return nil

}

func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	initialConfigOptions := &fluffycore_contracts_runtime.ConfigOptions{
		Destination: &contracts_config.InitialConfig{},
		RootConfig:  contracts_config.ConfigDefaultJSON,
		EnvPrefix:   "RAGE",
	}
	err := core_runtime.LoadConfig(initialConfigOptions)
	if err != nil {
		panic(err)
	}
	err = onLoadRageConfig(context.Background(), initialConfigOptions.Destination.(*contracts_config.InitialConfig).ConfigFiles.RagePath)
	if err != nil {
		panic(err)
	}
	s.config = &contracts_config.Config{}
	s.configOptions = &fluffycore_contracts_runtime.ConfigOptions{
		Destination: s.config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
		EnvPrefix:   "RAGE",
	}
	return s.configOptions
}
func (s *startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	_, err := fluffycore_utils_redact.CloneAndRedact(s.configOptions.Destination)
	if err != nil {
		panic(err)
	}
	config := s.configOptions.Destination.(*contracts_config.Config)
	config.DDProfilerConfig.ApplicationEnvironment = config.ApplicationEnvironment
	config.DDProfilerConfig.ServiceName = config.ApplicationName
	config.DDProfilerConfig.Version = internal_version.Version()
	di.AddInstance[*fluffycore_contracts_ddprofiler.Config](builder, config.DDProfilerConfig)

	services.ConfigureServices(ctx, config, builder)
	fluffycore_services_ddprofiler.AddSingletonIProfiler(builder)
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
	log.Info().Interface("config", config).Msg("config")

}
func (s *startup) Configure(ctx context.Context, rootContainer di.Container, unaryServerInterceptorBuilder fluffycore_contracts_middleware.IUnaryServerInterceptorBuilder, streamServerInterceptorBuilder fluffycore_contracts_middleware.IStreamServerInterceptorBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()

	// puts a zerlog logger into the request context
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor())
	log.Info().Msg("adding streamServerInterceptorBuilder: fluffycore_middleware_logging.EnsureContextLoggingStreamServerInterceptor")
	streamServerInterceptorBuilder.Use(fluffycore_middleware_logging.EnsureContextLoggingStreamServerInterceptor())

	// log correlation and spans
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_correlation.EnsureCorrelationIDUnaryServerInterceptor())
	// dicontext is responsible of create a scoped context for each request.
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_dicontext.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_dicontext.UnaryServerInterceptor(rootContainer))
	log.Info().Msg("adding streamServerInterceptorBuilder: fluffycore_middleware_dicontext.StreamServerInterceptor")
	streamServerInterceptorBuilder.Use(fluffycore_middleware_dicontext.StreamServerInterceptor(rootContainer))

	// auth
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_auth_jwt.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_auth_jwt.UnaryServerInterceptor(rootContainer))

	// Here the gating happens
	grpcEntrypointClaimsMap := internal_auth.BuildGrpcEntrypointPermissionsClaimsMap()
	// claims principal
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_claimsprincipal.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMapWithZeroTrustV2(grpcEntrypointClaimsMap))

	// last is the recovery middleware
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	unaryServerInterceptorBuilder.Use(grpc_recovery.UnaryServerInterceptor(opts...))

}

// OnPreServerStartup ...
func (s *startup) OnPreServerStartup(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerStartup").Logger()

	s.oidcserverRuntime = core_echo_runtime.New(oidcserver.NewStartup(
		oidcserver.WithConfigureServices(s.ext),
	))
	s.oidcserverFuture = fluffycore_async.ExecuteWithPromiseAsync(func(promise async.Promise[fluffycore_async.AsyncResponse]) {
		var err error
		defer func() {
			promise.Success(&fluffycore_async.AsyncResponse{
				Message: "End Serve - echoServer",
				Error:   err,
			})
		}()
		log.Info().Msg("echoServer starting up")
		err = s.oidcserverRuntime.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to start server")
		}
	})

	s.ddProfiler = di.Get[fluffycore_contracts_ddprofiler.IDataDogProfiler](s.RootContainer)
	s.ddProfiler.Start(ctx)
	return nil
}

// OnPreServerShutdown ...
func (s *startup) OnPreServerShutdown(ctx context.Context) {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerShutdown").Logger()

	log.Info().Msg("oidcserverRuntime shutting down")
	s.oidcserverRuntime.Stop()
	s.oidcserverFuture.Join()
	log.Info().Msg("oidcserverRuntime shutdown complete")

	log.Info().Msg("Stopping Datadog Tracer and Profiler")
	s.ddProfiler.Stop(ctx)
	log.Info().Msg("Datadog Tracer and Profiler stopped")
}
