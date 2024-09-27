package manifest

type (
	Page string
	IDP  struct {
		Slug string `json:"slug"`
	}
	Manifest struct {
		SocialIdps     []IDP        `json:"social_idps"`
		PasskeyEnabled bool         `json:"passkey_enabled"`
		LandingPage    *LandingPage `json:"landing_page"`
	}
	LandingPage struct {
		Code string `json:"code"`
		Page Page   `json:"page"`
	}
)

const (
	Login         Page = "Login"
	VerifyCode    Page = "VerifyCode"
	CreateAccount Page = "CreateAccount"
)
