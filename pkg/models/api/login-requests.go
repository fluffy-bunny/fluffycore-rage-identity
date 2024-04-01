package api

type (
	IDPLoginRequest struct {
		BaseRequest
		Slug string `param:"slug" query:"slug" form:"slug" json:"slug" xml:"slug" validate:"required"`
	}

	IDPLoginResponse struct {
		BaseResponse
		RedirectUri string `json:"redirect_uri"`
	}

	PasswordLoginRequest struct {
		BaseRequest
		Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"required"`
		Password string `param:"password" query:"password" form:"password" json:"password" xml:"password" validate:"required"`
	}
	PasswordLoginResponse struct {
		BaseResponse
		// NextPage to offer the user.
		// VerifyEmailPage
		//
		NextPage    string `json:"next_page,omitempty"`
		RedirectUri string `json:"redirect_uri,omitempty"`
	}
)
