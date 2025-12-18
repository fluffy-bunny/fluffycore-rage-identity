package models

type LinkedIdentity struct {
	Subject    string `json:"subject"`
	Provider   string `json:"provider"`
	Email      string `json:"email"`
	CreatedOn  int64  `json:"createdOn,omitempty"`
	LastUsedOn int64  `json:"lastUsedOn,omitempty"`
}

type LinkedAccountsResponse struct {
	Identities      []LinkedIdentity `json:"identities"`
	IsClaimedDomain bool             `json:"isClaimedDomain"`
}

type DeleteLinkedAccountRequest struct {
	Identity string `json:"identity" param:"identity"`
}

type DeleteLinkedAccountResponse struct {
	Success bool `json:"success"`
}
