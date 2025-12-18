package AuthorizationCodeClaimsAugmentor

/*
Default IEventSink.
Does Nothing
*/
import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
	proto_events_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/events/types"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)
var _ contracts_events.IEventSink = stemService

func (s *service) Ctor() (contracts_events.IEventSink, error) {
	return &service{}, nil
}

func AddSingletonIEventSink(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_events.IEventSink](cb, stemService.Ctor)
}
func (s *service) OnEvent(ctx context.Context, event *proto_events_types.Event) {

}
