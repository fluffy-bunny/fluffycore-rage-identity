package discovery

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	mocks_oauth2 "github.com/fluffy-bunny/fluffycore/mocks/oauth2"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		someUtil contracts_util.ISomeUtil
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

	signingKey, _ = mocks_oauth2.LoadSigningKey()
	jwksKeys = &mocks_oauth2.JWKSKeys{
		Keys: []mocks_oauth2.PublicJwk{
			signingKey.PublicJwk,
		},
	}
}

var signingKey *mocks_oauth2.SigningKey
var jwksKeys *mocks_oauth2.JWKSKeys

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
	return c.JSONPretty(http.StatusOK, jwksKeys, "  ")
}
