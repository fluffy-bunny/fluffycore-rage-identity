package base

import (
	"context"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	core_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
	core_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	BaseHandler struct {
		Container               di.Container
		Localizer               func() contracts_localizer.ILocalizer
		ClaimsPrincipal         func() fluffycore_contracts_common.IClaimsPrincipal
		EchoContextAccessor     func() fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
		IdpServiceServer        func() proto_oidc_idp.IFluffyCoreIDPServiceServer
		UserService             func() proto_oidc_user.IFluffyCoreUserServiceServer
		OIDCFlowStore           func() contracts_eko_gocache.IOIDCFlowStore
		ExternalOauth2FlowStore func() contracts_eko_gocache.IExternalOauth2FlowStore
		ScopedMemoryCache       func() fluffycore_contracts_common.IScopedMemoryCache
		EmailService            func() contracts_email.IEmailService

		localizer               contracts_localizer.ILocalizer
		claimsPrincipal         fluffycore_contracts_common.IClaimsPrincipal
		echoContextAccessor     fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
		idpServiceServer        proto_oidc_idp.IFluffyCoreIDPServiceServer
		userService             proto_oidc_user.IFluffyCoreUserServiceServer
		oidcFlowStore           contracts_eko_gocache.IOIDCFlowStore
		externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore
		scopedMemoryCache       fluffycore_contracts_common.IScopedMemoryCache
		emailService            contracts_email.IEmailService
	}
)

func NewBaseHandler(container di.Container) *BaseHandler {

	obj := &BaseHandler{Container: container}
	obj.Localizer = obj.getLocalizer
	obj.ClaimsPrincipal = obj.getClaimsPrincipal
	obj.EchoContextAccessor = obj.getEchoContextAccessor
	obj.IdpServiceServer = obj.getIdpServiceServer
	obj.UserService = obj.getUserService
	obj.OIDCFlowStore = obj.getOIDCFlowStore
	obj.ExternalOauth2FlowStore = obj.getExternalOauth2FlowStore
	obj.ScopedMemoryCache = obj.getScopedMemoryCache
	obj.EmailService = obj.getEmailService
	return obj

}
func (b *BaseHandler) getEmailService() contracts_email.IEmailService {
	if b.emailService == nil {
		b.emailService = di.Get[contracts_email.IEmailService](b.Container)
	}
	return b.emailService
}
func (b *BaseHandler) getScopedMemoryCache() fluffycore_contracts_common.IScopedMemoryCache {
	if b.scopedMemoryCache == nil {
		b.scopedMemoryCache = di.Get[fluffycore_contracts_common.IScopedMemoryCache](b.Container)
	}
	return b.scopedMemoryCache
}

func (b *BaseHandler) getLocalizer() contracts_localizer.ILocalizer {
	if b.localizer == nil {
		b.localizer = di.Get[contracts_localizer.ILocalizer](b.Container)
	}
	return b.localizer
}
func (b *BaseHandler) getClaimsPrincipal() fluffycore_contracts_common.IClaimsPrincipal {
	if b.claimsPrincipal == nil {
		b.claimsPrincipal = di.Get[fluffycore_contracts_common.IClaimsPrincipal](b.Container)
	}
	return b.claimsPrincipal
}
func (b *BaseHandler) getEchoContextAccessor() fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor {
	if b.echoContextAccessor == nil {
		b.echoContextAccessor = di.Get[fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor](b.Container)
	}
	return b.echoContextAccessor
}
func (b *BaseHandler) getIdpServiceServer() proto_oidc_idp.IFluffyCoreIDPServiceServer {
	if b.idpServiceServer == nil {
		b.idpServiceServer = di.Get[proto_oidc_idp.IFluffyCoreIDPServiceServer](b.Container)
	}
	return b.idpServiceServer
}
func (b *BaseHandler) getUserService() proto_oidc_user.IFluffyCoreUserServiceServer {
	if b.userService == nil {
		b.userService = di.Get[proto_oidc_user.IFluffyCoreUserServiceServer](b.Container)
	}
	return b.userService
}
func (b *BaseHandler) getOIDCFlowStore() contracts_eko_gocache.IOIDCFlowStore {
	if b.oidcFlowStore == nil {
		b.oidcFlowStore = di.Get[contracts_eko_gocache.IOIDCFlowStore](b.Container)
	}
	return b.oidcFlowStore
}
func (b *BaseHandler) getExternalOauth2FlowStore() contracts_eko_gocache.IExternalOauth2FlowStore {
	if b.externalOauth2FlowStore == nil {
		b.externalOauth2FlowStore = di.Get[contracts_eko_gocache.IExternalOauth2FlowStore](b.Container)
	}
	return b.externalOauth2FlowStore
}

func (b *BaseHandler) RenderAutoPost(c echo.Context, action string, formData []models.FormParam) error {
	data := map[string]interface{}{
		"form_params": formData,
		"action":      action,
	}
	return b.Render(c, http.StatusFound, "oidc/autopost/index", data)
}

func (b *BaseHandler) Render(c echo.Context, code int, name string, data map[string]interface{}) error {
	localizer := b.Localizer().GetLocalizer()
	data["LocalizeMessage"] = func(key string) string {
		message, _ := localizer.LocalizeMessage(&i18n.Message{ID: key})
		return message
	}
	data["isAuthenticated"] = func() bool {
		if b.ClaimsPrincipal == nil {
			return false
		}
		isAuthenticated := b.ClaimsPrincipal().HasClaimType(core_wellknown.ClaimTypeAuthenticated)
		return isAuthenticated
	}
	data["getUsername"] = func() string {
		claims := b.ClaimsPrincipal().GetClaimsByType("email")
		if len(claims) > 0 {
			return claims[0].Value
		}
		return "Account"
	}
	data["paths"] = wellknown_echo.NewPaths()
	data["username"] = "Account"
	if b.ClaimsPrincipal != nil {
		data["claims"] = b.ClaimsPrincipal().GetClaims()
		claims := b.ClaimsPrincipal().GetClaimsByType("email")
		if len(claims) > 0 {
			data["username"] = claims[0].Value
		}
	}

	return core_echo_templates.Render(c, code, name, data)

}

func (b *BaseHandler) GetIDPs(ctx context.Context) ([]*proto_oidc_models.IDP, error) {
	listIDPResponse, err := b.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			Enabled: &proto_types.BoolFilterExpression{
				Eq: true,
			},
			Hidden: &proto_types.BoolFilterExpression{
				Eq: false,
			},
		},
	})
	if err != nil {

		return nil, err
	}
	return listIDPResponse.Idps, nil
}
