package loginfinish

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	services_handlers_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	protocol "github.com/go-webauthn/webauthn/protocol"
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

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	webAuthN contracts_webauthn.IWebAuthN,
	cookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),
		webAuthN:    webAuthN,
		cookies:     cookies,
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

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	log.Info().Msg("WebAuthN_Login_Finish")
	getWebAuthNCookieResponse, err := s.cookies.GetWebAuthNCookie(c)
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
	webAuthNUser := services_handlers_webauthn.NewWebAuthNUser(getRageUserResponse.User)
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

	return c.JSON(http.StatusOK, credential)
}
