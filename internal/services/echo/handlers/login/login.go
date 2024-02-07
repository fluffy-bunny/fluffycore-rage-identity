package login

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		someUtil contracts_util.ISomeUtil
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil: someUtil,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.LoginPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LoginRequest struct {
	Code string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /login [get]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	c.SetCookie(&http.Cookie{
		Name:   "_code",
		Value:  model.Code,
		Path:   "/",
		Secure: true,
	})
	return c.Render(http.StatusOK, "views/login/index", map[string]interface{}{})
}
