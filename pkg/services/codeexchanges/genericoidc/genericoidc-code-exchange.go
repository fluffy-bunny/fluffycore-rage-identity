package genericoidc

import (
	"context"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_codeexchange "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/codeexchange"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oauth2factory"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
	fluffycore_services_claims "github.com/fluffy-bunny/fluffycore/services/claims"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config              *contracts_config.Config
		oauth2Factory       contracts_oauth2factory.IOAuth2Factory
		oidcProviderFactory contracts_oauth2factory.IOIDCProviderFactory
		tokenService        contracts_tokenservice.ITokenService
	}
	GithubUser struct {
		ID       int    `json:"id"`
		UserName string `json:"login"`
		Email    string `json:"email"`
	}
	GitHubEmail struct {
		Email      string `json:"email"`
		Verified   bool   `json:"verified"`
		Primary    bool   `json:"primary"`
		Visibility string `json:"visibility"`
	}
	GithubUserEmails []GitHubEmail
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_codeexchange.IGenericOIDCCodeExchange = stemService
}
func (s *service) Ctor(config *contracts_config.Config,
	tokenService contracts_tokenservice.ITokenService,
	oauth2Factory contracts_oauth2factory.IOAuth2Factory,
	oidcProviderFactory contracts_oauth2factory.IOIDCProviderFactory) (contracts_codeexchange.IGenericOIDCCodeExchange, error) {
	return &service{
		config:              config,
		oidcProviderFactory: oidcProviderFactory,
		tokenService:        tokenService,
		oauth2Factory:       oauth2Factory,
	}, nil
}

func AddSingletonIGenericOIDCCodeExchange(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_codeexchange.IGenericOIDCCodeExchange](cb, stemService.Ctor)
}
func (s *service) validateExchangeCodeRequest(request *contracts_codeexchange.ExchangeCodeRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.IDPHint) {
		return status.Error(codes.InvalidArgument, "IDPHint is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Code) {
		return status.Error(codes.InvalidArgument, "code is required")
	}

	return nil
}

var Duration30MinutesSeconds = 1800

func (s *service) ExchangeCode(ctx context.Context, request *contracts_codeexchange.ExchangeCodeRequest) (*contracts_codeexchange.ExchangeCodeResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	if err := s.validateExchangeCodeRequest(request); err != nil {
		return nil, err
	}
	// get the oidc provider
	oidcProviderResponse, err := s.oidcProviderFactory.GetOIDCProvider(ctx, &contracts_oauth2factory.GetOIDCProviderRequest{
		IDPHint: request.IDPHint,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get oidc provider")
		return nil, err
	}
	oidcProvider := oidcProviderResponse.OIDCProvider
	// get the config
	oidcConfig := &oidc.Config{
		ClientID: request.ClientID,
	}
	idTokenVerifier := oidcProvider.Verifier(oidcConfig)
	getConfigResponse, err := s.oauth2Factory.GetConfig(ctx, &contracts_oauth2factory.GetConfigRequest{
		IDPHint: request.IDPHint,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get config")
		return nil, err
	}
	config := getConfigResponse.Config
	authCodeOptions := []oauth2.AuthCodeOption{}
	if fluffycore_utils.IsNotEmptyOrNil(request.CodeVerifier) {
		authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam("code_verifier", request.CodeVerifier))
	}
	oauth2Token, err := config.Exchange(context.Background(),
		request.Code, authCodeOptions...)
	if err != nil {
		log.Error().Err(err).Msg("failed to exchange code")
		return nil, err
	}
	log.Debug().Interface("oauth2Token", oauth2Token).Msg("callbackPath")
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Error().Msg("failed to get id_token")
		return nil, status.Error(codes.Internal, "failed to get id_token")
	}
	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to verify id_token")
		return nil, err
	}
	if idToken.Nonce != request.Nonce {
		log.Error().Msg("nonce does not match")
		return nil, status.Error(codes.Internal, "nonce does not match")
	}
	var idTokenClaims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&idTokenClaims); err != nil {
		// handle error
		log.Error().Err(err).Msg("failed to get claims")
		return nil, err
	}
	claims := fluffycore_services_claims.NewClaims()
	claims.Set("sub", idToken.Subject)
	claims.Set("idp", request.IDPHint)
	claims.Set("email", idTokenClaims.Email)
	claims.Set("email_verified", idTokenClaims.EmailVerified)

	mintTokenResponse, err := s.tokenService.MintToken(ctx,
		&contracts_tokenservice.MintTokenRequest{
			Claims:                  claims,
			DurationLifeTimeSeconds: Duration30MinutesSeconds,
		})
	if err != nil {
		log.Error().Err(err).Msg("failed to mint token")
		return nil, err
	}
	return &contracts_codeexchange.ExchangeCodeResponse{
		IdToken: mintTokenResponse.Token,
	}, nil
}
