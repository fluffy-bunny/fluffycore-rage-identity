package api_login_current_user

import (
	"encoding/json"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
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
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_LoginCurrentUser,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_LoginCurrentUser_001 = "rg-logincurrentuser-001"
	InternalError_LoginCurrentUser_002 = "rg-logincurrentuser-002"
	InternalError_LoginCurrentUser_003 = "rg-logincurrentuser-003"
	InternalError_LoginCurrentUser_004 = "rg-logincurrentuser-004"
	InternalError_LoginCurrentUser_005 = "rg-logincurrentuser-005"
	InternalError_LoginCurrentUser_006 = "rg-logincurrentuser-006"
	InternalError_LoginCurrentUser_007 = "rg-logincurrentuser-007"
	InternalError_LoginCurrentUser_008 = "rg-logincurrentuser-008"
	InternalError_LoginCurrentUser_009 = "rg-logincurrentuser-009"
	InternalError_LoginCurrentUser_010 = "rg-logincurrentuser-010"
	InternalError_LoginCurrentUser_011 = "rg-logincurrentuser-011"
	InternalError_LoginCurrentUser_099 = "rg-logincurrentuser-099"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API Manifest godoc
// @Summary login current user.
// @Description login current user..
// @Tags root
// @Accept  json
// @Produce json
// @Success 200 {object} login_models.LoginCurrentUserResponse
// @Failure 401 {object} login_models.LoginCurrentUserErrorResponse
// @Failure 404 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/login-current-user [get]
func (s *service) Do(c echo.Context) error {
	rootPath := echo_utils.GetMyRootPath(c)

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	claimsPrincipal := s.ClaimsPrincipal()
	subjectClaims := claimsPrincipal.GetClaimsByType(fluffycore_echo_wellknown.ClaimTypeSubject)
	if fluffycore_utils.IsEmptyOrNil(subjectClaims) {
		response := &login_models.LoginCurrentUserErrorResponse{
			Reason: InternalError_LoginCurrentUser_001,
		}
		return c.JSON(http.StatusUnauthorized, response)
	}
	claim := subjectClaims[0]
	if fluffycore_utils.IsEmptyOrNil(claim.Value) {
		response := &login_models.LoginCurrentUserErrorResponse{
			Reason: InternalError_LoginCurrentUser_002,
		}
		return c.JSON(http.StatusUnauthorized, response)
	}
	subject := claim.Value

	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	user := getRageUserResponse.User

	session, err := s.getSession()
	if err != nil {
		log.Error().Err(err).Msg("getSession")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")

	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("session.Get")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	// check if multi factor is required
	// ---------------------------------
	if user.TOTP == nil {
		user.TOTP = &proto_oidc_models.TOTP{
			Enabled: false,
		}
	}
	acrClaims := claimsPrincipal.GetClaimsByType("acr")
	amrClaims := claimsPrincipal.GetClaimsByType("amr")

	acrSlice := []string{}
	err = json.Unmarshal([]byte(acrClaims[0].Value), &acrSlice)
	if err != nil {
		log.Error().Err(err).Msg("json.Unmarshal")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	amrSlice := []string{}
	err = json.Unmarshal([]byte(amrClaims[0].Value), &amrSlice)
	if err != nil {
		log.Error().Err(err).Msg("json.Unmarshal")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	// we can process the final state nowa
	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       user.RootIdentity.Subject,
				Email:         user.RootIdentity.Email,
				EmailVerified: user.RootIdentity.EmailVerified,
			},
			Acr: acrSlice,
			Amr: amrSlice,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().
		GetAuthorizationRequestState(ctx,
			&proto_oidc_flows.GetAuthorizationRequestStateRequest{
				State: authorizationRequest.State,
			})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	// pull the acr and amr from the user claims principal
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		Acr:     acrSlice,
		Amr:     amrSlice,
	}
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
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
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
		return c.JSONPretty(http.StatusInternalServerError,
			wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	// redirect to the client with the code.
	redirectUri := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + rootPath
	response := &login_models.LoginPasswordResponse{
		Email:     user.RootIdentity.Email,
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
