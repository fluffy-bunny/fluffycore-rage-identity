package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// VerifyCodeData holds data for the verify code partial.
type VerifyCodeData struct {
	*RenderContext
	Errors    []string
	Email     string
	Directive string
	Code      string
}

// VerifyCodePartial renders the verification code entry step.
func VerifyCodePartial(data VerifyCodeData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("verifycode"))),
		P(g.Text("A verification code has been emailed to " + data.Email + ".")),
		HtmxForm(data.Paths.HTMXVerifyCode, "verify-indicator",
			CsrfInput(data.CSRF),
			Input(Type("hidden"), Name("email"), Value(data.Email)),
			Input(Type("hidden"), Name("directive"), Value(data.Directive)),
			Div(Class("form-group"),
				Label(g.Attr("for", "code"), g.Text(data.L("code"))),
				Input(Type("text"), Class("verification-input"), ID("code"), Name("code"),
					Value(data.Code), g.Attr("required"), g.Attr("autofocus")),
			),
			ButtonGroup(
				SecondaryButton(data.L("cancel"), data.Paths.HTMXHome),
				PrimaryButton(data.L("next"), "verify-indicator"),
			),
		),
	})
}
