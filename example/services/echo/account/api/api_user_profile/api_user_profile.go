package api_user_profile

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
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
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_external_user.IFluffyCoreUserServiceServer,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
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
	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	PhoneNumber string `json:"phoneNumber"`
}

// API UserIdentityInfo godoc
// @Summary get the highlevel UserIdentityInfo post login.
// @Description get the highlevel UserIdentityInfo post login.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} Profile
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/user-profile [post]
// @Router /api/user-profile [get]
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
func (s *service) DoPost(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	memCache := s.ScopedMemoryCache()
	cachedItem, err := memCache.Get("rootIdentity")
	if err != nil {
		log.Error().Err(err).Msg("memCache.Get")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity := cachedItem.(*proto_oidc_models.Identity)
	if rootIdentity == nil {
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
			return c.JSONPretty(http.StatusNotFound, err.Error(), "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}

	user := getUserResponse.User

	model := &Profile{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}

	_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
		User: &proto_external_models.ExampleUserUpdate{
			Id: user.Id,
			Profile: &proto_external_models.ProfileUpdate{
				FamilyName: &wrapperspb.StringValue{Value: model.FamilyName},
				GivenName:  &wrapperspb.StringValue{Value: model.GivenName},
				PhoneNumbers: []*types.PhoneNumberDTOUpdate{
					{
						Id:     "0",
						Number: &wrapperspb.StringValue{Value: model.PhoneNumber},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	return c.JSONPretty(http.StatusOK, model, "  ")

}
func (s *service) DoGet(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, err := memCache.Get("rootIdentity")
	if err != nil {
		log.Error().Err(err).Msg("memCache.Get")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity := cachedItem.(*proto_oidc_models.Identity)
	if rootIdentity == nil {
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
			return c.JSONPretty(http.StatusNotFound, err.Error(), "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}

	user := getUserResponse.User
	profileResponse := &Profile{}
	if fluffycore_utils.IsNotNil(user.Profile) {
		profileResponse.FamilyName = user.Profile.FamilyName
		profileResponse.GivenName = user.Profile.GivenName
		if fluffycore_utils.IsNotEmptyOrNil(user.Profile.PhoneNumbers) {
			profileResponse.PhoneNumber = user.Profile.PhoneNumbers[0].Number
		}
	}
	return c.JSONPretty(http.StatusOK, profileResponse, "  ")
}
