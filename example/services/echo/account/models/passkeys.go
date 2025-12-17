package models

type PasskeyItem struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendlyName"`
	CreatedAt    int64  `json:"createdAt,omitempty"`
	LastUsedAt   int64  `json:"lastUsedAt,omitempty"`
}

type PasskeysResponse struct {
	Passkeys []PasskeyItem `json:"passkeys"`
}

type PasskeyRenameRequest struct {
	CredentialID string `json:"credentialID"`
	FriendlyName string `json:"friendlyName"`
}

type PasskeyRenameResponse struct {
	Success bool `json:"success"`
}

type PasskeyDeleteRequest struct {
	CredentialID string `json:"credentialID"`
}
type PasskeyDeleteResponse struct {
	Success bool `json:"success"`
}
