package signup

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	services_handlers_shared "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/shared"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config           *contracts_config.Config
		passwordHasher   contracts_identity.IPasswordHasher
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	passwordHasher contracts_identity.IPasswordHasher,
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_oidc_user.IFluffyCoreUserServiceServer,
) (*service, error) {

	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			// do auto post
			//contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.SignupPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type SignupGetRequest struct {
}
type ExternalIDPAuthRequest struct {
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}
type SignupPostRequest struct {
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	Type     string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
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
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return c.Redirect(http.StatusFound, "/error")
	}

	var rows []row
	//	rows = append(rows, row{Key: "code", Value: model.Code})

	return s.Render(c, http.StatusOK, "oidc/signup/index",
		map[string]interface{}{
			"errors":    rows,
			"idps":      idps,
			"directive": models.SignupDirective,
		})
}

func (s *service) validateSignupPostRequest(request *SignupPostRequest) ([]*services_handlers_shared.Error, error) {
	var err error
	errors := make([]*services_handlers_shared.Error, 0)

	if fluffycore_utils.IsEmptyOrNil(request.UserName) {

		errors = append(errors, services_handlers_shared.NewErrorF("username", "username is empty"))
	} else {
		_, ok := echo_utils.IsValidEmailAddress(request.UserName)
		if !ok {
			errors = append(errors, services_handlers_shared.NewErrorF("username", "username:%s is not a valid email address", request.UserName))
		}
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		errors = append(errors, services_handlers_shared.NewErrorF("password", "password is empty"))
	}

	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	idps, err := s.GetIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return c.Redirect(http.StatusFound, "/error")
	}
	doError := func(errors []*services_handlers_shared.Error) error {
		return s.Render(c, http.StatusBadRequest, "oidc/signup/index",
			map[string]interface{}{
				"errors":    errors,
				"idps":      idps,
				"directive": models.SignupDirective,
			})

	}
	// is the request get or post?
	model := &SignupPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Debug().Err(err).Msg("Bind")
		return doError([]*services_handlers_shared.Error{
			services_handlers_shared.NewErrorF("model", "model is invalid"),
		})
	}
	log.Info().Interface("model", model).Msg("model")
	if model.Type == "GET" {
		return s.DoGet(c)
	}
	errors, err := s.validateSignupPostRequest(model)
	if err != nil {
		return err
	}
	if len(errors) > 0 {
		return doError(errors)
	}
	model.UserName = strings.ToLower(model.UserName)
	// get the domain from the email
	parts := strings.Split(model.UserName, "@")
	domainPart := parts[1]
	// first lets see if this domain has been claimed.
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			ClaimedDomain: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		errors = append(errors, services_handlers_shared.NewErrorF("error", err.Error()))
		return doError(errors)
	}
	if len(listIDPRequest.Idps) > 0 {
		// an idp has claimed this domain.
		// post to the externalIDP
		return s.RenderAutoPost(c, wellknown_echo.ExternalIDPPath,
			[]models.FormParam{

				{
					Name:  "idp_hint",
					Value: listIDPRequest.Idps[0].Slug,
				},
				{
					Name:  "directive",
					Value: models.LoginDirective,
				},
			})
	}
	// does the user exist.
	listUserResponse, err := s.UserService().ListUser(ctx, &proto_oidc_user.ListUserRequest{
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
		return doError([]*services_handlers_shared.Error{
			services_handlers_shared.NewErrorF("username", "username:%s already exists", model.UserName),
		})
	}
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return c.Redirect(http.StatusFound, "/error")
	}
	user := &proto_oidc_models.User{
		RootIdentity: &proto_oidc_models.Identity{
			Subject:       xid.New().String(),
			Email:         model.UserName,
			IdpSlug:       models.RootIdp,
			EmailVerified: false,
		},
		Password: &proto_oidc_models.Password{
			Hash: hashPasswordResponse.HashedPassword,
		},
		State: proto_oidc_models.UserState_USER_STATE_PENDING,
	}
	_, err = s.UserService().CreateUser(ctx, &proto_oidc_user.CreateUserRequest{
		User: user,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	if s.config.EmailVerificationRequired {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:   model.UserName,
					Code:    verificationCode,
					Subject: user.RootIdentity.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return c.Redirect(http.StatusFound, "/error")
		}
		_, err = s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.UserName,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification..message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendSimpleEmail")
			return c.Redirect(http.StatusFound, "/error")
		}
		formParams := []models.FormParam{

			{
				Name:  "email",
				Value: model.UserName,
			},
			{
				Name:  "directive",
				Value: models.VerifyEmailDirective,
			},
			{
				Name:  "type",
				Value: "GET",
			},
		}

		if s.config.SystemConfig.DeveloperMode {
			formParams = append(formParams, models.FormParam{
				Name:  "code",
				Value: verificationCode,
			})

		}
		return s.RenderAutoPost(c, wellknown_echo.VerifyCodePath, formParams)

	}

	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
		[]models.FormParam{})

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