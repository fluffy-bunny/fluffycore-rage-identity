package api_user_linked_accounts

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	api_user_linked_accounts "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_user_linked_accounts"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
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

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_UserLinkedAccounts_001 = "rg-userlinkedaccounts-001"
	InternalError_UserLinkedAccounts_002 = "rg-userlinkedaccounts-002"
	InternalError_UserLinkedAccounts_003 = "rg-userlinkedaccounts-003"
	InternalError_UserLinkedAccounts_004 = "rg-userlinkedaccounts-004"
	InternalError_UserLinkedAccounts_005 = "rg-userlinkedaccounts-005"
	InternalError_UserLinkedAccounts_006 = "rg-userlinkedaccounts-006"
	InternalError_UserLinkedAccounts_007 = "rg-userlinkedaccounts-007"
	InternalError_UserLinkedAccounts_008 = "rg-userlinkedaccounts-008"
	InternalError_UserLinkedAccounts_009 = "rg-userlinkedaccounts-009"
	InternalError_UserLinkedAccounts_010 = "rg-userlinkedaccounts-010"
	InternalError_UserLinkedAccounts_011 = "rg-userlinkedaccounts-011"
	InternalError_UserLinkedAccounts_099 = "rg-userlinkedaccounts-099"
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
			contracts_handler.GET,
			contracts_handler.DELETE,
		},
		wellknown_echo.API_UserLinkedAccounts,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API UserLinkedAccounts godoc
// @Summary get the users linked accounts.
// @Description get the users linked accounts.
// @Tags root
// @Produce json
// @Success 200 {object} api_user_identity_info.UserLinkedAccounts
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/user-linked-accounts [get]
func (s *service) DoGet(c echo.Context) error {
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
	response := &api_user_linked_accounts.UserLinkedAccounts{}
	if fluffycore_utils.IsNotNil(user.LinkedIdentities) {
		response.Identities = make([]api_user_linked_accounts.Identity, 0)
		for _, identity := range user.LinkedIdentities.Identities {
			response.Identities = append(response.Identities, api_user_linked_accounts.Identity{
				Name: identity.IdpSlug,
			})
		}
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}

type DoDeleteRequest struct {
	Identity string `param:"identity" query:"identity" form:"identity" json:"identity" xml:"identity"`
}

// API UserLinkedAccounts Delete godoc
// @Summary delete a users linked identity.
// @Description delete a users linked identity.
// @Tags root
// @Produce json
// @Param        identity   path      string  true  "identity name"
// @Success 200 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/user-linked-accounts/{identity} [delete]
func (s *service) DoDelete(c echo.Context) error {
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

	model := &DoDeleteRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}
	log = log.With().Interface("model", model).Logger()

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

	var foundIdentity *models.Identity
	if fluffycore_utils.IsNotNil(getRageUserResponse.User.LinkedIdentities) {
		for _, identity := range getRageUserResponse.User.LinkedIdentities.Identities {
			isEqual := strings.EqualFold(identity.IdpSlug, model.Identity)
			if isEqual {
				foundIdentity = identity
				break
			}
		}
	}
	if fluffycore_utils.IsNil(foundIdentity) {
		return c.JSONPretty(http.StatusOK, "", "  ")
	}

	_, err = s.RageUserService().UnlinkRageUser(ctx,
		&proto_oidc_user.UnlinkRageUserRequest{
			RootSubject:      subject,
			ExternalIdentity: foundIdentity,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusOK, "", "  ")
		}
		log.Error().Err(err).Msg("UnlinkRageUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	return c.JSONPretty(http.StatusOK, "", "  ")
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodDelete:
		return s.DoDelete(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
