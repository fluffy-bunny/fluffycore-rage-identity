package discovery_endpoint

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/util"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		config   *contracts_config.Config
		someUtil contracts_util.ISomeUtil
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(config *contracts_config.Config, someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		config:   config,
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
		wellknown_echo.WellKnownOpenIDCOnfiguationPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /.well-known/openid-configuration [get]
func (s *service) Do(c echo.Context) error {
	rootPath := s.config.OIDCConfig.BaseUrl

	discovery := models.DiscoveryDocument{
		Issuer:                rootPath,
		TokenEndpoint:         rootPath + wellknown_echo.OAuth2TokenEndpointPath,
		JwksURI:               rootPath + wellknown_echo.WellKnownJWKS,
		UserinfoEndpoint:      rootPath + wellknown_echo.UserInfoPath,
		AuthorizationEndpoint: rootPath + wellknown_echo.OIDCAuthorizationEndpointPath,
		//	RevocationEndpoint:    rootPath + wellknown.OAuth2RevokePath,
		//	IntrospectionEndpoint: rootPath + wellknown.OAuth2IntrospectPath,
		GrantTypesSupported: []string{
			models.OAUTH2GrantType_AuthorizationCode,
		},
		ScopesSupported: []string{
			"openid",
			"email",
			"profile",
		},
		ResponseTypesSupported: []string{
			"code",
			"token",
			"id_token",
			"code token",
			"code id_token",
			"token id_token",
			"code token id_token",
			"none",
		},
		IDTokenSigningAlgValuesSupported: []string{
			"ES256",
		},
	}
	return c.JSONPretty(http.StatusOK, discovery, "  ")
}
