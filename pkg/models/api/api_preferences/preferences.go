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

	ClearSSOResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)
