package components

import (
	"fmt"
	"strings"
	"time"

	api_linked "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_linked_identities"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// LinkedPageData holds data for the linked accounts page.
type LinkedPageData struct {
	*RenderContext
	Identities      []api_linked.LinkedIdentity
	IsClaimedDomain bool
	Error           string
	Success         string
}

// LinkedAccountsPage renders the linked accounts page.
func LinkedAccountsPage(d *LinkedPageData) g.Node {
	children := []g.Node{
		Div(Class("profile-header"),
			H1(g.Text(d.L("mgmt_linked_accounts"))),
			P(Class("profile-subtitle"), g.Text(d.L("mgmt_manage_linked_accounts"))),
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
			g.If(len(d.Identities) == 0,
				linkedEmptyState(d),
			),
			g.If(len(d.Identities) > 0,
				linkedAccountsList(d),
			),
		),
	)

	return Div(Class("profile-container"), g.Group(children))
}

func linkedAccountsList(d *LinkedPageData) g.Node {
	cards := make([]g.Node, len(d.Identities))
	for i, identity := range d.Identities {
		cards[i] = linkedAccountCard(d, identity)
	}
	return Div(g.Attr("style", "display:flex;flex-direction:column;gap:1rem"), g.Group(cards))
}

func linkedAccountCard(d *LinkedPageData, identity api_linked.LinkedIdentity) g.Node {
	icon, providerName := providerBranding(identity.Provider)
	iconClass := providerIconClass(identity.Provider)

	// Build title like WASM: "Provider (January 2, 2006 at 3:04 PM)"
	title := providerName
	if identity.CreatedOn > 0 {
		title += " (" + formatUnixTime(identity.CreatedOn) + ")"
	}

	// Build description like WASM: "Last used on January 2, 2006 at 3:04 PM"
	description := formatLastUsed(identity.LastUsedOn)

	cardBody := []g.Node{}
	if !d.IsClaimedDomain {
		cardBody = append(cardBody,
			Div(Class("button-group"),
				Button(Type("button"), Class("btn-unlink"),
					g.Attr("hx-post", d.Paths.HTMXManagementLinked),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-vals", fmt.Sprintf(`{"action":"unlink","identity":"%s","csrf":"%s"}`, identity.Subject, d.CSRF)),
					g.Attr("hx-confirm", "Are you sure you want to unlink this account?"),
					g.Text(d.L("mgmt_unlink")),
				),
			),
		)
	}

	return ProfileCard(
		Div(Class("card-header"),
			Div(Class("card-header-content"),
				Div(Class("card-icon "+iconClass), g.Raw(icon)),
				Div(Class("card-title-group"),
					H2(g.Text(title)),
					P(Class("card-description"), g.Text(description)),
				),
			),
		),
		Div(Class("card-body"), g.Group(cardBody)),
	)
}

func linkedEmptyState(d *LinkedPageData) g.Node {
	return ProfileCard(
		Div(Class("card-body"),
			Div(g.Attr("style", "text-align:center;padding:40px 20px"),
				Div(Class("home-feature-icon accounts-icon"),
					g.Attr("style", "margin:0 auto 24px"),
					g.Raw(`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
						<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
					</svg>`),
				),
				H3(g.Attr("style", "margin:0 0 8px 0"), g.Text(d.L("mgmt_no_linked_accounts"))),
				P(Class("card-description"), g.Attr("style", "color:var(--text-secondary)"),
					g.Text(d.L("mgmt_no_linked_accounts_desc"))),
			),
		),
	)
}

// formatUnixTime converts unix timestamp to friendly date string with time
func formatUnixTime(unixTime int64) string {
	if unixTime == 0 {
		return ""
	}
	t := time.Unix(unixTime, 0)
	return t.Format("January 2, 2006 at 3:04 PM")
}

// formatLastUsed returns a friendly string for last used timestamp
func formatLastUsed(lastUsedOn int64) string {
	if lastUsedOn == 0 {
		return "Never been used"
	}
	t := time.Unix(lastUsedOn, 0)
	return "Last used on " + t.Format("January 2, 2006 at 3:04 PM")
}

func providerBranding(provider string) (string, string) {
	lower := strings.ToLower(provider)
	switch {
	case strings.Contains(lower, "google"):
		return GoogleIconSVG, "Google"
	case strings.Contains(lower, "microsoft"):
		return MicrosoftIconSVG, "Microsoft"
	case strings.Contains(lower, "github"):
		return GitHubIconSVG, "GitHub"
	default:
		return LinkIconSVG, provider
	}
}

func providerIconClass(provider string) string {
	lower := strings.ToLower(provider)
	switch {
	case strings.Contains(lower, "google"):
		return "icon-google"
	case strings.Contains(lower, "microsoft"):
		return "icon-microsoft"
	case strings.Contains(lower, "github"):
		return "icon-github"
	default:
		return "icon-purple"
	}
}
