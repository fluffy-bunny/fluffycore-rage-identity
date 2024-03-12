package webauthn

import (
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
	uuid "github.com/gofrs/uuid"
)

type (
	WebAuthNConfig struct {
		RPDisplayName string   `json:"rpDisplayName"`
		RPID          string   `json:"rpid"`
		RPOrigins     []string `json:"rpOrigins"`
	}
	IWebAuthN interface {
		GetWebAuthN() *go_webauthn.WebAuthn
		GetFriendlyNameByAAGUID(aaguid uuid.UUID) string
	}
)
