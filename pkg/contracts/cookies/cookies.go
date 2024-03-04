package cookies

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	echo "github.com/labstack/echo/v4"
)

type (
	SetExternalOauth2CookieRequest struct {
		ExternalOAuth2State *proto_oidc_models.ExternalOauth2State `json:"externalOAuth2State"`
	}
	GetExternalOauth2CookieResponse struct {
		ExternalOAuth2State *proto_oidc_models.ExternalOauth2State `json:"externalOAuth2State"`
	}
	SetVerificationCodeCookieRequest struct {
		VerificationCode *VerificationCode `json:"verificationCode"`
	}
	VerificationCode struct {
		Code    string `json:"code"`
		Email   string `json:"email"`
		Subject string `json:"subject"`
	}
	GetVerificationCodeCookieResponse struct {
		VerificationCode *VerificationCode `json:"verificationCode"`
	}
	PasswordReset struct {
		Subject string `json:"subject"`
	}
	SetPasswordResetCookieRequest struct {
		PasswordReset *PasswordReset `json:"passwordReset"`
	}
	GetPasswordResetCookieResponse struct {
		PasswordReset *PasswordReset `json:"passwordReset"`
	}
	AccountStateCookie struct {
		State string `json:"state"`
		Nonce string `json:"nonce"`
	}
	SetAccountStateCookieRequest struct {
		AccountStateCookie *AccountStateCookie `json:"accountStateCookie"`
	}
	GetAccountStateCookieResponse struct {
		AccountStateCookie *AccountStateCookie `json:"accountStateCookie"`
	}
	AuthCookie struct {
		Identity *proto_oidc_models.Identity `json:"identity"`
	}
	SetAuthCookieRequest struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	GetAuthCookieResponse struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	IWellknownCookies interface {
		// External OAuth2 Cookie
		//---------------------------------------------------------------------
		SetExternalOauth2Cookie(c echo.Context, request *SetExternalOauth2CookieRequest) error
		DeleteExternalOauth2CookieCookie(c echo.Context)
		GetExternalOauth2CookieCookie(c echo.Context) (*GetExternalOauth2CookieResponse, error)
		// Verification Code Cookie
		//---------------------------------------------------------------------
		SetVerificationCodeCookie(c echo.Context, request *SetVerificationCodeCookieRequest) error
		DeleteVerificationCodeCookie(c echo.Context)
		GetVerificationCodeCookie(c echo.Context) (*GetVerificationCodeCookieResponse, error)
		// Password Reset Cookie
		//---------------------------------------------------------------------
		SetPasswordResetCookie(c echo.Context, request *SetPasswordResetCookieRequest) error
		DeletePasswordResetCookie(c echo.Context)
		GetPasswordResetCookie(c echo.Context) (*GetPasswordResetCookieResponse, error)
		// Account State Cookie
		//---------------------------------------------------------------------
		SetAccountStateCookie(c echo.Context, request *SetAccountStateCookieRequest) error
		DeleteAccountStateCookie(c echo.Context)
		GetAccountStateCookie(c echo.Context) (*GetAccountStateCookieResponse, error)
		// Auth Cookie
		//---------------------------------------------------------------------
		SetAuthCookie(c echo.Context, request *SetAuthCookieRequest) error
		DeleteAuthCookie(c echo.Context)
		GetAuthCookie(c echo.Context) (*GetAuthCookieResponse, error)
		// Insecure Cookies
		//---------------------------------------------------------------------
		SetInsecureCookie(c echo.Context, name string, value interface{}) error
		DeleteInsecureCookie(c echo.Context, name string)
		GetInsecureCookie(c echo.Context, name string) (interface{}, error)
	}
)

const (
	CookieNameVerificationCode    = "_verificationCode"
	CookieNamePasswordReset       = "_passwordReset"
	CookieNameAccountState        = "_accountState"
	CookieNameAuth                = "_auth"
	LoginRequest                  = "_loginRequest"
	CookieNameExternalOauth2State = "_externalOauth2State"
)
