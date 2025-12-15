package api_user_passkeys

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
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
	}
)

var stemService = (*service)(nil)
var _ contracts_handler.IHandler = stemService

const (
	InternalError_Passkeys_001 = "rg-passkeys-001"
	InternalError_Passkeys_002 = "rg-passkeys-002"
	InternalError_Passkeys_003 = "rg-passkeys-003"
	InternalError_Passkeys_004 = "rg-passkeys-004"
	InternalError_Passkeys_005 = "rg-passkeys-005"
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_Passkeys,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type PasskeyResponse struct {
	ID           string   `json:"id"`
	FriendlyName string   `json:"friendlyName"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	Transport    []string `json:"transport,omitempty"`
}

type PasskeysListResponse struct {
	Passkeys []PasskeyResponse `json:"passkeys"`
}

// API GetPasskeys godoc
// @Summary get user's registered passkeys.
// @Description get user's registered passkeys.
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} PasskeysListResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/passkeys [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Get the user subject from claims principal
	claimsPrincipal := s.ClaimsPrincipal()
	subjectClaims := claimsPrincipal.GetClaimsByType(fluffycore_echo_wellknown.ClaimTypeSubject)
	if fluffycore_utils.IsEmptyOrNil(subjectClaims) {
		return c.JSONPretty(http.StatusUnauthorized, wellknown_echo.RestErrorResponse{Error: "Unauthorized"}, "  ")
	}
	claim := subjectClaims[0]
	if fluffycore_utils.IsEmptyOrNil(claim.Value) {
		return c.JSONPretty(http.StatusUnauthorized, wellknown_echo.RestErrorResponse{Error: "Unauthorized"}, "  ")
	}
	subject := claim.Value

	// Get the user from the store
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: "User not found"}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_Passkeys_001}, "  ")
	}

	user := getRageUserResponse.User
	var passkeys []PasskeyResponse

	if user.WebAuthN != nil && user.WebAuthN.Credentials != nil {
		for _, cred := range user.WebAuthN.Credentials {
			passkey := PasskeyResponse{
				ID:        base64.RawURLEncoding.EncodeToString(cred.ID),
				Transport: cred.Transport,
			}

			if cred.Authenticator != nil {
				passkey.FriendlyName = cred.Authenticator.FriendlyName
			}

			passkeys = append(passkeys, passkey)
		}
	}

	return c.JSONPretty(http.StatusOK, PasskeysListResponse{Passkeys: passkeys}, "  ")
}
