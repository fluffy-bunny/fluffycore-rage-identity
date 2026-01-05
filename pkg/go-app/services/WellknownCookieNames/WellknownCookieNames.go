package WellknownCookieNames

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
)

type (
	service struct {
		config      *contracts_cookies.WellknownCookieNamesConfig
		cookieNames map[contracts_cookies.CookieName]string
	}
)

var WellknownCookieNamesConfig *contracts_cookies.WellknownCookieNamesConfig

var stemService = (*service)(nil)
var _ contracts_cookies.IWellknownCookieNames = stemService

func (s *service) Ctor() (contracts_cookies.IWellknownCookieNames, error) {

	return &service{}, nil
}

func AddSingletonIWellknownCookieNames(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_cookies.IWellknownCookieNames](cb, stemService.Ctor)
}

func (s *service) ensureFastMap() {
	if s.cookieNames != nil {
		return // Already initialized
	}
	if WellknownCookieNamesConfig == nil {
		panic("WellknownCookieNames service not initialized properly: WellknownCookieNamesConfig is nil")
	}
	prefix := WellknownCookieNamesConfig.CookiePrefix
	cookieNames := map[contracts_cookies.CookieName]string{
		contracts_cookies.CookieName_VerificationCode:            prefix + "_verificationCode",
		contracts_cookies.CookieName_PasswordReset:               prefix + "_passwordReset",
		contracts_cookies.CookieName_AuthCompleted:               prefix + "_authCompleted",
		contracts_cookies.CookieName_AccountState:                prefix + "_accountState",
		contracts_cookies.CookieName_Auth:                        prefix + "_auth",
		contracts_cookies.CookieName_SSO:                         prefix + "_sso",
		contracts_cookies.CookieName_LoginRequest:                prefix + "_loginRequest",
		contracts_cookies.CookieName_ExternalOauth2StateTemplate: prefix + "_externalOauth2State_{state}",
		contracts_cookies.CookieName_WebAuthN:                    prefix + "_webAuthN",
		contracts_cookies.CookieName_SigninUserName:              prefix + "_signinUserName",
		contracts_cookies.CookieName_Error:                       prefix + "_error",
		contracts_cookies.CookieName_CSRF:                        "_csrf", // keep this for now, I think it is hard coded in fluffycore
		contracts_cookies.CookieName_AuthorizationState:          prefix + "_authorization_state",
		contracts_cookies.CookieName_AccountManagementSession:    prefix + "_account_management_session",
		contracts_cookies.CookieName_OIDCSession:                 prefix + "_oidc_session",
	}
	s.cookieNames = cookieNames

}
func (s *service) GetCookieName(cookieName contracts_cookies.CookieName) string {
	s.ensureFastMap()
	if name, ok := s.cookieNames[cookieName]; ok {
		return name
	}
	return ""
}
