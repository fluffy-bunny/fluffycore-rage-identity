package api_preferences

type (
	GetKeepSignedInPreferenceResponse struct {
		DoNotShowAgain bool `json:"doNotShowAgain"`
		KeepSignedIn   bool `json:"keepSignedIn"`
	}

	UpdateKeepSignedInPreferenceRequest struct {
		DoNotShowAgain bool `json:"doNotShowAgain"`
		KeepSignedIn   bool `json:"keepSignedIn"`
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
