package RequiresNoAuth

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_auth "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/auth"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
)

type (
	service struct {
		authMap map[string]bool
		config  *contracts_config.Config
	}
)

var _ contracts_auth.IRequiresNoAuth = (*service)(nil)
var stemService = (*service)(nil)

func (s *service) Ctor(
	config *contracts_config.Config,
) (contracts_auth.IRequiresNoAuth, error) {
	ss := &service{
		authMap: RequiresNoAuth(),
		config:  config,
	}
	for _, path := range config.RequiresNoAuthConfig.Remove {
		delete(ss.authMap, path)
	}
	for _, path := range config.RequiresNoAuthConfig.Add {
		ss.authMap[path] = true
	}
	return ss, nil
}

func AddSingletonIRequiresNoAuth(builder di.ContainerBuilder) {
	di.AddSingleton[contracts_auth.IRequiresNoAuth](builder, stemService.Ctor)
}
func (s *service) GetAuthMap() map[string]bool {
	return s.authMap
}

// everything requries auth unless otherwise documented here.
// -- this is a list of paths that do not require auth
func RequiresNoAuth() map[string]bool {

	return map[string]bool{
		wellknown_echo.API_Login: true,

		wellknown_echo.ManagementPath:      true,
		wellknown_echo.ManagementAllPath:   true,
		wellknown_echo.StaticPath:          true,
		wellknown_echo.AboutPath:           true,
		wellknown_echo.APIPath:             true,
		wellknown_echo.AccountCallbackPath: true,
		wellknown_echo.ErrorPath:           true,
		wellknown_echo.ExternalIDPPath:     true,
		wellknown_echo.ForgotPasswordPath:  true,
		wellknown_echo.HealthzPath:         true,
		//	wellknown_echo.HomePath:                      true,
		wellknown_echo.LoginPath:                     true,
		wellknown_echo.LogoutPath:                    true,
		wellknown_echo.OAuth2CallbackPath:            true,
		wellknown_echo.OAuth2TokenEndpointPath:       true,
		wellknown_echo.OIDCAuthorizationEndpointPath: true,
		wellknown_echo.OIDCLoginPath:                 true,
		wellknown_echo.OIDCLoginUIPath:               true,
		wellknown_echo.OIDCLoginUIStaticPath:         true,

		wellknown_echo.API_OIDCFlowAppConfig:      true,
		wellknown_echo.API_AppSettings:            true,
		wellknown_echo.API_Manifest:               true,
		wellknown_echo.API_StartOver:              true,
		wellknown_echo.API_Start_ExternalLogin:    true,
		wellknown_echo.API_VerifyUsername:         true,
		wellknown_echo.API_VerifyPasswordStrength: true,
		wellknown_echo.API_LoginPhaseOne:          true,
		wellknown_echo.API_LoginPassword:          true,
		wellknown_echo.API_VerifyCode:             true,
		wellknown_echo.API_VerifyCodeBegin:        true,
		wellknown_echo.API_Signup:                 true,
		wellknown_echo.API_Logout:                 true,
		wellknown_echo.API_PasswordResetStart:     true,
		wellknown_echo.API_PasswordResetFinish:    true,
		wellknown_echo.API_KeepSignedIn:           true,

		wellknown_echo.OIDCLoginPasskeyPath:            true,
		wellknown_echo.OIDCLoginPasswordPath:           true,
		wellknown_echo.OIDCLoginTOTPPath:               true,
		wellknown_echo.PasswordResetPath:               true,
		wellknown_echo.ReadyPath:                       true,
		wellknown_echo.SignupPath:                      true,
		wellknown_echo.SwaggerPath:                     true,
		wellknown_echo.UserInfoPath:                    true,
		wellknown_echo.VerifyCodePath:                  true,
		wellknown_echo.WellKnownJWKS:                   true,
		wellknown_echo.WellKnownOpenIDCOnfiguationPath: true,
		// WebAuthN Registrationhandlers: Must be authenticated
		//----------------------------------------------------
		//			wellknown_echo.WebAuthN_Register_Begin:  true,
		//			wellknown_echo.WebAuthN_Register_Finish: true,
		// WebAuthN Loginhandlers: Must NOT be authenticated
		//----------------------------------------------------
		wellknown_echo.WebAuthN_Login_Begin:  true,
		wellknown_echo.WebAuthN_Login_Finish: true,
	}

}
