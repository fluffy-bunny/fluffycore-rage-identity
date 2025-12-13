package api_linked_accounts

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	rage_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		wellknownCookies contracts_cookies.IWellknownCookies
		userService      proto_external_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_UserIdentityInfo_001 = "rg-userprofile-001"
	InternalError_UserIdentityInfo_002 = "rg-userprofile-002"
	InternalError_UserIdentityInfo_003 = "rg-userprofile-003"
	InternalError_UserIdentityInfo_004 = "rg-userprofile-004"
	InternalError_UserIdentityInfo_005 = "rg-userprofile-005"
	InternalError_UserIdentityInfo_006 = "rg-userprofile-006"
	InternalError_UserIdentityInfo_007 = "rg-userprofile-007"
	InternalError_UserIdentityInfo_008 = "rg-userprofile-008"
	InternalError_UserIdentityInfo_009 = "rg-userprofile-009"
	InternalError_UserIdentityInfo_010 = "rg-userprofile-010"
	InternalError_UserIdentityInfo_011 = "rg-userprofile-011"
	InternalError_UserIdentityInfo_099 = "rg-userprofile-099"
)

func (s *service) Ctor(
	container di.Container,
	config *rage_contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_external_user.IFluffyCoreUserServiceServer,

) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		wellknownCookies: wellknownCookies,
		userService:      userService,
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
		"/api/linked-accounts",
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	// is the request get or delete?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodDelete:
		return s.DoDelete(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}

// DoDelete godoc
// @Summary Delete a linked account
// @Description Removes a linked identity from the authenticated user's account
// @Tags account
// @Accept json
// @Produce json
// @Param identity query string true "Identity name to unlink"
// @Success 200 {object} models.DeleteLinkedAccountResponse "Successfully deleted linked account"
// @Failure 400 {object} string "Invalid request or missing identity parameter"
// @Failure 401 {object} string "Unauthorized - user not authenticated"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal server error"
// @Security CookieAuth
// @Router /api/linked-accounts [delete]
func (s *service) DoDelete(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Get identity from query parameter
	identity := c.QueryParam("identity")
	if identity == "" {
		log.Error().Msg("identity parameter is required")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"error": "identity parameter is required"}, "  ")
	}

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("memCache.Get")
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, "  ")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, "  ")
	}
	userService := s.userService

	// Get the user
	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"error": "user not found"}, "  ")
		}
		log.Error().Err(err).Msg("GetUser")
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, "  ")
	}

	user := getUserResponse.User

	// Find and remove the identity
	var foundIdentity *proto_oidc_models.Identity
	if user.RageUser != nil && user.RageUser.LinkedIdentities != nil && fluffycore_utils.IsNotEmptyOrNil(user.RageUser.LinkedIdentities.Identities) {
		for _, linkedIdentity := range user.RageUser.LinkedIdentities.Identities {
			if linkedIdentity.Subject == identity {
				foundIdentity = linkedIdentity
				break
			}
		}
	}

	if foundIdentity == nil {
		log.Warn().Str("identity", identity).Msg("Identity not found in user's linked accounts")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"error": "identity not found"}, "  ")
	}

	// Update user to remove the linked identity
	_, err = userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
		User: &proto_external_models.ExampleUserUpdate{
			Id: user.Id,
			RageUser: &proto_oidc_models.RageUserUpdate{
				LinkedIdentities: &proto_oidc_models.LinkedIdentitiesUpdate{
					Update: &proto_oidc_models.LinkedIdentitiesUpdate_Granular_{
						Granular: &proto_oidc_models.LinkedIdentitiesUpdate_Granular{
							Remove: []*proto_oidc_models.Identity{
								{
									Subject: foundIdentity.Subject,
									IdpSlug: foundIdentity.IdpSlug,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, "  ")
	}

	log.Info().Str("identity", identity).Msg("Successfully deleted linked account")
	return c.JSONPretty(http.StatusOK, &models.DeleteLinkedAccountResponse{Success: true}, "  ")
}

// DoGet godoc
// @Summary Get linked accounts
// @Description Retrieves the list of linked identities for the authenticated user
// @Tags account
// @Produce json
// @Success 200 {object} models.LinkedAccountsResponse "List of linked accounts"
// @Failure 401 {object} string "Unauthorized - user not authenticated"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal server error"
// @Security CookieAuth
// @Router /api/linked-accounts [get]
func (s *service) DoGet(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("memCache.Get")
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, "  ")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, "  ")
	}

	// Get the user
	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"error": "user not found"}, "  ")
		}
		log.Error().Err(err).Msg("GetUser")
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, "  ")
	}

	user := getUserResponse.User
	response := &models.LinkedAccountsResponse{
		Identities: []models.LinkedIdentity{},
	}

	// Extract linked identities from user
	if user.RageUser != nil && user.RageUser.LinkedIdentities != nil && fluffycore_utils.IsNotEmptyOrNil(user.RageUser.LinkedIdentities.Identities) {
		for _, identity := range user.RageUser.LinkedIdentities.Identities {
			// Get provider name from IdpSlug
			provider := identity.IdpSlug
			if provider != "" {
				// Capitalize first letter
				if provider[0] >= 'a' && provider[0] <= 'z' {
					provider = string(provider[0]-32) + provider[1:]
				}
			}

			// Get email from identity
			email := identity.Email

			// Format linked timestamp - TODO: add timestamp when field is available
			linkedAt := ""

			response.Identities = append(response.Identities, models.LinkedIdentity{
				Subject:  identity.Subject,
				Provider: provider,
				Email:    email,
				LinkedAt: linkedAt,
			})
		}
	}

	// Check if user has claimed domain ACR from ClaimsPrincipal
	cp := s.ClaimsPrincipal()
	isClaimedDomain := false
	acrClaims := cp.GetClaimsByType("acr")
	log.Info().Interface("acrClaims", acrClaims).Int("count", len(acrClaims)).Msg("Checking ACR claims from ClaimsPrincipal")

	for _, claim := range acrClaims {
		if claim.Value == "urn:rage:claimed-domain" {
			isClaimedDomain = true
			log.Info().Msg("Found claimed domain ACR!")
			break
		}
	}
	response.IsClaimedDomain = isClaimedDomain

	log.Info().Int("count", len(response.Identities)).Bool("isClaimedDomain", isClaimedDomain).Msg("Retrieved linked accounts")
	return c.JSONPretty(http.StatusOK, response, "  ")
}
