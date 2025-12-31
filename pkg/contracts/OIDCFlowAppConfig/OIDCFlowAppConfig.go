package OIDCFlowAppConfig

type (
	IDP struct {
		Slug string `json:"slug"`
	}
	OIDCFlowAppConfig struct {
		SocialIdps                  []IDP `json:"socialIdps"`
		PasskeyEnabled              bool  `json:"passkeyEnabled"`
		EnabledWebAuthN             bool  `json:"enabledWebAuthN"`
		EnabledTotp                 bool  `json:"enabledTotp"`
		DevelopmentMode             bool  `json:"developmentMode"`
		DisableLocalAccountCreation bool  `json:"disableLocalAccountCreation"`
		DisableSocialAccounts       bool  `json:"disableSocialAccounts"`
	}
)
