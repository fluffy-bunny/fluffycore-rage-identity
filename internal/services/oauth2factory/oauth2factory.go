package oauth2factory

import (
	"context"
	"sync"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/oauth2factory"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/models"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config           *contracts_config.Config
		idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer
		oidcProviders    map[string]*oidc.Provider
		oauth2Configs    map[string]*oauth2.Config
		lock             sync.Mutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_oauth2factory.IOAuth2Factory = stemService
}
func (s *service) Ctor(config *contracts_config.Config, idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer) (contracts_oauth2factory.IOAuth2Factory, error) {
	return &service{
		config:           config,
		idpServiceServer: idpServiceServer,
		oidcProviders:    make(map[string]*oidc.Provider),
		oauth2Configs:    make(map[string]*oauth2.Config),
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
	if fluffycore_utils.IsEmptyOrNil(request.IDPSlug) {
		return status.Error(codes.InvalidArgument, "IDPSlug is required")
	}
	return nil
}
func (s *service) GetConfig(ctx context.Context, request *contracts_oauth2factory.GetConfigRequest) (*contracts_oauth2factory.GetConfigResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetConfigRequest(request)
	if err != nil {
		return nil, err
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	getIDPBySlugResponse, err := s.idpServiceServer.GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: request.IDPSlug,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return nil, err
	}
	idp := getIDPBySlugResponse.Idp
	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch v := idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			oauth2Config, ok := s.oauth2Configs[request.IDPSlug]
			if !ok {
				config := oauth2.Config{
					ClientID:     v.Github.ClientId,
					ClientSecret: v.Github.ClientSecret,
					Scopes:       GithubScopes,
					RedirectURL:  s.config.BaseUrl + wellknown_echo.OAuth2CallbackPath,
					Endpoint: oauth2.Endpoint{
						AuthURL:  GithubAuthURL,
						TokenURL: GithubTokenURL,
					},
				}
				s.oauth2Configs[request.IDPSlug] = &config
				oauth2Config = &config
			}
			return &contracts_oauth2factory.GetConfigResponse{
				Config: oauth2Config,
			}, nil

		case *proto_oidc_models.Protocol_Oidc:
			{
				oidcProvider, ok := s.oidcProviders[request.IDPSlug]
				if !ok {
					provider, err := oidc.NewProvider(ctx, v.Oidc.Authority)
					if err != nil {
						log.Error().Err(err).Msg("oidc.NewProvider")
						return nil, err
					}
					s.oidcProviders[request.IDPSlug] = provider
					oidcProvider = provider
				}
				oauth2Config, ok := s.oauth2Configs[request.IDPSlug]
				if !ok {

					config := oauth2.Config{
						ClientID:     v.Oidc.ClientId,
						ClientSecret: v.Oidc.ClientSecret,
						Scopes:       GithubScopes,
						RedirectURL:  s.config.BaseUrl + wellknown_echo.OAuth2CallbackPath,
						Endpoint: oauth2.Endpoint{
							AuthURL:  oidcProvider.Endpoint().AuthURL,
							TokenURL: oidcProvider.Endpoint().TokenURL,
						},
					}
					s.oauth2Configs[request.IDPSlug] = &config
					oauth2Config = &config
				}
				return &contracts_oauth2factory.GetConfigResponse{
					Config: oauth2Config,
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.NotFound, "no oauth2 protocol found for IDPSlug: %s", request.IDPSlug)
}
