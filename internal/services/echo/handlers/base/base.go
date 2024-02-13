package base

import (
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/localizer"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	core_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
	core_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	BaseHandler struct {
		Localizer           contracts_localizer.ILocalizer
		ClaimsPrincipal     fluffycore_contracts_common.IClaimsPrincipal
		EchoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
	}
)

func (b BaseHandler) Render(c echo.Context, code int, name string, data map[string]interface{}) error {
	localizer := b.Localizer.GetLocalizer()
	data["LocalizeMessage"] = func(key string) string {
		loginMsg, _ := localizer.LocalizeMessage(&i18n.Message{ID: key})
		return loginMsg
	}
	data["isAuthenticated"] = func() bool {
		if b.ClaimsPrincipal == nil {
			return false
		}
		isAuthenticated := b.ClaimsPrincipal.HasClaimType(core_wellknown.ClaimTypeAuthenticated)
		return isAuthenticated
	}
	data["getUsername"] = func() string {
		claims := b.ClaimsPrincipal.GetClaimsByType("email")
		if len(claims) > 0 {
			return claims[0].Value
		}
		return "Account"
	}
	data["paths"] = wellknown_echo.NewPaths()
	data["username"] = "Account"
	if b.ClaimsPrincipal != nil {
		data["claims"] = b.ClaimsPrincipal.GetClaims()
		claims := b.ClaimsPrincipal.GetClaimsByType("email")
		if len(claims) > 0 {
			data["username"] = claims[0].Value
		}
	}

	return core_echo_templates.Render(c, code, name, data)

}
