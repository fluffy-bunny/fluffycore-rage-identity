package login_models

const (
	DIRECTIVE_Redirect                               = "redirect"
	DIRECTIVE_StartExternalLogin                     = "startExternalLogin"
	DIRECTIVE_LoginPhaseOne_UserDoesNotExist         = "userDoesNotExist"
	DIRECTIVE_LoginPhaseOne_DisplayPasswordPage      = "displayPasswordPage"
	DIRECTIVE_VerifyCode_DisplayVerifyCodePage       = "displayVerifyCodePage"
	DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage      = "displayLoginPhaseOnePage"
	DIRECTIVE_PasswordReset_DisplayPasswordResetPage = "displayPasswordResetPage"
)

type SignupErrorReason int

const (
	SignupErrorReason_NoError SignupErrorReason = iota
	SignupErrorReason_InvalidPassword
	SignupErrorReason_UserAlreadyExists
)

type PasswordResetErrorReason int

const (
	PasswordResetErrorReason_NoError PasswordResetErrorReason = iota
	PasswordResetErrorReason_InvalidPassword
)

type (
	LoginPhaseOneRequest struct {
		Email string `json:"email" validate:"required"`
	}
	DirectiveRedirect struct {
		RedirectURI string `json:"redirectUri"`
	}

	DirectiveStartExternalLogin struct {
		Slug string `json:"slug"`
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
		DirectiveStartExternalLogin  *DirectiveStartExternalLogin  `json:"directiveStartExternalLogin,omitempty"`
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
	SignupRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	SignupResponse struct {
		Email                       string                       `json:"email" validate:"required"`
		Directive                   string                       `json:"directive" validate:"required"`
		DirectiveRedirect           *DirectiveRedirect           `json:"directiveRedirect,omitempty"`
		DirectiveEmailCodeChallenge *DirectiveEmailCodeChallenge `json:"directiveEmailCodeChallenge,omitempty"`
		DirectiveStartExternalLogin *DirectiveStartExternalLogin `json:"directiveStartExternalLogin,omitempty"`
		Message                     string                       `json:"message,omitempty"`
		ErrorReason                 SignupErrorReason            `json:"errorReason,omitempty"`
	}
	LogoutRequest  struct{}
	LogoutResponse struct {
		Directive   string `json:"directive" validate:"required"`
		RedirectURL string `json:"redirectURL,omitempty"`
	}
	PasswordResetStartRequest struct {
		Email string `json:"email" validate:"required"`
	}
	PasswordResetStartResponse struct {
		Email                       string                       `json:"email" validate:"required"`
		Directive                   string                       `json:"directive" validate:"required"`
		DirectiveEmailCodeChallenge *DirectiveEmailCodeChallenge `json:"directiveEmailCodeChallenge,omitempty"`
	}
	PasswordResetFinishRequest struct {
		Password        string `json:"password" validate:"required"`
		PasswordConfirm string `json:"passwordConfirm" validate:"required"`
	}
	PasswordResetFinishResponse struct {
		Directive   string                   `json:"directive" validate:"required"`
		ErrorReason PasswordResetErrorReason `json:"errorReason,omitempty"`
	}
)
