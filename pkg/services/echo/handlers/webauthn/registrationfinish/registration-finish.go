package registrationfinish

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
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	webauthn_protocol "github.com/go-webauthn/webauthn/protocol"
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
			contracts_handler.POST,
		},
		wellknown_echo.WebAuthN_Register_Finish,
	)

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_WebAuthN_RegisterFinish_001 = "rg-webAuthN-RF-001"
	InternalError_WebAuthN_RegisterFinish_002 = "rg-webAuthN-RF-002"
	InternalError_WebAuthN_RegisterFinish_003 = "rg-webAuthN-RF-003"
	InternalError_WebAuthN_RegisterFinish_004 = "rg-webAuthN-RF-004"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	log.Info().Msg("WebAuthN_Register_Finish")

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

	// double match.
	//  the subecjt in the claims principal MUST match the one in the current session

	getWebAuthNCookieResponse, err := s.cookies.GetWebAuthNCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetWebAuthNCookie")
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	if getWebAuthNCookieResponse.Value.Identity.Subject != subject {
		log.Error().Msg("subject mismatch")
		s.cookies.DeleteWebAuthNCookie(c)
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	sessionData := getWebAuthNCookieResponse.Value.SessionData
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
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_RegisterFinish_001)
	}
	webAuthNUser := services_handlers_webauthn.NewWebAuthNUser(getRageUserResponse.User)

	body := c.Request().Body

	response, err := webauthn_protocol.ParseCredentialCreationResponseBody(body)
	if err != nil {
		log.Error().Err(err).Msg("ParseCredentialCreationResponseBody")
		return c.JSON(http.StatusInternalServerError, err)
	}
	webAuthN := s.webAuthN.GetWebAuthN()

	credential, err := webAuthN.CreateCredential(webAuthNUser, *sessionData, response)
	if err != nil {
		log.Error().Err(err).Msg("CreateCredential")
		return c.JSON(http.StatusInternalServerError, InternalError_WebAuthN_RegisterFinish_002)
	}

	// we need to add the credentials to the user.
	// TODO: Add creds to user database
	return c.JSON(http.StatusOK, credential)

}
