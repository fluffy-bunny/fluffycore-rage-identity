package authorization_endpoint

import (
	"net/http"
	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	models "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/models"
	clientauthorization "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/middleware/clientauthorization"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		someUtil          contracts_util.ISomeUtil
		scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache
		oidcFlowCache     contracts_eko_gocache.IOIDCFlowCache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	oidcFlowCache contracts_eko_gocache.IOIDCFlowCache,
	someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil:          someUtil,
		scopedMemoryCache: scopedMemoryCache,
		oidcFlowCache:     oidcFlowCache,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OAuth2AuthorizationEndpointPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		clientauthorization.AuthenticateOAuth2Client(),
	}
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Security BasicAuth
// @Param       response_type    		query     string  true  "response_type requested"
// @Param       scope            		query     string  true  "scope requested" default("openid profile email")
// @Param       state            		query     string  true  "state requested"
// @Param       redirect_uri     		query     string  true  "redirect_uri requested"
// @Param       audience     	 		query     string  false  "audience requested"
// @Param       code_challenge   		query     string  false  "PKCE challenge code"
// @Param       code_challenge_method 	query     string  false  "PKCE challenge method" default("S256")
// @Success 200 {object} string
// @Router /o/oauth2/v2/auth [get]
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
	err := s.oidcFlowCache.Set(ctx, code, model, store.WithExpiration(30*time.Minute))
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}

	mm, err := s.oidcFlowCache.Get(ctx, code)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("mm", mm).Msg("mm")
	// redirect to the server Auth login pages.
	//
	return c.Redirect(http.StatusTemporaryRedirect, "/login?code="+code)
}
