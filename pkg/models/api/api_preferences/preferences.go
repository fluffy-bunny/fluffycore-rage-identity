package api_preferences

type (
	GetKeepSignedInPreferenceResponse struct {
		HasPreference bool `json:"hasPreference"`
	}

	UpdateKeepSignedInPreferenceRequest struct {
		SkipKeepSignedInPage bool `json:"skipKeepSignedInPage"`
	}

	UpdateKeepSignedInPreferenceResponse struct {
		Success bool `json:"success"`
	}
)
