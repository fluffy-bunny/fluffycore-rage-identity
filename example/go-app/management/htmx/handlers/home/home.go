package home

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	management_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config    *contracts_config.Config
		appConfig *management_config.AppConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	appConfig *management_config.AppConfig,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		appConfig:   appConfig,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.HTMXManagementHomePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	// Non-HTMX GET requests (e.g. browser refresh) need the full shell page
	if !components.IsHTMXRequest(c) {
		return c.Redirect(http.StatusFound, wellknown_echo.HTMXManagementPath+"?redirect="+c.Request().URL.Path)
	}
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	rc.AppConfig = s.appConfig
	return components.RenderNode(c, http.StatusOK, components.HomePage(rc))
}
