package verify_username

type (
	VerifyUsernameRequest struct {
		UserName string `json:"userName"`
	}
	VerifyUsernameResponse struct {
		UserName         string `json:"userName"`
		PasskeyAvailable bool   `json:"passkeyAvailable"`
	}
)
