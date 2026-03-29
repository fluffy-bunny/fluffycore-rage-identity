package linked

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	pkg_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	api_linked "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_linked_identities"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
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
		wellknown_echo.HTMXManagementLinkedPath,
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

func (s *service) getLinkedData(c *echo.Context) (*components.LinkedPageData, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return &components.LinkedPageData{RenderContext: rc, Error: rc.L("mgmt_unexpected_error")}, nil
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return &components.LinkedPageData{RenderContext: rc, Error: rc.L("mgmt_unexpected_error")}, nil
	}

	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetUser")
		return &components.LinkedPageData{RenderContext: rc, Error: rc.L("mgmt_something_went_wrong")}, nil
	}

	user := getUserResponse.User
	var identities []api_linked.LinkedIdentity

	if user.RageUser != nil && user.RageUser.LinkedIdentities != nil &&
		fluffycore_utils.IsNotEmptyOrNil(user.RageUser.LinkedIdentities.Identities) {
		for _, identity := range user.RageUser.LinkedIdentities.Identities {
			provider := identity.IdpSlug
			listIDPResponse, err := s.IdpServiceServer().ListIDP(ctx,
				&proto_oidc_idp.ListIDPRequest{
					Filter: &proto_oidc_idp.Filter{
						Enabled: &proto_types.BoolFilterExpression{
							Eq: true,
						},
						Slug: &proto_types.StringFilterExpression{
							Eq: identity.IdpSlug,
						},
					},
				})
			if err != nil {
				log.Error().Err(err).Str("idpSlug", identity.IdpSlug).Msg("ListIDP")
			} else if listIDPResponse != nil && len(listIDPResponse.IDPs) > 0 {
				if fluffycore_utils.IsNotEmptyOrNil(listIDPResponse.IDPs[0].Name) {
					provider = listIDPResponse.IDPs[0].Name
				}
			}

			identities = append(identities, api_linked.LinkedIdentity{
				Subject:  identity.Subject,
				Provider: provider,
				Email:    identity.Email,
				CreatedOn: func() int64 {
					if identity.CreatedOn != nil {
						return identity.CreatedOn.AsTime().Unix()
					}
					return 0
				}(),
				LastUsedOn: func() int64 {
					if identity.LastUsedOn != nil {
						return identity.LastUsedOn.AsTime().Unix()
					}
					return 0
				}(),
			})
		}
	}

	// Check claimed domain
	cp := s.ClaimsPrincipal()
	isClaimedDomain := false
	acrClaims := cp.GetClaimsByType("acr")
	for _, claim := range acrClaims {
		if claim.Value == pkg_models.ACRClaimedDomain {
			isClaimedDomain = true
			break
		}
	}

	return &components.LinkedPageData{
		RenderContext:   rc,
		Identities:      identities,
		IsClaimedDomain: isClaimedDomain,
	}, nil
}

func (s *service) DoGet(c *echo.Context) error {
	data, _ := s.getLinkedData(c)
	return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(data))
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
		return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(&components.LinkedPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(&components.LinkedPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	switch action {
	case "unlink":
		identitySubject := c.FormValue("identity")

		getUserResponse, err := s.userService.GetUser(ctx,
			&proto_external_user.GetUserRequest{
				Subject: rootIdentity.Subject,
			})
		if err != nil {
			log.Error().Err(err).Msg("GetUser")
			return s.renderWithError(c, rc.L("mgmt_something_went_wrong"))
		}

		user := getUserResponse.User
		var foundIdentity *proto_oidc_models.Identity
		if user.RageUser != nil && user.RageUser.LinkedIdentities != nil {
			for _, li := range user.RageUser.LinkedIdentities.Identities {
				if li.Subject == identitySubject {
					foundIdentity = li
					break
				}
			}
		}

		if foundIdentity == nil {
			return s.renderWithError(c, rc.L("mgmt_something_went_wrong"))
		}

		_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
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
			log.Error().Err(err).Msg("UpdateUser unlink")
			return s.renderWithError(c, rc.L("mgmt_something_went_wrong"))
		}

		data, _ := s.getLinkedData(c)
		data.Success = rc.L("mgmt_success")
		return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(data))
	}

	data, _ := s.getLinkedData(c)
	return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(data))
}

func (s *service) renderWithError(c *echo.Context, errMsg string) error {
	data, _ := s.getLinkedData(c)
	data.Error = errMsg
	return components.RenderNode(c, http.StatusOK, components.LinkedAccountsPage(data))
}
