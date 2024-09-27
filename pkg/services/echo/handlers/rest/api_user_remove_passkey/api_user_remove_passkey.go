package api_user_remove_passkey

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_user_remove_passkey"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
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

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_UserIdentityInfo_001 = "rg-removepasskey-001"
	InternalError_UserIdentityInfo_002 = "rg-removepasskey-002"
	InternalError_UserIdentityInfo_003 = "rg-removepasskey-003"
	InternalError_UserIdentityInfo_004 = "rg-removepasskey-004"
	InternalError_UserIdentityInfo_005 = "rg-removepasskey-005"
	InternalError_UserIdentityInfo_006 = "rg-removepasskey-006"
	InternalError_UserIdentityInfo_007 = "rg-removepasskey-007"
	InternalError_UserIdentityInfo_008 = "rg-removepasskey-008"
	InternalError_UserIdentityInfo_009 = "rg-removepasskey-009"
	InternalError_UserIdentityInfo_010 = "rg-removepasskey-010"
	InternalError_UserIdentityInfo_011 = "rg-removepasskey-011"
	InternalError_UserIdentityInfo_099 = "rg-removepasskey-099"
)

func (s *service) Ctor(
	container di.Container,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_UserRemovePasskey,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
func (s *service) validateRemovePasskeyRequest(model *api_user_remove_passkey.RemovePasskeyRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.AAGUID) {
		return status.Error(codes.InvalidArgument, "model.Name is nil")
	}
	_, err := base64.StdEncoding.DecodeString(model.AAGUID)
	if err != nil {
		return status.Error(codes.InvalidArgument, "model.AAGUID is not base64")
	}
	return nil
}

// API UserIdentityInfo godoc
// @Summary get the highlevel UserIdentityInfo post login.
// @Description get the highlevel UserIdentityInfo post login.
// @Tags root
// @Produce json
// @Param		request body		api_user_remove_passkey.RemovePasskeyRequest	true	"RemovePasskeyRequest"
// @Success 200 {object} api_user_remove_passkey.RemovePasskeyResonse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {object} api_user_remove_passkey.RemovePasskeyResonse
// @Router /api/user-remove-passkey [post]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &api_user_remove_passkey.RemovePasskeyRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validateRemovePasskeyRequest(model); err != nil {
		log.Error().Err(err).Msg("validateRemovePasskeyRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	aaguid, _ := base64.StdEncoding.DecodeString(model.AAGUID)

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
	// passkey is only allowed for direct username/password users.
	response := &api_user_remove_passkey.RemovePasskeyResonse{
		AAGUID: model.AAGUID,
	}

	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		response.Error = InternalError_UserIdentityInfo_001
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	user := getRageUserResponse.User
	_, err = s.RageUserService().UpdateRageUser(ctx,
		&proto_oidc_user.UpdateRageUserRequest{
			User: &proto_oidc_models.RageUserUpdate{
				RootIdentity: &proto_oidc_models.IdentityUpdate{
					Subject: user.RootIdentity.Subject,
				},
				WebAuthN: &proto_oidc_models.WebAuthNUpdate{
					Credentials: &proto_types_webauthn.CredentialArrayUpdate{
						Update: &proto_types_webauthn.CredentialArrayUpdate_Granular_{
							Granular: &proto_types_webauthn.CredentialArrayUpdate_Granular{
								RemoveAAGUIDs: [][]byte{aaguid},
							},
						},
					},
				},
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("UpdateRageUser")
		response.Error = InternalError_UserIdentityInfo_002
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
