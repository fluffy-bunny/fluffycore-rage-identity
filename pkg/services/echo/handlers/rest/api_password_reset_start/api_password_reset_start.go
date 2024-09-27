package api_password_reset_start

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
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

		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
		oidcSession:      oidcSession,
	}, nil
}

// API_PasswordResetStart     = "/api/password-reset-start"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_PasswordResetStart,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validatePasswordResetStartRequest(model *login_models.PasswordResetStartRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		return status.Error(codes.InvalidArgument, "model.Email is nil")
	}
	_, ok := echo_utils.IsValidEmailAddress(model.Email)
	if !ok {
		return status.Error(codes.InvalidArgument, "model.Email is not a valid email address")
	}
	return nil
}

// API Manifest godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		login_models.PasswordResetStartRequest	true	"PasswordResetStartRequest"
// @Success 200 {object} login_models.PasswordResetStartResponse
// @Router /api/password-reset-start [post]
func (s *service) Do(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &login_models.PasswordResetStartRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validatePasswordResetStartRequest(model); err != nil {
		log.Error().Err(err).Msg("validatePasswordResetStartRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	response := &login_models.PasswordResetStartResponse{
		Email:     model.Email,
		Directive: login_models.DIRECTIVE_VerifyCode_DisplayVerifyCodePage,
	}

	model.Email = strings.ToLower(model.Email)
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
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
	}
	subject := "NA"
	if getRageUserResponse != nil {
		subject = getRageUserResponse.User.RootIdentity.Subject
	}
	verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
	err = s.wellknownCookies.SetVerificationCodeCookie(c, &contracts_cookies.SetVerificationCodeCookieRequest{
		VerificationCode: &contracts_cookies.VerificationCode{
			Email:             model.Email,
			Code:              verificationCode,
			Subject:           subject,
			VerifyCodePurpose: contracts_cookies.VerifyCode_PasswordReset,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetVerificationCodeCookie")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	message, err := localizer.LocalizeMessage(&i18n.Message{
		ID: "password.reset.message"})
	if err != nil {
		log.Error().Err(err).Msg("failed to localize message")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
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
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		if s.config.SystemConfig.DeveloperMode {
			response.DirectiveEmailCodeChallenge = &login_models.DirectiveEmailCodeChallenge{
				Code: verificationCode,
			}
		}
	} else {
		// no user found, is a probe.
		verificationCode = "NA"
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
