package clientauthorization

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	oauth2_server "github.com/go-oauth2/oauth2/v4/server"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

// AuthorizeOAuth2Client validates client_id and redirect_uri at the authorization
// endpoint (RFC 6749 §4.1.1).  No secret is required here — that happens at /token.
// On success, puts the client into the scoped memory cache (same key as AuthenticateOAuth2Client)
// so downstream handlers can read it consistently.
func AuthorizeOAuth2Client() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Request().Context()
			log := zerolog.Ctx(ctx).With().Str("middleware", "AuthorizeOAuth2Client").Logger()

			clientID := c.QueryParam("client_id")
			if clientID == "" {
				clientID = c.FormValue("client_id")
			}
			if clientID == "" {
				log.Warn().Msg("missing client_id")
				return c.Redirect(http.StatusFound, "/error?error=invalid_request&error_description=client_id+is+required")
			}

			redirectURI := c.QueryParam("redirect_uri")
			if redirectURI == "" {
				redirectURI = c.FormValue("redirect_uri")
			}
			if redirectURI == "" {
				log.Warn().Str("clientId", clientID).Msg("missing redirect_uri")
				return c.Redirect(http.StatusFound, "/error?error=invalid_request&error_description=redirect_uri+is+required")
			}

			scopedContainer := c.Get(echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			scopedMemoryCache := di.Get[contracts_cache.IScopedMemoryCache](scopedContainer)
			clientService := di.Get[proto_oidc_client.IFluffyCoreClientServiceServer](scopedContainer)

			getClientResponse, err := clientService.GetClient(ctx, &proto_oidc_client.GetClientRequest{
				ClientId: clientID,
			})
			if err != nil {
				log.Warn().Err(err).Str("clientId", clientID).Msg("client not found")
				return c.Redirect(http.StatusFound, "/error?error=unauthorized_client&error_description=client+not+found")
			}

			client := getClientResponse.Client
			// Validate redirect_uri against the registered list (exact match, RFC 6749 §3.1.2)
			validRedirect := false
			for _, allowed := range client.AllowedRedirectUris {
				if allowed == redirectURI {
					validRedirect = true
					break
				}
			}
			if !validRedirect {
				log.Warn().Str("clientId", clientID).Str("redirectUri", redirectURI).Msg("redirect_uri not registered")
				return c.Redirect(http.StatusFound, "/error?error=invalid_request&error_description=redirect_uri+mismatch")
			}

			scopedMemoryCache.Set("client", client)
			return next(c)
		}
	}
}

func AuthenticateOAuth2Client() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Request().Context()
			r := c.Request()
			scopedContainer := c.Get(echo_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			scopedMemoryCache := di.Get[contracts_cache.IScopedMemoryCache](scopedContainer)

			clientService := di.Get[proto_oidc_client.IFluffyCoreClientServiceServer](scopedContainer)
			clientID, clientSecret, err := oauth2_server.ClientBasicHandler(r)
			if err != nil {
				clientID, clientSecret, err = oauth2_server.ClientFormHandler(r)
			}
			if err != nil {
				return err
			}
			ValidateClientSecretResponse, err := clientService.ValidateClientSecret(ctx, &proto_oidc_client.ValidateClientSecretRequest{
				ClientId: clientID,
				Secret:   clientSecret,
			})
			if err != nil {
				return err
			}
			if !ValidateClientSecretResponse.Valid {
				return echo.ErrUnauthorized
			}
			getClientResponse, err := clientService.GetClient(ctx, &proto_oidc_client.GetClientRequest{
				ClientId: clientID,
			})
			if err != nil {
				return err
			}
			scopedMemoryCache.Set("client", getClientResponse.Client)
			return next(c)
		}
	}

}
