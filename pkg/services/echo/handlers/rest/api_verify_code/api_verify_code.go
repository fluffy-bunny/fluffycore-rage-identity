package api_verify_code

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
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

// 	API_VerifyCode             = "/api/verify-code"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_VerifyCode,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateVerifyCodeRequest(model *login_models.VerifyCodeRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Code) {
		return status.Error(codes.InvalidArgument, "model.Code is nil")
	}

	return nil
}

// API VerifyCode godoc
// @Summary verify code.
// @Description verify code
// @Tags root
// @Accept */*
// @Produce json
// @Param		request body		login_models.VerifyCodeRequest	true	"VerifyCodeRequest"
// @Success 200 {object} login_models.VerifyCodeResponse
// @Failure 401 {string} string
// @Router /api/verify-code [post]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &login_models.VerifyCodeRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if err := s.validateVerifyCodeRequest(model); err != nil {
		log.Error().Err(err).Msg("validateVerifyCodeRequest")
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}

	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetVerificationCodeCookie")
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	verificationCode := getVerificationCodeCookieResponse.VerificationCode
	code := verificationCode.Code
	if code != model.Code {
		return c.JSONPretty(http.StatusNotFound, "code does not match", "  ")
	}
	userService := s.RageUserService()
	getRageUserResponse, err := userService.GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: verificationCode.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("GetRageUser")
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, "User not found", "  ")
		}
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	rageUser := getRageUserResponse.User
	_, err = userService.UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: verificationCode.Subject,
				EmailVerified: &wrapperspb.BoolValue{
					Value: true,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	// one time only
	s.wellknownCookies.DeleteVerificationCodeCookie(c)

	switch verificationCode.VerifyCodePurpose {
	case contracts_cookies.VerifyCode_EmailVerification:
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	case contracts_cookies.VerifyCode_PasswordReset:
		err = s.wellknownCookies.SetPasswordResetCookie(c,
			&contracts_cookies.SetPasswordResetCookieRequest{
				PasswordReset: &contracts_cookies.PasswordReset{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_PassowrdReset_DisplayPasswordResetPage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	case contracts_cookies.VerifyCode_Challenge:
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_VerifyCode_Redirect,
		}
		err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
			AuthCookie: &contracts_cookies.AuthCookie{
				Identity: &proto_oidc_models.Identity{
					Subject:       rageUser.RootIdentity.Subject,
					Email:         rageUser.RootIdentity.Email,
					EmailVerified: rageUser.RootIdentity.EmailVerified,
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("SetAuthCookie")
			// redirect to error page
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		session, err := s.getSession()
		if err != nil {
			log.Error().Err(err).Msg("getSession")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		sessionRequest, err := session.Get("request")
		if err != nil {
			log.Error().Err(err).Msg("Get")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

		getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().
			GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
				State: authorizationRequest.State,
			})
		if err != nil {
			log.Error().Err(err).Msg("GetAuthorizationRequestState")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
		authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
			Subject: rageUser.RootIdentity.Subject,
			Email:   rageUser.RootIdentity.Email,
			Acr: []string{
				models.ACRPassword,
				models.ACRIdpRoot,
			},
			Amr: []string{
				models.AMRPassword,
				// always true, as we are the root idp
				models.AMRIdp,
				// this is a multifactor
				models.AMRMFA,
				models.AMREmailCode,
			},
		}
		// "urn:rage:idp:google", "urn:rage:idp:spacex", "urn:rage:idp:github-enterprise", etc.
		// "urn:rage:password", "urn:rage:2fa", "urn:rage:email", etc.
		// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
		_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationFinal.Request.Code,
			AuthorizationRequestState: authorizationFinal,
		})
		if err != nil {
			log.Error().Err(err).Msg("StoreAuthorizationRequestState")
			// redirect to error page
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		s.AuthorizationRequestStateStore().DeleteAuthorizationRequestState(ctx, &proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
			State: authorizationRequest.State,
		})
		_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationRequest.State,
			AuthorizationRequestState: authorizationFinal,
		})
		if err != nil {
			// redirect to error page
			log.Error().Err(err).Msg("StoreAuthorizationRequestState")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		rootPath := echo_utils.GetMyRootPath(c)

		// redirect to the client with the code.
		redirectUri := authorizationFinal.Request.RedirectUri +
			"?code=" + authorizationFinal.Request.Code +
			"&state=" + authorizationFinal.Request.State +
			"&iss=" + rootPath
		response.DirectiveRedirect = &login_models.DirectiveRedirect{
			VERB:        http.MethodGet,
			RedirectURI: redirectUri,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")

	}
	return c.JSONPretty(http.StatusInternalServerError, "Unknown VerifyCodePurpose", "  ")

}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
