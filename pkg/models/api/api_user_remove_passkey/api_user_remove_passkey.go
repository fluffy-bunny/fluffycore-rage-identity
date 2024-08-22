package api_user_remove_passkey

type (
	RemovePasskeyRequest struct {
		AAGUID string `json:"aaguid" validate:"required"`
	}

	RemovePasskeyResonse struct {
		AAGUID string `json:"aaguid"`
		Error  string `json:"error"`
	}
)
