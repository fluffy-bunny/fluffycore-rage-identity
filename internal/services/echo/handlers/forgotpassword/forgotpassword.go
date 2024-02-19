package forgotpassword

import (
	"fmt"
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	"github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	services_handlers_shared "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/shared"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
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
		emailService     contracts_email.IEmailService
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	emailService contracts_email.IEmailService,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		emailService:     emailService,
		wellknownCookies: wellknownCookies,
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
		wellknown_echo.ForgotPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ForgotPasswordGetRequest struct {
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
}

type ForgotPasswordPostRequest struct {
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
}

func (s *service) validateForgotPasswordGetRequest(model *ForgotPasswordGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.State) {
		return status.Error(codes.InvalidArgument, "State is empty")
	}
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
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validateForgotPasswordGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateForgotPasswordGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.Render(c, http.StatusOK, "views/forgotpassword/index",
		map[string]interface{}{
			"state": model.State,
			"email": model.Email,
		})
	return err
}

func (s *service) validateForgotPasswordPostRequest(request *ForgotPasswordPostRequest) ([]*services_handlers_shared.Error, error) {
	var err error
	errors := make([]*services_handlers_shared.Error, 0)
	if fluffycore_utils.IsEmptyOrNil(request.State) {
		errors = append(errors, services_handlers_shared.NewErrorF("state", "State is empty"))
	}
	if fluffycore_utils.IsEmptyOrNil(request.Email) {
		errors = append(errors, services_handlers_shared.NewErrorF("email", "Email is empty"))
	}
	_, ok := echo_utils.IsValidEmailAddress(request.Email)
	if !ok {
		errors = append(errors, services_handlers_shared.NewErrorF("email", "Email:%s is not a valid email address", request.Email))
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
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	errors, err := s.validateForgotPasswordPostRequest(model)
	if err != nil {
		return s.Render(c, http.StatusBadRequest, "views/forgotpassword/index",
			map[string]interface{}{
				"state": model.State,
				"email": model.Email,
				"defs":  errors,
			})
	}

	// NOTE: We don't want to give bots the ability to probe our service to see if an email exists.
	// we check here and we redirect to the enter code in all cases.
	// we just don't send the email, but we drop the cookie with a verification code just for the fun of it.
	listUserResponse, err := s.UserService().ListUser(ctx,
		&proto_oidc_user.ListUserRequest{
			Filter: &proto_oidc_user.Filter{
				RootIdentity: &proto_oidc_user.IdentityFilter{
					Email: &proto_types.StringFilterExpression{
						Eq: model.Email,
					},
				},
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, "/Error")
	}
	subject := ""
	if listUserResponse != nil && len(listUserResponse.Users) > 0 {
		subject = listUserResponse.Users[0].RootIdentity.Subject
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
		return c.Redirect(http.StatusFound, "/error")
	}
	localizer := s.Localizer().GetLocalizer()
	message, err := localizer.LocalizeMessage(&i18n.Message{ID: "password.reset.message"})
	if err != nil {
		log.Error().Err(err).Msg("failed to localize message")
		return c.Redirect(http.StatusFound, "/error")
	}
	message = strings.ReplaceAll(message, "{code}", verificationCode)
	if len(listUserResponse.Users) > 0 {
		// send the email
		_, err = s.emailService.SendEmail(ctx,
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
			return c.Redirect(http.StatusFound, "/error")
		}
	}

	redirectURL := fmt.Sprintf("%s?state=%s&email=%s&directive=%s",
		wellknown_echo.VerifyCodePath,
		model.State,
		model.Email,
		models.PasswordResetDirective,
	)
	return c.Redirect(http.StatusFound, redirectURL)

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
