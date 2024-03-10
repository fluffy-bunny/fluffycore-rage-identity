package cookies

import (
	"encoding/json"
	"strings"
	"time"

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
func (s *service) validateSetVerificationCodeCookieRequest(c echo.Context, request *contracts_cookies.SetVerificationCodeCookieRequest) error {
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
	err := s.validateSetVerificationCodeCookieRequest(c, request)
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
func (s *service) validateSetPasswordResetCookieRequest(c echo.Context, request *contracts_cookies.SetPasswordResetCookieRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}

	if fluffycore_utils.IsEmptyOrNil(request.PasswordReset.Subject) {
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	return nil
}
func (s *service) SetPasswordResetCookie(c echo.Context, request *contracts_cookies.SetPasswordResetCookieRequest) error {
	err := s.validateSetPasswordResetCookieRequest(c, request)
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
func (s *service) validateSetAccountStateCookieRequest(c echo.Context, request *contracts_cookies.SetAccountStateCookieRequest) error {
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
func (s *service) SetAccountStateCookie(c echo.Context, request *contracts_cookies.SetAccountStateCookieRequest) error {
	err := s.validateSetAccountStateCookieRequest(c, request)
	if err != nil {
		return err
	}
	return SetCookie(c, s.cookieConfig, s.secureCookies, contracts_cookies.CookieNameAccountState, request.AccountStateCookie)
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

func (s *service) validateSetAuthCookieRequest(c echo.Context, request *contracts_cookies.SetAuthCookieRequest) error {
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
func (s *service) SetAuthCookie(c echo.Context, request *contracts_cookies.SetAuthCookieRequest) error {
	// TODO: Configurable expiration
	err := s.validateSetAuthCookieRequest(c, request)
	if err != nil {
		return err
	}
	b, err := json.Marshal(request.AuthCookie)
	if err != nil {
		return err
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	_, err = s.secureCookies.SetCookie(c,
		&fluffycore_contracts_cookies.SetCookieRequest{
			Name:     contracts_cookies.CookieNameAuth,
			Value:    value,
			HttpOnly: false,
			Expires:  time.Now().Add(30 * time.Minute),
			Path:     "/",
			Domain:   s.cookieConfig.Domain,
		})
	if err != nil {
		return err
	}
	return nil
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
	return SetCookie(c, s.cookieConfig, s.insecureCookies, name, value)
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
func (s *service) validateSetExternalOauth2CookieRequest(c echo.Context, request *contracts_cookies.SetExternalOauth2CookieRequest) error {
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
	err := s.validateSetExternalOauth2CookieRequest(c, request)
	if err != nil {
		return err
	}
	b, err := json.Marshal(request.ExternalOAuth2State)
	if err != nil {
		return err
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	cookieName := s.makeExternalOAuth2CookieName(request.State)
	_, err = s.secureCookies.SetCookie(c,
		&fluffycore_contracts_cookies.SetCookieRequest{
			Name:     cookieName,
			Value:    value,
			HttpOnly: false,
			Expires:  time.Now().Add(30 * time.Minute),
			Path:     "/",
			Domain:   s.cookieConfig.Domain,
		})
	if err != nil {
		return err
	}
	return nil
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

func (s *service) validateSetWebAuthNCookieRequest(c echo.Context, request *contracts_cookies.SetWebAuthNCookieRequest) error {
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
	err := s.validateSetWebAuthNCookieRequest(c, request)
	if err != nil {
		return err
	}
	b, err := json.Marshal(request.Value)
	if err != nil {
		return err
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	_, err = s.secureCookies.SetCookie(c,
		&fluffycore_contracts_cookies.SetCookieRequest{
			Name:     contracts_cookies.CookieNameWebAuthN,
			Value:    value,
			HttpOnly: false,
			Expires:  time.Now().Add(30 * time.Minute),
			Path:     "/",
			Domain:   s.cookieConfig.Domain,
		})
	if err != nil {
		return err
	}
	return nil
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
