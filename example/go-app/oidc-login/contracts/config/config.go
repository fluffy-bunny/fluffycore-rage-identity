package config

type (
	BannerBranding struct {
		Title             string `json:"title,omitempty"`
		LogoURL           string `json:"logoUrl,omitempty"`
		ShowBannerVersion bool   `json:"showBannerVersion,omitempty"`
	}
	AppConfig struct {
		BaseHREF       string         `json:"basehref,omitempty"`
		RageBaseURL    string         `json:"rageBaseUrl,omitempty"`
		BannerBranding BannerBranding `json:"bannerBranding,omitempty"`
	}
)
