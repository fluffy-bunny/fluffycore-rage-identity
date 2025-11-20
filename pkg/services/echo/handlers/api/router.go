package api

import (
	"bytes"
	"io"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models_api "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		config           *contracts_config.Config
		oidcSession      contracts_oidc_session.IOIDCSession
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	oidcSession contracts_oidc_session.IOIDCSession,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {

	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		oidcSession:      oidcSession,
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			// everthing is a POST
			contracts_handler.POST,
		},
		wellknown_echo.APIPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// make a copy of the original reqest so that next handler can call c.Bind
	b, _ := io.ReadAll(c.Request().Body)
	// we will put it back so that the handler can read it as well
	c.Request().Body = io.NopCloser(bytes.NewReader(b))
	request := &models_api.BaseRequest{}
	err := c.Bind(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to bind")
		return c.JSON(400, "failed to bind")
	}
	// Restore the io.ReadCloser to its original state
	c.Request().Body = io.NopCloser(bytes.NewReader(b))

	switch request.RequestType {
	case "InitialPageRequest":
		return s.DoInitialPageRequest(c)
	}
	log.Debug().Msg("api")
	return c.JSON(500, &models_api.BaseResponse{
		Errors: []string{"unknown request type"},
	})
}
func (s *service) DoInitialPageRequest(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	request := &models_api.InitialPageRequest{}
	err := c.Bind(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to bind")
		return c.JSON(400, "failed to bind")
	}
	response := &models_api.InitialPageResponse{}

	errors := []string{}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		errors = append(errors, err.Error())
		response.Errors = errors
		return c.JSON(500, response)
	}

	for _, idp := range idps {
		response.IDPs = append(response.IDPs, models_api.IDP{
			Name: idp.Name,
			Slug: idp.Slug,
		})
	}
	return c.JSON(200, response)
}
