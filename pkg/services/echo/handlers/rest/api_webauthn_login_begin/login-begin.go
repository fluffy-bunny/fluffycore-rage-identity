package api_webauthn_login_begin

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	services_handlers_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
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
		webAuthN contracts_webauthn.IWebAuthN
		cookies  contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	webAuthN contracts_webauthn.IWebAuthN,
	cookies contracts_cookies.IWellknownCookies,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		webAuthN:    webAuthN,
		cookies:     cookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.WebAuthN_Login_Begin,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_WebAuthN_LoginBegin_001 = "rg-webAuthN-LB-001"
	InternalError_WebAuthN_LoginBegin_002 = "rg-webAuthN-LB-002"
	InternalError_WebAuthN_LoginBegin_003 = "rg-webAuthN-LB-003"
	InternalError_WebAuthN_LoginBegin_004 = "rg-webAuthN-LB-004"
	InternalError_WebAuthN_LoginBegin_005 = "rg-webAuthN-LB-005"
	InternalError_WebAuthN_LoginBegin_006 = "rg-webAuthN-LB-006"
	InternalError_WebAuthN_LoginBegin_007 = "rg-webAuthN-LB-007"
	InternalError_WebAuthN_LoginBegin_008 = "rg-webAuthN-LB-008"
	InternalError_WebAuthN_LoginBegin_009 = "rg-webAuthN-LB-009"
	InternalError_WebAuthN_LoginBegin_010 = "rg-webAuthN-LB-010"

	InternalError_WebAuthN_LoginBegin_099 = "rg-webAuthN-LB-099"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Try to get the signin cookie, but don't fail if it doesn't exist
	// This supports both flows:
	// 1. User enters email first (cookie exists) - specific user authentication
	// 2. User clicks passkey directly (no cookie) - discoverable credential authentication
	signinResponse, err := s.cookies.GetSigninUserNameCookie(c)

	var webAuthNUser *services_handlers_webauthn.WebAuthNUser
	var userIdentity *proto_oidc_models.Identity
	var isDiscoverableFlow bool

	if err != nil {
		// No signin cookie - use discoverable credentials (resident keys)
		log.Info().Msg("No signin cookie found, using discoverable credentials flow")
		isDiscoverableFlow = true
		userIdentity = nil
	} else {
		// Signin cookie exists - get the specific user
		log.Info().Str("email", signinResponse.Value.Email).Msg("Signin cookie found, authenticating specific user")
		isDiscoverableFlow = false

		// get the user from the store
		getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
			&proto_oidc_user.GetRageUserRequest{
				By: &proto_oidc_user.GetRageUserRequest_Email{
					Email: signinResponse.Value.Email,
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
			return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_002)
		}
		webAuthNUser = services_handlers_webauthn.NewWebAuthNUser(getRageUserResponse.User)
		userIdentity = getRageUserResponse.User.RootIdentity
	}

	var credentialAssertion *protocol.CredentialAssertion
	var webAuthNSession *go_webauthn.SessionData

	if isDiscoverableFlow {
		// For discoverable credentials, use BeginDiscoverableLogin which doesn't require a user
		var err error
		credentialAssertion, webAuthNSession, err = s.webAuthN.GetWebAuthN().BeginDiscoverableLogin()
		if err != nil {
			log.Error().Err(err).Msg("BeginDiscoverableLogin")
			return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_003)
		}
	} else {
		// For specific user, use regular BeginLogin
		var err error
		credentialAssertion, webAuthNSession, err = s.webAuthN.GetWebAuthN().BeginLogin(webAuthNUser)
		if err != nil {
			log.Error().Err(err).Msg("BeginLogin")
			return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_003)
		}
	}
	cookieValue := &contracts_cookies.WebAuthNCookie{
		Identity:    userIdentity, // Will be nil for discoverable credentials
		SessionData: webAuthNSession,
	}
	err = s.cookies.SetWebAuthNCookie(c, &contracts_cookies.SetWebAuthNCookieRequest{
		Value: cookieValue,
	})
	if err != nil {
		log.Error().Err(err).Msg("SetWebAuthNCookie")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_004)
	}
	return c.JSON(http.StatusOK, credentialAssertion)
}
