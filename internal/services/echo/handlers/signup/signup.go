package signup

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/util"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		services_echo_handlers_base.BaseHandler

		config           *contracts_config.Config
		container        di.Container
		oidcFlowStore    contracts_eko_gocache.IOIDCFlowStore
		idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer
		someUtil         contracts_util.ISomeUtil
		userService      proto_oidc_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	config *contracts_config.Config,
	someUtil contracts_util.ISomeUtil,
	userService proto_oidc_user.IFluffyCoreUserServiceServer,
	container di.Container,
	oidcFlowStore contracts_eko_gocache.IOIDCFlowStore,
	claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal,
	idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer,
	localizer contracts_localizer.ILocalizer,
	echoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor) (*service, error) {

	return &service{
		BaseHandler: services_echo_handlers_base.BaseHandler{
			ClaimsPrincipal:     claimsPrincipal,
			EchoContextAccessor: echoContextAccessor,
			Localizer:           localizer,
		},
		config:           config,
		container:        container,
		someUtil:         someUtil,
		idpServiceServer: idpServiceServer,
		oidcFlowStore:    oidcFlowStore,
		userService:      userService,
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
		wellknown_echo.SignupPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type SignupGetRequest struct {
	WizardMode  bool   `param:"wizard_mode" query:"wizard_mode" form:"wizard_mode" json:"wizard_mode" xml:"wizard_mode"`
	RedirectURL string `param:"redirect_url" query:"redirect_url" form:"redirect_url" json:"redirect_url" xml:"redirect_url"`
}
type ExternalIDPAuthRequest struct {
	IDPSlug string `param:"idp_slug" query:"idp_slug" form:"idp_slug" json:"idp_slug" xml:"idp_slug"`
}
type SignupPostRequest struct {
	WizardMode bool   `param:"wizard_mode" query:"wizard_mode" form:"wizard_mode" json:"wizard_mode" xml:"wizard_mode"`
	Code       string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	UserName   string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password   string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
}

func (s *service) getIDPs(ctx context.Context) ([]*proto_oidc_models.IDP, error) {
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

		return nil, err
	}

	return listIDPResponse.Idps, nil
}
func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &SignupGetRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	type row struct {
		Key   string
		Value string
	}
	idps, err := s.getIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return c.Redirect(http.StatusFound, "/error")
	}

	echo_utils.SetCookieInterface(c, &http.Cookie{
		Name:     "_signup_request",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}, model)

	var rows []row
	//	rows = append(rows, row{Key: "code", Value: model.Code})

	return s.Render(c, http.StatusOK, "views/signup/index",
		map[string]interface{}{
			"defs": rows,
			"idps": idps,
			"isWizardMode": func() bool {
				return model.WizardMode
			},
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
func (s *service) validateSignupPostRequest(request *SignupPostRequest) ([]*Error, error) {
	var err error
	errors := make([]*Error, 0)

	if fluffycore_utils.IsEmptyOrNil(request.UserName) {

		errors = append(errors, NewErrorF("username", "username is empty"))
	} else {
		_, ok := echo_utils.IsValidEmailAddress(request.UserName)
		if !ok {
			errors = append(errors, NewErrorF("username", "username:%s is not a valid email address", request.UserName))
		}
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		errors = append(errors, NewErrorF("password", "password is empty"))
	}

	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	idps, err := s.getIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return c.Redirect(http.StatusFound, "/error")
	}
	// is the request get or post?
	model := &SignupPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Debug().Err(err).Msg("Bind")
		return s.Render(c, http.StatusBadRequest, "views/signup/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("model", "model is invalid")},
				"isWizardMode": func() bool {
					return model.WizardMode
				},
				"idps": idps,
			})
	}
	log.Info().Interface("model", model).Msg("model")
	errors, err := s.validateSignupPostRequest(model)
	if err != nil {
		return err
	}
	if len(errors) > 0 {
		return s.Render(c, http.StatusBadRequest, "views/signup/index",
			map[string]interface{}{
				"defs": errors,
				"isWizardMode": func() bool {
					return model.WizardMode
				},
				"idps": idps,
			})
	}
	model.UserName = strings.ToLower(model.UserName)
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
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	if len(listUserResponse.Users) > 0 {
		return s.Render(c, http.StatusBadRequest, "views/signup/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("username", "username:%s already exists", model.UserName)},
				"isWizardMode": func() bool {
					return model.WizardMode
				},
				"idps": idps,
			})
	}
	hash, err := echo_utils.GeneratePasswordHash(model.Password)
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return c.Redirect(http.StatusFound, "/error")
	}
	user := &proto_oidc_models.User{
		RootIdentity: &proto_oidc_models.Identity{
			Subject:       xid.New().String(),
			Email:         model.UserName,
			IdpSlug:       "root-idp",
			EmailVerified: false,
		},
		Password: &proto_oidc_models.Password{
			Hash: hash,
		},
		State: proto_oidc_models.UserState_USER_STATE_PENDING,
	}
	_, err = s.userService.CreateUser(ctx, &proto_oidc_user.CreateUserRequest{
		User: user,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	if s.config.EmailVerificationRequired {
		// send email
		// TODO implement email service.
		log.Error().Msg("TODO: send email")

	}
	var signupGetRequest *SignupGetRequest = &SignupGetRequest{}
	err = echo_utils.GetCookieInterface(c, "_signup_request", signupGetRequest)

	if err != nil {
		log.Error().Err(err).Msg("Unmarshal")
		return c.Redirect(http.StatusFound, "/error")
	}
	if !fluffycore_utils.IsEmptyOrNil(signupGetRequest.RedirectURL) {
		return c.Redirect(http.StatusFound, signupGetRequest.RedirectURL)
	}

	return c.Redirect(http.StatusFound, "/login")

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
