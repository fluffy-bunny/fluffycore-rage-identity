package models

// DiscoveryDocument is the OpenID Connect Provider Metadata document served at
// /.well-known/openid-configuration.
// Spec: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type DiscoveryDocument struct {
	// Issuer is the base URL of the identity provider. Must match the iss claim
	// in issued tokens exactly.
	Issuer string `json:"issuer,omitempty"`

	// JwksURI is the URL of the JSON Web Key Set containing the public keys used
	// to verify id_token signatures.
	JwksURI string `json:"jwks_uri,omitempty"`

	// AuthorizationEndpoint is the URL of the OAuth 2.0 authorization endpoint
	// where clients begin the authorization code flow.
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`

	// TokenEndpoint is the URL of the OAuth 2.0 token endpoint where
	// authorization codes and refresh tokens are exchanged for tokens.
	TokenEndpoint string `json:"token_endpoint,omitempty"`

	// UserinfoEndpoint is the URL of the UserInfo endpoint. Clients may call it
	// with a bearer access token to retrieve claims about the authenticated user.
	UserinfoEndpoint string `json:"userinfo_endpoint,omitempty"`

	// EndSessionEndpoint is the URL of the RP-Initiated Logout endpoint.
	// Clients redirect the browser here (or load it in a hidden iframe for
	// front-channel logout) to clear the SSO session on the identity server.
	// Spec: https://openid.net/specs/openid-connect-rpinitiated-1_0.html
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`

	// CheckSessionIframe is the URL of an iframe endpoint used for session
	// management via the OpenID Connect Session Management spec.
	// Spec: https://openid.net/specs/openid-connect-session-1_0.html
	CheckSessionIframe string `json:"check_session_iframe,omitempty"`

	// RevocationEndpoint is the URL of the OAuth 2.0 token revocation endpoint.
	// Spec: https://www.rfc-editor.org/rfc/rfc7009
	RevocationEndpoint string `json:"revocation_endpoint,omitempty"`

	// IntrospectionEndpoint is the URL of the OAuth 2.0 token introspection
	// endpoint used by resource servers to validate tokens.
	// Spec: https://www.rfc-editor.org/rfc/rfc7662
	IntrospectionEndpoint string `json:"introspection_endpoint,omitempty"`

	// DeviceAuthorizationEndpoint is the URL of the OAuth 2.0 Device Authorization
	// endpoint for the device code flow.
	// Spec: https://www.rfc-editor.org/rfc/rfc8628
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint,omitempty"`

	// FrontchannelLogoutSupported indicates that this provider supports
	// front-channel logout: clients may load EndSessionEndpoint in a hidden
	// iframe to clear the SSO cookie without a full-page redirect.
	// Spec: https://openid.net/specs/openid-connect-frontchannel-1_0.html
	FrontchannelLogoutSupported bool `json:"frontchannel_logout_supported,omitempty"`

	// FrontchannelLogoutSessionSupported indicates that the provider passes a
	// sid (session ID) parameter in the front-channel logout request so clients
	// can match it to a specific session. Currently false — sid is not validated.
	FrontchannelLogoutSessionSupported bool `json:"frontchannel_logout_session_supported,omitempty"`

	// BackchannelLogoutSupported indicates that this provider supports
	// back-channel logout via a server-to-server POST with a logout token.
	// Spec: https://openid.net/specs/openid-connect-backchannel-1_0.html
	BackchannelLogoutSupported bool `json:"backchannel_logout_supported,omitempty"`

	// BackchannelLogoutSessionSupported indicates that the provider includes a
	// sid claim in back-channel logout tokens.
	BackchannelLogoutSessionSupported bool `json:"backchannel_logout_session_supported,omitempty"`

	// ScopesSupported lists the OAuth 2.0 scope values this provider supports.
	ScopesSupported []string `json:"scopes_supported,omitempty"`

	// ClaimsSupported lists the claim names this provider may include in tokens
	// and the UserInfo response.
	ClaimsSupported []string `json:"claims_supported,omitempty"`

	// GrantTypesSupported lists the OAuth 2.0 grant type values this provider
	// supports.
	GrantTypesSupported []string `json:"grant_types_supported,omitempty"`

	// ResponseTypesSupported lists the OAuth 2.0 response_type values this
	// provider supports at the authorization endpoint.
	ResponseTypesSupported []string `json:"response_types_supported,omitempty"`

	// ResponseModesSupported lists the OAuth 2.0 response_mode values this
	// provider supports (e.g. "query", "fragment", "form_post").
	ResponseModesSupported []string `json:"response_modes_supported,omitempty"`

	// TokenEndpointAuthMethodsSupported lists the client authentication methods
	// supported at the token endpoint (e.g. "client_secret_basic").
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`

	// SubjectTypesSupported lists the subject identifier types supported
	// (e.g. "public", "pairwise").
	SubjectTypesSupported []string `json:"subject_types_supported,omitempty"`

	// IDTokenSigningAlgValuesSupported lists the JWS signing algorithms supported
	// for the id_token (e.g. "RS256", "ES256").
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`

	// CodeChallengeMethodsSupported lists the PKCE code challenge methods
	// supported at the authorization endpoint (e.g. "S256").
	// Spec: https://www.rfc-editor.org/rfc/rfc7636
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported,omitempty"`
}
