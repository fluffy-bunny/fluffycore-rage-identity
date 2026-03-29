package passkey

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	management_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	pkg_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types/webauthn"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config    *contracts_config.Config
		appConfig *management_config.AppConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	appConfig *management_config.AppConfig,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		appConfig:   appConfig,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXManagementPasskeyPath,
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

func (s *service) getSubject() string {
	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return ""
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return ""
	}
	return rootIdentity.Subject
}

func (s *service) getPasskeys(c *echo.Context) ([]api_passkey.PasskeyItem, bool, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	subject := s.getSubject()
	if subject == "" {
		return nil, false, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("GetRageUser")
		return nil, false, err
	}

	user := getRageUserResponse.User
	var passkeys []api_passkey.PasskeyItem
	if user.WebAuthN != nil && user.WebAuthN.Credentials != nil {
		for _, cred := range user.WebAuthN.Credentials {
			pk := api_passkey.PasskeyItem{
				ID: base64.RawURLEncoding.EncodeToString(cred.ID),
			}
			if cred.Authenticator != nil {
				pk.FriendlyName = cred.Authenticator.FriendlyName
			}
			if cred.CreatedOn != nil {
				pk.CreatedAt = cred.CreatedOn.AsTime().Unix()
			}
			if cred.LastUsedOn != nil {
				pk.LastUsedAt = cred.LastUsedOn.AsTime().Unix()
			}
			passkeys = append(passkeys, pk)
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

	return passkeys, isClaimedDomain, nil
}

func (s *service) DoGet(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	passkeys, isClaimedDomain, err := s.getPasskeys(c)
	if err != nil {
		return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_something_went_wrong"),
		}))
	}

	renameID := c.QueryParam("rename")

	return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
		RenderContext:   rc,
		Passkeys:        passkeys,
		IsClaimedDomain: isClaimedDomain,
		EnabledWebAuthN: s.appConfig.EnabledWebAuthN,
		RenameID:        renameID,
	}))
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	action := c.FormValue("action")

	subject := s.getSubject()
	if subject == "" {
		return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	switch action {
	case "rename":
		credentialID := c.FormValue("credentialId")
		friendlyName := c.FormValue("friendlyName")

		credentialIDBytes, err := base64.RawURLEncoding.DecodeString(credentialID)
		if err != nil {
			log.Error().Err(err).Msg("DecodeString credentialId")
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}

		// Get the user to find and update the credential
		getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
			&proto_oidc_user.GetRageUserRequest{
				By: &proto_oidc_user.GetRageUserRequest_Subject{
					Subject: subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("GetRageUser")
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}

		user := getRageUserResponse.User
		found := false
		if user.WebAuthN != nil && user.WebAuthN.Credentials != nil {
			for _, cred := range user.WebAuthN.Credentials {
				if string(cred.ID) == string(credentialIDBytes) {
					if cred.Authenticator != nil {
						cred.Authenticator.FriendlyName = friendlyName
						cred.UpdatedOn = timestamppb.Now()
						found = true
						break
					}
				}
			}
		}

		if !found {
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}

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
			log.Error().Err(err).Msg("UpdateRageUser rename passkey")
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}
		return s.renderWithSuccess(c, rc, rc.L("mgmt_success"))

	case "delete":
		credentialID := c.FormValue("credentialId")

		credentialIDBytes, err := base64.RawURLEncoding.DecodeString(credentialID)
		if err != nil {
			log.Error().Err(err).Msg("DecodeString credentialId")
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}

		_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
			User: &proto_oidc_models.RageUserUpdate{
				RootIdentity: &proto_oidc_models.IdentityUpdate{
					Subject: subject,
				},
				WebAuthN: &proto_oidc_models.WebAuthNUpdate{
					Credentials: &proto_types_webauthn.CredentialArrayUpdate{
						Update: &proto_types_webauthn.CredentialArrayUpdate_Granular_{
							Granular: &proto_types_webauthn.CredentialArrayUpdate_Granular{
								RemoveAAGUIDs: [][]byte{credentialIDBytes},
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("UpdateRageUser delete passkey")
			return s.renderWithError(c, rc, rc.L("mgmt_something_went_wrong"))
		}
		return s.renderWithSuccess(c, rc, rc.L("mgmt_success"))
	}

	return s.renderDefault(c, rc)
}

func (s *service) renderWithError(c *echo.Context, rc *components.RenderContext, errMsg string) error {
	passkeys, isClaimedDomain, _ := s.getPasskeys(c)
	return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
		RenderContext:   rc,
		Passkeys:        passkeys,
		IsClaimedDomain: isClaimedDomain,
		EnabledWebAuthN: s.appConfig.EnabledWebAuthN,
		Error:           errMsg,
	}))
}

func (s *service) renderWithSuccess(c *echo.Context, rc *components.RenderContext, successMsg string) error {
	passkeys, isClaimedDomain, _ := s.getPasskeys(c)
	return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
		RenderContext:   rc,
		Passkeys:        passkeys,
		IsClaimedDomain: isClaimedDomain,
		EnabledWebAuthN: s.appConfig.EnabledWebAuthN,
		Success:         successMsg,
	}))
}

func (s *service) renderDefault(c *echo.Context, rc *components.RenderContext) error {
	passkeys, isClaimedDomain, _ := s.getPasskeys(c)
	return components.RenderNode(c, http.StatusOK, components.PasskeyPage(&components.PasskeyPageData{
		RenderContext:   rc,
		Passkeys:        passkeys,
		IsClaimedDomain: isClaimedDomain,
		EnabledWebAuthN: s.appConfig.EnabledWebAuthN,
	}))
}
