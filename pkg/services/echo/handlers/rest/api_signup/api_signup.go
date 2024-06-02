package api_signup

import (
	"net/http"

	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
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
		userIdGenerator  contracts_identity.IUserIdGenerator
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
	userIdGenerator contracts_identity.IUserIdGenerator,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
		userIdGenerator:  userIdGenerator,
	}, nil
}

// 	API_VerifyCode             = "/api/signup"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_Signup,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateSignupRequest(model *login_models.SignupRequest) error {
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
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return status.Error(codes.InvalidArgument, "model.Password is nil")
	}
	return nil
}

// API VerifyCode godoc
// @Summary verify code.
// @Description verify code
// @Tags root
// @Accept */*
// @Produce json
// @Param		request body		login_models.SignupRequest	true	"SignupRequest"
// @Success 200 {object} login_models.SignupResponse
// @Failure 302 {string} login_models.SignupResponse
// @Failure 400 {string} login_models.SignupResponse
// @Failure 500 {string} string
// @Router /api/signup [post]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &login_models.SignupRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if err := s.validateSignupRequest(model); err != nil {
		log.Error().Err(err).Msg("validateSignupRequest")
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}
	response := &login_models.SignupResponse{
		Email:       model.Email,
		ErrorReason: login_models.SignupErrorReason_NoError,
	}

	session, err := s.getSession()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")

	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	model.Email = strings.ToLower(model.Email)
	// get the domain from the email
	parts := strings.Split(model.Email, "@")
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
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if len(listIDPRequest.Idps) > 0 {
		// this domain is claimed.
		response.Directive = login_models.DIRECTIVE_Redirect
		response.DirectiveRedirect = &login_models.DirectiveRedirect{
			RedirectURI: wellknown_echo.ExternalIDPPath,
			VERB:        http.MethodPost,
			FormParams: []models.FormParam{
				{
					Name:  "state",
					Value: authorizationRequest.State,
				},
				{
					Name:  "idp_hint",
					Value: listIDPRequest.Idps[0].Slug,
				},
				{
					Name:  "directive",
					Value: models.LoginDirective,
				},
			},
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: strings.ToLower(model.Email),
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			err = nil
		} else {
			log.Error().Err(err).Msg("GetRageUser")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
	}
	if getRageUserResponse != nil {
		response.Message = "User already exists"
		return c.JSONPretty(http.StatusFound, response, "  ")
	}

	subjectId := s.userIdGenerator.GenerateUserId()
	user := &proto_oidc_models.RageUser{
		RootIdentity: &proto_oidc_models.Identity{
			Subject:       subjectId,
			Email:         model.Email,
			IdpSlug:       models.RootIdp,
			EmailVerified: false,
		},
		State: proto_oidc_models.RageUserState_USER_STATE_PENDING,
	}

	//  check password strength
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		response.ErrorReason = login_models.SignupErrorReason_InvalidPassword
		return c.JSONPretty(http.StatusBadRequest, response, "  ")
	}
	user.Password = &proto_oidc_models.Password{
		Hash: hashPasswordResponse.HashedPassword,
	}

	_, err = s.RageUserService().CreateRageUser(ctx, &proto_oidc_user.CreateRageUserRequest{
		User: user,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if s.config.EmailVerificationRequired {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             model.Email,
					Code:              verificationCode,
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_EmailVerification,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		_, err = s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.Email,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification..message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendSimpleEmail")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		if s.config.SystemConfig.DeveloperMode {
			response.DirectiveEmailCodeChallenge = &login_models.DirectiveEmailCodeChallenge{
				Code: verificationCode,
			}
		}
		response.Directive = login_models.DIRECTIVE_VerifyCode_DisplayVerifyCodePage
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage
	response.Message = "User created"
	return c.JSONPretty(http.StatusOK, response, "  ")

}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
