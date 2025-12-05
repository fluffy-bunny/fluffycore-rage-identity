package account_management_server

import (
	"fmt"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	models_api "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
)

var csrfSkipperPaths map[string]bool

func CSRFSkipperPaths() map[string]bool {
	// needs to be a func as some of these are configured in.
	if csrfSkipperPaths == nil {
		csrfSkipperPaths = map[string]bool{
			wellknown_echo.StaticPath:                      true,
			wellknown_echo.WellKnownJWKS:                   true,
			wellknown_echo.ErrorPath:                       true,
			wellknown_echo.HealthzPath:                     true,
			wellknown_echo.ReadyPath:                       true,
			wellknown_echo.SwaggerPath:                     true,
			wellknown_echo.WellKnownOpenIDCOnfiguationPath: true,
			wellknown_echo.OAuth2TokenEndpointPath:         true,
			wellknown_echo.UserInfoPath:                    true,
			wellknown_echo.API_AppSettings:                 true,
			wellknown_echo.API_Manifest:                    true,
			wellknown_echo.API_StartOver:                   true,
			wellknown_echo.API_VerifyCodeBegin:             true,
			wellknown_echo.API_UserIdentityInfo:            true,
			wellknown_echo.API_UserProfilePath:             true,
			wellknown_echo.OIDCLoginUIStaticPath:           true,
		}
	}
	return csrfSkipperPaths
}

var requiresNoAuthPaths map[string]bool

// everything requries auth unless otherwise documented here.
// -- this is a list of paths that do not require auth
func RequiresNoAuth() map[string]bool {
	// needs to be a func as some of these are configured in.
	if requiresNoAuthPaths == nil {
		requiresNoAuthPaths = map[string]bool{
			wellknown_echo.API_AppConfig:       true,
			wellknown_echo.API_Login:           true,
			wellknown_echo.AccountCallbackPath: true,
		}
	}
	return requiresNoAuthPaths
}

// EnsureAuth ...
func EnsureAuth(_ di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get path
			path := c.Path()

			//ctx := c.Request().Context()
			subContainer, ok := c.Get(fluffycore_echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			if !ok {
				return next(c)
			}

			if _, ok := RequiresNoAuth()[path]; ok {
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
			if path == wellknown_echo.HomePath {
				redirectUrl := fmt.Sprintf("%s?returnUrl=%s", wellknown_echo.LoginPath, path)
				return c.Redirect(http.StatusFound, redirectUrl)
			}

			// return StatusUnauthorized
			return c.JSON(http.StatusUnauthorized, models_api.UnautorizedResponse{Path: path})
		}
	}
}
