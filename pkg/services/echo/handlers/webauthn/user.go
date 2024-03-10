package webauthn

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
)

type (
	WebAuthNUser struct {
		RageUser *proto_oidc_models.RageUser
	}
)

func init() {
	var _ go_webauthn.User = (*WebAuthNUser)(nil)
}

func NewWebAuthNUser(rageUser *proto_oidc_models.RageUser) *WebAuthNUser {
	return &WebAuthNUser{
		RageUser: rageUser,
	}
}

// WebAuthnID provides the user handle of the user account. A user handle is an opaque byte sequence with a maximum
// size of 64 bytes, and is not meant to be displayed to the user.
//
// To ensure secure operation, authentication and authorization decisions MUST be made on the basis of this id
// member, not the displayName nor name members. See Section 6.1 of [RFC8266].
//
// It's recommended this value is completely random and uses the entire 64 bytes.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dom-publickeycredentialuserentity-id)
func (s *WebAuthNUser) WebAuthnID() []byte {
	subject := s.RageUser.RootIdentity.Subject
	return []byte(subject)
}

// WebAuthnName provides the name attribute of the user account during registration and is a human-palatable name for the user
// account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party SHOULD let the user
// choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dictdef-publickeycredentialuserentity)
func (s *WebAuthNUser) WebAuthnName() string {
	return s.RageUser.RootIdentity.Email
}

// WebAuthnDisplayName provides the name attribute of the user account during registration and is a human-palatable
// name for the user account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party
// SHOULD let the user choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://www.w3.org/TR/webauthn/#dom-publickeycredentialuserentity-displayname)
func (s *WebAuthNUser) WebAuthnDisplayName() string {
	return s.RageUser.RootIdentity.Email
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (s *WebAuthNUser) WebAuthnCredentials() []go_webauthn.Credential {
	return nil
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (s *WebAuthNUser) WebAuthnIcon() string {
	return ""
}
