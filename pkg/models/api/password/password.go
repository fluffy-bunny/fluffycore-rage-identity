package password

type (
	VerifyPasswordStrengthRequest struct {
		Password string `json:"password" validate:"required"`
	}
	VerifyPasswordStrengthResponse struct {
		Valid bool `json:"valid"`
	}
)
