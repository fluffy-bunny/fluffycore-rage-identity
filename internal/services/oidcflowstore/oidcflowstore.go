package oidcflowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	"github.com/rs/zerolog"
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

func AddSingletonIOIDCFlowStore(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_eko_gocache.IOIDCFlowStore](cb, stemService.Ctor)
}

func (s *service) StoreAuthorizationFinal(ctx context.Context, state string, value *models.AuthorizationFinal) error {
	log := zerolog.Ctx(ctx).With().Str("state", state).Logger()
	err := s.oidcFlowCache.Set(ctx, state, value, store.WithExpiration(30*time.Minute))
	log.Info().Err(err).Interface("request", value).Msg("StoreAuthorizationFinal")
	return err
}
func (s *service) GetAuthorizationFinal(ctx context.Context, state string) (*models.AuthorizationFinal, error) {
	log := zerolog.Ctx(ctx).With().Str("state", state).Logger()
	mm, err := s.oidcFlowCache.Get(ctx, state)
	if err != nil {
		// redirect to error page
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		return nil, err
	}
	var value *models.AuthorizationFinal = new(models.AuthorizationFinal)
	mmB, err := json.Marshal(mm)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		return nil, err
	}
	err = json.Unmarshal(mmB, value)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		return nil, err
	}
	return value, nil
}
func (s *service) DeleteAuthorizationFinal(ctx context.Context, state string) error {
	log := zerolog.Ctx(ctx).With().Str("state", state).Logger()
	err := s.oidcFlowCache.Delete(ctx, state)
	log.Info().Err(err).Msg("DeleteAuthorizationFinal")

	return err
}
