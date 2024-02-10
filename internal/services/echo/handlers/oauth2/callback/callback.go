package callback

/*
reference: https://github.com/go-oauth2/oauth2/blob/master/example/client/client.go
*/
import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/client"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		someUtil                contracts_util.ISomeUtil
		scopedMemoryCache       fluffycore_contracts_common.IScopedMemoryCache
		externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore
		clientServiceServer     proto_oidc_client.IFluffyCoreClientServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore,
	someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil:                someUtil,
		scopedMemoryCache:       scopedMemoryCache,
		externalOauth2FlowStore: externalOauth2FlowStore,
		clientServiceServer:     clientServiceServer,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OAuth2CallbackPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		//clientauthorization.AuthenticateOAuth2Client(),
	}
}

type CallbackRequest struct {
	Code  string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code    				query     string  true  "code requested"
// @Param       state            		query     string  true  "state requested"
// @Success 200 {object} string
// @Router /oauth2/callback [get]
func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &CallbackRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	return c.Redirect(http.StatusTemporaryRedirect, "/login?code="+model.Code)
}
