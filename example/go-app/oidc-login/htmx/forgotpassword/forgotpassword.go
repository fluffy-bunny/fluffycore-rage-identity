package forgotpassword

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/components"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v5"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
		oidcSession:      oidcSession,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXForgotPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ForgotPasswordPostRequest struct {
	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
}

func (s *service) Do(c *echo.Context) error {
	r := c.Request()
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	return c.NoContent(http.StatusNotFound)
}

func (s *service) renderForgotPassword(c *echo.Context, code int, errors []string, email string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, code, components.ForgotPasswordPartial(components.ForgotPasswordData{
		RenderContext: rc,
		Errors:        errors,
		Email:         email,
	}))
}

func (s *service) renderError(c *echo.Context, errorCode, errorMessage string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.ErrorPartial(components.ErrorData{
		RenderContext: rc,
		ErrorCode:     errorCode,
		ErrorMessage:  errorMessage,
	}))
}

func (s *service) DoGet(c *echo.Context) error {
	return s.renderForgotPassword(c, http.StatusOK, nil, "")
}

func (s *service) DoPost(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &ForgotPasswordPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderForgotPassword(c, http.StatusBadRequest, []string{"Invalid request"}, "")
	}

	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		msg := utils.LocalizeSimple(localizer, "username.is.empty")
		return s.renderForgotPassword(c, http.StatusBadRequest, []string{msg}, "")
	}

	model.Email = strings.ToLower(model.Email)
	_, ok := echo_utils.IsValidEmailAddress(model.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "username.not.valid", map[string]string{"username": model.Email})
		return s.renderForgotPassword(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	// Always show verify-code page to prevent email enumeration
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
			log.Error().Err(err).Msg("GetRageUser")
			return s.renderError(c, "htmx-forgot-001", err.Error())
		}
	}

	subject := "NA"
	if getRageUserResponse != nil {
		subject = getRageUserResponse.User.RootIdentity.Subject
	}

	codeResult, err := echo_utils.GenerateHashedVerificationCode(ctx, s.passwordHasher, 6)
	if err != nil {
		log.Error().Err(err).Msg("GenerateHashedVerificationCode")
		return s.renderError(c, "htmx-forgot-002", err.Error())
	}

	plainCode := ""
	if s.config.SystemConfig.DeveloperMode {
		plainCode = codeResult.PlainCode
	}

	err = s.wellknownCookies.SetVerificationCodeCookie(c,
		&contracts_cookies.SetVerificationCodeCookieRequest{
			VerificationCode: &contracts_cookies.VerificationCode{
				Email:             model.Email,
				CodeHash:          codeResult.HashedCode,
				PlainCode:         plainCode,
				Subject:           subject,
				VerifyCodePurpose: contracts_cookies.VerifyCode_PasswordReset,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetVerificationCodeCookie")
		return s.renderError(c, "htmx-forgot-003", err.Error())
	}

	// Send email if user exists
	if getRageUserResponse != nil {
		message, err := localizer.LocalizeMessage(&i18n.Message{
			ID: "password.reset.message"})
		if err != nil {
			log.Error().Err(err).Msg("failed to localize message")
		} else {
			message = strings.ReplaceAll(message, "{code}", codeResult.PlainCode)
			s.EmailService().SendEmail(ctx,
				&contracts_email.SendEmailRequest{
					ToEmail:      model.Email,
					SubjectId:    "forgotpassword.email.subject",
					HtmlTemplate: "emails/generic/index",
					TextTemplate: "emails/generic/txt",
					Data: map[string]interface{}{
						"body": message,
					},
				})
		}
	}

	devCode := ""
	if s.config.SystemConfig.DeveloperMode {
		devCode = codeResult.PlainCode
	}

	localizer2 := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer2)
	return components.RenderNode(c, http.StatusOK, components.VerifyCodePartial(components.VerifyCodeData{
		RenderContext: rc,
		Email:         model.Email,
		Directive:     "passwordReset",
		Code:          devCode,
		Errors:        []string{},
	}))
}
