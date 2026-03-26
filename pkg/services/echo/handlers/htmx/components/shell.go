package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// ShellData holds data for the full HTML shell page.
type ShellData struct {
	*RenderContext
	BrandTitle string
}

// ShellPage renders the full HTML document shell with wizard-container and app-header.
func ShellPage(data ShellData) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    "Login",
		Language: "en",
		Head: []g.Node{
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1, shrink-to-fit=no")),
			Meta(Name("description"), g.Attr("content", "OIDC Login")),
			Link(g.Attr("rel", "icon"), Type("image/x-icon"), Href("/static/assets/favicon.ico")),
			Link(g.Attr("rel", "stylesheet"), Href("/static/go-app/oidc-login/static_output/web/styles.css?v="+data.CacheBustVersion)),
			Script(Src("https://unpkg.com/htmx.org@2.0.4"),
				g.Attr("integrity", "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"),
				g.Attr("crossorigin", "anonymous")),
			Meta(Name("htmx-config"), g.Attr("content", `{"responseHandling":[{"code":".*", "swap": true}]}`)),
			StyleEl(g.Raw(`.htmx-indicator { display: none; }
.htmx-request .htmx-indicator, .htmx-request.htmx-indicator { display: inline-block; }`)),
		},
		Body: []g.Node{
			// Unregister any lingering WASM service workers from previous visits
			Script(g.Raw(`if("serviceWorker"in navigator){navigator.serviceWorker.getRegistrations().then(function(r){r.forEach(function(reg){reg.unregister()})})}`)),
			Div(Class("wizard-container"),
				Div(Class("app-header"),
					Div(Class("header-content"),
						Div(Class("logo-title-group"),
							Img(Src("/static/go-app/oidc-login/static_output/web/m_logo.svg"), Alt("Logo"), Class("app-logo")),
							Div(Class("title-version-group"),
								Div(Class("app-title"), g.Text(data.BrandTitle)),
							),
						),
					),
				),
				Div(ID("main-content"), Class("step-container"),
					g.Attr("hx-get", data.Paths.HTMXHome),
					g.Attr("hx-trigger", "load"),
					g.Attr("hx-swap", "innerHTML"),
					P(g.Attr("style", "text-align:center;padding:40px 0;"), g.Text("Loading...")),
				),
			),
			CookieBanner(data.RenderContext),
		},
	})
}
