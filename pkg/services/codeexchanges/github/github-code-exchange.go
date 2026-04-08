package github

/*
https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28#get-the-authenticated-user
https://docs.github.com/en/rest/users/emails?apiVersion=2022-11-28#list-email-addresses-for-the-authenticated-user
*/
import (
	"context"
	"strconv"

	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
		config        *contracts_config.Config
		oauth2Factory contracts_oauth2factory.IOAuth2Factory
		tokenService  contracts_tokenservice.ITokenService
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

const (
	// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#1-request-a-users-github-identity
	GithubAuthURL = "https://github.com/login/oauth/authorize"
	// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#2-users-are-redirected-back-to-your-site-by-github
	GithubTokenURL         = "https://github.com/login/oauth/access_token"
	GithubUserInfoEndpoint = "https://api.github.com/user"
	GitHubEmailsEndpoint   = "https://api.github.com/user/emails"
)

var GithubScopes = []string{"user:email"}

var stemService = (*service)(nil)
var _ contracts_codeexchange.IGithubCodeExchange = stemService

func (s *service) Ctor(config *contracts_config.Config,
	tokenService contracts_tokenservice.ITokenService,
	oauth2Factory contracts_oauth2factory.IOAuth2Factory) (contracts_codeexchange.IGithubCodeExchange, error) {
	return &service{
		config:        config,
		oauth2Factory: oauth2Factory,
		tokenService:  tokenService,
	}, nil
}

func AddSingletonIGithubCodeExchange(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_codeexchange.IGithubCodeExchange](cb, stemService.Ctor)
}
func (s *service) validateExchangeCodeRequest(request *contracts_codeexchange.ExchangeCodeRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Code) {
		return status.Error(codes.InvalidArgument, "code is required")
	}
	// GitHub OAuth2 does not support PKCE, so CodeVerifier is not required
	return nil
}

var Duration30MinutesSeconds = 1800

func (s *service) ExchangeCode(ctx context.Context, request *contracts_codeexchange.ExchangeCodeRequest) (*contracts_codeexchange.ExchangeCodeResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	if err := s.validateExchangeCodeRequest(request); err != nil {
		return nil, err
	}
	// get the config
	getConfigRequest := &contracts_oauth2factory.GetConfigRequest{
		IDPHint: request.IDPHint,
	}
	getConfigResponse, err := s.oauth2Factory.GetConfig(ctx, getConfigRequest)
	if err != nil {
		log.Error().Err(err).Msg("failed to get config")
		return nil, err
	}
	config := getConfigResponse.Config
	// GitHub now enforces PKCE: if code_challenge was sent, code_verifier is required
	log.Debug().
		Str("redirectURL", config.RedirectURL).
		Str("clientID", config.ClientID).
		Str("tokenURL", config.Endpoint.TokenURL).
		Int("codeLen", len(request.Code)).
		Bool("hasPKCE", fluffycore_utils.IsNotEmptyOrNil(request.CodeVerifier)).
		Msg("exchanging code with GitHub")
	var exchangeOpts []oauth2.AuthCodeOption
	if fluffycore_utils.IsNotEmptyOrNil(request.CodeVerifier) {
		exchangeOpts = append(exchangeOpts, oauth2.SetAuthURLParam("code_verifier", request.CodeVerifier))
	}
	token, err := config.Exchange(context.Background(), request.Code, exchangeOpts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to exchange code")
		return nil, err
	}
	doGitHubGet := func(url string, result interface{}) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, result)
	}

	githubUser := &GithubUser{}
	if err := doGitHubGet(GithubUserInfoEndpoint, githubUser); err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		return nil, status.Error(codes.Internal, "failed to get user info")
	}
	githubEmails := make(GithubUserEmails, 0)
	if err := doGitHubGet(GitHubEmailsEndpoint, &githubEmails); err != nil {
		log.Error().Err(err).Msg("failed to get user emails")
		return nil, status.Error(codes.Internal, "failed to get user emails")
	}
	var primaryEmail *GitHubEmail
	for _, e := range githubEmails {
		if fluffycore_utils.IsEmptyOrNil(e.Email) {
			continue
		}
		if e.Primary {
			primaryEmail = &e
		}
	}
	claims := fluffycore_services_claims.NewClaims()
	claims.Set("sub", strconv.Itoa(githubUser.ID))
	claims.Set("idp", request.IDPHint)

	if primaryEmail != nil {
		claims.Set("email", primaryEmail.Email)
		claims.Set("email_verified", primaryEmail.Verified)
	} else {
		// if we don't have an email, we can't mint a token
		return nil, status.Error(codes.Internal, "failed to get user info, can't get emails")
	}

	mintTokenResponse, err := s.tokenService.MintToken(ctx,
		&contracts_tokenservice.MintTokenRequest{
			DurationLifeTimeSeconds: Duration30MinutesSeconds,
			Claims:                  claims,
		})
	if err != nil {
		log.Error().Err(err).Msg("failed to mint token")
		return nil, err
	}
	return &contracts_codeexchange.ExchangeCodeResponse{
		IdToken: mintTokenResponse.Token,
	}, nil
}
