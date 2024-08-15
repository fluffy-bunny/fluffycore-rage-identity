package api_user_identity_info

type (
	Passkey struct {
		Name   string `json:"name"`
		AAGUID string `json:"aaguid"`
	}
	LinkedIdentity struct {
		Name string `json:"name"`
	}
	UserIdentityInfo struct {
		Email            string           `json:"email"`
		PasskeyEligible  bool             `json:"passkeyEligible"`
		Passkeys         []Passkey        `json:"passkeys"`
		LinkedIdentities []LinkedIdentity `json:"linkedIdentities"`
	}
)
