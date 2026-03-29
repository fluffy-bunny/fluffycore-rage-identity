package shell

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	management_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	example_version "github.com/fluffy-bunny/fluffycore-rage-identity/example/version"
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
		wellknown_echo.HTMXManagementPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	rc.CacheBustVersion = s.config.CacheBustVersion
	rc.AppConfig = s.appConfig
	rc.AppVersion = example_version.Version()

	// Get user info from claims principal
	cp := s.ClaimsPrincipal()
	if cp != nil {
		emailClaims := cp.GetClaimsByType("email")
		if len(emailClaims) > 0 {
			rc.UserEmail = emailClaims[0].Value
		}
		nameClaims := cp.GetClaimsByType("name")
		if len(nameClaims) > 0 {
			rc.UserName = nameClaims[0].Value
		}
		subClaims := cp.GetClaimsByType("sub")
		if len(subClaims) > 0 {
			rc.UserSubject = subClaims[0].Value
		}
	}

	// Support deep linking: if ?redirect= is present, load that page initially
	if redirect := c.QueryParam("redirect"); redirect != "" {
		// Only allow paths under /management/ to prevent open redirect
		if strings.HasPrefix(redirect, "/management/") {
			rc.DeepLinkPath = redirect
		}
	}

	return components.RenderNode(c, http.StatusOK, components.ShellPage(rc))
}
