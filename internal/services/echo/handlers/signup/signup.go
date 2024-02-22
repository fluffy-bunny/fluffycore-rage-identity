package signup

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
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
	WizardMode bool   `param:"wizard_mode" query:"wizard_mode" form:"wizard_mode" json:"wizard_mode" xml:"wizard_mode"`
	State      string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
}
type ExternalIDPAuthRequest struct {
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}
type SignupPostRequest struct {
	WizardMode bool   `param:"wizard_mode" query:"wizard_mode" form:"wizard_mode" json:"wizard_mode" xml:"wizard_mode"`
	State      string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	UserName   string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password   string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	Type       string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
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

	echo_utils.SetCookieInterface(c, &http.Cookie{
		Name:     "_signup_request",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}, model)

	var rows []row
	//	rows = append(rows, row{Key: "code", Value: model.Code})

	return s.Render(c, http.StatusOK, "oidc/signup/index",
		map[string]interface{}{
			"defs": rows,
			"idps": idps,
			"isWizardMode": func() bool {
				return model.WizardMode
			},
			"state": model.State,
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

	idps, err := s.GetIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return c.Redirect(http.StatusFound, "/error")
	}
	// is the request get or post?
	model := &SignupPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Debug().Err(err).Msg("Bind")
		return s.Render(c, http.StatusBadRequest, "oidc/signup/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("model", "model is invalid")},
				"isWizardMode": func() bool {
					return model.WizardMode
				},
				"idps":  idps,
				"state": model.State,
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
		return s.Render(c, http.StatusBadRequest, "oidc/signup/index",
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
		return s.Render(c, http.StatusBadRequest, "oidc/signup/index",
			map[string]interface{}{
				"defs": []*Error{NewErrorF("username", "username:%s already exists", model.UserName)},
				"isWizardMode": func() bool {
					return model.WizardMode
				},
				"idps": idps,
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
				Name:  "state",
				Value: model.State,
			},
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
	var signupGetRequest *SignupGetRequest = &SignupGetRequest{}
	err = echo_utils.GetCookieInterface(c, "_signup_request", signupGetRequest)

	if err != nil {
		log.Error().Err(err).Msg("Unmarshal")
		return c.Redirect(http.StatusFound, "/error")
	}

	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
		[]models.FormParam{
			{
				Name:  "state",
				Value: model.State,
			},
		})

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
