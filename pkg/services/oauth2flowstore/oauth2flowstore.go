package oauth2flowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/eko_gocache"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		proto_oidc_flows.UnimplementedExternalOauth2FlowStoreServer

		externalOAuth2Cache contracts_eko_gocache.IExternalOAuth2Cache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_flows.IFluffyCoreExternalOauth2FlowStoreServer = stemService
}
func (s *service) Ctor(externalOAuth2Cache contracts_eko_gocache.IExternalOAuth2Cache) (proto_oidc_flows.IFluffyCoreExternalOauth2FlowStoreServer, error) {
	return &service{
		externalOAuth2Cache: externalOAuth2Cache,
	}, nil
}

func AddSingletonIExternalOauth2FlowStore(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_flows.IFluffyCoreExternalOauth2FlowStoreServer](cb, stemService.Ctor)
}

func (s *service) StoreExternalOauth2Final(ctx context.Context, request *proto_oidc_flows.StoreExternalOauth2FinalRequest) (*proto_oidc_flows.StoreExternalOauth2FinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.externalOAuth2Cache.Set(ctx, request.State, request.ExternalOauth2Final, store.WithExpiration(30*time.Minute))
	log.Info().Err(err).Msg("externalOAuth2Cache.Set")
	return &proto_oidc_flows.StoreExternalOauth2FinalResponse{}, err
}
func (s *service) GetExternalOauth2Final(ctx context.Context, request *proto_oidc_flows.GetExternalOauth2FinalRequest) (*proto_oidc_flows.GetExternalOauth2FinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()

	mm, err := s.externalOAuth2Cache.Get(ctx, request.State)
	if err != nil {
		log.Error().Err(err).Msg("externalOAuth2Cache.Get")
		// redirect to error page
		return nil, err
	}
	var value *proto_oidc_models.ExternalOauth2Final = new(proto_oidc_models.ExternalOauth2Final)
	mmB, err := json.Marshal(mm)
	if err != nil {
		log.Error().Err(err).Msg("marshal")
		return nil, err
	}
	err = json.Unmarshal(mmB, value)
	if err != nil {
		log.Error().Err(err).Msg("unmarshal")
		return nil, err
	}
	return &proto_oidc_flows.GetExternalOauth2FinalResponse{
		ExternalOauth2Final: value,
	}, nil
}
func (s *service) DeleteExternalOauth2Final(ctx context.Context, request *proto_oidc_flows.DeleteExternalOauth2FinalRequest) (*proto_oidc_flows.DeleteExternalOauth2FinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.externalOAuth2Cache.Delete(ctx, request.State)
	log.Info().Err(err).Msg("externalOAuth2Cache.Delete")
	return &proto_oidc_flows.DeleteExternalOauth2FinalResponse{}, err
}
