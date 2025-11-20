package loginfinish

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	services_handlers_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	protocol "github.com/go-webauthn/webauthn/protocol"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		webAuthN         contracts_webauthn.IWebAuthN
		wellknownCookies contracts_cookies.IWellknownCookies
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	webAuthN contracts_webauthn.IWebAuthN,
	cookies contracts_cookies.IWellknownCookies,
	oidcSession contracts_oidc_session.IOIDCSession,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		webAuthN:         webAuthN,
		wellknownCookies: cookies,
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
		wellknown_echo.WebAuthN_Login_Finish,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_WebAuthN_LoginFinish_001 = "rg-webAuthN-LF-001"
	InternalError_WebAuthN_LoginFinish_002 = "rg-webAuthN-LF-002"
	InternalError_WebAuthN_LoginFinish_003 = "rg-webAuthN-LF-003"
	InternalError_WebAuthN_LoginFinish_004 = "rg-webAuthN-LF-004"
	InternalError_WebAuthN_LoginFinish_005 = "rg-webAuthN-LF-005"
	InternalError_WebAuthN_LoginFinish_006 = "rg-webAuthN-LF-006"
	InternalError_WebAuthN_LoginFinish_007 = "rg-webAuthN-LF-007"
	InternalError_WebAuthN_LoginFinish_008 = "rg-webAuthN-LF-008"
	InternalError_WebAuthN_LoginFinish_009 = "rg-webAuthN-LF-009"
	InternalError_WebAuthN_LoginFinish_010 = "rg-webAuthN-LF-010"

	InternalError_WebAuthN_LoginFinish_099 = "rg-webAuthN-LF-099"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}

type SucessResonseJson struct {
	RedirectUrl string                  `json:"redirectUrl"`
	Credential  *go_webauthn.Credential `json:"credential"`
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	rootPath := echo_utils.GetMyRootPath(c)

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	log.Debug().Msg("WebAuthN_Login_Finish")
	getWebAuthNCookieResponse, err := s.wellknownCookies.GetWebAuthNCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetWebAuthNCookie")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_001)
	}
	sessionData := getWebAuthNCookieResponse.Value.SessionData
	// get the user from the store
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: getWebAuthNCookieResponse.Value.Identity.Subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				return c.JSON(http.StatusNotFound, "User not found")
			}
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_002)
	}
	user := getRageUserResponse.User
	webAuthNUser := services_handlers_webauthn.NewWebAuthNUser(user)
	body := r.Body
	parsedCredentialAssertionData, err := protocol.ParseCredentialRequestResponseBody(body)
	if err != nil {
		// Handle Error and return.
		log.Error().Err(err).Msg("ParseCredentialRequestResponseBody")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_003)
	}
	credential, err := s.webAuthN.GetWebAuthN().ValidateLogin(webAuthNUser, *sessionData, parsedCredentialAssertionData)
	if err != nil {
		// Handle Error and return.
		log.Error().Err(err).Msg("ValidateLogin")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_004)
	}
	session, err := s.getSession()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_005)
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_006)
	}

	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().
		GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
			State: authorizationRequest.State,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_007)
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		Acr: []string{
			models.ACRPasskey,
			models.ACRIdpRoot,
		},
		Amr: []string{
			models.AMRPasskey,
			// always true, as we are the root idp
			models.AMRIdp,
		},
	}
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationFinal.Request.Code,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		log.Error().Err(err).Msg("StoreAuthorizationRequestState")
		// redirect to error page
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_008)
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
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_009)
	}

	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       user.RootIdentity.Subject,
				Email:         user.RootIdentity.Email,
				EmailVerified: user.RootIdentity.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginFinish_005)
	}
	// redirect to the client with the code.
	redirectUrl := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + rootPath
	successResponse := &SucessResonseJson{
		RedirectUrl: redirectUrl,
		Credential:  credential,
	}
	return c.JSON(http.StatusOK, successResponse)
}
