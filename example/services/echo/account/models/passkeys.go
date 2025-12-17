package models

type PasskeyItem struct {
	ID           string   `json:"id"`
	FriendlyName string   `json:"friendlyName"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	Transport    []string `json:"transport,omitempty"`
}

type PasskeysResponse struct {
	Passkeys []PasskeyItem `json:"passkeys"`
}

type PasskeyRenameRequest struct {
	CredentialID string `json:"credentialID"S`
	FriendlyName string `json:"friendlyName"`
}

type PasskeyRenameResponse struct {
	Success bool `json:"success"`
}

type PasskeyDeleteRequest struct {
	CredentialID string `json:"credentialID"S`
}
type PasskeyDeleteResponse struct {
	Success bool `json:"success"`
}
