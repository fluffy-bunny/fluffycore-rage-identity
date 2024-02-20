package cookies

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
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
	return SetCookie(c, s.secureCookies, "verificationCode", request.VerificationCode)

}
func (s *service) DeleteVerificationCodeCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c, "verificationCode")
}
func (s *service) GetVerificationCodeCookie(c echo.Context) (*contracts_cookies.GetVerificationCodeCookieResponse, error) {

	var value contracts_cookies.VerificationCode
	err := GetCookie(c, s.secureCookies, "verificationCode", &value)
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
	return SetCookie(c, s.secureCookies, "passwordReset", request.PasswordReset)
}
func (s *service) DeletePasswordResetCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c, "passwordReset")
}
func (s *service) GetPasswordResetCookie(c echo.Context) (*contracts_cookies.GetPasswordResetCookieResponse, error) {

	var value contracts_cookies.PasswordReset
	err := GetCookie(c, s.secureCookies, "passwordReset", &value)
	if err != nil {
		return nil, err
	}
	return &contracts_cookies.GetPasswordResetCookieResponse{
		PasswordReset: &value,
	}, nil
}
