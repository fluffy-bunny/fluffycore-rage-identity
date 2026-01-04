package login_models

type KeepSignedInRequest struct {
	KeepSignedIn   bool `json:"keepSignedIn"`
	DoNotShowAgain bool `json:"doNotShowAgain"`
}

type KeepSignedInResponse struct {
	Directive         string             `json:"directive"`
	DirectiveRedirect *DirectiveRedirect `json:"directiveRedirect,omitempty"`
}

type KeepSignedInErrorResponse struct {
	Reason string `json:"reason"`
}
