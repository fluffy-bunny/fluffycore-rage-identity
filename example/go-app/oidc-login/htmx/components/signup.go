package components

import (
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// SignupData holds data for the signup partial.
type SignupData struct {
	*RenderContext
	Errors     []string
	Email      string
	SocialIdps []*proto_oidc_models.IDP
}

// SignupPartial renders the signup (create account) step.
func SignupPartial(data SignupData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("signup"))),
		P(g.Text("Create a new account")),
		HtmxForm(data.Paths.HTMXSignup, "signup-indicator",
			CsrfInput(data.CSRF),
			FormGroupField(data.L("email"), "email", "email", "email", data.Email,
				g.Attr("required"), g.Attr("autofocus")),
			FormGroupField(data.L("password"), "password", "password", "password", "",
				g.Attr("required")),
			ButtonGroup(
				SecondaryButton(data.L("cancel"), data.Paths.HTMXHome),
				PrimaryButton(data.L("next"), "signup-indicator"),
			),
		),
		SocialIdpButtons(data.SocialIdps, data.CSRF, data.Paths.HTMXHome, data.L("or_signin_with")),
	})
}
