package oauth2factory

import (
	"context"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oauth2factory"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config                *contracts_config.Config
		idpServiceServer      proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
		oauth2ProviderFactory contracts_oauth2factory.IOIDCProviderFactory
	}
)

var stemService = (*service)(nil)
var _ contracts_oauth2factory.IOAuth2Factory = stemService

func (s *service) Ctor(config *contracts_config.Config,
	oauth2ProviderFactory contracts_oauth2factory.IOIDCProviderFactory,
	idpServiceServer proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer) (contracts_oauth2factory.IOAuth2Factory, error) {
	return &service{
		config:                config,
		idpServiceServer:      idpServiceServer,
		oauth2ProviderFactory: oauth2ProviderFactory,
	}, nil
}

func AddSingletonIOAuth2Factory(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_oauth2factory.IOAuth2Factory](cb, stemService.Ctor)
}

const (
	// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#1-request-a-users-github-identity
	GithubAuthURL = "https://github.com/login/oauth/authorize"
	// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#2-users-are-redirected-back-to-your-site-by-github
	GithubTokenURL         = "https://github.com/login/oauth/access_token"
	GithubUserInfoEndpoint = "https://api.github.com/user"
	GitHubEmailsEndpoint   = "https://api.github.com/user/emails"
)

var GithubScopes = []string{"user:email"}

func (s *service) validateGetConfigRequest(request *contracts_oauth2factory.GetConfigRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.IDPHint) {
		return status.Error(codes.InvalidArgument, "IDPHint is required")
	}
	return nil
}
func (s *service) GetConfig(ctx context.Context, request *contracts_oauth2factory.GetConfigRequest) (*contracts_oauth2factory.GetConfigResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetConfigRequest(request)
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
		log.Debug().Interface("listIDPResponse", listIDPResponse).Msg("listIDPResponse")
		switch v := idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			// Always create fresh config with current credentials
			config := oauth2.Config{
				ClientID:     v.Github.ClientId,
				ClientSecret: v.Github.ClientSecret,
				Scopes:       GithubScopes,
				RedirectURL:  s.config.OIDCConfig.BaseUrl + s.config.OIDCConfig.OAuth2CallbackPath,
				Endpoint: oauth2.Endpoint{
					AuthURL:  GithubAuthURL,
					TokenURL: GithubTokenURL,
				},
			}
			return &contracts_oauth2factory.GetConfigResponse{
				Config: &config,
			}, nil

		case *proto_oidc_models.Protocol_Oidc:
			{
				getOIDCProviderResponse, err := s.oauth2ProviderFactory.GetOIDCProvider(ctx,
					&contracts_oauth2factory.GetOIDCProviderRequest{
						IDPHint: request.IDPHint,
					})
				if err != nil {
					log.Error().Err(err).Msg("Failed to get oidcProvider")
					return nil, err
				}
				oidcProvider := getOIDCProviderResponse.OIDCProvider

				// Always create fresh config with current credentials
				scopes := strings.Split(v.Oidc.Scope, " ")
				config := oauth2.Config{
					ClientID:     v.Oidc.ClientId,
					ClientSecret: v.Oidc.ClientSecret,
					Scopes:       scopes,
					RedirectURL:  s.config.OIDCConfig.BaseUrl + s.config.OIDCConfig.OAuth2CallbackPath,
					Endpoint: oauth2.Endpoint{
						AuthURL:  oidcProvider.Endpoint().AuthURL,
						TokenURL: oidcProvider.Endpoint().TokenURL,
					},
				}
				return &contracts_oauth2factory.GetConfigResponse{
					Config: &config,
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.NotFound, "no oauth2 protocol found for IDPHint: %s", request.IDPHint)
}
