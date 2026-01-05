package error

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
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

func (s *service) Ctor(container di.Container, config *contracts_config.Config) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.ErrorPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// HealthCheck godoc
// @Summary get the error page.
// @Description get the error page.
// @Tags root
// @Accept */*
// @Produce json
// @Param error query string false "Error code"
// @Success 200 {object} string
// @Router /error [get]
func (s *service) Do(c echo.Context) error {
	errorCode := c.QueryParam("error")
	errorMessage := "An error occurred"

	// Map error codes to user-friendly messages
	switch errorCode {
	case "invalid_idp_hint":
		errorMessage = "Invalid identity provider specified. The requested IDP does not exist or is not enabled."
	case "invalid_root_candidate":
		errorMessage = "Invalid user candidate specified. The requested user does not exist."
	case "store_failed":
		errorMessage = "Failed to store authorization request. Please try again."
	case "state_not_found":
		errorMessage = "Authorization state not found. Please try again."
	default:
		if errorCode != "" {
			errorMessage = "Error: " + errorCode
		}
	}

	return s.Render(c, http.StatusOK, "oidc/error/index", map[string]interface{}{
		"error":   errorCode,
		"message": errorMessage,
	})
}
