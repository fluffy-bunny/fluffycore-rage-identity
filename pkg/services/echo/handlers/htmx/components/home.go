package components

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomeData holds data for the home (email entry) partial.
type HomeData struct {
	*RenderContext
	Errors          []string
	Email           string
	SocialIdps      []*proto_oidc_models.IDP
	DisableSignup   bool
	EnabledWebAuthN bool
}

// HomePartial renders the email entry step matching the WASM Home layout.
func HomePartial(data HomeData) g.Node {
	return g.Group([]g.Node{
		// Error banner
		ErrorMessages(data.Errors),
		// Email form
		HtmxForm(data.Paths.HTMXHome, "home-indicator",
			CsrfInput(data.CSRF),
			Div(Class("form-group"),
				Label(g.Attr("for", "email"), g.Text(data.L("email_address"))),
				Input(
					Type("email"),
					ID("email"),
					Name("email"),
					Value(data.Email),
					g.Attr("placeholder", data.L("email_address")),
					g.Attr("required"),
					g.Attr("autofocus"),
				),
			),
			ButtonGroup(
				PrimaryButton(data.L("continue"), "home-indicator"),
			),
		),
		// Passkey login section (conditional)
		g.If(data.EnabledWebAuthN,
			PasskeyLoginSection(data.CSRF, data.Paths.HTMXKeepSignedIn, data.L("signin_with_passkey")),
		),
		// Create account + forgot password links
		g.If(!data.DisableSignup,
			Div(
				Div(Class("create-account"),
					Span(g.Text(data.L("dont_have_account")+" ")),
					A(Href("#"),
						g.Attr("hx-get", data.Paths.HTMXSignup),
						g.Attr("hx-target", "#main-content"),
						g.Attr("hx-swap", "innerHTML"),
						g.Attr("hx-push-url", "true"),
						g.Text(data.L("create_one")),
					),
				),
				Div(Class("forgot-password"),
					A(Href("#"),
						g.Attr("hx-get", data.Paths.HTMXForgotPassword),
						g.Attr("hx-target", "#main-content"),
						g.Attr("hx-swap", "innerHTML"),
						g.Attr("hx-push-url", "true"),
						g.Text(data.L("forgot_password")),
					),
				),
			),
		),
		// Social login buttons
		g.If(!data.DisableSignup,
			SocialIdpButtons(data.SocialIdps, data.CSRF, data.Paths.HTMXHome, data.L("or_signin_with")),
		),
	})
}
