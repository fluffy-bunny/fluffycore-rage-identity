package nil

/*
Default IAuditStore.
Does Nothing
*/
import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)
var _ contracts_events.IAuditStore = stemService

func (s *service) Ctor() (contracts_events.IAuditStore, error) {
	return &service{}, nil
}

func AddSingletonIAuditStore(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_events.IAuditStore](cb, stemService.Ctor)
}
func (s *service) Submit(ctx context.Context, request *contracts_events.SubmitRequest) (*contracts_events.SubmitResponse, error) {
	return &contracts_events.SubmitResponse{}, nil
}
