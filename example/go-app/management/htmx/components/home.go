package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomePage renders the management dashboard home page content.
func HomePage(rc *RenderContext) g.Node {
	features := []g.Node{
		featureCard(ShieldIconLargeSVG, "security-icon",
			rc.L("mgmt_secure_authentication"),
			rc.L("mgmt_secure_authentication_desc"),
			rc.Paths.HTMXManagementPassword),
	}
	if rc.AppConfig != nil && rc.AppConfig.EnabledWebAuthN {
		features = append(features, featureCard(PasskeyLargeIconSVG, "passkey-icon",
			rc.L("mgmt_passkeys"),
			rc.L("mgmt_manage_passkeys"),
			rc.Paths.HTMXManagementPasskey))
	}
	features = append(features,
		featureCard(LinkLargeIconSVG, "accounts-icon",
			rc.L("mgmt_linked_accounts"),
			rc.L("mgmt_linked_accounts_desc"),
			rc.Paths.HTMXManagementLinked),
		featureCard(LockLargeIconSVG, "privacy-icon",
			rc.L("mgmt_privacy_matters"),
			rc.L("mgmt_privacy_matters_desc"),
			rc.Paths.HTMXManagementProfile),
	)

	return Div(Class("home-container"),
		// Hero section
		Div(Class("home-hero"),
			Div(Class("home-hero-content"),
				H1(Class("home-title"), g.Text(rc.L("mgmt_welcome"))),
				P(Class("home-subtitle"), g.Text(rc.L("mgmt_welcome_subtitle"))),
			),
		),
		// Feature cards
		Div(Class("home-features"), g.Group(features)),
		// CTA section
		Div(Class("home-cta"),
			H2(g.Text(rc.L("mgmt_ready_to_start"))),
			P(g.Text(rc.L("mgmt_sign_in_to_access"))),
		),
	)
}

func featureCard(iconSVG, iconClass, title, description, navPath string) g.Node {
	return Div(Class("home-feature-card"),
		g.Attr("hx-get", navPath),
		g.Attr("hx-target", "#dashboard-main"),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-push-url", "true"),
		g.Attr("style", "cursor:pointer"),
		Div(Class("home-feature-icon "+iconClass), g.Raw(iconSVG)),
		H3(Class("home-feature-title"), g.Text(title)),
		P(Class("home-feature-description"), g.Text(description)),
	)
}
