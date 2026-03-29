package signup

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
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
		oidcSession      contracts_oidc_session.IOIDCSession
		userIdGenerator  contracts_identity.IUserIdGenerator
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
	userIdGenerator contracts_identity.IUserIdGenerator,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
		oidcSession:      oidcSession,
		userIdGenerator:  userIdGenerator,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXSignupPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type SignupPostRequest struct {
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

func (s *service) renderSignup(c *echo.Context, code int, errors []string, email string) error {
	ctx := c.Request().Context()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	idps, _ := s.GetIDPs(ctx)
	return components.RenderNode(c, code, components.SignupPartial(components.SignupData{
		RenderContext: rc,
		Errors:        errors,
		Email:         email,
		SocialIdps:    idps,
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
	return s.renderSignup(c, http.StatusOK, nil, "")
}

func (s *service) DoPost(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &SignupPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderSignup(c, http.StatusBadRequest, []string{"Invalid request"}, "")
	}

	// Validate
	var errors []string
	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		errors = append(errors, utils.LocalizeSimple(localizer, "username.is.empty"))
	} else {
		_, ok := echo_utils.IsValidEmailAddress(model.Email)
		if !ok {
			errors = append(errors, utils.LocalizeWithInterperlate(localizer, "username.is.not.valid", map[string]string{"username": model.Email}))
		}
	}
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		errors = append(errors, utils.LocalizeSimple(localizer, "password.is.empty"))
	}
	if len(errors) > 0 {
		return s.renderSignup(c, http.StatusBadRequest, errors, model.Email)
	}

	model.Email = strings.ToLower(model.Email)

	// Check if domain is claimed
	parts := strings.Split(model.Email, "@")
	domainPart := parts[1]
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			ClaimedDomains: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		return s.renderSignup(c, http.StatusInternalServerError, []string{err.Error()}, model.Email)
	}
	if len(listIDPRequest.IDPs) > 0 {
		c.Response().Header().Set("HX-Redirect", wellknown_echo.ExternalIDPPath+"?idp_hint="+listIDPRequest.IDPs[0].Slug+"&directive=login")
		return c.NoContent(http.StatusOK)
	}

	// Check if user already exists
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
		} else {
			log.Error().Err(err).Msg("GetRageUser")
			return s.renderError(c, "htmx-signup-001", err.Error())
		}
	}
	if getRageUserResponse != nil {
		msg := utils.LocalizeWithInterperlate(localizer, "username.already.exists", map[string]string{"username": model.Email})
		return s.renderSignup(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	// Check password strength
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		msg := utils.LocalizeWithInterperlate(localizer, "password.is.not.acceptable", map[string]string{"username": model.Email})
		return s.renderSignup(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return s.renderError(c, "htmx-signup-002", "Failed to hash password")
	}

	subjectId := s.userIdGenerator.GenerateUserId()
	now := timestamppb.Now()
	user := &proto_oidc_models.RageUser{
		RootIdentity: &proto_oidc_models.Identity{
			Subject:       subjectId,
			Email:         model.Email,
			IdpSlug:       "root",
			EmailVerified: false,
			CreatedOn:     now,
			UpdatedOn:     now,
		},
		Password: &proto_oidc_models.Password{
			Hash: hashPasswordResponse.HashedPassword,
		},
		State: proto_oidc_models.RageUserState_USER_STATE_PENDING,
	}

	_, err = s.RageUserService().CreateRageUser(ctx, &proto_oidc_user.CreateRageUserRequest{
		User: user,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return s.renderError(c, "htmx-signup-003", "Failed to create user")
	}
	if err := s.SubmitAuditEvent(ctx,
		"com.fluffybunny.identity.user.created",
		user.RootIdentity.Subject,
		map[string]string{"email": user.RootIdentity.Email, "idp_slug": user.RootIdentity.IdpSlug},
		map[string]string{"mutation": "create_user", "handler": "htmx.signup"}); err != nil {
		log.Error().Err(err).Msg("SubmitAuditEvent")
		return s.renderError(c, "htmx-signup-003-audit", "Failed to write audit record")
	}

	if s.config.EmailVerificationRequired {
		codeResult, err := echo_utils.GenerateHashedVerificationCode(ctx, s.passwordHasher, 6)
		if err != nil {
			log.Error().Err(err).Msg("GenerateHashedVerificationCode")
			return s.renderError(c, "htmx-signup-004", "Failed to generate verification code")
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
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_EmailVerification,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return s.renderError(c, "htmx-signup-005", "Failed to set verification code")
		}
		s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.Email,
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
			Email:         model.Email,
			Directive:     "verifyEmailDirective",
			Code:          code,
			Errors:        []string{},
		}))
	}

	// No email verification - go back to login
	rc2 := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.HomePartial(components.HomeData{
		RenderContext: rc2,
		Errors:        []string{},
		Email:         model.Email,
	}))
}
