package oauth2flowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
)

type (
	service struct {
		externalOAuth2Cache contracts_eko_gocache.IExternalOAuth2Cache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_eko_gocache.IExternalOauth2FlowStore = stemService
}
func (s *service) Ctor(externalOAuth2Cache contracts_eko_gocache.IExternalOAuth2Cache) (contracts_eko_gocache.IExternalOauth2FlowStore, error) {
	return &service{
		externalOAuth2Cache: externalOAuth2Cache,
	}, nil
}

func AddSingletonIExternalOauth2FlowStore(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_eko_gocache.IExternalOauth2FlowStore](cb, stemService.Ctor)
}

func (s *service) StoreExternalOauth2Final(ctx context.Context, state string, value *models.ExternalOauth2Final) error {
	err := s.externalOAuth2Cache.Set(ctx, state, value, store.WithExpiration(30*time.Minute))
	return err
}
func (s *service) GetExternalOauth2Final(ctx context.Context, state string) (*models.ExternalOauth2Final, error) {
	mm, err := s.externalOAuth2Cache.Get(ctx, state)
	if err != nil {
		// redirect to error page
		return nil, err
	}
	var value *models.ExternalOauth2Final = new(models.ExternalOauth2Final)
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
func (s *service) DeleteExternalOauth2Final(ctx context.Context, state string) error {
	err := s.externalOAuth2Cache.Delete(ctx, state)
	return err
}
