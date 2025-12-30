package oidcproviderfactory

import (
	"context"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oauth2factory"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config           *contracts_config.Config
		idpServiceServer proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
	}
)

var stemService = (*service)(nil)
var _ contracts_oauth2factory.IOIDCProviderFactory = stemService

func (s *service) Ctor(config *contracts_config.Config, idpServiceServer proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer) (contracts_oauth2factory.IOIDCProviderFactory, error) {
	return &service{
		config:           config,
		idpServiceServer: idpServiceServer,
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

	// Always fetch fresh IDP data from database - no caching
	listIDPResponse, err := s.idpServiceServer.ListIDP(ctx,
		&proto_oidc_idp.ListIDPRequest{
			Filter: &proto_oidc_idp.Filter{
				Enabled: &proto_types.BoolFilterExpression{
					Eq: true,
				},
				Slug: &proto_types.StringFilterExpression{
					Eq: request.IDPHint,
				},
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("ListIDP")
		return nil, err
	}
	if listIDPResponse == nil || len(listIDPResponse.IDPs) == 0 {
		return nil, status.Errorf(codes.NotFound, "no idp found for IDPHint: %s", request.IDPHint)
	}
	idp := listIDPResponse.IDPs[0]
	if idp.Protocol != nil {
		log.Debug().Interface("idp", idp).Msg("listIDPResponse")
		switch v := idp.Protocol.Value.(type) {

		case *proto_oidc_models.Protocol_Oidc:
			{
				// Always create fresh provider with current config
				provider, err := oidc.NewProvider(ctx, v.Oidc.Authority)
				if err != nil {
					log.Error().Err(err).Msg("oidc.NewProvider")
					return nil, err
				}
				return &contracts_oauth2factory.GetOIDCProviderResponse{
					OIDCProvider: provider,
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.NotFound, "no oauth2 protocol found for IDPHint: %s", request.IDPHint)
}
