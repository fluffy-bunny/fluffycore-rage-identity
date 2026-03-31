package components

import (
	"fmt"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// --- About Page ---

// AboutRouteRow holds data for a single route row in the about page.
type AboutRouteRow struct {
	Verbs string
	Path  string
}

// AboutPage renders the about/routes page.
func AboutPage(rc *RenderContext, rows []AboutRouteRow) g.Node {
	tableRows := make([]g.Node, 0, len(rows))
	for idx, r := range rows {
		tableRows = append(tableRows,
			Tr(
				Th(Class("text-start"), g.Attr("scope", "row"), g.Text(fmt.Sprintf("%d", idx))),
				Td(Class("text-start"), g.Text(r.Verbs)),
				Td(Class("text-start"), g.Text(r.Path)),
			),
		)
	}

	return PageShellWithNavbar(rc,
		Div(Class("container"),
			Div(Class("text-center mt-5"), g.Attr("role", "alert"),
				H1(g.Text("Perfect Corp.")),
				P(Class("lead"), g.Text("Everything good, nothing bad")),
				Div(Class("mt-5 alert alert-success"), g.Attr("role", "alert"),
					Table(Class("table table-striped"),
						THead(
							Tr(
								Th(Class("text-start"), g.Attr("scope", "col"), g.Text("#")),
								Th(Class("text-start"), g.Attr("scope", "col"), g.Text("Verbs")),
								Th(Class("text-start"), g.Attr("scope", "col"), g.Text("Path")),
							),
						),
						TBody(g.Group(tableRows)),
					),
				),
			),
		),
	)
}

// --- Profile Page (React SPA shell) ---

// ProfilePage renders the SPA shell for the React-based account management.
func ProfilePage(rc *RenderContext) g.Node {
	return PageShell(
		[]g.Node{Class("bg-light d-flex align-items-center min-vh-100")},
		g.El("noscript", g.Text("You need to enable JavaScript to run this app.")),
		Div(ID("root")),
		Script(g.Attr("defer", "defer"), Src("/static/account-management/static/js/main.js")),
	)
}

// --- TOTP Management ---

// TOTPManagementData holds data for the TOTP management page.
type TOTPManagementData struct {
	ReturnUrl  string
	FormAction string
	PngQRCode  string
	Verified   bool
	Enabled    bool
}

// TOTPManagementPage renders the TOTP management page.
func TOTPManagementPage(rc *RenderContext, data TOTPManagementData) g.Node {
	var formContent g.Node
	if data.Verified {
		actionValue := "enable"
		buttonLabel := rc.L("totp_enable")
		if data.Enabled {
			actionValue = "disable"
			buttonLabel = rc.L("totp_disable")
		}
		formContent = FormEl(g.Attr("action", data.FormAction), Method("post"),
			CsrfInput(rc.CSRF),
			Input(Type("hidden"), Name("returnUrl"), Value(data.ReturnUrl)),
			Input(Type("hidden"), Name("action"), Value(actionValue)),
			Button(Type("submit"), Class("btn btn-primary btn-block"), g.Text(buttonLabel)),
		)
	} else {
		formContent = FormEl(g.Attr("action", data.FormAction), Method("post"),
			CsrfInput(rc.CSRF),
			Input(Type("hidden"), Name("returnUrl"), Value(data.ReturnUrl)),
			Input(Type("hidden"), Name("action"), Value("enroll")),
			Div(Class("mb-3"),
				Label(g.Attr("for", "code"), Class("form-label"), g.Text(rc.L("code"))),
				Input(Type("text"), Class("form-control"), ID("code"), Name("code"),
					g.Attr("placeholder", rc.L("totp_enter_placeholder")), g.Attr("required")),
			),
			Button(Type("submit"), Class("btn btn-primary btn-block"), g.Text(rc.L("totp_enroll"))),
		)
	}

	return PageShellWithNavbar(rc,
		Div(Class("container"),
			Div(Class("text-center mt-5"),
				H1(g.Text(rc.L("totp_management"))),
			),
			Div(Class("mb-3"),
				Img(Src("data:image/png;base64,"+data.PngQRCode),
					Alt("QR Code"),
					g.Attr("style", "max-width: 100%; max-height: 100%;")),
			),
			formContent,
		),
	)
}

// --- Passkey Management ---

// PasskeyManagementData holds data for the passkey management page.
type PasskeyManagementData struct {
	ReturnUrl string
}

// PasskeyManagementPage renders the passkey management page.
func PasskeyManagementPage(rc *RenderContext, data PasskeyManagementData) g.Node {
	return PageShellWithNavbar(rc,
		Div(Class("container"),
			Div(Class("text-center mt-5"),
				H1(g.Text(rc.L("passkey_management"))),
				Button(Class("btn btn-outline-primary"),
					g.Attr("onclick", fmt.Sprintf("registerUser(%s)", data.ReturnUrl)),
					g.Text(rc.L("register")),
				),
			),
		),
		Script(Src("/static/js/webauthn.js")),
	)
}

// --- Personal Information Page ---

// PersonalInformationPage renders the personal information page with navbar.
func PersonalInformationPage(rc *RenderContext, panelData PersonalInformationPanelData) g.Node {
	return PageShellWithNavbar(rc,
		Div(Class("container"),
			Div(Class("text-center mt-5"),
				H1(g.Text(rc.L("personal_information"))),
				Div(Class("row justify-content-center"),
					Div(Class("col-md-6"),
						PersonalInformationPanel(rc, panelData),
					),
				),
			),
		),
	)
}
