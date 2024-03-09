package oidcflowstore

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
		proto_oidc_flows.UnimplementedAuthorizationRequestStateStoreServer

		oidcFlowCache contracts_eko_gocache.IAuthorizationRequestStateCache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer = stemService
}
func (s *service) Ctor(oidcFlowCache contracts_eko_gocache.IAuthorizationRequestStateCache) (proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer, error) {
	return &service{
		oidcFlowCache: oidcFlowCache,
	}, nil
}

func AddSingletonAuthorizationRequestStateStoreServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer](cb, stemService.Ctor)
}

func (s *service) StoreAuthorizationRequestState(ctx context.Context, request *proto_oidc_flows.StoreAuthorizationRequestStateRequest) (*proto_oidc_flows.StoreAuthorizationRequestStateResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	err := s.oidcFlowCache.Set(ctx, request.State, request.AuthorizationRequestState, store.WithExpiration(30*time.Minute))
	log.Info().Err(err).Interface("request", request).Msg("StoreAuthorizationRequestState")
	return &proto_oidc_flows.StoreAuthorizationRequestStateResponse{}, err
}
func (s *service) GetAuthorizationRequestState(ctx context.Context, request *proto_oidc_flows.GetAuthorizationRequestStateRequest) (*proto_oidc_flows.GetAuthorizationRequestStateResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	mm, err := s.oidcFlowCache.Get(ctx, request.State)
	if err != nil {
		// redirect to error page
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return nil, err
	}
	var value *proto_oidc_models.AuthorizationRequestState = new(proto_oidc_models.AuthorizationRequestState)
	mmB, err := json.Marshal(mm)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return nil, err
	}
	err = json.Unmarshal(mmB, value)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return nil, err
	}
	return &proto_oidc_flows.GetAuthorizationRequestStateResponse{
		AuthorizationRequestState: value,
	}, nil
}
func (s *service) DeleteAuthorizationRequestState(ctx context.Context, request *proto_oidc_flows.DeleteAuthorizationRequestStateRequest) (*proto_oidc_flows.DeleteAuthorizationRequestStateResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	err := s.oidcFlowCache.Delete(ctx, request.State)
	log.Info().Err(err).Msg("DeleteAuthorizationRequestState")

	return &proto_oidc_flows.DeleteAuthorizationRequestStateResponse{}, err
}
