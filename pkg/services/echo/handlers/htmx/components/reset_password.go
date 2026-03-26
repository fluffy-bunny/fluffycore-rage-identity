package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ResetPasswordData holds data for the reset password partial.
type ResetPasswordData struct {
	*RenderContext
	Errors []string
	Email  string
}

// ResetPasswordPartial renders the reset password step.
func ResetPasswordPartial(data ResetPasswordData) g.Node {
	return g.Group([]g.Node{
		ErrorMessages(data.Errors),
		H2(g.Text(data.L("reset_password"))),
		P(g.Text(data.L("enter_new_password"))),
		HtmxForm(data.Paths.HTMXResetPassword, "reset-indicator",
			CsrfInput(data.CSRF),
			FormGroupField(data.L("new_password"), "password", "password", "password", "",
				g.Attr("required"), g.Attr("autofocus")),
			FormGroupField(data.L("confirm_password"), "password", "confirmPassword", "confirmPassword", ""),
			ButtonGroup(
				SecondaryButton(data.L("cancel"), data.Paths.HTMXHome),
				PrimaryButton(data.L("reset_password"), "reset-indicator"),
			),
		),
	})
}
