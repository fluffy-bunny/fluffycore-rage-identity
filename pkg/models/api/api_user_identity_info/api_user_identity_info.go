package api_user_identity_info

type (
	Passkeys struct {
		Name string `json:"name"`
	}
	UserIdentityInfo struct {
		Email           string     `json:"email"`
		PasskeyEligible bool       `json:"passkeyEligible"`
		Passkeys        []Passkeys `json:"passkeys"`
	}
)
