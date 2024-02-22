package models

const (
	RootIdp string = "root-idp"
)
const (
	LoginDirective         string = "login"
	SignupDirective        string = "signup"
	PasswordResetDirective string = "password-reset"
	VerifyEmailDirective   string = "verify-email"
)
const (
	InternalError            string = "internal-error"
	ExternalIDPNotLinked     string = "external-idp-not-linked"
	UsernamePasswordNotFound string = "username-password-not-found"
	IdentityFound            string = "identity-found"
)

// urn prefixes
const (
	URNIdpPrefix     string = "urn:mastodon:idp:{idp_hint}"
	URLRootCandidate string = "urn:mastodon:root_candidate:{user_id}"
)

type (
	AuthorizationRequest struct {
		ClientId            string `param:"client_id" query:"client_id" form:"client_id" json:"client_id" xml:"client_id"`
		ResponseType        string `param:"response_type" query:"response_type" form:"response_type" json:"response_type" xml:"response_type"`
		Scope               string `param:"scope" query:"scope" form:"scope" json:"scope" xml:"scope"`
		State               string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
		RedirectURI         string `param:"redirect_uri" query:"redirect_uri" form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri"`
		Audience            string `param:"audience" query:"audience" form:"audience" json:"audience" xml:"audience"`
		CodeChallenge       string `param:"code_challenge" query:"code_challenge" form:"code_challenge" json:"code_challenge" xml:"code_challenge"`
		CodeChallengeMethod string `param:"code_challenge_method" query:"code_challenge_method" form:"code_challenge_method" json:"code_challenge_method" xml:"code_challenge_method"`
		ACRValues           string `param:"acr_values" query:"acr_values" form:"acr_values" json:"acr_values" xml:"acr_values"`
		Nonce               string `param:"nonce" query:"nonce" form:"nonce" json:"nonce" xml:"nonce"`
		Code                string // this is the internal code that will be returned to the OIDC client
		// IDPHint is the idp_hint of the external idp that the authorization must authenticate against
		IDPHint string
		// CandidateUserID is the user_id of the candidate user that if the external IDP has no link should be linked to
		// The candidate user must exist.
		CandidateUserID string
	}

	ExternalOauth2Request struct {
		IDPHint               string `json:"idp_hint,omitempty"`
		ClientID              string `json:"client_id,omitempty"`
		CodeChallenge         string `json:"code_challenge,omitempty"`
		CodeChallengeMethod   string `json:"code_challenge_method,omitempty"`
		State                 string `json:"state,omitempty"`
		CodeChallengeVerifier string `json:"code_challenge_verifier,omitempty"`
		Nonce                 string `json:"nonce,omitempty"`
		Directive             string `json:"directive,omitempty"`
		ParentState           string `json:"parent_state,omitempty"`
	}
	Identity struct {
		Subject       string
		Email         string
		ACR           []string
		AMR           []string
		EmailVerified bool
	}
	AuthorizationFinal struct {
		Request          *AuthorizationRequest
		Identity         *Identity
		ExternalIdentity *Identity
		Directive        string
	}
	ExternalOauth2Final struct {
		Request  *ExternalOauth2Request
		Identity *Identity
	}
)
