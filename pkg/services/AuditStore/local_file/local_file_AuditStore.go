package local_file

/*
Development IAuditStore.
Persists CloudEvents to tmp/auditstore.jsonl (NDJSON, append-only)
*/
import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

type (
	service struct {
		mu       sync.Mutex
		filePath string
	}
)

var stemService = (*service)(nil)
var _ contracts_events.IAuditStore = stemService

func (s *service) Ctor() (contracts_events.IAuditStore, error) {
	return &service{
		filePath: filepath.Join("tmp", "auditstore.jsonl"),
	}, nil
}

func AddSingletonIAuditStore(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_events.IAuditStore](cb, stemService.Ctor)
}
func (s *service) Submit(ctx context.Context, request *contracts_events.SubmitRequest) (*contracts_events.SubmitResponse, error) {
	_ = ctx
	if request == nil || request.CloudEvent == nil {
		return nil, errors.New("request.CloudEvent is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	eventJSON, err := protojson.Marshal(request.CloudEvent)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Write(eventJSON); err != nil {
		return nil, err
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return nil, err
	}

	return &contracts_events.SubmitResponse{}, nil
}
