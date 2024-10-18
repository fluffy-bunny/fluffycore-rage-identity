package api_login_password

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
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
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
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_LoginPassword,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_LoginPassword_001 = "rg-loginpassword-001"
	InternalError_LoginPassword_002 = "rg-loginpassword-002"
	InternalError_LoginPassword_003 = "rg-loginpassword-003"
	InternalError_LoginPassword_004 = "rg-loginpassword-004"
	InternalError_LoginPassword_005 = "rg-loginpassword-005"
	InternalError_LoginPassword_006 = "rg-loginpassword-006"
	InternalError_LoginPassword_007 = "rg-loginpassword-007"
	InternalError_LoginPassword_008 = "rg-loginpassword-008"
	InternalError_LoginPassword_009 = "rg-loginpassword-009"
	InternalError_LoginPassword_010 = "rg-loginpassword-010"
	InternalError_LoginPassword_011 = "rg-loginpassword-011"
	InternalError_LoginPassword_099 = "rg-loginpassword-099"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateLoginPasswordRequest(model *login_models.LoginPasswordRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return status.Error(codes.InvalidArgument, "model.Password is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		return status.Error(codes.InvalidArgument, "model.Email is nil")
	}
	model.Email = strings.ToLower(model.Email)

	return nil
}

// API Manifest godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept  json
// @Produce json
// @Param		request body		login_models.LoginPasswordRequest	true	"LoginPasswordRequest"
// @Success 200 {object} login_models.LoginPasswordResponse
// @Failure 401 {object} login_models.LoginPasswordErrorResponse
// @Failure 500 {object} login_models.LoginPasswordErrorResponse
// @Router /api/login-password [post]
func (s *service) Do(c echo.Context) error {
	rootPath := echo_utils.GetMyRootPath(c)

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &login_models.LoginPasswordRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: "bind error",
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	if err := s.validateLoginPasswordRequest(model); err != nil {
		log.Error().Err(err).Msg("validateLoginPasswordRequest")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: "invalid request",
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	session, err := s.getSession()
	if err != nil {
		log.Error().Err(err).Msg("getSession")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_001,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("session.Get")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_002,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.Email,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			response := &login_models.LoginPasswordErrorResponse{
				Email:  model.Email,
				Reason: "user not found",
			}
			return c.JSONPretty(http.StatusNotFound, response, "  ")
		}
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_003,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	user := getRageUserResponse.User
	doEmailVerification := func(purpose contracts_cookies.VerifyCodePurpose) (string, error) {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             model.Email,
					Code:              verificationCode,
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: purpose,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return "", err
		}
		s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.Email,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification.message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})

		return verificationCode, nil

	}
	doPasswordVerification := func() error {
		if fluffycore_utils.IsNil(user.Password) {
			return status.Error(codes.NotFound, "Password not found")
		}
		if fluffycore_utils.IsEmptyOrNil(user.Password.Hash) {
			return status.Error(codes.NotFound, "Password hash not found")
		}
		return s.passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password:       model.Password,
			HashedPassword: user.Password.Hash,
		})
	}
	err = doPasswordVerification()
	if err != nil {
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: "wrong_password",
		}
		return c.JSONPretty(http.StatusUnauthorized, response, "  ")
	}
	// check if multi factor is required
	// ---------------------------------
	if user.TOTP == nil {
		user.TOTP = &proto_oidc_models.TOTP{
			Enabled: false,
		}
	}
	if (s.config.EmailVerificationRequired && !user.RootIdentity.EmailVerified) || s.config.MultiFactorRequiredByEmailCode {
		vCode, err := doEmailVerification(contracts_cookies.VerifyCode_Challenge)
		if err != nil {
			log.Error().Err(err).Msg("doEmailVerification")
			response := &login_models.LoginPasswordErrorResponse{
				Email:  model.Email,
				Reason: InternalError_LoginPassword_004,
			}
			return c.JSONPretty(http.StatusInternalServerError, response, "  ")
		}
		response := &login_models.LoginPasswordResponse{
			Email:     model.Email,
			Directive: login_models.DIRECTIVE_VerifyCode_DisplayVerifyCodePage,
		}
		if s.config.SystemConfig.DeveloperMode {
			response.DirectiveEmailCodeChallenge = &login_models.DirectiveEmailCodeChallenge{
				Code: vCode,
			}
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	oidcIdentity := &proto_oidc_models.OIDCIdentity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		Acr: []string{
			models.ACRPassword,
			models.ACRIdpRoot,
		},
		Amr: []string{
			models.AMRPassword,
			// always true, as we are the root idp
			models.AMRIdp,
		},
	}

	// we can process the final state now
	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       user.RootIdentity.Subject,
				Email:         user.RootIdentity.Email,
				EmailVerified: user.RootIdentity.EmailVerified,
			},
			Acr: oidcIdentity.Acr,
			Amr: oidcIdentity.Amr,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_005,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().
		GetAuthorizationRequestState(ctx,
			&proto_oidc_flows.GetAuthorizationRequestStateRequest{
				State: authorizationRequest.State,
			})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_006,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = oidcIdentity
	// "urn:rage:idp:google", "urn:rage:idp:spacex", "urn:rage:idp:github-enterprise", etc.
	// "urn:rage:password", "urn:rage:2fa", "urn:rage:email", etc.
	// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx,
		&proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationFinal.Request.Code,
			AuthorizationRequestState: authorizationFinal,
		})
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationRequestState")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_007,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	s.AuthorizationRequestStateStore().DeleteAuthorizationRequestState(ctx, &proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
		State: authorizationRequest.State,
	})
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx,
		&proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationRequest.State,
			AuthorizationRequestState: authorizationFinal,
		})
	if err != nil {
		// redirect to error page
		log.Error().Err(err).Msg("StoreAuthorizationRequestState")
		response := &login_models.LoginPasswordErrorResponse{
			Email:  model.Email,
			Reason: InternalError_LoginPassword_008,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	// redirect to the client with the code.
	redirectUri := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + rootPath
	response := &login_models.LoginPasswordResponse{
		Email:     model.Email,
		Directive: login_models.DIRECTIVE_Redirect,
		DirectiveRedirect: &login_models.DirectiveRedirect{
			RedirectURI: redirectUri,
		},
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
