package oidcserver

import (
	"fmt"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
)

type (
	AuthPath struct {
		Path string `json:"path"`
	}
)

// everything requries auth unless otherwise documented here.
// -- this is a list of paths that do not require auth
var (
	RequiresNoAuth = map[string]bool{
		wellknown_echo.AboutPath:                       true,
		wellknown_echo.AccountCallbackPath:             true,
		wellknown_echo.ErrorPath:                       true,
		wellknown_echo.ExternalIDPPath:                 true,
		wellknown_echo.ForgotPasswordPath:              true,
		wellknown_echo.HealthzPath:                     true,
		wellknown_echo.HomePath:                        true,
		wellknown_echo.LoginPath:                       true,
		wellknown_echo.LogoutPath:                      true,
		wellknown_echo.OAuth2CallbackPath:              true,
		wellknown_echo.OAuth2TokenEndpointPath:         true,
		wellknown_echo.OIDCAuthorizationEndpointPath:   true,
		wellknown_echo.OIDCLoginPath:                   true,
		wellknown_echo.OIDCLoginPasswordPath:           true,
		wellknown_echo.PasswordResetPath:               true,
		wellknown_echo.ReadyPath:                       true,
		wellknown_echo.SignupPath:                      true,
		wellknown_echo.SwaggerPath:                     true,
		wellknown_echo.UserInfoPath:                    true,
		wellknown_echo.VerifyCodePath:                  true,
		wellknown_echo.WellKnownJWKS:                   true,
		wellknown_echo.WellKnownOpenIDCOnfiguationPath: true,
	}
)

// EnsureAuth ...
func EnsureAuth(_ di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//ctx := c.Request().Context()
			subContainer, ok := c.Get(fluffycore_echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			if !ok {
				return next(c)
			}

			// get path
			path := c.Path()
			if _, ok := RequiresNoAuth[path]; ok {
				return next(c)
			}
			claimsPrincipal := di.Get[fluffycore_contracts_common.IClaimsPrincipal](subContainer)
			isAuthenticated := claimsPrincipal.HasClaim(fluffycore_contracts_common.Claim{
				Type:  fluffycore_echo_wellknown.ClaimTypeAuthenticated,
				Value: "true",
			})
			if isAuthenticated {
				return next(c)
			}
			// redirect to root
			redirectUrl := fmt.Sprintf("%s?returnUrl=%s", wellknown_echo.LoginPath, path)
			return c.Redirect(http.StatusFound, redirectUrl)
		}
	}
}
