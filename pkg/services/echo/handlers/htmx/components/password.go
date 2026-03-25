package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PasswordData holds data for the password entry partial.
type PasswordData struct {
	*RenderContext
	Errors []string
	Email  string
}

// PasswordPartial renders the password entry step.
func PasswordPartial(data PasswordData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("password"))),
		P(g.Text(data.Email)),
		HtmxForm(data.Paths.HTMXPassword, "password-indicator",
			CsrfInput(data.CSRF),
			Input(Type("hidden"), Name("email"), Value(data.Email)),
			FormGroupField(data.L("password"), "password", "password", "password", "",
				g.Attr("required"), g.Attr("autofocus")),
			ButtonGroup(
				SecondaryButton(data.L("cancel"), data.Paths.HTMXHome),
				PrimaryButton(data.L("next"), "password-indicator"),
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
	})
}
