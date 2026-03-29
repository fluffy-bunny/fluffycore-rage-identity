package events

import (
	"context"

	proto_events_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/events/types"
)

type (
	IAuditStore interface {
		Submit(ctx context.Context, request *SubmitRequest) (*SubmitResponse, error)
	}
	SubmitRequest struct {
		CloudEvent *proto_events_types.CloudEvent `json:"cloudEvent,omitempty"`
	}
	SubmitResponse struct{}
	LoginEvent     struct {
		Subject   string   `json:"subject,omitempty"`
		Email     string   `json:"email,omitempty"`
		ClientID  string   `json:"client_id,omitempty"`
		ACR       []string `json:"acr,omitempty"`
		AMR       []string `json:"amr,omitempty"`
		IDP       []string `json:"idp,omitempty"`
		LoginType string   `json:"login_type,omitempty"`
		Outcome   string   `json:"outcome,omitempty"`
	}
)
