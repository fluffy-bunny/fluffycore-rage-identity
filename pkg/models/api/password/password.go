package password

type (
	VerifyPasswordStrengthRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	VerifyPasswordStrengthResponse struct {
		Valid bool `json:"valid"`
	}
)
