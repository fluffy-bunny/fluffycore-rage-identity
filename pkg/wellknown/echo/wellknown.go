package echo

import (
	"fmt"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	echo "github.com/labstack/echo/v4"
	oauth2 "golang.org/x/oauth2"
)

var (
	AboutPath           = "/about"
	AccountCallbackPath = "/auth/callback"
	ErrorPath           = "/error"
	ExternalIDPPath     = "/external-idp"
	ForgotPasswordPath  = "/forgot-password"
	HealthzPath         = "/healthz"
	HomePath            = "/"
	LoginPath           = "/login"
	LogoutPath          = "/logout"
	StaticPath          = "/static*"
	//OAuth2CallbackPath                               = "/oauth2/callback"
	OAuth2CallbackPath                               = "@@OAuth2CallbackPath@@"
	OAuth2TokenEndpointPath                          = "/token"
	OIDCAuthorizationEndpointPath                    = "/oidc/v1/auth"
	OIDCLoginPath                                    = "/oidc-login"
	OIDCLoginPasswordPath                            = "/oidc-login-password"
	PasswordResetPath                                = "/password-reset"
	ProfilePath                                      = "/profile"
	PersonalInformationPath                          = "/profile/personal-information"
	PasskeyManagementPath                            = "/passkey-management"
	ReadyPath                                        = "/ready"
	SignupPath                                       = "/signup"
	SwaggerPath                                      = "/swagger/*"
	UserInfoPath                                     = "/v1/userinfo"
	VerifyCodePath                                   = "/verify-code"
	WellKnownJWKS                                    = "/.well-known/jwks"
	WellKnownOpenIDCOnfiguationPath                  = "/.well-known/openid-configuration"
	WebAuthN_Register_GetCredentialCreationOptions   = "/webauthn/register/get_credential_creation_options"
	WebAuthN_Register_ProcessRegistrationAttestation = "/webauthn/register/process_registration_attestation"
	WebAuthN_Login_GetCredentialRequestOptions       = "/webauthn/login/get_credential_request_options"
	WebAuthN_Login_ProcessLoginAssertion             = "/webauthn/login/process_login_assertion"
	WebAuthN_Register_Begin                          = "/webauthn/register/begin"
	WebAuthN_Register_Finish                         = "/webauthn/register/finish"
	WebAuthN_Login_Begin                             = "/webauthn/login/begin"
	WebAuthN_Login_Finish                            = "/webauthn/login/finish"
)

type Paths struct {
	About               string
	ExternalIDP         string
	Home                string
	Login               string
	Logout              string
	PasskeyManagement   string
	PersonalInformation string
	Profile             string
	OIDCLogin           string
	OIDCLoginPassword   string
	Signup              string
	ForgotPassword      string
	VerifyCode          string
	PasswordReset       string
}

func NewPaths() *Paths {
	return &Paths{
		About:               AboutPath,
		ExternalIDP:         ExternalIDPPath,
		Home:                HomePath,
		Login:               LoginPath,
		Logout:              LogoutPath,
		PasskeyManagement:   PasskeyManagementPath,
		PersonalInformation: PersonalInformationPath,
		Profile:             ProfilePath,
		OIDCLogin:           OIDCLoginPath,
		OIDCLoginPassword:   OIDCLoginPasswordPath,
		Signup:              SignupPath,
		ForgotPassword:      ForgotPasswordPath,
		VerifyCode:          VerifyCodePath,
		PasswordReset:       PasswordResetPath,
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
