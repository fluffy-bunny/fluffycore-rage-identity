package api_isauthorized

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	api "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),
	}, nil
}

// AddScopedIHandler registers the *service.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_IsAuthorized,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type IsAuthorizedResponse struct{}

// API Manifest godoc
// @Summary is authorized
// @Description is authorized
// @Tags root
// @Produce json
// @Success 200 {object} api.AuthorizedResponse
// @Failure 401 {object} api.UnauthorizedResponse
// @Router /api/is-authorized [get]
func (s *service) Do(c echo.Context) error {

	response := api.AuthorizedResponse{}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
