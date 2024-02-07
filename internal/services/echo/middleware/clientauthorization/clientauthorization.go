package clientauthorization

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	oauth2_server "github.com/go-oauth2/oauth2/v4/server"
	echo "github.com/labstack/echo/v4"
)

func AuthenticateOAuth2Client() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			scopedContainer := c.Get(echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			scopedMemoryCache := di.Get[fluffycore_contracts_common.IScopedMemoryCache](scopedContainer)

			clientID, clientSecret, err := oauth2_server.ClientBasicHandler(r)
			if err != nil {
				clientID, clientSecret, err = oauth2_server.ClientFormHandler(r)
			}
			if err != nil {
				return err
			}
			scopedMemoryCache.Set("client_id", clientID)
			scopedMemoryCache.Set("client_secret", clientSecret)
			return next(c)
		}
	}

}
