package token_endpoint

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
	clientauthorization "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/middleware/clientauthorization"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	oauth2 "github.com/go-oauth2/oauth2/v4"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		scopedMemoryCache              contracts_cache.IScopedMemoryCache
		authorizationRequestStateStore proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer
		tokenService                   contracts_tokenservice.ITokenService
		claimsaugmentor                contracts_tokenservice.IAuthorizationCodeClaimsAugmentor
		eventSink                      contracts_events.IEventSink
		idpServiceServer               proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	scopedMemoryCache contracts_cache.IScopedMemoryCache,
	authorizationRequestStateStore proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer,
	tokenService contracts_tokenservice.ITokenService,
	claimsaugmentor contracts_tokenservice.IAuthorizationCodeClaimsAugmentor,
	eventSink contracts_events.IEventSink,
	idpServiceServer proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer,
) (*service, error) {
	return &service{
		scopedMemoryCache:              scopedMemoryCache,
		authorizationRequestStateStore: authorizationRequestStateStore,
		tokenService:                   tokenService,
		claimsaugmentor:                claimsaugmentor,
		eventSink:                      eventSink,
		idpServiceServer:               idpServiceServer,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.OAuth2TokenEndpointPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		// this will pull the client_id and client_secret from the Authorization header or form and put the client into the scoped memory cache
		clientauthorization.AuthenticateOAuth2Client(),
	}
}

type TokenEndpointRequest struct {
	GrantType string `param:"grant_type" query:"grant_type" form:"grant_type" json:"grant_type" xml:"grant_type"`
}

// HealthCheck godoc
// @Summary OAuth2 token endpoint.
// @Description OAuth2 token endpoint.
// @Tags root
// @Accept */*
// @Produce json
// @Security BasicAuth
// @Param       response_type    query     string  true  "response_type requested"
// @Param       scope            query     string  true  "scope requested" default("openid profile email")
// @Param       state            query     string  true  "state requested"
// @Param       redirect_uri     query     string  true  "redirect_uri requested"
// @Success 200 {object} string
// @Router /token [post]
func (s *service) Do(c echo.Context) error {
	tokenEndpointRequest := &TokenEndpointRequest{}
	if err := c.Bind(tokenEndpointRequest); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	switch tokenEndpointRequest.GrantType {
	case string(oauth2.AuthorizationCode):
		return s.handleAuthorizationCode(c)
	}
	return c.String(http.StatusBadRequest, "Bad Request")
}
