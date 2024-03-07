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

const (
	ACRPassword = "urn:mastodon:password"
	ACRIdpRoot  = "urn:mastodon:idp:root"
	ACR2FA      = "urn:rage:loa:2fa"
	ACRIdp      = "urn:rage:loa:idp:{idp}"
)
const (
	AMRPassword = "pwd"
	AMR2FA      = "mfa"
	AMRIdp      = "idp"
)

const (
	WellknownIdpRoot      = "root"
	WellknownIdpGoogle    = "google"
	WellknownIdpGithub    = "github"
	WellknownIdpMicrosoft = "microsoft"
	WellknownIdpApple     = "apple"
)

const (
	ClaimTypeAcr   = "acr"
	ClaimTypeSub   = "sub"
	ClaimTypeIat   = "iat"
	ClaimTypeExp   = "exp"
	ClaimTypeNbf   = "nbf"
	ClaimTypeJti   = "jti"
	ClaimTypeIss   = "iss"
	ClainTypeAud   = "aud"
	ClaimTypeAzp   = "azp"
	ClaimTypeNonce = "nonce"
	ClaimTypeAmr   = "amr"
)
