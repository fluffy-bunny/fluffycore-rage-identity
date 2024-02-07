package echo

const (
	HomePath                        = "/"
	LoginPath                       = "/login"
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
)

type Paths struct {
	Home  string
	About string
}

func NewPaths() *Paths {
	return &Paths{
		Home:  HomePath,
		About: AboutPath,
	}
}
