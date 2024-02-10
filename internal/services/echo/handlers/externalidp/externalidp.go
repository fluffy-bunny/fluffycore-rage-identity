package externalidp

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_echo_handlers_base.BaseHandler
		container     di.Container
		oidcFlowStore contracts_eko_gocache.IOIDCFlowStore

		someUtil contracts_util.ISomeUtil
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(someUtil contracts_util.ISomeUtil,
	container di.Container,
	oidcFlowStore contracts_eko_gocache.IOIDCFlowStore,
	claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal,
	echoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor) (*service, error) {

	return &service{
		BaseHandler: services_echo_handlers_base.BaseHandler{
			ClaimsPrincipal: claimsPrincipal, EchoContextAccessor: echoContextAccessor},
		container:     container,
		someUtil:      someUtil,
		oidcFlowStore: oidcFlowStore,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.ExternalIDPPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ExternalIDPAuthRequest struct {
	IDPSlug string `param:"idp_slug" query:"idp_slug" form:"idp_slug" json:"idp_slug" xml:"idp_slug"`
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	//rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &ExternalIDPAuthRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")
	return c.Redirect(http.StatusFound, "/error")

}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code            		query     string  true  "code"
// @Success 200 {object} string
// @Router /login [get,post]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {

	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
