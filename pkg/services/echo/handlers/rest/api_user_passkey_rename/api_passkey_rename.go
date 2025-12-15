package api_user_passkey_rename

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types/webauthn"
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
	InternalError_PasskeyRename_001 = "rg-passkey-rename-001"
	InternalError_PasskeyRename_002 = "rg-passkey-rename-002"
	InternalError_PasskeyRename_003 = "rg-passkey-rename-003"
	InternalError_PasskeyRename_004 = "rg-passkey-rename-004"
	InternalError_PasskeyRename_005 = "rg-passkey-rename-005"
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
			contracts_handler.PATCH,
		},
		wellknown_echo.API_PasskeyRename,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type RenamePasskeyRequest struct {
	CredentialID string `param:"credentialId" json:"credentialId"`
	FriendlyName string `json:"friendlyName"`
}

// API RenamePasskey godoc
// @Summary rename a user's passkey.
// @Description update the friendly name of a user's passkey.
// @Tags account
// @Accept json
// @Produce json
// @Param credentialId path string true "Credential ID (base64 encoded)"
// @Param request body RenamePasskeyRequest true "Rename request"
// @Success 200 {object} wellknown_echo.RestSuccessResponse
// @Failure 400 {object} wellknown_echo.RestErrorResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 404 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/passkeys/{credentialId} [patch]
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

	// Bind the request
	req := &RenamePasskeyRequest{}
	if err := c.Bind(req); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Invalid request"}, "  ")
	}

	if fluffycore_utils.IsEmptyOrNil(req.CredentialID) {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Credential ID is required"}, "  ")
	}

	if fluffycore_utils.IsEmptyOrNil(req.FriendlyName) {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Friendly name is required"}, "  ")
	}

	// Decode the base64 credential ID
	credentialIDBytes, err := base64.RawURLEncoding.DecodeString(req.CredentialID)
	if err != nil {
		log.Error().Err(err).Str("credentialId", req.CredentialID).Msg("Failed to decode credential ID")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Invalid credential ID format"}, "  ")
	}

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
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_PasskeyRename_001}, "  ")
	}

	user := getRageUserResponse.User

	// Find the credential and update its friendly name
	found := false
	if user.WebAuthN != nil && user.WebAuthN.Credentials != nil {
		for _, cred := range user.WebAuthN.Credentials {
			if string(cred.ID) == string(credentialIDBytes) {
				if cred.Authenticator != nil {
					cred.Authenticator.FriendlyName = req.FriendlyName
					found = true
					break
				}
			}
		}
	}

	if !found {
		return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: "Passkey not found"}, "  ")
	}

	// Update the user with the modified credentials
	_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: subject,
			},
			WebAuthN: &proto_oidc_models.WebAuthNUpdate{
				Credentials: &proto_types_webauthn.CredentialArrayUpdate{
					Update: &proto_types_webauthn.CredentialArrayUpdate_Granular_{
						Granular: &proto_types_webauthn.CredentialArrayUpdate_Granular{
							Add: user.WebAuthN.Credentials,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_PasskeyRename_002}, "  ")
	}

	return c.JSONPretty(http.StatusOK, wellknown_echo.RestSuccessResponse{Message: "Passkey renamed successfully"}, "  ")
}
