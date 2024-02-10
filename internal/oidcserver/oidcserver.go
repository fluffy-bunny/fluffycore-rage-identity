package oidcserver

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/config"
	services "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services"
	services_handlers_about "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/about"
	services_handlers_authorization_endpoint "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/authorization_endpoint"
	services_handlers_discovery_endpoint "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/discovery_endpoint"
	services_handlers_error "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/error"
	services_handlers_externalidp "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/externalidp"
	services_handlers_healthz "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/healthz"
	services_handlers_home "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/home"
	services_handlers_jwks_endpoint "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/jwks_endpoint"
	services_handlers_login "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/login"
	services_handlers_oauth2_callback "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/oauth2/callback"
	services_handlers_swagger "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/swagger"
	services_handlers_token_endpoint "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/token_endpoint"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	contracts_startup "github.com/fluffy-bunny/fluffycore/echo/contracts/startup"
	services_startup "github.com/fluffy-bunny/fluffycore/echo/services/startup"
	echo "github.com/labstack/echo/v4"
	log "github.com/rs/zerolog/log"
)

type (
	startup struct {
		services_startup.StartupBase
		config *contracts_config.Config
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
	log.Info().Msg("PreStartHook")
	return nil
}
func (s *startup) PostBuildHook(container di.Container) error {
	log.Info().Msg("PostBuildHook")
	return nil
}
func (s *startup) PreShutdownHook(echo *echo.Echo) error {
	log.Info().Msg("PreShutdownHook")
	return nil
}
func (s *startup) addAppHandlers(builder di.ContainerBuilder) {
	// add your handlers here
	services_handlers_healthz.AddScopedIHandler(builder)
	services_handlers_home.AddScopedIHandler(builder)
	services_handlers_about.AddScopedIHandler(builder)
	services_handlers_error.AddScopedIHandler(builder)
	services_handlers_login.AddScopedIHandler(builder)
	services_handlers_externalidp.AddScopedIHandler(builder)
	services_handlers_swagger.AddScopedIHandler(builder)
	services_handlers_discovery_endpoint.AddScopedIHandler(builder)
	services_handlers_jwks_endpoint.AddScopedIHandler(builder)
	services_handlers_authorization_endpoint.AddScopedIHandler(builder)
	services_handlers_token_endpoint.AddScopedIHandler(builder)
	services_handlers_oauth2_callback.AddScopedIHandler(builder)
}
func (s *startup) RegisterStaticRoutes(e *echo.Echo) error {
	// i.e. e.Static("/css", "./css")
	e.Static("/static", "./static")
	return nil
}
