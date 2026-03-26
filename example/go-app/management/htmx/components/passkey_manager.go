package components

import (
	"fmt"

	api_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_passkey"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PasskeyPageData holds data for the passkey manager page.
type PasskeyPageData struct {
	*RenderContext
	Passkeys        []api_passkey.PasskeyItem
	IsClaimedDomain bool
	EnabledWebAuthN bool
	Error           string
	Success         string
	RenameID        string // credential ID being renamed, if any
}

// PasskeyPage renders the passkey manager page.
func PasskeyPage(d *PasskeyPageData) g.Node {
	children := []g.Node{
		Div(Class("profile-header"),
			H1(g.Text(d.L("mgmt_passkeys"))),
			P(Class("profile-subtitle"), g.Text(d.L("mgmt_manage_passkeys"))),
		),
	}

	if d.IsClaimedDomain || !d.EnabledWebAuthN {
		children = append(children,
			ProfileCard(
				CardHeader(PasskeyIconSVG, d.L("mgmt_passkeys_not_available"), "", "passkey-icon"),
				Div(Class("card-body"),
					P(g.Text(d.L("mgmt_passkeys_not_available_desc"))),
				),
			),
		)
		return Div(Class("profile-container"), g.Group(children))
	}

	if d.Error != "" {
		children = append(children, ErrorBanner(d.Error))
	}
	if d.Success != "" {
		children = append(children, SuccessBanner(d.Success))
	}

	// All cards go into a single profile-cards container for consistent spacing
	cards := []g.Node{
		// Passkey usage definition card
		ProfileCard(
			CardHeader(PasskeyIconSVG, d.L("mgmt_passkeys"),
				d.L("mgmt_manage_passkeys"), "passkey-icon"),
			Div(Class("card-body"),
				P(g.Attr("style", "margin-bottom:16px;color:var(--text-secondary);"),
					g.Text(d.L("mgmt_passkey_usage_definition"))),
				Div(Class("card-actions"),
					Button(Type("button"), Class("btn-primary"), ID("add-passkey-btn"),
						g.Raw(PasskeyIconSVG),
						Span(g.Text(" "+d.L("mgmt_add_passkey"))),
					),
				),
			),
		),
	}

	// Individual passkey cards
	for _, pk := range d.Passkeys {
		cards = append(cards, passkeyCard(d, pk))
	}

	children = append(children,
		Div(Class("profile-cards"), g.Group(cards)),
	)

	// WebAuthn registration script
	children = append(children, passkeyRegistrationScript(d))

	return Div(Class("profile-container"), g.Group(children))
}

func passkeyCard(d *PasskeyPageData, pk api_passkey.PasskeyItem) g.Node {
	if d.RenameID == pk.ID {
		return passkeyRenameCard(d, pk)
	}
	return ProfileCard(
		Div(Class("card-header"),
			Div(Class("card-header-content"),
				Div(Class("card-icon passkey-icon"), g.Raw(PasskeyIconSVG)),
				Div(Class("card-title-group"),
					H2(g.Text(pk.FriendlyName)),
				),
			),
		),
		Div(Class("card-body"),
			Div(Class("card-actions"),
				Button(Type("button"), Class("btn-secondary"),
					g.Attr("hx-get", d.Paths.HTMXManagementPasskey+fmt.Sprintf("?rename=%s", pk.ID)),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text(d.L("mgmt_rename")),
				),
				Button(Type("button"), Class("btn-unlink"),
					g.Attr("hx-post", d.Paths.HTMXManagementPasskey),
					g.Attr("hx-target", "#dashboard-main"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-vals", fmt.Sprintf(`{"action":"delete","credentialId":"%s","csrf":"%s"}`, pk.ID, d.CSRF)),
					g.Attr("hx-confirm", d.L("mgmt_confirm_delete_passkey")),
					g.Text(d.L("mgmt_delete")),
				),
			),
		),
	)
}

func passkeyRenameCard(d *PasskeyPageData, pk api_passkey.PasskeyItem) g.Node {
	return ProfileCard(
		CardHeader(PasskeyIconSVG, d.L("mgmt_rename"), "", "passkey-icon"),
		Div(Class("card-body"),
			HtmxForm(d.Paths.HTMXManagementPasskey, "rename-indicator",
				CsrfInput(d.CSRF),
				Input(Type("hidden"), Name("action"), Value("rename")),
				Input(Type("hidden"), Name("credentialId"), Value(pk.ID)),
				FormGroup(d.L("mgmt_rename"), "text", "friendlyName", "friendlyName", pk.FriendlyName,
					g.Attr("required"),
				),
				ButtonGroup(
					PrimaryButton(d.L("mgmt_save"), "rename-indicator"),
					SecondaryButton(d.L("mgmt_cancel"), d.Paths.HTMXManagementPasskey),
				),
			),
		),
	)
}

func passkeyRegistrationScript(d *PasskeyPageData) g.Node {
	return Script(g.Raw(`document.getElementById("add-passkey-btn").addEventListener("click",function(){
  this.disabled=true;
  var btn=this;
  registerPasskey("").then(function(success){
    btn.disabled=false;
    if(success){
      htmx.ajax("GET","` + d.Paths.HTMXManagementPasskey + `",{target:"#dashboard-main",swap:"innerHTML"});
    }
  }).catch(function(err){
    btn.disabled=false;
    console.error("Passkey registration error:",err);
  });
});`))
}
