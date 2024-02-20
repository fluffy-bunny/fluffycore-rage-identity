package cookies

import (
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
	IWellknownCookies interface {
		SetVerificationCodeCookie(c echo.Context, request *SetVerificationCodeCookieRequest) error
		DeleteVerificationCodeCookie(c echo.Context)
		GetVerificationCodeCookie(c echo.Context) (*GetVerificationCodeCookieResponse, error)
		SetPasswordResetCookie(c echo.Context, request *SetPasswordResetCookieRequest) error
		DeletePasswordResetCookie(c echo.Context)
		GetPasswordResetCookie(c echo.Context) (*GetPasswordResetCookieResponse, error)
	}
)
