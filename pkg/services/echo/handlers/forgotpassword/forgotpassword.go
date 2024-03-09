package forgotpassword

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		wellknownCookies contracts_cookies.IWellknownCookies
		config           *contracts_config.Config
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_ForgotPassword_001 = "rg-forgot-password-001"
	InternalError_ForgotPassword_002 = "rg-forgot-password-002"
	InternalError_ForgotPassword_003 = "rg-forgot-password-003"
	InternalError_ForgotPassword_004 = "rg-forgot-password-004"
	InternalError_ForgotPassword_005 = "rg-forgot-password-005"
	InternalError_ForgotPassword_006 = "rg-forgot-password-006"
	InternalError_ForgotPassword_007 = "rg-forgot-password-007"
	InternalError_ForgotPassword_008 = "rg-forgot-password-008"
	InternalError_ForgotPassword_009 = "rg-forgot-password-009"
	InternalError_ForgotPassword_010 = "rg-forgot-password-010"
	InternalError_ForgotPassword_011 = "rg-forgot-password-011"
	InternalError_ForgotPassword_099 = "rg-forgot-password-099" // 99 is a bind problem
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies: wellknownCookies,
		config:           config,
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
		wellknown_echo.ForgotPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ForgotPasswordGetRequest struct {
	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
}

type ForgotPasswordPostRequest struct {
	Email  string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Type   string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
	Action string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
}

func (s *service) validateForgotPasswordGetRequest(request *ForgotPasswordGetRequest) error {

	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &ForgotPasswordGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return s.TeleportBackToLogin(c, InternalError_ForgotPassword_099)
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validateForgotPasswordGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateForgotPasswordGetRequest")
		return s.TeleportBackToLogin(c, InternalError_ForgotPassword_002)
	}

	err = s.Render(c, http.StatusOK, "oidc/forgotpassword/index",
		map[string]interface{}{
			"email":  model.Email,
			"errors": make([]string, 0),
		})
	return err
}

func (s *service) validateForgotPasswordPostRequest(request *ForgotPasswordPostRequest) ([]string, error) {
	var err error
	localizer := s.Localizer().GetLocalizer()
	errors := make([]string, 0)

	if fluffycore_utils.IsEmptyOrNil(request.Email) {
		msg := utils.LocalizeWithInterperlate(localizer, "email.is.empty", nil)
		errors = append(errors, msg)
	}
	_, ok := echo_utils.IsValidEmailAddress(request.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "email.is.not.valid", map[string]string{"email": request.Email})
		errors = append(errors, msg)
	}
	request.Email = strings.ToLower(request.Email)
	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &ForgotPasswordPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLogin(c, InternalError_ForgotPassword_099)
	}
	log.Info().Interface("model", model).Msg("model")

	errors, err := s.validateForgotPasswordPostRequest(model)
	if err != nil {
		return s.Render(c, http.StatusBadRequest, "oidc/forgotpassword/index",
			map[string]interface{}{
				"email":  model.Email,
				"errors": errors,
			})
	}
	if model.Type == "GET" {
		return s.DoGet(c)
	}
	if model.Action == "cancel" {
		return s.TeleportToPath(c, wellknown_echo.OIDCLoginPath)
	}
	// NOTE: We don't want to give bots the ability to probe our service to see if an email exists.
	// we check here and we redirect to the enter code in all cases.
	// we just don't send the email, but we drop the cookie with a verification code just for the fun of it.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.Email,
			},
		})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			err = nil
		}
		if err != nil {
			log.Error().Err(err).Msg("ListUser")
			return s.TeleportBackToLogin(c, InternalError_ForgotPassword_003)
		}
	}
	subject := "NA"
	if getRageUserResponse != nil {
		subject = getRageUserResponse.User.RootIdentity.Subject
	}

	verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
	err = s.wellknownCookies.SetVerificationCodeCookie(c, &contracts_cookies.SetVerificationCodeCookieRequest{
		VerificationCode: &contracts_cookies.VerificationCode{
			Email:   model.Email,
			Code:    verificationCode,
			Subject: subject,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetVerificationCodeCookie")
		return s.TeleportBackToLogin(c, InternalError_ForgotPassword_004)
	}
	localizer := s.Localizer().GetLocalizer()
	message, err := localizer.LocalizeMessage(&i18n.Message{ID: "password.reset.message"})
	if err != nil {
		log.Error().Err(err).Msg("failed to localize message")
		return s.TeleportBackToLogin(c, InternalError_ForgotPassword_005)
	}
	message = strings.ReplaceAll(message, "{code}", verificationCode)
	if getRageUserResponse != nil {
		// send the email
		_, err = s.EmailService().SendEmail(ctx,
			&contracts_email.SendEmailRequest{
				ToEmail:      model.Email,
				SubjectId:    "forgotpassword.email.subject",
				HtmlTemplate: "emails/generic/index",
				TextTemplate: "emails/generic/txt",
				Data: map[string]interface{}{
					"body": message,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendEmail")
			return s.TeleportBackToLogin(c, InternalError_ForgotPassword_006)
		}
	} else {
		// no user found, is a probe.
		verificationCode = "NA"
	}
	formParams := []models.FormParam{

		{
			Name:  "email",
			Value: model.Email,
		},
		{
			Name:  "directive",
			Value: models.PasswordResetDirective,
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
