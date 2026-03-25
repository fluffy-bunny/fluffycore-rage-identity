package components

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomeData holds data for the home (email entry) partial.
type HomeData struct {
	*RenderContext
	Errors        []string
	Email         string
	SocialIdps    []*proto_oidc_models.IDP
	DisableSignup bool
}

// HomePartial renders the email entry step.
func HomePartial(data HomeData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("signin"))),
		P(g.Text("Enter your email to get started")),
		HtmxForm(data.Paths.HTMXHome, "home-indicator",
			CsrfInput(data.CSRF),
			FormGroupField(data.L("email"), "email", "email", "email", data.Email,
				g.Attr("required"), g.Attr("autofocus")),
			ButtonGroup(
				PrimaryButton(data.L("next"), "home-indicator"),
			),
		),
		g.If(!data.DisableSignup,
			Div(Class("create-account"),
				A(Href("#"),
					g.Attr("hx-get", data.Paths.HTMXSignup),
					g.Attr("hx-target", "#main-content"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-push-url", "true"),
					g.Text(data.L("signup")),
				),
			),
		),
		SocialIdpButtons(data.SocialIdps, data.CSRF, data.Paths.HTMXHome, data.L("or_signin_with")),
	})
}
