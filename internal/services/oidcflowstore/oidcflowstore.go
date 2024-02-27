package oidcflowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/eko_gocache"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		proto_oidc_flows.UnimplementedOIDCFlowStoreServer

		oidcFlowCache contracts_eko_gocache.IOIDCFlowCache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_flows.IFluffyCoreOIDCFlowStoreServer = stemService
}
func (s *service) Ctor(oidcFlowCache contracts_eko_gocache.IOIDCFlowCache) (proto_oidc_flows.IFluffyCoreOIDCFlowStoreServer, error) {
	return &service{
		oidcFlowCache: oidcFlowCache,
	}, nil
}

func AddSingletonIOIDCFlowStore(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_flows.IFluffyCoreOIDCFlowStoreServer](cb, stemService.Ctor)
}

func (s *service) StoreAuthorizationFinal(ctx context.Context, request *proto_oidc_flows.StoreAuthorizationFinalRequest) (*proto_oidc_flows.StoreAuthorizationFinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	err := s.oidcFlowCache.Set(ctx, request.State, request.AuthorizationFinal, store.WithExpiration(30*time.Minute))
	log.Info().Err(err).Interface("request", request).Msg("StoreAuthorizationFinal")
	return &proto_oidc_flows.StoreAuthorizationFinalResponse{}, err
}
func (s *service) GetAuthorizationFinal(ctx context.Context, request *proto_oidc_flows.GetAuthorizationFinalRequest) (*proto_oidc_flows.GetAuthorizationFinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	mm, err := s.oidcFlowCache.Get(ctx, request.State)
	if err != nil {
		// redirect to error page
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		return nil, err
	}
	var value *proto_oidc_models.AuthorizationFinal = new(proto_oidc_models.AuthorizationFinal)
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
	return &proto_oidc_flows.GetAuthorizationFinalResponse{
		AuthorizationFinal: value,
	}, nil
}
func (s *service) DeleteAuthorizationFinal(ctx context.Context, request *proto_oidc_flows.DeleteAuthorizationFinalRequest) (*proto_oidc_flows.DeleteAuthorizationFinalResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	err := s.oidcFlowCache.Delete(ctx, request.State)
	log.Info().Err(err).Msg("DeleteAuthorizationFinal")

	return &proto_oidc_flows.DeleteAuthorizationFinalResponse{}, err
}
