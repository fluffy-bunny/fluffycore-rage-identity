package api_user_identity_info

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_user_identity_info"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
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

		webAuthNConfig *contracts_webauthn.WebAuthNConfig
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_UserIdentityInfo_001 = "rg-useridentityinfo-001"
	InternalError_UserIdentityInfo_002 = "rg-useridentityinfo-002"
	InternalError_UserIdentityInfo_003 = "rg-useridentityinfo-003"
	InternalError_UserIdentityInfo_004 = "rg-useridentityinfo-004"
	InternalError_UserIdentityInfo_005 = "rg-useridentityinfo-005"
	InternalError_UserIdentityInfo_006 = "rg-useridentityinfo-006"
	InternalError_UserIdentityInfo_007 = "rg-useridentityinfo-007"
	InternalError_UserIdentityInfo_008 = "rg-useridentityinfo-008"
	InternalError_UserIdentityInfo_009 = "rg-useridentityinfo-009"
	InternalError_UserIdentityInfo_010 = "rg-useridentityinfo-010"
	InternalError_UserIdentityInfo_011 = "rg-useridentityinfo-011"
	InternalError_UserIdentityInfo_099 = "rg-useridentityinfo-099"
)

func (s *service) Ctor(
	container di.Container,
	webAuthNConfig *contracts_webauthn.WebAuthNConfig,
) (*service, error) {
	return &service{
		BaseHandler:    services_echo_handlers_base.NewBaseHandler(container),
		webAuthNConfig: webAuthNConfig,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_UserIdentityInfo,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API UserIdentityInfo godoc
// @Summary get the highlevel UserIdentityInfo post login.
// @Description get the highlevel UserIdentityInfo post login.
// @Tags root
// @Produce json
// @Success 200 {object} api_user_identity_info.UserIdentityInfo
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/user-identity-info [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

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

	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, err.Error(), "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}

	user := getRageUserResponse.User

	// passkey is only allowed for direct username/password users.
	response := &api_user_identity_info.UserIdentityInfo{
		Email: user.RootIdentity.Email,
	}

	if fluffycore_utils.IsNotNil(user.LinkedIdentities) {
		response.LinkedIdentities = make([]api_user_identity_info.LinkedIdentity, 0)
		for _, linkedIdentity := range user.LinkedIdentities.Identities {
			response.LinkedIdentities = append(response.LinkedIdentities, api_user_identity_info.LinkedIdentity{
				Name: linkedIdentity.IdpSlug,
			})
		}
	}
	if s.webAuthNConfig.Enabled {
		response.PasskeyEligible = (user.Password != nil && fluffycore_utils.IsNotEmptyOrNil(user.Password.Hash))
		if fluffycore_utils.IsNotNil(user.WebAuthN) {
			response.Passkeys = make([]api_user_identity_info.Passkey, 0)
			for _, webAuthNCreds := range user.WebAuthN.Credentials {
				name := "Unknown"
				if fluffycore_utils.IsNotEmptyOrNil(webAuthNCreds.Authenticator.FriendlyName) {
					name = webAuthNCreds.Authenticator.FriendlyName
				}
				aaGUID := base64.StdEncoding.EncodeToString(webAuthNCreds.Authenticator.AAGUID)
				response.Passkeys = append(response.Passkeys, api_user_identity_info.Passkey{
					AAGUID: aaGUID,
					Name:   name,
				})
			}
		}
	}
	return c.JSONPretty(http.StatusOK, response, "  ")
}
