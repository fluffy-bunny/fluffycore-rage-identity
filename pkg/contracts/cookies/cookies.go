package cookies

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
	echo "github.com/labstack/echo/v4"
)

type (
	SetExternalOauth2CookieRequest struct {
		State               string                                 `json:"state"`
		ExternalOAuth2State *proto_oidc_models.ExternalOauth2State `json:"externalOAuth2State"`
	}
	DeleteExternalOauth2CookieRequest struct {
		State string `json:"state"`
	}
	GetExternalOauth2CookieRequest struct {
		State string `json:"state"`
	}
	GetExternalOauth2CookieResponse struct {
		State               string                                 `json:"state"`
		ExternalOAuth2State *proto_oidc_models.ExternalOauth2State `json:"externalOAuth2State"`
	}
	SetVerificationCodeCookieRequest struct {
		VerificationCode *VerificationCode `json:"verificationCode"`
	}
	VerificationCode struct {
		Code             string `json:"code"`
		Email            string `json:"email"`
		Subject          string `json:"subject"`
		PasswordVerified bool   `json:"passwordVerified"`
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
	WebAuthNCookie struct {
		Identity    *proto_oidc_models.Identity `json:"identity"`
		SessionData *go_webauthn.SessionData    `json:"sessionData"`
	}
	SetWebAuthNCookieRequest struct {
		Value *WebAuthNCookie `json:"webAuthNCookie"`
	}
	GetWebAuthNCookieResponse struct {
		Value *WebAuthNCookie `json:"webAuthNCookie"`
	}
	SigninUserNameCookie struct {
		Email      string `json:"email"`
		HasPasskey bool   `json:"hasPasskey"`
	}
	SetSigninUserNameCookieRequest struct {
		Value *SigninUserNameCookie `json:"signinUserNameCookie"`
	}
	GetSigninUserNameCookieResponse struct {
		Value *SigninUserNameCookie `json:"signinUserNameCookie"`
	}
	IWellknownCookies interface {
		// External OAuth2 Cookie
		//---------------------------------------------------------------------
		SetExternalOauth2Cookie(c echo.Context, request *SetExternalOauth2CookieRequest) error
		DeleteExternalOauth2Cookie(c echo.Context, request *DeleteExternalOauth2CookieRequest) error
		GetExternalOauth2Cookie(c echo.Context, request *GetExternalOauth2CookieRequest) (*GetExternalOauth2CookieResponse, error)
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
		// WebAuthN Cookie
		//---------------------------------------------------------------------
		SetWebAuthNCookie(c echo.Context, request *SetWebAuthNCookieRequest) error
		DeleteWebAuthNCookie(c echo.Context)
		GetWebAuthNCookie(c echo.Context) (*GetWebAuthNCookieResponse, error)

		// SigninUserName Cookie
		//---------------------------------------------------------------------
		SetSigninUserNameCookie(c echo.Context, request *SetSigninUserNameCookieRequest) error
		DeleteSigninUserNameCookie(c echo.Context)
		GetSigninUserNameCookie(c echo.Context) (*GetSigninUserNameCookieResponse, error)
	}
)

const (
	CookieNameVerificationCode            = "_verificationCode"
	CookieNamePasswordReset               = "_passwordReset"
	CookieNameAccountState                = "_accountState"
	CookieNameAuth                        = "_auth"
	LoginRequest                          = "_loginRequest"
	CookieNameExternalOauth2StateTemplate = "_externalOauth2State_{state}"
	CookieNameWebAuthN                    = "_webAuthN"
	CookieNameSigninUserName              = "_signinUserName"
)
