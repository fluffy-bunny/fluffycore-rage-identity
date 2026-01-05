package api_clear_sso

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	models_api_preferences "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_preferences"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
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
	config *contracts_config.Config,
	container di.Container,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_ClearSSO,
	)
}

const (
	InternalError_ClearSSO_001 = "rg-clearsso-001"
	InternalError_ClearSSO_002 = "rg-clearsso-002"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API ClearSSO godoc
// @Summary Clear SSO session cookie.
// @Description Clears the SSO cookie to end the single sign-on session
// @Tags preferences
// @Accept json
// @Produce json
// @Success 200 {object} models_api_preferences.ClearSSOResponse
// @Failure 500 {object} models_api_preferences.ErrorResponse
// @Router /api/clear-sso [post]
func (s *service) Do(c echo.Context) error {

	// Clear the SSO cookie
	s.WellknownCookies().DeleteSSOCookie(c)
	return c.JSON(http.StatusOK, models_api_preferences.ClearSSOResponse{
		Success: true,
		Message: "SSO session cleared successfully",
	})
}
