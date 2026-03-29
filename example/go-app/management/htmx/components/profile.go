package components

import (
	api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ProfilePageData holds data for the profile page.
type ProfilePageData struct {
	*RenderContext
	Profile *api_profile.Profile
	Editing string // "personal" or "contact" or "" for view mode
	Error   string
	Success string
}

// ProfilePage renders the profile page with personal info and contact cards.
func ProfilePage(d *ProfilePageData) g.Node {
	children := []g.Node{
		Div(Class("profile-header"),
			H1(g.Text(d.L("mgmt_my_profile"))),
			P(Class("profile-subtitle"), g.Text(d.L("mgmt_manage_personal_info"))),
		),
	}
	if d.Error != "" {
		children = append(children, ErrorBanner(d.Error))
	}
	if d.Success != "" {
		children = append(children, SuccessBanner(d.Success))
	}
	children = append(children,
		Div(Class("profile-cards"),
			personalInfoCard(d),
			contactInfoCard(d),
			passwordCard(d),
		),
	)
	return Div(Class("profile-container"), g.Group(children))
}

func personalInfoCard(d *ProfilePageData) g.Node {
	if d.Editing == "personal" {
		return personalInfoEditCard(d)
	}

	return ProfileCard(
		CardHeader(PersonIconSVG, d.L("mgmt_personal_information"),
			d.L("mgmt_your_name_and_info"), "personal-info-icon"),
		Div(Class("card-body"),
			Div(Class("info-rows"),
				InfoRow(d.L("mgmt_first_name"), d.Profile.GivenName),
				InfoRow(d.L("mgmt_last_name"), d.Profile.FamilyName),
			),
			Div(Class("card-actions"),
				Button(Type("button"), Class("btn-edit"),
					g.Attr("hx-get", d.Paths.HTMXManagementProfile+"?edit=personal"),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text(d.L("mgmt_edit")),
				),
			),
		),
	)
}

func personalInfoEditCard(d *ProfilePageData) g.Node {
	return ProfileCard(
		CardHeader(PersonIconSVG, d.L("mgmt_personal_information"),
			d.L("mgmt_your_name_and_info"), "personal-info-icon"),
		Div(Class("card-body"),
			HtmxForm(d.Paths.HTMXManagementProfile, "save-personal-indicator",
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("save-personal")),
				FormGroup(d.L("mgmt_first_name"), "text", "givenName", "givenName", d.Profile.GivenName),
				FormGroup(d.L("mgmt_last_name"), "text", "familyName", "familyName", d.Profile.FamilyName),
				ButtonGroup(
					PrimaryButton(d.L("mgmt_save"), "save-personal-indicator"),
					SecondaryButton(d.L("mgmt_cancel"), d.Paths.HTMXManagementProfile),
				),
			),
		),
	)
}

func contactInfoCard(d *ProfilePageData) g.Node {
	if d.Editing == "contact" {
		return contactInfoEditCard(d)
	}

	return ProfileCard(
		CardHeader(
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>`,
			d.L("mgmt_contact_information"),
			d.L("mgmt_how_we_reach_you"), "contact-info-icon"),
		Div(Class("card-body"),
			Div(Class("info-rows"),
				InfoRow(d.L("mgmt_email_address"), d.Profile.Email),
				InfoRow(d.L("mgmt_phone_number"), phoneOrEmpty(d.Profile.PhoneNumber)),
			),
			Div(Class("card-actions"),
				Button(Type("button"), Class("btn-edit"),
					g.Attr("hx-get", d.Paths.HTMXManagementProfile+"?edit=contact"),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text(d.L("mgmt_edit")),
				),
			),
		),
	)
}

func contactInfoEditCard(d *ProfilePageData) g.Node {
	return ProfileCard(
		CardHeader(
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>`,
			d.L("mgmt_contact_information"),
			d.L("mgmt_how_we_reach_you"), "contact-info-icon"),
		Div(Class("card-body"),
			HtmxForm(d.Paths.HTMXManagementProfile, "save-contact-indicator",
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("save-contact")),
				FormGroup(d.L("mgmt_phone_number"), "tel", "phoneNumber", "phoneNumber", d.Profile.PhoneNumber),
				ButtonGroup(
					PrimaryButton(d.L("mgmt_save"), "save-contact-indicator"),
					SecondaryButton(d.L("mgmt_cancel"), d.Paths.HTMXManagementProfile),
				),
			),
		),
	)
}

func passwordCard(d *ProfilePageData) g.Node {
	return ProfileCard(
		CardHeader(LockIconSVG, d.L("mgmt_password_manager"),
			d.L("mgmt_reset_password_desc"), "password-icon"),
		Div(Class("card-body"),
			Div(Class("card-actions"),
				Button(Type("button"), Class("btn-secondary"),
					g.Attr("hx-get", d.Paths.HTMXManagementPassword),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-push-url", "true"),
					g.Text(d.L("mgmt_reset_password")),
				),
			),
		),
	)
}

func phoneOrEmpty(phone string) string {
	if phone == "" {
		return "—"
	}
	return phone
}
