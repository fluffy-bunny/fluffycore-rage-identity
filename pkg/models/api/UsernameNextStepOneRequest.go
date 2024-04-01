package api

type (
	UsernameNextStepOneRequest struct {
		BaseRequest
		Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"required"`
	}
	UsernameNextStepOneResponse struct {
		// UserExists is true if the user exists in the system.
		// Show the proper error to the user
		UserExists bool `json:"user_exists"`
		// NextPage to offer the user.
		// EnterPasswordPage
		NextPage string `json:"next_page,omitempty"`
		// IDPs is a list of IDPs that the user can use to login.
		IDPs []IDP `json:"idps,omitempty"`
		// RedirectUri is used to redirect the user to the next step.
		// This will happen if the username/email is claimed by an external IDP.
		RedirectUri string `json:"redirect_uri,omitempty"`
	}
)
