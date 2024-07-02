package api_user_identity_info

type (
	UserIdentityInfo struct {
		Email           string `json:"email"`
		PasskeyEligible bool   `json:"passkeyEligible"`
	}
)
