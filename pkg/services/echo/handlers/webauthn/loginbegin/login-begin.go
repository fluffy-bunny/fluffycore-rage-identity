package loginbegin

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
	signinResponse, err := s.cookies.GetSigninUserNameCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetSigninUserNameCookie")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_001)
	}
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
	webAuthNUser := services_handlers_webauthn.NewWebAuthNUser(getRageUserResponse.User)
	credentialAssertion, webAuthNSession, err := s.webAuthN.GetWebAuthN().BeginLogin(webAuthNUser)
	if err != nil {
		// Handle Error and return.
		log.Error().Err(err).Msg("BeginLogin")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_LoginBegin_003)
	}
	cookieValue := &contracts_cookies.WebAuthNCookie{
		Identity:    getRageUserResponse.User.RootIdentity,
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
