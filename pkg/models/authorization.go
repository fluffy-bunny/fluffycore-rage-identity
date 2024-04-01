package models

import (
	"encoding/gob"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
)

func init() {
	gob.Register(&proto_oidc_models.AuthorizationRequest{})
	gob.Register(&proto_oidc_models.ExternalOauth2Request{})
	gob.Register(&proto_oidc_models.OIDCIdentity{})
	gob.Register(&proto_oidc_models.AuthorizationRequestState{})
	gob.Register(&proto_oidc_models.ExternalOauth2State{})
	gob.Register(&FormParam{})

}

const (
	RootIdp string = "root-idp"
)
const (
	LoginDirective           string = "login"
	SignupDirective          string = "signup"
	PasswordResetDirective   string = "password-reset"
	VerifyEmailDirective     string = "verify-email"
	MFA_VerifyEmailDirective string = "mfa-verify-email"

	CancelDirective string = "cancel"
)
const (
	InternalError            string = "internal-error"
	ExternalIDPNotLinked     string = "external-idp-not-linked"
	UsernamePasswordNotFound string = "username-password-not-found"
	IdentityFound            string = "identity-found"
)

// urn prefixes
const (
	URNIdpPrefix     string = "urn:rage:idp:{idp_hint}"
	URLRootCandidate string = "urn:rage:root_candidate:{user_id}"
)

const (
	OIDCSessionName = "_oidc_session"
)

type (
	FormParam struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)
