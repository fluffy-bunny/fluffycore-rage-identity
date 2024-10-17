package cookies

import (
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_contracts_cookies "github.com/fluffy-bunny/fluffycore/echo/contracts/cookies"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		CustomCookieBase
		insecureCookies fluffycore_contracts_cookies.ICookies
		secureCookies   fluffycore_contracts_cookies.ICookies
		config          *contracts_config.EchoConfig
		cookieConfig    *contracts_config.CookieConfig
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_cookies.IWellknownCookies = stemService
}
func (s *service) Ctor(
	insecureCookies fluffycore_contracts_cookies.ICookies,
	secureCookies fluffycore_contracts_cookies.ISecureCookies,
	config *contracts_config.EchoConfig,
	cookieConfig *contracts_config.CookieConfig,
) (contracts_cookies.IWellknownCookies, error) {

	var secureCookesService fluffycore_contracts_cookies.ICookies
	secureCookesService = secureCookies
	if config.DisableSecureCookies {
		secureCookesService = insecureCookies
	}

	return &service{
		CustomCookieBase: CustomCookieBase{},
		insecureCookies:  insecureCookies,
		secureCookies:    secureCookesService,
		config:           config,
		cookieConfig:     cookieConfig,
	}, nil
}

func AddSingletonIWellknownCookies(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_cookies.IWellknownCookies](cb, stemService.Ctor)
}
func (s *service) validateSetVerificationCodeCookieRequest(request *contracts_cookies.SetVerificationCodeCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.VerificationCode == nil {
		return status.Error(codes.InvalidArgument, "request.VerificationCode is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.VerificationCode.Code) {
		return status.Error(codes.InvalidArgument, "Code is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.VerificationCode.Email) {
		return status.Error(codes.InvalidArgument, "Email is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.VerificationCode.Subject) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	return nil
}
func (s *service) SetVerificationCodeCookie(c echo.Context, request *contracts_cookies.SetVerificationCodeCookieRequest) error {
	err := s.validateSetVerificationCodeCookieRequest(request)
	if err != nil {
		return err
	}
	return SetCookie(c, s.cookieConfig, s.secureCookies, contracts_cookies.CookieNameVerificationCode, request.VerificationCode)

}
func (s *service) DeleteVerificationCodeCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameVerificationCode,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetVerificationCodeCookie(c echo.Context) (*contracts_cookies.GetVerificationCodeCookieResponse, error) {

	var value contracts_cookies.VerificationCode
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNameVerificationCode, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetVerificationCodeCookieResponse{
		VerificationCode: &value,
	}, nil
}
func (s *service) validateSetPasswordResetCookieRequest(request *contracts_cookies.SetPasswordResetCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}

	if fluffycore_utils.IsEmptyOrNil(request.PasswordReset.Subject) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	return nil
}
func (s *service) SetPasswordResetCookie(c echo.Context, request *contracts_cookies.SetPasswordResetCookieRequest) error {
	err := s.validateSetPasswordResetCookieRequest(request)
	if err != nil {
		return err
	}
	return SetCookie(c, s.cookieConfig, s.secureCookies, contracts_cookies.CookieNamePasswordReset, request.PasswordReset)
}
func (s *service) DeletePasswordResetCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNamePasswordReset,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})

}
func (s *service) GetPasswordResetCookie(c echo.Context) (*contracts_cookies.GetPasswordResetCookieResponse, error) {

	var value contracts_cookies.PasswordReset
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNamePasswordReset, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetPasswordResetCookieResponse{
		PasswordReset: &value,
	}, nil
}
func (s *service) validateSetAccountStateCookieRequest(request *contracts_cookies.SetAccountStateCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.AccountStateCookie == nil {
		return status.Error(codes.InvalidArgument, "request.AccountStateCookie is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.AccountStateCookie.State) {
		return status.Error(codes.InvalidArgument, "State is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.AccountStateCookie.Nonce) {
		return status.Error(codes.InvalidArgument, "Nonce is empty")
	}
	return nil

}
func tPtr[T any](t T) *T {
	return &t
}
func (s *service) SetAccountStateCookie(c echo.Context, request *contracts_cookies.SetAccountStateCookieRequest) error {
	err := s.validateSetAccountStateCookieRequest(request)
	if err != nil {
		return err
	}
	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: contracts_cookies.CookieNameAccountState,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest, request.AccountStateCookie)

}
func (s *service) DeleteAccountStateCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameAccountState,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})

}
func (s *service) GetAccountStateCookie(c echo.Context) (*contracts_cookies.GetAccountStateCookieResponse, error) {
	var value contracts_cookies.AccountStateCookie
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNameAccountState, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetAccountStateCookieResponse{
		AccountStateCookie: &value,
	}, nil
}

func (s *service) validateSetAuthCookieRequest(request *contracts_cookies.SetAuthCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.AuthCookie == nil {
		return status.Error(codes.InvalidArgument, "request.AuthCookie is nil")
	}
	if request.AuthCookie.Identity == nil {
		return status.Error(codes.InvalidArgument, "request.AuthCookie.Identity is nil")
	}
	return nil
}
func (s *service) SetAuthCookie(c echo.Context,
	request *contracts_cookies.SetAuthCookieRequest) error {
	// TODO: Configurable expiration
	err := s.validateSetAuthCookieRequest(request)
	if err != nil {
		return err
	}
	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: contracts_cookies.CookieNameAuth,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest,
		request.AuthCookie)
}
func (s *service) DeleteAuthCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameAuth,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetAuthCookie(c echo.Context) (*contracts_cookies.GetAuthCookieResponse, error) {
	var value contracts_cookies.AuthCookie
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNameAuth, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetAuthCookieResponse{
		AuthCookie: &value,
	}, nil
}
func (s *service) SetInsecureCookie(c echo.Context, name string, value interface{}) error {

	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: name,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.insecureCookies, setCookieRequest,
		value)

}
func (s *service) DeleteInsecureCookie(c echo.Context, name string) {
	s.insecureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   name,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetInsecureCookie(c echo.Context, name string) (interface{}, error) {
	var value interface{}
	err := GetCookie(c, s.insecureCookies, name, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (s *service) validateSetExternalOauth2CookieRequest(request *contracts_cookies.SetExternalOauth2CookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.State) {
		return status.Error(codes.InvalidArgument, "request.State is empty")
	}
	if request.ExternalOAuth2State == nil {
		return status.Error(codes.InvalidArgument, "request.ExternalOAuth2State is nil")
	}
	if request.ExternalOAuth2State.Request == nil {
		return status.Error(codes.InvalidArgument, "request.ExternalOAuth2State.Request is nil")
	}
	return nil
}
func (s *service) makeExternalOAuth2CookieName(state string) string {
	if fluffycore_utils.IsEmptyOrNil(state) {
		panic("state is empty")
	}
	result := strings.ReplaceAll(contracts_cookies.CookieNameExternalOauth2StateTemplate, "{state}", state)
	return result
}
func (s *service) SetExternalOauth2Cookie(c echo.Context, request *contracts_cookies.SetExternalOauth2CookieRequest) error {
	err := s.validateSetExternalOauth2CookieRequest(request)
	if err != nil {
		return err
	}
	cookieName := s.makeExternalOAuth2CookieName(request.State)

	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: cookieName,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest,
		request.ExternalOAuth2State)

}
func (s *service) validateDeleteExternalOauth2CookieRequest(request *contracts_cookies.DeleteExternalOauth2CookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.State) {
		return status.Error(codes.InvalidArgument, "request.State is empty")
	}
	return nil
}
func (s *service) DeleteExternalOauth2Cookie(c echo.Context, request *contracts_cookies.DeleteExternalOauth2CookieRequest) error {
	err := s.validateDeleteExternalOauth2CookieRequest(request)
	if err != nil {
		return err
	}
	cookieName := s.makeExternalOAuth2CookieName(request.State)

	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   cookieName,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
	return nil
}
func (s *service) validateGetExternalOauth2CookieRequest(request *contracts_cookies.GetExternalOauth2CookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.State) {
		return status.Error(codes.InvalidArgument, "request.State is empty")
	}
	return nil
}
func (s *service) GetExternalOauth2Cookie(c echo.Context, request *contracts_cookies.GetExternalOauth2CookieRequest) (*contracts_cookies.GetExternalOauth2CookieResponse, error) {
	err := s.validateGetExternalOauth2CookieRequest(request)
	if err != nil {
		return nil, err
	}
	var value proto_oidc_models.ExternalOauth2State
	cookieName := s.makeExternalOAuth2CookieName(request.State)
	err = GetCookie(c, s.secureCookies, cookieName, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetExternalOauth2CookieResponse{
		ExternalOAuth2State: &value,
	}, nil
}

// WebAuthN Cookie
// ---------------------------------------------------------------------
func (s *service) validateSetWebAuthNCookieRequest(request *contracts_cookies.SetWebAuthNCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.Value == nil {
		return status.Error(codes.InvalidArgument, "request.Value is nil")
	}
	if request.Value.Identity == nil {
		return status.Error(codes.InvalidArgument, "request.Value.Identity is nil")
	}
	return nil
}
func (s *service) SetWebAuthNCookie(c echo.Context, request *contracts_cookies.SetWebAuthNCookieRequest) error {
	// TODO: Configurable expiration
	err := s.validateSetWebAuthNCookieRequest(request)
	if err != nil {
		return err
	}
	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: contracts_cookies.CookieNameWebAuthN,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest,
		request.Value)

}
func (s *service) DeleteWebAuthNCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameWebAuthN,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetWebAuthNCookie(c echo.Context) (*contracts_cookies.GetWebAuthNCookieResponse, error) {
	var value contracts_cookies.WebAuthNCookie
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNameWebAuthN, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetWebAuthNCookieResponse{
		Value: &value,
	}, nil
}

// SigninUserName Cookie
// ---------------------------------------------------------------------
func (s *service) validateSetSigninUserNameCookieRequest(request *contracts_cookies.SetSigninUserNameCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.Value == nil {
		return status.Error(codes.InvalidArgument, "request.Value is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Value.Email) {
		return status.Error(codes.InvalidArgument, "Email is empty")

	}
	return nil
}
func (s *service) SetSigninUserNameCookie(c echo.Context, request *contracts_cookies.SetSigninUserNameCookieRequest) error {
	// TODO: Configurable expiration
	err := s.validateSetSigninUserNameCookieRequest(request)
	if err != nil {
		return err
	}
	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: contracts_cookies.CookieNameSigninUserName,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest,
		request.Value)
}
func (s *service) DeleteSigninUserNameCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameSigninUserName,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetSigninUserNameCookie(c echo.Context) (*contracts_cookies.GetSigninUserNameCookieResponse, error) {
	var value contracts_cookies.SigninUserNameCookie
	err := GetCookie(c, s.secureCookies, contracts_cookies.CookieNameSigninUserName, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetSigninUserNameCookieResponse{
		Value: &value,
	}, nil
}

// SigninUserName Cookie
// ---------------------------------------------------------------------
func (s *service) validateSetErrorCookieRequest(request *contracts_cookies.SetErrorCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if request.Value == nil {
		return status.Error(codes.InvalidArgument, "request.Value is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Value.Code) {
		return status.Error(codes.InvalidArgument, "Code is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Value.Error) {
		return status.Error(codes.InvalidArgument, "Error is empty")
	}
	return nil
}

func (s *service) SetErrorCookie(c echo.Context, request *contracts_cookies.SetErrorCookieRequest) error {
	// TODO: Configurable expiration
	err := s.validateSetErrorCookieRequest(request)
	if err != nil {
		return err
	}

	setCookieRequest := &fluffycore_contracts_cookies.SetCookieRequest{
		Name: contracts_cookies.CookieNameErrorName,
	}
	if s.config.DisableSecureCookies {
		setCookieRequest.HttpOnly = true
		setCookieRequest.SameSite = 0
		setCookieRequest.Secure = tPtr(false)
	}
	return SetCookieByRequest(c, s.cookieConfig, s.secureCookies, setCookieRequest,
		request.Value)

}
func (s *service) DeleteErrorCookie(c echo.Context) {
	s.insecureCookies.DeleteCookie(c,
		&fluffycore_contracts_cookies.DeleteCookieRequest{
			Name:   contracts_cookies.CookieNameErrorName,
			Path:   "/",
			Domain: s.cookieConfig.Domain,
		})
}
func (s *service) GetErrorCookie(c echo.Context) (*contracts_cookies.GetErrorCookieResponse, error) {
	var value contracts_cookies.ErrorCookie
	err := GetCookie(c, s.insecureCookies, contracts_cookies.CookieNameErrorName, &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetErrorCookieResponse{
		Value: &value,
	}, nil
}
