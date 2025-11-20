package logout

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.LogoutPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LogoutGetRequest struct {
	RedirectURL string `param:"redirect_url" query:"redirect_url" form:"redirect_url" json:"redirect_url" xml:"redirect_url"`
}

func (s *service) validateLoginGetRequest(model *LogoutGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.RedirectURL) {
		model.RedirectURL = "/"
	}
	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LogoutGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Debug().Interface("model", model).Msg("model")
	err := s.validateLoginGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateLoginGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}
	s.wellknownCookies.DeleteAuthCookie(c)

	return s.Render(c, http.StatusOK, "oidc/logout/index",
		map[string]interface{}{
			"url": model.RedirectURL,
		})
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       redirect_url            		query     string  true  "redirect url"
// @Success 200 {object} string
// @Router /logout [get,post]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)

	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
