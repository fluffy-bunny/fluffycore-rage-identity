package jwks_endpoint

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	fluffycore_contracts_jwtminter "github.com/fluffy-bunny/fluffycore/contracts/jwtminter"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		jwtMinter fluffycore_contracts_jwtminter.IJWTMinter
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(jwtMinter fluffycore_contracts_jwtminter.IJWTMinter) (*service, error) {
	return &service{
		jwtMinter: jwtMinter,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.WellKnownJWKS,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// HealthCheck godoc
// @Summary get the public keys of the server.
// @Description get the public keys of the server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /.well-known/jwks [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	publickKeys, err := s.jwtMinter.PublicKeys(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	return c.JSONPretty(http.StatusOK, publickKeys, "  ")
}
