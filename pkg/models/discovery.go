package models

type DiscoveryDocument struct {
	Issuer                             string   `json:"issuer,omitempty"`
	JwksURI                            string   `json:"jwks_uri,omitempty"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint,omitempty"`
	TokenEndpoint                      string   `json:"token_endpoint,omitempty"`
	UserinfoEndpoint                   string   `json:"userinfo_endpoint,omitempty"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint,omitempty"`
	CheckSessionIframe                 string   `json:"check_session_iframe,omitempty"`
	RevocationEndpoint                 string   `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint,omitempty"`
	DeviceAuthorizationEndpoint        string   `json:"device_authorization_endpoint,omitempty"`
	FrontchannelLogoutSupported        bool     `json:"frontchannel_logout_supported,omitempty"`
	FrontchannelLogoutSessionSupported bool     `json:"frontchannel_logout_session_supported,omitempty"`
	BackchannelLogoutSupported         bool     `json:"backchannel_logout_supported,omitempty"`
	BackchannelLogoutSessionSupported  bool     `json:"backchannel_logout_session_supported,omitempty"`
	ScopesSupported                    []string `json:"scopes_supported,omitempty"`
	ClaimsSupported                    []string `json:"claims_supported,omitempty"`
	GrantTypesSupported                []string `json:"grant_types_supported,omitempty"`
	ResponseTypesSupported             []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported             []string `json:"response_modes_supported,omitempty"`
	TokenEndpointAuthMethodsSupported  []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	SubjectTypesSupported              []string `json:"subject_types_supported,omitempty"`
	IDTokenSigningAlgValuesSupported   []string `json:"id_token_signing_alg_values_supported,omitempty"`
	CodeChallengeMethodsSupported      []string `json:"code_challenge_methods_supported,omitempty"`
}
