package shell

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx/components"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.HTMXOIDCLoginPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	rc.CacheBustVersion = s.config.CacheBustVersion

	// Determine initial page based on session landing page directive
	initialPage := rc.Paths.HTMXHome
	session, err := s.OIDCSession().GetSession()
	if err == nil {
		landingPageI, err := session.Get("landingPage")
		if err == nil && landingPageI != nil {
			if lp, ok := landingPageI.(*models_api_manifest.LandingPage); ok && lp != nil {
				switch lp.Page {
				case models_api_manifest.PageVerifyCode:
					initialPage = rc.Paths.HTMXVerifyCode
				case models_api_manifest.PageKeepSignedIn:
					initialPage = rc.Paths.HTMXKeepSignedIn
				case models_api_manifest.PagePasswordEntry:
					initialPage = rc.Paths.HTMXPassword
				}
				if initialPage != rc.Paths.HTMXHome {
					log.Info().Str("landingPage", string(lp.Page)).Str("initialPage", initialPage).Msg("shell: routing to session landing page")
				}
				// Consume the landing page so it doesn't fire again on refresh
				session.Set("landingPage", nil)
				session.Save()
			}
		}
	}

	return components.RenderNode(c, http.StatusOK, components.ShellPage(components.ShellData{
		RenderContext: rc,
		BrandTitle:    "MAPPED Identity",
		InitialPage:   initialPage,
	}))
}
