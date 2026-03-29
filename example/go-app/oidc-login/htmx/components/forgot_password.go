package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ForgotPasswordData holds data for the forgot password partial.
type ForgotPasswordData struct {
	*RenderContext
	Errors []string
	Email  string
}

// ForgotPasswordPartial renders the forgot password step.
func ForgotPasswordPartial(data ForgotPasswordData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("forgot_password"))),
		P(g.Text(data.L("enter_email_for_reset"))),
		HtmxForm(data.Paths.HTMXForgotPassword, "forgot-indicator",
			CsrfInput(data.CSRF),
			FormGroupField(data.L("email"), "email", "email", "email", data.Email,
				g.Attr("required"), g.Attr("autofocus")),
			ButtonGroup(
				SecondaryButton(data.L("cancel"), data.Paths.HTMXHome),
				PrimaryButton(data.L("next"), "forgot-indicator"),
			),
		),
	})
}
