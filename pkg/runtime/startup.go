package runtime

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	oidcserver "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/oidcserver"
	services "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services"
	services_health "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/health"
	pkg_types "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/types"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidcuser "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_async "github.com/fluffy-bunny/fluffycore/async"
	fluffycore_contracts_GRPCClientFactory "github.com/fluffy-bunny/fluffycore/contracts/GRPCClientFactory"
	contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_contracts_middleware "github.com/fluffy-bunny/fluffycore/contracts/middleware"
	fluffycore_contracts_otel "github.com/fluffy-bunny/fluffycore/contracts/otel"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	core_echo_runtime "github.com/fluffy-bunny/fluffycore/echo/runtime"
	fluffycore_middleware_correlation "github.com/fluffy-bunny/fluffycore/middleware/correlation"
	fluffycore_middleware_dicontext "github.com/fluffy-bunny/fluffycore/middleware/dicontext"
	fluffycore_middleware_logging "github.com/fluffy-bunny/fluffycore/middleware/logging"
	core_runtime "github.com/fluffy-bunny/fluffycore/runtime"
	fluffycore_runtime_otel "github.com/fluffy-bunny/fluffycore/runtime/otel"
	fluffycore_services_GRPCClientFactory "github.com/fluffy-bunny/fluffycore/services/GRPCClientFactory"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	fluffycore_utils_redact "github.com/fluffy-bunny/fluffycore/utils/redact"
	status "github.com/gogo/status"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	async "github.com/reugn/async"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

type (
	startup struct {
		*fluffycore_runtime_otel.FluffyCoreOTELStartup

		configOptions *fluffycore_contracts_runtime.ConfigOptions
		config        *contracts_config.Config

		oidcserverFuture  async.Future[*fluffycore_async.AsyncResponse]
		oidcserverRuntime *core_echo_runtime.Runtime

		ext pkg_types.ConfigureServices
	}
)
type WithOption func(startup *startup)

func WithConfigureServices(ext pkg_types.ConfigureServices) WithOption {
	return func(startup *startup) {
		startup.ext = ext
	}
}
func emptyEntrypointConfigs() map[string]contracts_common.IEntryPointConfig {
	return map[string]contracts_common.IEntryPointConfig{}
}
func NewStartup(options ...WithOption) fluffycore_contracts_runtime.IStartup {
	var s = &startup{
		FluffyCoreOTELStartup: fluffycore_runtime_otel.NewFluffyCoreOTELStartup(&fluffycore_runtime_otel.FluffyCoreOTELStartupConfig{
			FuncAuthGetEntryPointConfigs: emptyEntrypointConfigs,
		}),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func onLoadRageConfig(ctx context.Context, ragePath string) error {
	log := zerolog.Ctx(ctx).With().Str("method", "onLoadRageConfig").Logger()
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
	// need to set the OTEL Config in the base startup
	if config.OTELConfig == nil {
		config.OTELConfig = &fluffycore_contracts_otel.OTELConfig{}
	}
	config.OTELConfig.ServiceName = config.ApplicationName
	s.FluffyCoreOTELStartup.SetConfig(config.OTELConfig)
	// add grpcclient factory that is config aware.  Will make sure that you get one that has otel tracing if enabled.
	fluffycore_contracts_GRPCClientFactory.AddGRPCClientConfig(builder,
		&fluffycore_contracts_GRPCClientFactory.GRPCClientConfig{
			OTELTracingEnabled: config.OTELConfig.TracingConfig.Enabled,
		})
	fluffycore_services_GRPCClientFactory.AddSingletonIGRPCClientFactory(builder)

	wellknown_echo.OAuth2CallbackPath = config.OIDCConfig.OAuth2CallbackPath

	services.ConfigureServices(ctx, config, builder)
	services_health.AddHealthService(builder)

	log.Info().Interface("config", config).Msg("config")

}
func (s *startup) ConfigureOld(ctx context.Context, rootContainer di.Container, unaryServerInterceptorBuilder fluffycore_contracts_middleware.IUnaryServerInterceptorBuilder, streamServerInterceptorBuilder fluffycore_contracts_middleware.IStreamServerInterceptorBuilder) {
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

	// last is the recovery middleware
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	unaryServerInterceptorBuilder.Use(grpc_recovery.UnaryServerInterceptor(opts...))

}
func (s *startup) OnLoadSeedUsers(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnLoadSeedUsers").Logger()
	config := s.configOptions.Destination.(*contracts_config.Config)

	fileContent, err := os.ReadFile(config.ConfigFiles.SeedUsersPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read OIDCClientPath - may not be a problem if clients are comming from a DB")
		return nil
	}
	rageUsers := &proto_oidc_models.RageUsers{}

	err = protojson.Unmarshal(fileContent, rageUsers)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	log.Info().Interface("rageUsers", rageUsers).Msg("rageUsers")

	rageUserService := di.Get[proto_oidcuser.IFluffyCoreRageUserServiceServer](s.RootContainer)
	for _, rageUser := range rageUsers.Users {
		_, err := rageUserService.CreateRageUser(ctx, &proto_oidcuser.CreateRageUserRequest{
			User: rageUser,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to CreateRageUser")
			return err
		}
	}
	return nil
}

// OnPreServerStartup ...
func (s *startup) OnPreServerStartup(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerStartup").Logger()

	err := s.FluffyCoreOTELStartup.OnPreServerStartup(ctx)
	if err != nil {
		return err
	}
	//s.OnLoadSeedUsers(ctx)
	s.oidcserverRuntime = core_echo_runtime.New(oidcserver.NewStartup(
		oidcserver.WithConfigureServices(s.ext),
	))
	s.oidcserverFuture = fluffycore_async.ExecuteWithPromiseAsync(func(promise async.Promise[*fluffycore_async.AsyncResponse]) {
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

	return nil
}

// OnPreServerShutdown ...
func (s *startup) OnPreServerShutdown(ctx context.Context) {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerShutdown").Logger()

	log.Info().Msg("oidcserverRuntime shutting down")
	s.oidcserverRuntime.Stop()
	s.oidcserverFuture.Join()
	log.Info().Msg("oidcserverRuntime shutdown complete")

	log.Info().Msg("FluffyCoreOTELStartup stopped")
	s.FluffyCoreOTELStartup.OnPreServerShutdown(ctx)

}
