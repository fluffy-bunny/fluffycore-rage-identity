package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PreferencesPageData holds data for the preferences page.
type PreferencesPageData struct {
	*RenderContext
	KeepSignedIn  bool
	DontShowAgain bool
	Error         string
	Success       string
}

// PreferencesPage renders the preferences page.
func PreferencesPage(d *PreferencesPageData) g.Node {
	children := []g.Node{
		Div(Class("profile-header"),
			H1(g.Text(d.L("mgmt_preferences"))),
			P(Class("profile-subtitle"), g.Text(d.L("mgmt_preferences"))),
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
			Div(g.Attr("style", "display:flex;flex-direction:column;gap:1rem"),
				ProfileCard(
					CardHeader(SettingsIconSVG, d.L("mgmt_preferences"), "", "personal-info-icon"),
					Div(Class("card-body"),
						HtmxForm(d.Paths.HTMXManagementPrefs, "prefs-indicator",
							CsrfInput(d.CSRF),
							Input(Type("hidden"), Name("action"), Value("save-prefs")),
							ToggleSwitch("keepSignedIn", "keepSignedIn", d.KeepSignedIn,
								d.L("mgmt_keep_signed_in")),
							ToggleSwitch("dontShowAgain", "dontShowAgain", d.DontShowAgain,
								d.L("mgmt_dont_show_again")),
							ButtonGroup(
								PrimaryButton(d.L("mgmt_save"), "prefs-indicator"),
							),
						),
					),
				),
				// Clear SSO card
				ProfileCard(
					CardHeader(SettingsIconSVG, d.L("mgmt_clear_sso"), "", "password-icon"),
					Div(Class("card-body"),
						HtmxForm(d.Paths.HTMXManagementPrefs, "clear-sso-indicator",
							CsrfInput(d.CSRF),
							Input(Type("hidden"), Name("action"), Value("clear-sso")),
							ButtonGroup(
								PrimaryButton(d.L("mgmt_clear_sso"), "clear-sso-indicator"),
							),
						),
					),
				),
			),
		),
	)

	return Div(Class("profile-container"), g.Group(children))
}
