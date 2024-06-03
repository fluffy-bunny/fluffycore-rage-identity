package external_idp

// StartExternalIDPLoginRequest example
type StartExternalIDPLoginRequest struct {
	Slug      string `json:"slug" validate:"required"`
	Directive string `json:"directive" validate:"required"`
}

// StartExternalIDPLoginResponse example
type StartExternalIDPLoginResponse struct {
	RedirectURI string `json:"redirectUri"`
}
