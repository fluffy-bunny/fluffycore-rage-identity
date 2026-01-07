package api_user_profile

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	pkg_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_external_user.IFluffyCoreUserServiceServer,
	config *contracts_config.Config,
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
			contracts_handler.POST,
		},
		wellknown_echo.API_UserProfilePath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type Profile struct {
	Email           string `json:"email"` // not editable
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	PhoneNumber     string `json:"phoneNumber"`
	IsClaimedDomain bool   `json:"isClaimedDomain"`
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}

// API GetUserProfile godoc
// @Summary get user profile.
// @Description get user profile.
// @Tags root
// @Accept 	json
// @Produce json
// @Param		request body		Profile	true	"Profile"
// @Success 200 {object} Profile
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 404 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/user-profile [post]
func (s *service) DoPost(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("rootIdentity not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}

	// get the user
	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	user := getUserResponse.User

	profile := &Profile{
		Email: rootIdentity.Email,
	}
	if err := c.Bind(profile); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
		User: &proto_external_models.ExampleUserUpdate{
			Id: user.Id,
			Profile: &proto_external_models.ProfileUpdate{
				FamilyName: &wrapperspb.StringValue{Value: profile.FamilyName},
				GivenName:  &wrapperspb.StringValue{Value: profile.GivenName},
				PhoneNumbers: []*types.PhoneNumberDTOUpdate{
					{
						Id:     "0",
						Number: &wrapperspb.StringValue{Value: profile.PhoneNumber},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	return c.JSONPretty(http.StatusOK, profile, "  ")

}

// API PutUserProfile godoc
// @Summary set user profile.
// @Description set user profile.
// @Tags root
// @Produce json
// @Success 200 {object} Profile
// @Failure 401 {string} string
// @Failure 404 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/user-profile [get]
func (s *service) DoGet(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("rootIdentity not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}

	// get the user
	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	user := getUserResponse.User
	profileResponse := &Profile{
		Email: rootIdentity.Email,
	}
	if fluffycore_utils.IsNotNil(user.Profile) {
		profileResponse.FamilyName = user.Profile.FamilyName
		profileResponse.GivenName = user.Profile.GivenName
		if fluffycore_utils.IsNotEmptyOrNil(user.Profile.PhoneNumbers) {
			profileResponse.PhoneNumber = user.Profile.PhoneNumbers[0].Number
		}
	}

	// Check if user has claimed domain ACR from ClaimsPrincipal
	cp := s.ClaimsPrincipal()
	isClaimedDomain := false
	acrClaims := cp.GetClaimsByType("acr")
	log.Info().Interface("acrClaims", acrClaims).Int("count", len(acrClaims)).Msg("Checking ACR claims for claimed domain")

	for _, claim := range acrClaims {
		if claim.Value == pkg_models.ACRClaimedDomain {
			isClaimedDomain = true
			log.Info().Msg("Found claimed domain ACR")
			break
		}
	}

	profileResponse.IsClaimedDomain = isClaimedDomain
	log.Info().Bool("isClaimedDomain", isClaimedDomain).Msg("Profile response prepared")

	return c.JSONPretty(http.StatusOK, profileResponse, "  ")
}
