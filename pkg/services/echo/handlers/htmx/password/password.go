package password

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx/components"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v5"
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
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
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
		wellknown_echo.HTMXPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type PasswordPostRequest struct {
	Email    string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
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

func (s *service) renderPassword(c *echo.Context, code int, errors []string, email string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, code, components.PasswordPartial(components.PasswordData{
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
	signinResponse, err := s.wellknownCookies.GetSigninUserNameCookie(c)
	if err != nil {
		return s.renderError(c, "htmx-password-001", "Session expired. Please start over.")
	}
	return s.renderPassword(c, http.StatusOK, nil, signinResponse.Value.Email)
}

func (s *service) DoPost(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	r := c.Request()
	rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &PasswordPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderError(c, "htmx-password-099", "Invalid request")
	}

	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return s.renderPassword(c, http.StatusBadRequest,
			[]string{utils.LocalizeSimple(localizer, "password.is.empty")}, model.Email)
	}

	session, err := s.oidcSession.GetSession()
	if err != nil {
		log.Error().Err(err).Msg("GetSession")
		return s.renderError(c, "htmx-password-002", "Session error")
	}

	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("Get request from session")
		return s.renderError(c, "htmx-password-003", "Session error")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	model.Email = strings.ToLower(model.Email)

	// Get user
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.Email,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			msg := utils.LocalizeWithInterperlate(localizer, "user.not.found", map[string]string{"username": model.Email})
			return s.renderPassword(c, http.StatusBadRequest, []string{msg}, model.Email)
		}
		log.Error().Err(err).Msg("GetRageUser")
		return s.renderError(c, "htmx-password-004", err.Error())
	}
	user := getRageUserResponse.User

	// Email verification check
	if s.config.EmailVerificationRequired && !user.RootIdentity.EmailVerified {
		return s.handleEmailVerification(c, ctx, model.Email, user, models.VerifyEmailDirective, contracts_cookies.VerifyCode_EmailVerification)
	}

	if user.Password == nil {
		msg := utils.LocalizeWithInterperlate(localizer, "username.does.not.have.password", map[string]string{"username": model.Email})
		return s.renderPassword(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	// Verify password
	err = s.passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
		Password:       model.Password,
		HashedPassword: user.Password.Hash,
	})
	if err != nil {
		log.Warn().Err(err).Msg("ComparePasswordHash")
		msg := utils.LocalizeWithInterperlate(localizer, "password.is.invalid", nil)
		return s.renderPassword(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	// MFA checks
	if user.TOTP == nil {
		user.TOTP = &proto_oidc_models.TOTP{Enabled: false}
	}
	if s.config.MultiFactorRequired || s.config.MultiFactorRequiredByEmailCode {
		return s.handleEmailVerification(c, ctx, model.Email, user, models.MFA_VerifyEmailDirective, contracts_cookies.VerifyCode_Challenge)
	}

	// Process final authentication
	result, err := s.ProcessFinalAuthenticationState(ctx, c, &services_echo_handlers_base.ProcessFinalAuthenticationStateRequest{
		AuthorizationRequest: authorizationRequest,
		Identity: &proto_oidc_models.OIDCIdentity{
			Subject: user.RootIdentity.Subject,
			Email:   user.RootIdentity.Email,
			IdpSlug: models.RootIdp,
			Acr: []string{
				models.ACRPassword,
				models.ACRIdpRoot,
			},
			Amr: []string{
				models.AMRPassword,
				models.AMRIdp,
			},
		},
		RootPath: rootPath,
	})
	if err != nil {
		log.Error().Err(err).Msg("ProcessFinalAuthenticationState")
		return s.renderError(c, "htmx-password-005", err.Error())
	}

	// Use HX-Redirect for the final OAuth redirect
	c.Response().Header().Set("HX-Redirect", result.RedirectURI)
	return c.NoContent(http.StatusOK)
}

func (s *service) handleEmailVerification(c *echo.Context, ctx interface{ Value(interface{}) interface{} }, email string, user *proto_oidc_models.RageUser, directive string, purpose contracts_cookies.VerifyCodePurpose) error {
	log := zerolog.Ctx(c.Request().Context()).With().Logger()
	goCtx := c.Request().Context()

	codeResult, err := echo_utils.GenerateHashedVerificationCode(goCtx, s.passwordHasher, 6)
	if err != nil {
		log.Error().Err(err).Msg("GenerateHashedVerificationCode")
		return s.renderError(c, "htmx-password-006", "Failed to generate verification code")
	}

	plainCode := ""
	if s.config.SystemConfig.DeveloperMode {
		plainCode = codeResult.PlainCode
	}

	err = s.wellknownCookies.SetVerificationCodeCookie(c,
		&contracts_cookies.SetVerificationCodeCookieRequest{
			VerificationCode: &contracts_cookies.VerificationCode{
				Email:             email,
				CodeHash:          codeResult.HashedCode,
				PlainCode:         plainCode,
				Subject:           user.RootIdentity.Subject,
				VerifyCodePurpose: purpose,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetVerificationCodeCookie")
		return s.renderError(c, "htmx-password-007", "Failed to set verification code")
	}

	s.EmailService().SendSimpleEmail(goCtx,
		&contracts_email.SendSimpleEmailRequest{
			ToEmail:   email,
			SubjectId: "email.verification.subject",
			BodyId:    "email.verification.message",
			Data: map[string]string{
				"code": codeResult.PlainCode,
			},
		})

	code := ""
	if s.config.SystemConfig.DeveloperMode {
		code = codeResult.PlainCode
	}

	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.VerifyCodePartial(components.VerifyCodeData{
		RenderContext: rc,
		Email:         email,
		Directive:     directive,
		Code:          code,
		Errors:        []string{},
	}))
}
