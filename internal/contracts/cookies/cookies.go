package cookies

import (
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	echo "github.com/labstack/echo/v4"
)

type (
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
		Identity *models.Identity `json:"identity"`
	}
	SetAuthCookieRequest struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	GetAuthCookieResponse struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	IWellknownCookies interface {
		SetVerificationCodeCookie(c echo.Context, request *SetVerificationCodeCookieRequest) error
		DeleteVerificationCodeCookie(c echo.Context)
		GetVerificationCodeCookie(c echo.Context) (*GetVerificationCodeCookieResponse, error)
		SetPasswordResetCookie(c echo.Context, request *SetPasswordResetCookieRequest) error
		DeletePasswordResetCookie(c echo.Context)
		GetPasswordResetCookie(c echo.Context) (*GetPasswordResetCookieResponse, error)
		SetAccountStateCookie(c echo.Context, request *SetAccountStateCookieRequest) error
		DeleteAccountStateCookie(c echo.Context)
		GetAccountStateCookie(c echo.Context) (*GetAccountStateCookieResponse, error)
		SetAuthCookie(c echo.Context, request *SetAuthCookieRequest) error
		DeleteAuthCookie(c echo.Context)
		GetAuthCookie(c echo.Context) (*GetAuthCookieResponse, error)

		SetInsecureCookie(c echo.Context, name string, value interface{}) error
		DeleteInsecureCookie(c echo.Context, name string)
		GetInsecureCookie(c echo.Context, name string) (interface{}, error)
	}
)

const (
	CookieNameVerificationCode = "verificationCode"
	CookieNamePasswordReset    = "passwordReset"
	CookieNameAccountState     = "accountState"
	CookieNameAuth             = "_auth"
	LoginRequest               = "_loginRequest"
)
