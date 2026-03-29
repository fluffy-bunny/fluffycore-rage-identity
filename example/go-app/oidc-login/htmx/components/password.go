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

// PasswordPartial renders the password entry step matching the WASM layout.
func PasswordPartial(data PasswordData) g.Node {
	return g.Group([]g.Node{
		H2(g.Text(data.L("enter_your_password"))),
		P(g.Text(data.LF("enter_your_password_for", map[string]string{"email": data.Email}))),
		HtmxForm(data.Paths.HTMXPassword, "password-indicator",
			CsrfInput(data.CSRF),
			Input(Type("hidden"), Name("email"), Value(data.Email)),
			Div(Class("form-group"),
				Label(g.Attr("for", "password"), g.Text(data.L("password"))),
				Input(
					Type("password"),
					ID("password"),
					Name("password"),
					g.Attr("placeholder", data.L("enter_your_password")),
					g.Attr("required"),
					g.Attr("autofocus"),
				),
				ErrorMessages(data.Errors),
			),
			ButtonGroup(
				SecondaryButton(data.L("back"), data.Paths.HTMXHome),
				PrimaryButton(data.L("signin"), "password-indicator"),
			),
		),
	})
}
