package authorization_endpoint

/*
reference: https://developers.onelogin.com/openid-connect/api/authorization-code
*/
import (
	"fmt"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/util"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/client"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		someUtil            contracts_util.ISomeUtil
		scopedMemoryCache   fluffycore_contracts_common.IScopedMemoryCache
		oidcFlowStore       contracts_eko_gocache.IOIDCFlowStore
		clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	oidcFlowStore contracts_eko_gocache.IOIDCFlowStore,
	someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil:            someUtil,
		scopedMemoryCache:   scopedMemoryCache,
		oidcFlowStore:       oidcFlowStore,
		clientServiceServer: clientServiceServer,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OIDCAuthorizationEndpointPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		//clientauthorization.AuthenticateOAuth2Client(),
	}
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       client_id    			query     string  true  "client_id requested"
// @Param       response_type    		query     string  true  "response_type requested"
// @Param       scope            		query     string  true  "scope requested" default("openid profile email")
// @Param       state            		query     string  true  "state requested"
// @Param       redirect_uri     		query     string  true  "redirect_uri requested"
// @Param       audience     	 		query     string  false  "audience requested"
// @Param       code_challenge   		query     string  false  "PKCE challenge code"
// @Param       code_challenge_method 	query     string  false  "PKCE challenge method" default("S256")
// @Success 200 {object} string
// @Router /oidc/v1/auth [get]
func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &models.AuthorizationRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	// TODO: validate the request
	// does the client have the permissions to do this?
	code := xid.New().String()
	// store the model in the cache.  Redis in production.
	authorizationFinal := &models.AuthorizationFinal{
		Request: model,
	}

	err := s.oidcFlowStore.StoreAuthorizationFinal(ctx, code, authorizationFinal)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}

	mm, err := s.oidcFlowStore.GetAuthorizationFinal(ctx, code)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("mm", mm).Msg("mm")
	// redirect to the server Auth login pages.
	//
	finalOIDCPath := fmt.Sprintf("%s?code=%s", wellknown_echo.OIDCLoginPath, code)
	redirectPath := fmt.Sprintf("%s?redirect_uri=%s", wellknown_echo.LoginPath, finalOIDCPath)
	return c.Redirect(http.StatusFound, redirectPath)
}
