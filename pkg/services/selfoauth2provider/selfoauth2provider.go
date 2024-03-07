package selfoauth2provider

import (
	"context"
	"sync"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_selfoauth2provider "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/selfoauth2provider"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config       *contracts_config.SelfIDPConfig
		oauth2Config *oauth2.Config
		lock         sync.Mutex
		verifier     *oidc.IDTokenVerifier
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_selfoauth2provider.ISelfOAuth2Provider = stemService
}
func (s *service) Ctor(
	config *contracts_config.SelfIDPConfig,
) (contracts_selfoauth2provider.ISelfOAuth2Provider, error) {
	return &service{
		config: config,
	}, nil
}

func AddSingletonISelfOAuth2Provider(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_selfoauth2provider.ISelfOAuth2Provider](cb, stemService.Ctor)
}

func (s *service) GetConfig(ctx context.Context) (*contracts_selfoauth2provider.GetConfigResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.oauth2Config != nil {
		return &contracts_selfoauth2provider.GetConfigResponse{
			Config:   s.oauth2Config,
			Verifier: s.verifier,
		}, nil
	}
	provider, err := oidc.NewProvider(ctx, s.config.Authority)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query provider.")
		return nil, status.Error(codes.Internal, "Failed to query provider.")
	}

	oauth2Config := oauth2.Config{
		ClientID:     s.config.ClientID,
		ClientSecret: s.config.ClientSecret,
		Scopes:       s.config.Scopes,
		RedirectURL:  s.config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.Endpoint().AuthURL,
			TokenURL: provider.Endpoint().TokenURL,
		},
	}
	s.oauth2Config = &oauth2Config
	oidcConfig := &oidc.Config{
		ClientID: s.config.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)
	s.verifier = verifier

	return &contracts_selfoauth2provider.GetConfigResponse{
		Config:   s.oauth2Config,
		Verifier: s.verifier,
	}, nil
}
