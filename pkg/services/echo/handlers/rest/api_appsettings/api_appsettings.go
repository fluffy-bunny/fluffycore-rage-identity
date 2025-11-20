package api_appsettings

import (
	"encoding/json"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service[T any] struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config

		AppSettings *T
	}
)

var stemService = (*service[any])(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

// AddScopedIHandler registers the *service as a scope.
func AddScopedIHandler[T any](builder di.ContainerBuilder, ptr *T) {
	contracts_handler.AddScopedIHandleWithMetadata[*service[T]](builder,
		func(container di.Container, config *contracts_config.Config) (*service[T], error) {
			return &service[T]{
				BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
				AppSettings: ptr,
				config:      config,
			}, nil
		},
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_AppSettings,
	)

}

func (s *service[T]) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API Manifest godoc
// @Summary get the app settings
// @Description This is the configuration of the server..
// @Tags root
// @Produce json
// @Success 200 {object} any
// @Router /api/appsettings [get]
func (s *service[T]) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	jsonData, err := json.Marshal(s.AppSettings)
	if err != nil {
		log.Error().Err(err).Msg("json.Marshal failed")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	return c.JSONBlob(http.StatusOK, jsonData)
}
