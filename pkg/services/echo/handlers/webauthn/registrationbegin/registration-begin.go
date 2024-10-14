package registrationbegin

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	services_handlers_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/webauthn"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
		wellknown_echo.WebAuthN_Register_Begin,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_WebAuthN_RegisterBegin_001 = "rg-webAuthN-RB-001"
	InternalError_WebAuthN_RegisterBegin_002 = "rg-webAuthN-RB-002"
	InternalError_WebAuthN_RegisterBegin_003 = "rg-webAuthN-RB-003"
	InternalError_WebAuthN_RegisterBegin_004 = "rg-webAuthN-RB-004"
	InternalError_WebAuthN_RegisterBegin_005 = "rg-webAuthN-RB-005"
	InternalError_WebAuthN_RegisterBegin_006 = "rg-webAuthN-RB-006"
	InternalError_WebAuthN_RegisterBegin_007 = "rg-webAuthN-RB-007"
	InternalError_WebAuthN_RegisterBegin_008 = "rg-webAuthN-RB-008"
	InternalError_WebAuthN_RegisterBegin_009 = "rg-webAuthN-RB-009"
	InternalError_WebAuthN_RegisterBegin_010 = "rg-webAuthN-RB-010"

	InternalError_WebAuthN_RegisterBegin_099 = "rg-webAuthN-RB-099"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

/*
Requirments.
1. The user must be authenticated, and all that information is in the claims principal
2. Pull the subject and get the user from the store
3. Put the user in the WebAuthNUser wrapper, which pull the username/email to generate the challenge.
*/
func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// the the user subject from claims principal
	claimsPrincipal := s.ClaimsPrincipal()
	subjectClaims := claimsPrincipal.GetClaimsByType(fluffycore_echo_wellknown.ClaimTypeSubject)
	if fluffycore_utils.IsEmptyOrNil(subjectClaims) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	claim := subjectClaims[0]
	if fluffycore_utils.IsEmptyOrNil(claim.Value) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	subject := claim.Value

	log.Debug().Msg("WebAuthN_Register_Begin")

	// get the user from the store
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
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
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_RegisterBegin_001)
	}
	webAuthNUser := services_handlers_webauthn.NewWebAuthNUser(getRageUserResponse.User)
	credentialCreation, webAuthNSession, err := s.webAuthN.GetWebAuthN().BeginRegistration(webAuthNUser)
	if err != nil {
		log.Error().Err(err).Msg("BeginRegistration")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_RegisterBegin_002)
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
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_RegisterBegin_003)
	}
	// store the WebAuthNSession in a cookie.
	return c.JSON(http.StatusOK, credentialCreation)
}
