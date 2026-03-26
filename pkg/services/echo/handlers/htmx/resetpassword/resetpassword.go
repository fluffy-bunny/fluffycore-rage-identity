package resetpassword

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx/components"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
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
		wellknown_echo.HTMXResetPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ResetPasswordPostRequest struct {
	Email           string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Password        string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	ConfirmPassword string `param:"confirmPassword" query:"confirmPassword" form:"confirmPassword" json:"confirmPassword" xml:"confirmPassword"`
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

func (s *service) renderResetPassword(c *echo.Context, code int, errors []string, email string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, code, components.ResetPasswordPartial(components.ResetPasswordData{
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
	getPasswordResetCookieResponse, err := s.wellknownCookies.GetPasswordResetCookie(c)
	if err != nil || getPasswordResetCookieResponse == nil || getPasswordResetCookieResponse.PasswordReset == nil {
		return s.renderError(c, "htmx-reset-001", "Password reset session expired. Please start over.")
	}
	return s.renderResetPassword(c, http.StatusOK, nil, "")
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &ResetPasswordPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderError(c, "htmx-reset-099", "Invalid request")
	}

	// Validate
	var errors []string
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		errors = append(errors, "Password is required")
	}
	if fluffycore_utils.IsEmptyOrNil(model.ConfirmPassword) {
		errors = append(errors, "Password confirmation is required")
	}
	if model.Password != model.ConfirmPassword {
		errors = append(errors, "Passwords do not match")
	}
	if len(errors) > 0 {
		return s.renderResetPassword(c, http.StatusBadRequest, errors, model.Email)
	}

	// Get password reset cookie
	getPasswordResetCookieResponse, err := s.wellknownCookies.GetPasswordResetCookie(c)
	if err != nil || getPasswordResetCookieResponse == nil || getPasswordResetCookieResponse.PasswordReset == nil {
		return s.renderError(c, "htmx-reset-002", "Password reset session expired. Please start over.")
	}

	if fluffycore_utils.IsEmptyOrNil(getPasswordResetCookieResponse.PasswordReset.Subject) {
		s.wellknownCookies.DeletePasswordResetCookie(c)
		return s.renderError(c, "htmx-reset-003", "Invalid password reset session.")
	}

	// Get user
	getUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("GetRageUser")
		return s.renderError(c, "htmx-reset-004", err.Error())
	}

	// Check password strength
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		return s.renderResetPassword(c, http.StatusBadRequest, []string{"Password does not meet requirements"}, model.Email)
	}

	// Hash password
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("HashPassword")
		return s.renderError(c, "htmx-reset-005", "Failed to hash password")
	}

	// Update user
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
		log.Error().Err(err).Msg("UpdateRageUser")
		return s.renderError(c, "htmx-reset-006", "Failed to update password")
	}
	if err := s.SubmitAuditEvent(ctx,
		"com.fluffybunny.identity.user.password.updated",
		getPasswordResetCookieResponse.PasswordReset.Subject,
		map[string]string{"operation": "password_reset"},
		map[string]string{"mutation": "update_password", "handler": "htmx.resetpassword"}); err != nil {
		log.Error().Err(err).Msg("SubmitAuditEvent")
		return s.renderError(c, "htmx-reset-006-audit", "Failed to write audit record")
	}

	// Send confirmation email
	s.EmailService().SendSimpleEmail(ctx,
		&contracts_email.SendSimpleEmailRequest{
			ToEmail:   getUserResponse.User.RootIdentity.Email,
			SubjectId: "password.reset.changed.subject",
			BodyId:    "password.reset.changed.message",
		})

	s.wellknownCookies.DeletePasswordResetCookie(c)

	// Go back to login
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.HomePartial(components.HomeData{
		RenderContext: rc,
		Errors:        []string{},
		Email:         getUserResponse.User.RootIdentity.Email,
	}))
}
