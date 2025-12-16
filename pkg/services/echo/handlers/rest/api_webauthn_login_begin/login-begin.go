package api_webauthn_login_begin

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	protocol "github.com/go-webauthn/webauthn/protocol"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
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

	// For passkey login, we should always use discoverable credentials
	// Delete any existing signin cookie to prevent conflicts
	s.cookies.DeleteSigninUserNameCookie(c)
	log.Info().Msg("Cleared signin cookie for passkey login")

	// Always use discoverable credentials for passkey login
	// This allows the user to select any passkey for any account
	log.Info().Msg("Using discoverable credentials flow")
	var userIdentity *proto_oidc_models.Identity = nil

	var credentialAssertion *protocol.CredentialAssertion
	var webAuthNSession *go_webauthn.SessionData

	// Always use BeginDiscoverableLogin for passkey authentication
	credentialAssertion, webAuthNSession, err := s.webAuthN.GetWebAuthN().BeginDiscoverableLogin()
	if err != nil {
		log.Error().Err(err).Msg("BeginDiscoverableLogin")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_003)
	}

	// Log the challenge being sent to the client
	log.Info().
		Str("challenge", string(webAuthNSession.Challenge)).
		Str("challenge_in_assertion", string(credentialAssertion.Response.Challenge)).
		Msg("Sending credential assertion to client")

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
