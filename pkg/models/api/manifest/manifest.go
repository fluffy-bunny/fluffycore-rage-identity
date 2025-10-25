package manifest

type (
	Page string
	IDP  struct {
		Slug string `json:"slug"`
	}
	Manifest struct {
		SocialIdps                  []IDP        `json:"social_idps"`
		PasskeyEnabled              bool         `json:"passkey_enabled"`
		LandingPage                 *LandingPage `json:"landing_page"`
		DevelopmentMode             bool         `json:"development_mode"`
		DisableLocalAccountCreation bool         `json:"disable_local_account_creation"`
		DisableSocialAccounts       bool         `json:"disable_social_accounts"`
	}
	LandingPage struct {
		Page Page `json:"page"`
	}
)

const (
	Login         Page = "Login"
	VerifyCode    Page = "VerifyCode"
	CreateAccount Page = "CreateAccount"
)
