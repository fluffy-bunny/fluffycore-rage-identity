package manifest

type (
	Page string

	CurrentManifestError int
)

const (
	PageLogin         Page = "Login"
	PageVerifyCode    Page = "VerifyCode"
	PageCreateAccount Page = "CreateAccount"
	PagePasswordEntry Page = "PasswordEntry"
	PageUsernameEntry Page = "UsernameEntry"
)

type (
	IDP struct {
		Slug string `json:"slug"`
	}

	Manifest struct {
		SessionId                   string       `json:"session_id"`
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
