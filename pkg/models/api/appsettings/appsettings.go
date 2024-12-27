package appsettings

type (
	OIDCUIAppSettings struct {
		ApplicationEnvironment string `json:"ApplicationEnvironment"`
		BaseApiUrl             string `json:"BaseApiUrl"`
	}
	AccountAppSettings struct {
		ApplicationEnvironment string `json:"ApplicationEnvironment"`
		BaseApiUrl             string `json:"BaseApiUrl"`
	}
	ApiAppSettings struct {
		ApplicationEnvironment string `json:"ApplicationEnvironment"`
		BaseApiUrl             string `json:"BaseApiUrl"`
		PrivacyPolicyUrl       string `json:"PrivacyPolicyUrl"`
		CookiePolicyUrl        string `json:"CookiePolicyUrl"`
	}
)
