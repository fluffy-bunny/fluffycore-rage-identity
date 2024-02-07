package oidcflowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	models "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/models"
)

type (
	service struct {
		oidcFlowCache contracts_eko_gocache.IOIDCFlowCache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_eko_gocache.IOIDCFlowStore = stemService
}
func (s *service) Ctor(oidcFlowCache contracts_eko_gocache.IOIDCFlowCache) (contracts_eko_gocache.IOIDCFlowStore, error) {
	return &service{
		oidcFlowCache: oidcFlowCache,
	}, nil
}

func AddSingletonIOIDCFlowCache(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_eko_gocache.IOIDCFlowStore](cb, stemService.Ctor)
}

func (s *service) StoreAuthorizationFinal(ctx context.Context, code string, value *models.AuthorizationFinal) error {
	err := s.oidcFlowCache.Set(ctx, code, value, store.WithExpiration(30*time.Minute))
	return err
}
func (s *service) GetAuthorizationFinal(ctx context.Context, code string) (*models.AuthorizationFinal, error) {
	mm, err := s.oidcFlowCache.Get(ctx, code)
	if err != nil {
		// redirect to error page
		return nil, err
	}
	var value *models.AuthorizationFinal = new(models.AuthorizationFinal)
	mmB, err := json.Marshal(mm)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(mmB, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
