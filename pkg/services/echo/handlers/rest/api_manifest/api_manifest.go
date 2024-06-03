package api_manifest

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	models_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	container di.Container,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_Manifest,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API Manifest godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /api/manifest [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()

	idps, err := s.GetIDPs(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	response := &models_manifest.Manifest{}
	for _, idp := range idps {
		if idp.Enabled && !idp.Hidden {
			response.SocialIdps = append(response.SocialIdps, models_manifest.IDP{
				Slug: idp.Slug,
			})
		}
	}
	return c.JSONPretty(http.StatusOK, response, "  ")
}
