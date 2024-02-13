package oidclogin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/util"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_echo_handlers_base.BaseHandler
		container        di.Container
		oidcFlowStore    contracts_eko_gocache.IOIDCFlowStore
		idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer
		userService      proto_oidc_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(someUtil contracts_util.ISomeUtil,
	container di.Container,
	localizer contracts_localizer.ILocalizer,
	userService proto_oidc_user.IFluffyCoreUserServiceServer,
	oidcFlowStore contracts_eko_gocache.IOIDCFlowStore,
	claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal,
	idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer,
	echoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor) (*service, error) {

	return &service{
		BaseHandler: services_echo_handlers_base.BaseHandler{
			ClaimsPrincipal:     claimsPrincipal,
			EchoContextAccessor: echoContextAccessor,
			Localizer:           localizer,
		},
		container:        container,
		userService:      userService,
		idpServiceServer: idpServiceServer,
		oidcFlowStore:    oidcFlowStore,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.OIDCLoginPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LoginGetRequest struct {
	Code string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
}
type ExternalIDPAuthRequest struct {
	IDPSlug string `param:"idp_slug" query:"idp_slug" form:"idp_slug" json:"idp_slug" xml:"idp_slug"`
}
type LoginPostRequest struct {
	Code     string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginGetRequest{}
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
	type row struct {
		Key   string
		Value string
	}

	listIDPResponse, err := s.idpServiceServer.ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			Enabled: &proto_types.BoolFilterExpression{
				Eq: true,
			},
			Metadata: &proto_types.StringMapStringFilterExpression{
				Key: "hidden",
				Value: &proto_types.StringFilterExpression{
					Eq: "false",
				},
			},
		},
	})
	if err != nil {
		return err
	}
	var rows []row
	rows = append(rows, row{Key: "code", Value: model.Code})

	return s.Render(c, http.StatusOK, "views/oidclogin/index",
		map[string]interface{}{
			"defs": rows,
			"idps": listIDPResponse.Idps,
			"code": model.Code,
		})
}

type Error struct {
	Key   string `json:"key"`
	Value string `json:"msg"`
}

func NewErrorF(key string, value string, args ...interface{}) *Error {
	return &Error{
		Key:   key,
		Value: fmt.Sprintf(value, args...),
	}
}
func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginPostRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	model.UserName = strings.ToLower(model.UserName)

	// get the code from the cookie
	cookie, err := c.Cookie("_code")
	if err != nil {
		return err
	}
	code := cookie.Value
	log.Info().Str("code", code).Msg("code")

	// does the user exist.
	listUserResponse, err := s.userService.ListUser(ctx, &proto_oidc_user.ListUserRequest{
		Filter: &proto_oidc_user.Filter{
			RootIdentity: &proto_oidc_user.IdentityFilter{
				Email: &proto_types.StringFilterExpression{
					Eq: model.UserName,
				},
			},
		},
	})
	if len(listUserResponse.Users) == 0 {
		return s.Render(c, http.StatusBadRequest, "views/oidclogin/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("username", "username:%s does not exist", model.UserName)},
			})

	}
	if err != nil {
		log.Warn().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", wellknown_echo.OIDCLoginPath, code))
	}
	user := listUserResponse.Users[0]
	if user.Password == nil {
		return s.Render(c, http.StatusBadRequest, "views/oidclogin/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("username", "username:%s does not have a password", model.UserName)},
			})
	}

	passwordValid, err := echo_utils.ComparePasswordHash(model.Password, user.Password.Hash)
	if err != nil && !passwordValid {
		log.Warn().Err(err).Msg("ComparePasswordHash")
		redirectUrl := fmt.Sprintf("%s%s?code=%s", rootPath, wellknown_echo.OIDCLoginPath, code)
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	echo_utils.SetCookieInterface(c, &http.Cookie{
		Name:     "_auth",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}, user.RootIdentity)

	mm, err := s.oidcFlowStore.GetAuthorizationFinal(ctx, code)
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationFinal")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	mm.Identity = &models.Identity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		ACR:     []string{"urn:mastodon:password", "urn:mastodon:2fa", "urn:mastodon:idp:root"},
	}

	// "urn:mastodon:idp:google", "urn:mastodon:idp:spacex", "urn:mastodon:idp:github-enterprise", etc.
	// "urn:mastodon:password", "urn:mastodon:2fa", "urn:mastodon:email", etc.
	err = s.oidcFlowStore.StoreAuthorizationFinal(ctx, code, mm)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	// redirect to the client with the code.
	redirectUri := mm.Request.RedirectURI +
		"?code=" + code +
		"&state=" + mm.Request.State +
		"&iss=" + rootPath
	return c.Redirect(http.StatusFound, redirectUri)

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
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
