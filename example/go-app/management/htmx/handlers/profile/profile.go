package profile

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	pkg_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config      *contracts_config.Config
		userService proto_external_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	userService proto_external_user.IFluffyCoreUserServiceServer,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		userService: userService,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXManagementProfilePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	// Non-HTMX GET requests (e.g. browser refresh) need the full shell page
	if c.Request().Method == http.MethodGet && !components.IsHTMXRequest(c) {
		return c.Redirect(http.StatusFound, wellknown_echo.HTMXManagementPath+"?redirect="+c.Request().URL.Path)
	}
	r := c.Request()
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	return c.NoContent(http.StatusNotFound)
}

func (s *service) getProfile(c *echo.Context) (*api_profile.Profile, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetUser")
		return nil, err
	}

	user := getUserResponse.User
	profile := &api_profile.Profile{
		Email:   rootIdentity.Email,
		Subject: rootIdentity.Subject,
	}
	if fluffycore_utils.IsNotNil(user.Profile) {
		profile.FamilyName = user.Profile.FamilyName
		profile.GivenName = user.Profile.GivenName
		if fluffycore_utils.IsNotEmptyOrNil(user.Profile.PhoneNumbers) {
			profile.PhoneNumber = user.Profile.PhoneNumbers[0].Number
		}
	}

	// Check claimed domain
	cp := s.ClaimsPrincipal()
	acrClaims := cp.GetClaimsByType("acr")
	for _, claim := range acrClaims {
		if claim.Value == pkg_models.ACRClaimedDomain {
			profile.IsClaimedDomain = true
			break
		}
	}

	return profile, nil
}

func (s *service) DoGet(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	profile, err := s.getProfile(c)
	if err != nil {
		return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
			RenderContext: rc,
			Profile:       &api_profile.Profile{},
			Error:         rc.L("mgmt_failed_load_profile"),
		}))
	}

	editing := c.QueryParam("edit")

	return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
		RenderContext: rc,
		Profile:       profile,
		Editing:       editing,
	}))
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	action := c.FormValue("action")

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
			RenderContext: rc,
			Profile:       &api_profile.Profile{},
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
			RenderContext: rc,
			Profile:       &api_profile.Profile{},
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetUser")
		return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
			RenderContext: rc,
			Profile:       &api_profile.Profile{},
			Error:         rc.L("mgmt_failed_load_profile"),
		}))
	}

	user := getUserResponse.User

	switch action {
	case "save-personal":
		givenName := c.FormValue("givenName")
		familyName := c.FormValue("familyName")
		_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
			User: &proto_external_models.ExampleUserUpdate{
				Id: user.Id,
				Profile: &proto_external_models.ProfileUpdate{
					GivenName:  &wrapperspb.StringValue{Value: givenName},
					FamilyName: &wrapperspb.StringValue{Value: familyName},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("UpdateUser")
			profile, _ := s.getProfile(c)
			if profile == nil {
				profile = &api_profile.Profile{}
			}
			return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
				RenderContext: rc,
				Profile:       profile,
				Error:         rc.L("mgmt_failed_save_profile"),
			}))
		}

	case "save-contact":
		phoneNumber := c.FormValue("phoneNumber")
		_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
			User: &proto_external_models.ExampleUserUpdate{
				Id: user.Id,
				Profile: &proto_external_models.ProfileUpdate{
					PhoneNumbers: []*proto_types.PhoneNumberDTOUpdate{
						{
							Id:     "0",
							Number: &wrapperspb.StringValue{Value: phoneNumber},
						},
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("UpdateUser")
			profile, _ := s.getProfile(c)
			if profile == nil {
				profile = &api_profile.Profile{}
			}
			return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
				RenderContext: rc,
				Profile:       profile,
				Error:         rc.L("mgmt_failed_save_profile"),
			}))
		}
	}

	// Reload and render with success
	profile, _ := s.getProfile(c)
	if profile == nil {
		profile = &api_profile.Profile{}
	}
	return components.RenderNode(c, http.StatusOK, components.ProfilePage(&components.ProfilePageData{
		RenderContext: rc,
		Profile:       profile,
		Success:       rc.L("mgmt_success"),
	}))
}
