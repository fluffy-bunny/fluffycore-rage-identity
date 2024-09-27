package api_password_reset_finish

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

// 	API_PasswordResetFinish    = "/api/password-reset-finish"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_PasswordResetFinish,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validatePasswordResetFinishRequest(model *login_models.PasswordResetFinishRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return status.Error(codes.InvalidArgument, "model.Password is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.PasswordConfirm) {
		return status.Error(codes.InvalidArgument, "model.PasswordConfirm is nil")
	}
	if model.Password != model.PasswordConfirm {
		return status.Error(codes.InvalidArgument, "model.Password and model.PasswordConfirm do not match")
	}

	return nil
}

// API Manifest godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		login_models.PasswordResetFinishRequest	true	"PasswordResetStartRequest"
// @Success 200 {object} login_models.PasswordResetFinishResponse
// @Failure 400 {string} login_models.PasswordResetFinishResponse
// @Router /api/password-reset-finish [post]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &login_models.PasswordResetFinishRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validatePasswordResetFinishRequest(model); err != nil {
		log.Error().Err(err).Msg("validatePasswordResetFinishRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	response := &login_models.PasswordResetFinishResponse{}
	getPasswordResetCookieResponse, err := s.wellknownCookies.GetPasswordResetCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetPasswordResetCookie")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if getPasswordResetCookieResponse == nil {
		response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	if getPasswordResetCookieResponse.PasswordReset == nil {
		response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage
	}
	if fluffycore_utils.IsEmptyOrNil(getPasswordResetCookieResponse.PasswordReset.Subject) {
		s.wellknownCookies.DeletePasswordResetCookie(c)
		response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage
	}
	getUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	// do password acceptablity check
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		response.Directive = login_models.DIRECTIVE_PasswordReset_DisplayPasswordResetPage
		response.ErrorReason = login_models.PasswordResetErrorReason_InvalidPassword
		return c.JSONPretty(http.StatusBadRequest, response, "  ")
	}
	// hash the password
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx,
		&contracts_identity.HashPasswordRequest{
			Password: model.Password,
		})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
			Password: &proto_oidc_models.PasswordUpdate{
				Hash: &wrapperspb.StringValue{
					Value: hashPasswordResponse.HashedPassword,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	// send the email
	sendSimpleEmailRequest := &contracts_email.SendSimpleEmailRequest{
		ToEmail:   getUserResponse.User.RootIdentity.Email,
		SubjectId: "password.reset.changed.subject",
		BodyId:    "password.reset.changed.message",
	}
	log.Debug().Interface("sendSimpleEmailRequest", sendSimpleEmailRequest).Err(err).Msg("SendEmail Reset Password")
	_, err = s.EmailService().SendSimpleEmail(ctx, sendSimpleEmailRequest)
	if err != nil {
		log.Error().Interface("sendSimpleEmailRequest", sendSimpleEmailRequest).Err(err).Msg("SendEmail")
		// eat it since we have already updated the password
		err = nil
	}
	s.wellknownCookies.DeletePasswordResetCookie(c)
	response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage

	return c.JSONPretty(http.StatusOK, response, "  ")
}
