package cookies

import (
	"encoding/json"
	"time"

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
		insecureCookies: insecureCookies,
		secureCookies:   secureCookesService,
		config:          config,
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
	b, err := json.Marshal(request)
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
			Name:     "verificationCode",
			Value:    value,
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(30 * time.Minute),
			Path:     "/",
		})
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteVerificationCodeCookie(c echo.Context) {
	s.secureCookies.DeleteCookie(c, "verificationCode")
}
func (s *service) GetVerificationCodeCookie(c echo.Context) (*contracts_cookies.GetVerificationCodeCookieResponse, error) {
	getCookieResponse, err := s.secureCookies.GetCookie(c, "verificationCode")
	if err != nil {
		return nil, err
	}
	if getCookieResponse.Value == nil {
		return nil, status.Error(codes.NotFound, "verificationCode not found")
	}

	bb, err := json.Marshal(getCookieResponse.Value)
	if err != nil {
		return nil, err
	}
	var value contracts_cookies.GetVerificationCodeCookieResponse
	err = json.Unmarshal(bb, &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}
