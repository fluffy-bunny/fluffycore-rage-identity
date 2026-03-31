package oidcflowstore

import (
	"context"
	"encoding/json"

	"time"

	store "github.com/eko/gocache/lib/v4/store"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/eko_gocache"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	zerolog "github.com/rs/zerolog"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type (
	service struct {
		proto_oidc_flows.UnimplementedAuthorizationRequestStateStoreServer

		oidcFlowCache contracts_eko_gocache.IAuthorizationRequestStateCache
		ttl           time.Duration
		codeTTL       time.Duration
	}
)

var stemService = (*service)(nil)

var _ proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	oidcFlowCache contracts_eko_gocache.IAuthorizationRequestStateCache,
) (proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer, error) {
	ttlMinutes := 30
	if config.OIDCConfig != nil && config.OIDCConfig.AuthorizationStateTTLMinutes > 0 {
		ttlMinutes = config.OIDCConfig.AuthorizationStateTTLMinutes
	}
	codeTTLSeconds := 600
	if config.OIDCConfig != nil && config.OIDCConfig.AuthorizationCodeTTLSeconds > 0 {
		codeTTLSeconds = config.OIDCConfig.AuthorizationCodeTTLSeconds
	}
	return &service{
		oidcFlowCache: oidcFlowCache,
		ttl:           time.Duration(ttlMinutes) * time.Minute,
		codeTTL:       time.Duration(codeTTLSeconds) * time.Second,
	}, nil
}

func AddSingletonAuthorizationRequestStateStoreServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer](cb, stemService.Ctor)
}

func (s *service) StoreAuthorizationRequestState(ctx context.Context, request *proto_oidc_flows.StoreAuthorizationRequestStateRequest) (*proto_oidc_flows.StoreAuthorizationRequestStateResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("state", request.State).Logger()
	request.AuthorizationRequestState.Updated = timestamppb.Now()
	// Use the shorter code TTL when identity is set (authorization code entry),
	// otherwise use the longer session TTL (login flow state).
	ttl := s.ttl
	if request.AuthorizationRequestState.Identity != nil {
		ttl = s.codeTTL
	}
	err := s.oidcFlowCache.Set(ctx, request.State, request.AuthorizationRequestState, store.WithExpiration(ttl))
	log.Debug().Err(err).Dur("ttl", ttl).Interface("request", request).Msg("StoreAuthorizationRequestState")
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
	log.Debug().Err(err).Msg("DeleteAuthorizationRequestState")

	return &proto_oidc_flows.DeleteAuthorizationRequestStateResponse{}, err
}
