package OIDCFlowAppConfig

type (
	IDP struct {
		Slug string `json:"slug"`
	}
	OIDCFlowAppConfig struct {
		SocialIdps                  []IDP `json:"social_idps"`
		PasskeyEnabled              bool  `json:"passkey_enabled"`
		DevelopmentMode             bool  `json:"development_mode"`
		DisableLocalAccountCreation bool  `json:"disable_local_account_creation"`
		DisableSocialAccounts       bool  `json:"disable_social_accounts"`
	}
)
