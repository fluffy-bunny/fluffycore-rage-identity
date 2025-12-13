package models

type LinkedIdentity struct {
	Subject  string `json:"subject"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
	LinkedAt string `json:"linkedAt"`
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
