package oidcproviderfactory

import (
	"context"
	"sync"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oauth2factory"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config           *contracts_config.Config
		idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer
		oidcProviders    map[string]*oidc.Provider
		lock             sync.Mutex
	}
)

var stemService = (*service)(nil)
var _ contracts_oauth2factory.IOIDCProviderFactory = stemService

func (s *service) Ctor(config *contracts_config.Config, idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer) (contracts_oauth2factory.IOIDCProviderFactory, error) {
	return &service{
		config:           config,
		idpServiceServer: idpServiceServer,
		oidcProviders:    make(map[string]*oidc.Provider),
	}, nil
}

func AddSingletonIOIDCProviderFactory(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_oauth2factory.IOIDCProviderFactory](cb, stemService.Ctor)
}

func (s *service) validateGetOIDCProviderRequest(request *contracts_oauth2factory.GetOIDCProviderRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.IDPHint) {
		return status.Error(codes.InvalidArgument, "IDPHint is required")
	}
	return nil
}
func (s *service) GetOIDCProvider(ctx context.Context, request *contracts_oauth2factory.GetOIDCProviderRequest) (*contracts_oauth2factory.GetOIDCProviderResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetOIDCProviderRequest(request)
	if err != nil {
		return nil, err
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	getIDPBySlugResponse, err := s.idpServiceServer.GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: request.IDPHint,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return nil, err
	}
	idp := getIDPBySlugResponse.Idp
	if idp.Protocol != nil {
		log.Debug().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch v := idp.Protocol.Value.(type) {

		case *proto_oidc_models.Protocol_Oidc:
			{
				oidcProvider, ok := s.oidcProviders[request.IDPHint]
				if !ok {
					provider, err := oidc.NewProvider(ctx, v.Oidc.Authority)
					if err != nil {
						log.Error().Err(err).Msg("oidc.NewProvider")
						return nil, err
					}
					s.oidcProviders[request.IDPHint] = provider
					oidcProvider = provider
				}
				return &contracts_oauth2factory.GetOIDCProviderResponse{
					OIDCProvider: oidcProvider,
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.NotFound, "no oauth2 protocol found for IDPHint: %s", request.IDPHint)
}
