package verify_code

type (
	VerifyCodeBeginResponse struct {
		// Code is supplied only if in development mode
		Code string `json:"code"`
		// Valid is true if we truely are doing a code verification
		Valid bool `json:"valid"`

		Email string `json:"email"`
	}
)
