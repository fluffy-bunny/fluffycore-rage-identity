package startover

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.HTMXStartOverPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Clean up cookies
	s.wellknownCookies.DeleteVerificationCodeCookie(c)
	s.wellknownCookies.DeletePasswordResetCookie(c)
	s.wellknownCookies.DeleteAuthCompletedCookie(c)

	// Clear session landing page
	session, err := s.oidcSession.GetSession()
	if err == nil {
		session.Set("landing_page", nil)
		session.Save()
	} else {
		log.Debug().Err(err).Msg("Failed to get session during start-over")
	}

	return s.Render(c, http.StatusOK, "oidc/htmx/_partials/home", map[string]interface{}{
		"errors": []string{},
		"email":  "",
	})
}
