package cookies

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
	echo "github.com/labstack/echo/v4"
)

type VerifyCodePurpose int

const (
	VerifyCode_PasswordReset VerifyCodePurpose = iota + 1
	VerifyCode_EmailVerification
	VerifyCode_Challenge
)

type CookieName int

const (
	CookieName_VerificationCode CookieName = iota + 1
	CookieName_PasswordReset
	CookieName_AuthCompleted
	CookieName_AccountState
	CookieName_Auth
	CookieName_SSO
	CookieName_SkipKeepSignedIn
	CookieName_LoginRequest
	CookieName_ExternalOauth2StateTemplate
	CookieName_WebAuthN
	CookieName_SigninUserName
	CookieName_Error
	CookieName_CSRF
	CookieName_AuthorizationState
	CookieName_AccountManagementSession
	CookieName_OIDCSession
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
		CodeHash          string            `json:"codeHash"`
		PlainCode         string            `json:"plainCode,omitempty"`
		Email             string            `json:"email"`
		Subject           string            `json:"subject"`
		VerifyCodePurpose VerifyCodePurpose `json:"verifyCodePurpose"`
		DevelopmentMode   bool              `json:"developmentMode"`
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
	AuthCompleted struct {
		Subject string `json:"subject"`
	}
	SetAuthCompletedCookieRequest struct {
		AuthCompleted *AuthCompleted `json:"authCompleted"`
	}
	GetAuthCompletedCookieResponse struct {
		AuthCompleted *AuthCompleted `json:"authCompleted"`
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
		Acr      []string                    `json:"acr,omitempty"`
		Amr      []string                    `json:"amr,omitempty"`
	}
	SetAuthCookieRequest struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	GetAuthCookieResponse struct {
		AuthCookie *AuthCookie `json:"authCookie"`
	}
	SSOCookie struct {
		Identity *proto_oidc_models.Identity `json:"identity"`
		Acr      []string                    `json:"acr,omitempty"`
		Amr      []string                    `json:"amr,omitempty"`
	}
	SetSSOCookieRequest struct {
		SSOCookie *SSOCookie `json:"ssoCookie"`
	}
	GetSSOCookieResponse struct {
		SSOCookie *SSOCookie `json:"ssoCookie"`
	}
	KeepSigninPreferencesCookie struct {
		PreferenceValue bool `json:"preferenceValue"`
	}
	SetKeepSigninPreferencesCookieRequest struct {
		Subject                     string                       `json:"subject"`
		KeepSigninPreferencesCookie *KeepSigninPreferencesCookie `json:"keepSigninPreferencesCookie"`
	}
	GetKeepSigninPreferencesCookieRequest struct {
		Subject string `json:"subject"`
	}
	GetKeepSigninPreferencesCookieResponse struct {
		KeepSigninPreferencesCookie *KeepSigninPreferencesCookie `json:"keepSigninPreferencesCookie"`
	}
	DeleteKeepSigninPreferencesCookieRequest struct {
		Subject string `json:"subject"`
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
	ErrorCookie struct {
		Code   string            `json:"code"`
		Error  string            `json:"error"`
		Params map[string]string `json:"params"`
	}
	SetErrorCookieRequest struct {
		Value *ErrorCookie `json:"errorCookie"`
	}
	GetErrorCookieResponse struct {
		Value *ErrorCookie `json:"errorCookie"`
	}
	WellknownCookieNamesConfig struct {
		CookiePrefix string `json:"cookiePrefix"`
	}

	IWellknownCookieNames interface {
		// Cookie Name
		//---------------------------------------------------------------------
		GetCookieName(cookieName CookieName) string
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
		// Auth Completed Cookie
		//---------------------------------------------------------------------
		SetAuthCompletedCookie(c echo.Context, request *SetAuthCompletedCookieRequest) error
		DeleteAuthCompletedCookie(c echo.Context)
		GetAuthCompletedCookie(c echo.Context) (*GetAuthCompletedCookieResponse, error)
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
		// SSO Cookie
		//---------------------------------------------------------------------
		SetSSOCookie(c echo.Context, request *SetSSOCookieRequest) error
		DeleteSSOCookie(c echo.Context)
		GetSSOCookie(c echo.Context) (*GetSSOCookieResponse, error)
		// KeepSigninPreferences Cookie
		//---------------------------------------------------------------------
		SetKeepSigninPreferencesCookie(c echo.Context, request *SetKeepSigninPreferencesCookieRequest) error
		DeleteKeepSigninPreferencesCookie(c echo.Context, request *DeleteKeepSigninPreferencesCookieRequest)
		GetKeepSigninPreferencesCookie(c echo.Context, request *GetKeepSigninPreferencesCookieRequest) (*GetKeepSigninPreferencesCookieResponse, error)
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

		// SetErrorCookie Cookie
		//---------------------------------------------------------------------
		SetErrorCookie(c echo.Context, request *SetErrorCookieRequest) error
		DeleteErrorCookie(c echo.Context)
		GetErrorCookie(c echo.Context) (*GetErrorCookieResponse, error)
		// Legacy string constants - deprecated, use CookieName enum with GetCookieName() instead
	}
)

/*
const (
	CookieNameVerificationCode            = "_rage_verificationCode"
	CookieNamePasswordReset               = "_rage_passwordReset"
	CookieNameAuthCompleted               = "_rage_authCompleted"
	CookieNameAccountState                = "_rage_accountState"
	CookieNameAuth                        = "_rage_auth"
	CookieNameSSO                         = "_rage_sso"
	LoginRequest                          = "_rage_loginRequest"
	CookieNameExternalOauth2StateTemplate = "_rage_externalOauth2State_{state}"
	CookieNameWebAuthN                    = "_rage_webAuthN"
	CookieNameSigninUserName              = "_rage_signinUserName"
	CookieNameErrorName                   = "_rage_error"
	CookieNameCSRF                        = "_csrf" // keep this for now, I think it is hard coded in fluffycore
	CookieNameAuthorizationState          = "_rage_authorization_state"
	CookieNameAccountManagementSession    = "_rage_account_management_session"
	CookieNameOIDCSession                 = "_rage_oidc_session"
)
*/
