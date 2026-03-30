package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// ShellPage renders the full HTML document shell for the management dashboard.
func ShellPage(rc *RenderContext) g.Node {
	// Determine the initial page to load
	initialPage := rc.Paths.HTMXManagementHome
	if rc.DeepLinkPath != "" {
		initialPage = rc.DeepLinkPath
	}

	// Derive HTML <title> from config branding
	pageTitle := "Account Management"
	if rc.AppConfig != nil && rc.AppConfig.BannerBranding.Title != "" {
		pageTitle = rc.AppConfig.BannerBranding.Title
	}

	return c.HTML5(c.HTML5Props{
		Title:    pageTitle,
		Language: "en",
		Head: []g.Node{
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1, shrink-to-fit=no")),
			Meta(Name("description"), g.Attr("content", "Account Management")),
			Link(g.Attr("rel", "icon"), Type("image/x-icon"), Href("/static/assets/favicon.ico?v="+rc.CacheBustVersion)),
			Link(g.Attr("rel", "stylesheet"), Href("/static/go-app/management/htmx/styles.css?v="+rc.CacheBustVersion)),
			Script(Src("https://unpkg.com/htmx.org@2.0.4"),
				g.Attr("integrity", "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"),
				g.Attr("crossorigin", "anonymous")),
			Script(Src("/static/go-app/oidc-login/htmx/webauthn.js?v=" + rc.CacheBustVersion)),
			Meta(Name("htmx-config"), g.Attr("content", `{"responseHandling":[{"code":".*", "swap": true}]}`)),
			StyleEl(g.Raw(`.htmx-indicator { display: none; }
.htmx-request .htmx-indicator, .htmx-request.htmx-indicator { display: inline-block; }
.field-error { color: #e74c3c; font-size: 0.85rem; margin-top: 0.25rem; }`)),
		},
		Body: []g.Node{
			// Unregister stale service workers from prior WASM deployments
			Script(g.Raw(`if("serviceWorker"in navigator){navigator.serviceWorker.getRegistrations().then(function(r){r.forEach(function(reg){reg.unregister()})})}`)),
			Div(Class("dashboard-layout"),
				DashboardHeader(rc),
				Div(Class("dashboard-body"),
					Sidebar(rc),
					Div(Class("dashboard-main"),
						Div(ID("dashboard-main"),
							g.Attr("hx-get", initialPage),
							g.Attr("hx-trigger", "load"),
							g.Attr("hx-swap", "innerHTML"),
							g.Attr("hx-push-url", "true"),
							P(g.Attr("style", "text-align:center;padding:40px 0;"), g.Text("Loading...")),
						),
					),
				),
			),
			CookieBanner(rc),
			// Handle browser back/forward navigation
			Script(g.Raw(`window.addEventListener("popstate",function(){
  var p=location.pathname;
  if(p.startsWith("/management/")){
    htmx.ajax("GET",p,{target:"#dashboard-main",swap:"innerHTML"});
  }
});`)),
			// Sidebar toggle for mobile (checkbox+CSS approach)
			Script(g.Raw(`document.addEventListener("DOMContentLoaded",function(){
  var links=document.querySelectorAll(".sidebar-link");
  var toggle=document.getElementById("sidebar-toggle");
  links.forEach(function(link){
    link.addEventListener("click",function(){
      if(toggle&&toggle.checked){toggle.checked=false;}
    });
  });
});`)),
		},
	})
}
