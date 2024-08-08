package api_user_identity_info

type (
	Identity struct {
		Name string `json:"name"`
	}
	UserLinkedAccounts struct {
		Identities []Identity `json:"identities"`
	}
)
