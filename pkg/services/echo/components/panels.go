package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PersonalInformationPanelData holds data for the personal information form panel.
type PersonalInformationPanelData struct {
	FormAction  string
	Action      string
	ReturnUrl   string
	Email       string
	GivenName   string
	FamilyName  string
	PhoneNumber string
	DisplayOnly bool
}

// PersonalInformationPanel renders a card with form fields for personal info.
func PersonalInformationPanel(rc *RenderContext, data PersonalInformationPanelData) g.Node {
	submitLabel := rc.L("next")
	if data.DisplayOnly {
		submitLabel = rc.L("edit")
	}
	readonlyAttr := func() g.Node {
		if data.DisplayOnly {
			return g.Attr("readonly")
		}
		return nil
	}

	return Div(Class("card shadow"),
		Div(Class("card-body p-4"),
			H2(Class("card-title text-center mb-4"), g.Text(rc.L("personal_information"))),
			FormEl(g.Attr("action", data.FormAction), Method("post"),
				CsrfInput(rc.CSRF),
				Input(Type("hidden"), Name("action"), Value(data.Action)),
				Input(Type("hidden"), Name("returnUrl"), Value(data.ReturnUrl)),
				Div(Class("mb-3"),
					Label(g.Attr("for", "email"), Class("form-label"), g.Text(rc.L("email"))),
					Input(Type("email"), Class("form-control"), ID("email"), Name("email"),
						Value(data.Email), g.Attr("required"), g.Attr("readonly")),
				),
				Div(Class("mb-3"),
					Label(g.Attr("for", "given_name"), Class("form-label"), g.Text(rc.L("given_name"))),
					Input(Type("text"), Class("form-control"), ID("given_name"), Name("given_name"),
						Value(data.GivenName), readonlyAttr()),
				),
				Div(Class("mb-3"),
					Label(g.Attr("for", "family_name"), Class("form-label"), g.Text(rc.L("family_name"))),
					Input(Type("text"), Class("form-control"), ID("family_name"), Name("family_name"),
						Value(data.FamilyName), readonlyAttr()),
				),
				Div(Class("mb-3"),
					Label(g.Attr("for", "phone_number"), Class("form-label"), g.Text(rc.L("phone_number"))),
					Input(Type("text"), Class("form-control"), ID("phone_number"), Name("phone_number"),
						Value(data.PhoneNumber), readonlyAttr()),
				),
				Button(Type("submit"), Class("btn btn-primary btn-block"), g.Text(submitLabel)),
			),
		),
	)
}

// PasswordResetPanelData holds data for the password reset form panel.
type PasswordResetPanelData struct {
	ReturnUrl string
}

// PasswordResetPanel renders a card with password/confirm password form.
func PasswordResetPanel(rc *RenderContext, data PasswordResetPanelData) g.Node {
	return Div(Class("card shadow"),
		Div(Class("card-body p-4"),
			H2(Class("card-title text-center mb-4"), g.Text(rc.L("password_reset"))),
			FormEl(g.Attr("action", rc.Paths.PasswordReset), Method("post"),
				CsrfInput(rc.CSRF),
				Input(Type("hidden"), Name("returnUrl"), Value(data.ReturnUrl)),
				Div(Class("mb-3"),
					Label(g.Attr("for", "password"), Class("form-label"), g.Text(rc.L("password"))),
					Input(Type("password"), Class("form-control"), ID("password"), Name("password"), g.Attr("required")),
				),
				Div(Class("mb-3"),
					Label(g.Attr("for", "confirmPassword"), Class("form-label"), g.Text(rc.L("confirm_password"))),
					Input(Type("password"), Class("form-control"), ID("confirmPassword"), Name("confirmPassword"), g.Attr("required")),
				),
				Div(Class("d-flex justify-content-between"),
					Button(Type("submit"), Class("btn btn-outline-primary"), Name("action"), Value("cancel"), g.Attr("formnovalidate"), g.Text(rc.L("cancel"))),
					Div(Class("btn-group"),
						Button(Type("submit"), Class("btn btn-primary"), Name("action"), Value("next"), g.Text(rc.L("next"))),
					),
				),
			),
		),
	)
}
