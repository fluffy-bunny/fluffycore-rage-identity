package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PasswordStage tracks the current stage of the password reset flow.
type PasswordStage string

const (
	PasswordStageInitial PasswordStage = "initial"
	PasswordStageVerify  PasswordStage = "verify"
	PasswordStageReset   PasswordStage = "reset"
	PasswordStageSuccess PasswordStage = "success"
)

// PasswordPageData holds data for the password manager page.
type PasswordPageData struct {
	*RenderContext
	Stage           PasswordStage
	Email           string
	IsClaimedDomain bool
	Error           string
	DevCode         string // verification code shown in dev mode
}

// PasswordPage renders the password manager page with the appropriate stage.
func PasswordPage(d *PasswordPageData) g.Node {
	children := []g.Node{
		Div(Class("profile-header"),
			H1(g.Text(d.L("mgmt_password_manager"))),
			P(Class("profile-subtitle"), g.Text(d.L("mgmt_reset_password_desc"))),
		),
	}

	if d.IsClaimedDomain {
		children = append(children,
			ProfileCard(
				CardHeader(LockIconSVG, d.L("mgmt_password_not_available"), "", "password-icon"),
				Div(Class("card-body"),
					P(g.Text(d.L("mgmt_password_not_available_desc"))),
					Div(Class("card-actions"),
						SecondaryButton(d.L("mgmt_back"), d.Paths.HTMXManagementProfile),
					),
				),
			),
		)
		return Div(Class("profile-container"), g.Group(children))
	}

	if d.Error != "" {
		children = append(children, ErrorBanner(d.Error))
	}

	switch d.Stage {
	case PasswordStageVerify:
		children = append(children, passwordVerifyStage(d))
	case PasswordStageReset:
		children = append(children, passwordResetStage(d))
	case PasswordStageSuccess:
		children = append(children, passwordSuccessStage(d))
	default:
		children = append(children, passwordInitialStage(d))
	}

	return Div(Class("profile-container"), g.Group(children))
}

func passwordInitialStage(d *PasswordPageData) g.Node {
	return ProfileCard(
		CardHeader(LockIconSVG, d.L("mgmt_reset_password"),
			d.L("mgmt_reset_password_desc"), "password-icon"),
		Div(Class("card-body"),
			Div(Class("info-rows"),
				InfoRow(d.L("mgmt_email_address"), d.Email),
			),
			HtmxForm(d.Paths.HTMXManagementPassword, "send-code-indicator",
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("send-code")),
				ButtonGroup(
					PrimaryButton(d.L("mgmt_send_verification_code"), "send-code-indicator"),
				),
			),
		),
	)
}

func passwordVerifyStage(d *PasswordPageData) g.Node {
	bodyChildren := []g.Node{
		HtmxForm(d.Paths.HTMXManagementPassword, "verify-code-indicator",
			CsrfInput(d.CSRF),
			Input(Type("hidden"), Name("action"), Value("verify-code")),
			FormGroup(d.L("mgmt_enter_verification_code"), "text", "code", "code", "",
				g.Attr("placeholder", "000000"),
				g.Attr("maxlength", "6"),
				g.Attr("required"),
				g.Attr("autocomplete", "one-time-code"),
			),
			ButtonGroup(
				SecondaryButton(d.L("mgmt_cancel"), d.Paths.HTMXManagementPassword),
				PrimaryButton(d.L("mgmt_verify_code"), "verify-code-indicator"),
			),
		),
		Div(Class("resend-link"),
			FormEl(
				g.Attr("hx-post", d.Paths.HTMXManagementPassword),
				g.Attr("hx-target", "#dashboard-main"),
				g.Attr("hx-swap", "innerHTML"),
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("send-code")),
				Button(Type("submit"), Class("btn-link"),
					g.Text(d.L("mgmt_didnt_receive_code")),
				),
			),
		),
	}

	// Show the verification code in dev mode
	if d.DevCode != "" {
		bodyChildren = append([]g.Node{
			Div(Class("dev-code-banner"),
				g.Attr("style", "background:#fff3cd;border:1px solid #ffc107;border-radius:8px;padding:12px 16px;margin-bottom:16px;font-family:monospace;"),
				Span(g.Attr("style", "font-weight:600;color:#856404;"), g.Text("DEV CODE: ")),
				Span(g.Attr("style", "font-size:18px;font-weight:700;color:#856404;letter-spacing:2px;"), g.Text(d.DevCode)),
			),
		}, bodyChildren...)
	}

	return ProfileCard(
		CardHeader(LockIconSVG, d.L("mgmt_verify_code"),
			d.LF("mgmt_enter_code_sent_to", map[string]string{"email": d.Email}), "password-icon"),
		Div(Class("card-body"), g.Group(bodyChildren)),
	)
}

func passwordResetStage(d *PasswordPageData) g.Node {
	return ProfileCard(
		CardHeader(LockIconSVG, d.L("mgmt_reset_password"),
			d.L("mgmt_enter_new_password_below"), "password-icon"),
		Div(Class("card-body"),
			HtmxForm(d.Paths.HTMXManagementPassword, "reset-password-indicator",
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("reset-password")),
				FormGroup(d.L("mgmt_new_password"), "password", "newPassword", "newPassword", "",
					g.Attr("placeholder", d.L("mgmt_enter_new_password")),
					g.Attr("required"),
					g.Attr("autocomplete", "new-password"),
				),
				FormGroup(d.L("mgmt_confirm_password"), "password", "confirmPassword", "confirmPassword", "",
					g.Attr("placeholder", d.L("mgmt_confirm_new_password")),
					g.Attr("required"),
					g.Attr("autocomplete", "new-password"),
				),
				ButtonGroup(
					PrimaryButton(d.L("mgmt_reset_password"), "reset-password-indicator"),
					SecondaryButton(d.L("mgmt_cancel"), d.Paths.HTMXManagementProfile),
				),
			),
		),
	)
}

func passwordSuccessStage(d *PasswordPageData) g.Node {
	return ProfileCard(
		CardHeader(LockIconSVG, d.L("mgmt_success"),
			d.L("mgmt_password_reset_success"), "success-icon"),
		Div(Class("card-body"),
			P(g.Text(d.L("mgmt_password_reset_success"))),
			ButtonGroup(
				SecondaryButton(d.L("mgmt_done"), d.Paths.HTMXManagementProfile),
			),
		),
	)
}
