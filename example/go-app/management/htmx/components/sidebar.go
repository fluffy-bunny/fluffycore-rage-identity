package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type navItem struct {
	Icon   string
	Label  string
	Path   string
	PageID string
}

// Sidebar renders the navigation sidebar with active state highlighting.
func Sidebar(rc *RenderContext) g.Node {
	items := []navItem{
		{HomeIconSVG, rc.L("mgmt_home"), rc.Paths.HTMXManagementHome, "home"},
		{PersonIconSVG, rc.L("mgmt_profile"), rc.Paths.HTMXManagementProfile, "profile"},
		{LockIconSVG, rc.L("mgmt_password_manager"), rc.Paths.HTMXManagementPassword, "password"},
	}
	if rc.AppConfig != nil && rc.AppConfig.EnabledWebAuthN {
		items = append(items, navItem{PasskeyIconSVG, rc.L("mgmt_passkeys"), rc.Paths.HTMXManagementPasskey, "passkey"})
	}
	items = append(items,
		navItem{LinkIconSVG, rc.L("mgmt_linked_accounts"), rc.Paths.HTMXManagementLinked, "linked"},
		navItem{SettingsIconSVG, rc.L("mgmt_preferences"), rc.Paths.HTMXManagementPrefs, "prefs"},
	)

	navLinks := make([]g.Node, len(items))
	for i, item := range items {
		activeClass := "sidebar-link"
		if item.PageID == rc.ActivePage {
			activeClass = "sidebar-link active"
		}
		navLinks[i] = A(Class(activeClass),
			g.Attr("hx-get", item.Path),
			g.Attr("hx-target", "#dashboard-main"),
			g.Attr("hx-swap", "innerHTML"),
			g.Attr("hx-push-url", "true"),
			g.Raw(item.Icon),
			Span(g.Text(item.Label)),
		)
	}

	return g.Group([]g.Node{
		// Hidden checkbox for mobile sidebar toggle (pure CSS approach)
		Input(Type("checkbox"), ID("sidebar-toggle"), Class("sidebar-toggle-input"),
			g.Attr("style", "display:none")),
		Nav(Class("dashboard-sidebar"),
			Div(Class("sidebar-content"),
				Div(Class("sidebar-section"),
					Div(Class("sidebar-section-title"), g.Text(rc.L("mgmt_account"))),
					g.Group(navLinks),
				),
			),
		),
	})
}
