package echo

import (
	"fmt"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	echo "github.com/labstack/echo/v4"
	oauth2 "golang.org/x/oauth2"
)

const (
	HomePath                        = "/"
	LoginPath                       = "/login"
	LogoutPath                      = "/logout"
	ProfilePath                     = "/profile"
	OIDCLoginPath                   = "/oidc-login"
	SignupPath                      = "/signup"
	ExternalIDPPath                 = "/external-idp"
	SwaggerPath                     = "/swagger/*"
	HealthzPath                     = "/healthz"
	ErrorPath                       = "/error"
	AboutPath                       = "/about"
	ReadyPath                       = "/ready"
	WellKnownOpenIDCOnfiguationPath = "/.well-known/openid-configuration"
	WellKnownJWKS                   = "/.well-known/jwks"
	OAuth2TokenEndpointPath         = "/token"
	OIDCAuthorizationEndpointPath   = "/oidc/v1/auth"
	UserInfoPath                    = "/v1/userinfo"
	OAuth2CallbackPath              = "/oauth2/callback"
)

type Paths struct {
	Home    string
	About   string
	Login   string
	Logout  string
	Profile string
}

func NewPaths() *Paths {
	return &Paths{
		Home:    HomePath,
		About:   AboutPath,
		Login:   LoginPath,
		Logout:  LogoutPath,
		Profile: ProfilePath,
	}
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

func GetMyRootPath(c echo.Context) string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
}
func GetGithubConfig(c echo.Context, protocol *proto_oidc_models.GithubOAuth2Protocol) *oauth2.Config {
	rootPath := GetMyRootPath(c)
	config := oauth2.Config{
		ClientID:     protocol.ClientId,
		ClientSecret: protocol.ClientSecret,
		Scopes:       GithubScopes,
		RedirectURL:  rootPath + OAuth2CallbackPath,
		Endpoint: oauth2.Endpoint{
			AuthURL:  GithubAuthURL,
			TokenURL: GithubTokenURL,
		},
	}
	return &config
}
