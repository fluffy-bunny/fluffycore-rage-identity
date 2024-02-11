package base

import (
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	core_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
	core_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
)

type (
	BaseHandler struct {
		ClaimsPrincipal     fluffycore_contracts_common.IClaimsPrincipal
		EchoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
	}
)

func (b BaseHandler) Render(c echo.Context, code int, name string, data map[string]interface{}) error {
	data["isAuthenticated"] = func() bool {
		if b.ClaimsPrincipal == nil {
			return false
		}
		return b.ClaimsPrincipal.HasClaimType(core_wellknown.ClaimTypeAuthenticated)
	}
	data["paths"] = wellknown_echo.NewPaths()
	if b.ClaimsPrincipal != nil {
		data["claims"] = b.ClaimsPrincipal.GetClaims()
	}

	return core_echo_templates.Render(c, code, name, data)

}
