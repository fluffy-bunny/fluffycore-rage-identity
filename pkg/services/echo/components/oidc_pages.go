package components

import (
	"fmt"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// --- OIDC Login (React SPA shell) ---

// OIDCLoginPage renders the SPA shell for the React-based OIDC login flow.
func OIDCLoginPage(rc *RenderContext) g.Node {
	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		g.El("noscript", g.Text("You need to enable JavaScript to run this app.")),
		Div(ID("root")),
		Script(g.Attr("defer", "defer"), Src("/static/oidc-flows/static/js/main.js")),
	)
}

// --- OIDC Login Password ---

// OIDCLoginPasswordData holds data for the password login page.
type OIDCLoginPasswordData struct {
	Errors     []string
	Email      string
	Directive  string
	HasPasskey bool
	IDPs       []*proto_oidc_models.IDP
}

// OIDCLoginPasswordPage renders the password login form.
func OIDCLoginPasswordPage(rc *RenderContext, data OIDCLoginPasswordData) g.Node {
	idpButtons := make([]g.Node, 0, len(data.IDPs))
	for _, idp := range data.IDPs {
		idpButtons = append(idpButtons,
			FormEl(g.Attr("action", rc.Paths.ExternalIDP), Method("post"),
				CsrfInput(rc.CSRF),
				Input(Type("hidden"), Name("directive"), Value(data.Directive)),
				Input(Type("hidden"), Name("idp_hint"), Value(idp.GetSlug())),
				Button(Type("submit"), Class("btn btn-outline-primary me-2"), g.Text(idp.GetName())),
			),
		)
	}

	var passkeySection g.Node
	if data.HasPasskey {
		passkeySection = g.Group([]g.Node{
			Hr(),
			Div(Class("d-flex justify-content-center"),
				FormEl(g.Attr("action", rc.Paths.OIDCLoginPasskey), Method("post"),
					CsrfInput(rc.CSRF),
					Button(Type("submit"), Class("btn btn-outline-primary me-2"), g.Text(rc.L("passkey"))),
				),
			),
		})
	}

	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					Div(Class("card shadow"),
						Div(Class("card-body p-4"),
							H2(Class("card-title text-center mb-4"), g.Text(rc.L("login"))),
							FormEl(g.Attr("action", rc.Paths.OIDCLoginPassword), Method("post"),
								CsrfInput(rc.CSRF),
								Div(Class("mb-3"),
									Label(g.Attr("for", "username"), Class("form-label"), g.Text("Email address")),
									Input(Type("email"), Class("form-control"), ID("username"), Name("username"),
										Value(data.Email), g.Attr("required"), g.Attr("readonly")),
								),
								Div(Class("mb-3"),
									Label(g.Attr("for", "password"), Class("form-label"), g.Text("Password")),
									Input(Type("password"), Class("form-control"), ID("password"), Name("password"), g.Attr("required")),
								),
								Button(Type("submit"), Class("btn btn-primary btn-block"), g.Text(rc.L("next"))),
							),
							P(Class("mt-0 text-center"),
								FormEl(g.Attr("action", rc.Paths.ForgotPassword), Method("post"),
									CsrfInput(rc.CSRF),
									Input(Type("hidden"), Name("type"), Value("GET")),
									Button(Type("submit"), Class("btn btn-link text-muted"), g.Text(rc.L("forgot_password"))),
								),
								FormEl(g.Attr("action", rc.Paths.Signup), Method("post"),
									CsrfInput(rc.CSRF),
									Input(Type("hidden"), Name("type"), Value("GET")),
									Button(Type("submit"), Class("btn btn-link text-muted"), g.Text(rc.L("signup"))),
								),
							),
							passkeySection,
							Hr(),
							P(Class("text-center"), g.Text(rc.L("or_signin_with"))),
							Div(Class("d-flex justify-content-center"),
								g.Group(idpButtons),
							),
						),
					),
				),
			),
		),
	)
}

// --- OIDC Login Passkey ---

// OIDCLoginPasskeyData holds data for the passkey login page.
type OIDCLoginPasskeyData struct {
	ReturnFailedUrl string
}

// OIDCLoginPasskeyPage renders the passkey login page (JS-triggered WebAuthn).
func OIDCLoginPasskeyPage(rc *RenderContext, data OIDCLoginPasskeyData) g.Node {
	return PageShell(nil,
		Script(g.Raw(fmt.Sprintf("window.onload = function() { LoginUser(%s); };", data.ReturnFailedUrl))),
		Script(Src("/static/js/webauthn.js")),
	)
}

// --- OIDC Login TOTP ---

// OIDCLoginTOTPData holds data for the TOTP login page.
type OIDCLoginTOTPData struct {
	Errors    []string
	Email     string
	PngQRCode string
	Verified  bool
}

// OIDCLoginTOTPPage renders the TOTP authentication form.
func OIDCLoginTOTPPage(rc *RenderContext, data OIDCLoginTOTPData) g.Node {
	var qrSection g.Node
	if !data.Verified {
		qrSection = Div(Class("mb-3"),
			Img(Src("data:image/png;base64,"+data.PngQRCode),
				Alt("QR Code"),
				g.Attr("style", "max-width: 100%; max-height: 100%;")),
		)
	}

	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					Div(Class("card shadow"),
						Div(Class("card-body p-4"),
							H2(Class("card-title text-center mb-4"), g.Text(rc.L("totp_authenticator_app_login"))),
							FormEl(g.Attr("action", rc.Paths.OIDCLoginTOTP), Method("post"),
								CsrfInput(rc.CSRF),
								qrSection,
								Div(Class("mb-3"),
									Label(g.Attr("for", "username"), Class("form-label"), g.Text("Email address")),
									Input(Type("email"), Class("form-control"), ID("username"), Name("username"),
										Value(data.Email), g.Attr("required"), g.Attr("readonly")),
								),
								Div(Class("mb-3"),
									Label(g.Attr("for", "code"), Class("form-label"), g.Text("Code")),
									Input(Type("text"), Class("form-control"), ID("code"), Name("code"), g.Attr("required")),
								),
								Div(Class("d-flex justify-content-between"),
									Div(Class("btn-group"),
										Button(Type("submit"), Class("btn btn-primary"), Name("action"), Value("next"), g.Text(rc.L("next"))),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

// --- Signup ---

// SignupData holds data for the signup page.
type SignupData struct {
	Errors    []string
	Email     string
	Directive string
	IDPs      []*proto_oidc_models.IDP
}

// SignupPage renders the signup form.
func SignupPage(rc *RenderContext, data SignupData) g.Node {
	idpButtons := make([]g.Node, 0, len(data.IDPs))
	for _, idp := range data.IDPs {
		idpButtons = append(idpButtons,
			FormEl(g.Attr("action", rc.Paths.ExternalIDP), Method("post"),
				Input(Type("hidden"), Name("directive"), Value(data.Directive)),
				Input(Type("hidden"), Name("idp_hint"), Value(idp.GetSlug())),
				Button(Type("submit"), Class("btn btn-outline-primary me-2"), g.Text(idp.GetName())),
			),
		)
	}

	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					Div(Class("card shadow"),
						Div(Class("card-body p-4"),
							H2(Class("card-title text-center mb-4"), g.Text(rc.L("signup"))),
							FormEl(g.Attr("action", rc.Paths.Signup), Method("post"),
								CsrfInput(rc.CSRF),
								Div(Class("mb-3"),
									Label(g.Attr("for", "username"), Class("form-label"), g.Text("Email address")),
									Input(Type("email"), Class("form-control"), ID("username"), Name("username"),
										Value(data.Email), g.Attr("required")),
								),
								Div(Class("mb-3"),
									Label(g.Attr("for", "password"), Class("form-label"), g.Text("Password")),
									Input(Type("password"), Class("form-control"), ID("password"), Name("password"), g.Attr("required")),
								),
								Div(Class("d-flex justify-content-between"),
									Button(Type("submit"), Class("btn btn-outline-primary"), Name("action"), Value("cancel"), g.Attr("formnovalidate"), g.Text(rc.L("cancel"))),
									Div(Class("btn-group"),
										Button(Type("submit"), Class("btn btn-primary"), Name("action"), Value("next"), g.Text(rc.L("next"))),
									),
								),
							),
							Hr(),
							P(Class("text-center"), g.Text(rc.L("or_signin_with"))),
							Div(Class("d-flex justify-content-center"),
								g.Group(idpButtons),
							),
						),
					),
				),
			),
		),
	)
}

// --- Verify Code ---

// VerifyCodeData holds data for the email verification code page.
type VerifyCodeData struct {
	Errors    []string
	Email     string
	Directive string
	Code      string
}

// VerifyCodePage renders the email verification code form.
func VerifyCodePage(rc *RenderContext, data VerifyCodeData) g.Node {
	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					Div(Class("card shadow"),
						Div(Class("card-body p-4"),
							H2(Class("card-title text-center mb-4"), g.Text(rc.L("verifycode"))),
							P(g.Textf("A verification code has be emailed to %s If an account exists. ", data.Email)),
							FormEl(g.Attr("action", rc.Paths.VerifyCode), Method("post"),
								CsrfInput(rc.CSRF),
								Input(Type("hidden"), Name("email"), Value(data.Email)),
								Input(Type("hidden"), Name("directive"), Value(data.Directive)),
								Div(Class("mb-3"),
									Label(g.Attr("for", "code"), Class("form-label"), g.Text("Code")),
									Input(Type("text"), Class("form-control"), ID("code"), Name("code"),
										Value(data.Code), g.Attr("required")),
								),
								Div(Class("d-flex justify-content-between"),
									Button(Type("submit"), Class("btn btn-outline-primary"), Name("action"), Value("cancel"), g.Attr("formnovalidate"), g.Text(rc.L("cancel"))),
									Div(Class("btn-group"),
										Button(Type("submit"), Class("btn btn-primary"), Name("action"), Value("next"), g.Text(rc.L("next"))),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

// --- Forgot Password ---

// ForgotPasswordData holds data for the forgot password page.
type ForgotPasswordData struct {
	Errors []string
	Email  string
}

// ForgotPasswordPage renders the forgot password form.
func ForgotPasswordPage(rc *RenderContext, data ForgotPasswordData) g.Node {
	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					Div(Class("card shadow"),
						Div(Class("card-body p-4"),
							H2(Class("card-title text-center mb-4"), g.Text(rc.L("forgot_password"))),
							FormEl(g.Attr("action", rc.Paths.ForgotPassword), Method("post"),
								CsrfInput(rc.CSRF),
								Div(Class("mb-3"),
									Label(g.Attr("for", "email"), Class("form-label"), g.Text("Email address")),
									Input(Type("email"), Class("form-control"), ID("email"), Name("email"),
										g.Attr("placeholder", "Enter your email"),
										Value(data.Email), g.Attr("required")),
								),
								Div(Class("d-flex justify-content-between"),
									Button(Type("submit"), Class("btn btn-outline-primary"), Name("action"), Value("cancel"), g.Attr("formnovalidate"), g.Text(rc.L("cancel"))),
									Div(Class("btn-group"),
										Button(Type("submit"), Class("btn btn-primary"), Name("action"), Value("next"), g.Text(rc.L("next"))),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

// --- Password Reset ---

// PasswordResetData holds data for the password reset page.
type PasswordResetData struct {
	Errors    []string
	ReturnUrl string
}

// PasswordResetPage renders the password reset page (errors + panel).
func PasswordResetPage(rc *RenderContext, data PasswordResetData) g.Node {
	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		Div(Class("container"),
			Div(Class("row justify-content-center"),
				Div(Class("col-md-6"),
					ErrorList(data.Errors),
					PasswordResetPanel(rc, PasswordResetPanelData{ReturnUrl: data.ReturnUrl}),
				),
			),
		),
	)
}

// --- Error Page ---

// ErrorPageData holds data for the error page.
type ErrorPageData struct {
	Message string
	Error   string
}

// ErrorPage renders the OIDC error page.
func ErrorPage(rc *RenderContext, data ErrorPageData) g.Node {
	messageNode := Div(Strong(g.Text("An error occurred")))
	if data.Message != "" {
		messageNode = Div(Strong(g.Text(data.Message)))
	}
	var errorNode g.Node
	if data.Error != "" {
		errorNode = Div(Class("mt-2"), Small(g.Textf("Error code: %s", data.Error)))
	}

	return PageShellWithNavbar(rc,
		Div(Class("container"),
			Div(Class("text-center mt-5"),
				H1(g.Text("Error")),
				Div(Class("alert alert-danger"), g.Attr("role", "alert"),
					messageNode,
					errorNode,
				),
				Div(Class("mt-4"),
					A(Href("/"), Class("btn btn-primary"), g.Text("Return to Home")),
				),
			),
		),
	)
}
