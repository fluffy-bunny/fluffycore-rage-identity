package events

import (
	"context"

	proto_events_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/events/types"
)

type (
	IEventSink interface {
		OnEvent(ctx context.Context, event *proto_events_types.Event)
	}
)
