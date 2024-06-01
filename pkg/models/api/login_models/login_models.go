package login_models

import (
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
)

const (
	DIRECTIVE_LoginPhaseOne_Redirect                      = "redirect"
	DIRECTIVE_LoginPhaseOne_UserDoesNotExist              = "userDoesNotExist"
	DIRECTIVE_LoginPhaseOne_DisplayPasswordPage           = "displayPasswordPage"
	DIRECTIVE_LoginPhaseOne_DisplayEmailVerificationPage  = "displayEmailVerificationPage"
	DIRECTIVE_LoginPassword_DisplayEmailCodeChallengePage = "displayEmailCodeChallengePage"
	DIRECTIVE_LoginPassword_Redirect                      = "redirect"
	DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage           = "displayLoginPhaseOnePage"
	DIRECTIVE_PassowrdReset_DisplayPasswordResetPage      = "displayPasswordResetPage"

	DIRECTIVE_VerifyCode_Redirect = "redirect"
)

type (
	LoginPhaseOneRequest struct {
		Email string `json:"email" validate:"required"`
	}
	DirectiveRedirect struct {
		RedirectURI string             `json:"redirectUri"`
		VERB        string             `json:"verb"`
		FormParams  []models.FormParam `json:"formParams"`
	}
	DirectiveDisplayPasswordPage struct {
		Email      string `json:"email"`
		HasPasskey bool   `json:"hasPasskey"`
	}

	DirectiveEmailCodeChallenge struct {
		Code string `json:"code"`
	}

	LoginPhaseOneResponse struct {
		Email                        string                        `json:"email" validate:"required"`
		Directive                    string                        `json:"directive" validate:"required"`
		DirectiveRedirect            *DirectiveRedirect            `json:"directiveRedirect,omitempty"`
		DirectiveDisplayPasswordPage *DirectiveDisplayPasswordPage `json:"directiveDisplayPasswordPage,omitempty"`
		DirectiveEmailCodeChallenge  *DirectiveEmailCodeChallenge  `json:"directiveEmailCodeChallenge,omitempty"`
	}

	LoginPasswordRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	LoginPasswordResponse struct {
		Email                       string                       `json:"email" validate:"required"`
		Directive                   string                       `json:"directive" validate:"required"`
		DirectiveRedirect           *DirectiveRedirect           `json:"directiveRedirect,omitempty"`
		DirectiveEmailCodeChallenge *DirectiveEmailCodeChallenge `json:"directiveEmailCodeChallenge,omitempty"`
	}
	VerifyCodeRequest struct {
		Code string `json:"code" validate:"required"`
	}
	VerifyCodeResponse struct {
		Directive         string             `json:"directive" validate:"required"`
		DirectiveRedirect *DirectiveRedirect `json:"directiveRedirect,omitempty"`
	}
)
