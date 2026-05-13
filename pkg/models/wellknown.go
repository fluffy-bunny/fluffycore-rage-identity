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

// ACR (Authentication Context Class Reference) values communicated via acr_values
// in an authorization request. They declare the required authentication strength
// or behaviour to the identity server.
const (
	// ACRPassword requires the user to authenticate with a password.
	ACRPassword = "urn:rage:password"
	// ACRIdpRoot requires authentication against the root (local) identity provider.
	ACRIdpRoot = "urn:rage:idp:root"
	// ACR2FA requires a second factor of authentication in addition to the primary credential.
	ACR2FA = "urn:rage:loa:2fa"
	// ACRIdp is a template value; replace {idp} with the IDP slug to require
	// authentication through a specific external identity provider
	// (e.g. "urn:rage:loa:idp:google").
	ACRIdp = "urn:rage:loa:idp:{idp}"
	// ACRPasskey requires the user to authenticate with a passkey (WebAuthn).
	ACRPasskey = "urn:rage:loa:passkey"
	// ACRClaimedDomain requires the user to authenticate against a claimed domain.
	ACRClaimedDomain = "urn:rage:claimed-domain"
	// ACRNoSSO instructs the authorization endpoint to skip any existing SSO
	// cookie and require the user to authenticate fresh. Use this when the
	// relying party lives on a different domain and cannot clear the SSO cookie
	// itself (e.g. a white-labelled management portal).
	ACRNoSSO = "urn:rage:no-sso"
)

// AMR (Authentication Methods References) values recorded in the id_token amr
// claim. They describe the actual methods used during authentication.
const (
	// AMRPassword indicates the user authenticated with a password.
	AMRPassword = "pwd"
	// AMRMFA indicates the user completed multi-factor authentication.
	AMRMFA = "mfa"
	// AMRIdp indicates the user authenticated via an external identity provider.
	AMRIdp = "idp"
	// AMRPasskey indicates the user authenticated with a passkey (WebAuthn).
	AMRPasskey = "passkey"
	// AMRTOTP indicates the user authenticated with a TOTP authenticator app.
	AMRTOTP = "totp"
	// AMREmailCode indicates the user authenticated with an emailed verification code.
	AMREmailCode = "emailcode"
)

// Wellknown IDP slug identifiers used with ACRIdp and acr_values.
const (
	// WellknownIdpRoot is the built-in local (root) identity provider.
	WellknownIdpRoot = "root"
	// WellknownIdpGoogle is the Google external identity provider.
	WellknownIdpGoogle = "google"
	// WellknownIdpGithub is the GitHub external identity provider.
	WellknownIdpGithub = "github"
	// WellknownIdpMicrosoft is the Microsoft external identity provider.
	WellknownIdpMicrosoft = "microsoft"
	// WellknownIdpApple is the Apple external identity provider.
	WellknownIdpApple = "apple"
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
