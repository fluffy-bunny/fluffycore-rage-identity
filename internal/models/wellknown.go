package models

const (
	OAuth2TokenType_Bearer            = "bearer"
	OAuth2TokenType_JWT               = "urn:ietf:params:oauth:token-type:jwt"
	OAuth2TokenType_IDToken           = "urn:ietf:params:oauth:token-type:id_token"
	OAuth2TokenType_RefreshToken      = "urn:ietf:params:oauth:token-type:refresh_token"
	OAuth2TokenType_AccessToken       = "urn:urn:ietf:params:oauth:token-type:access_token"
	OAuth2GrantType_ClientCredentials = "client_credentials"
	OAuth2GrantType_RefreshToken      = "refresh_token"
	OAuth2GrantType_TokenExchange     = "urn:ietf:params:oauth:grant-type:token-exchange"
	OAUTH2GrantType_AuthorizationCode = "authorization_code"
)
