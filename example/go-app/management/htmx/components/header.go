package components

import (
	"strings"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// DashboardHeader renders the dashboard header with logo, title, and user menu.
func DashboardHeader(rc *RenderContext) g.Node {
	displayName := rc.UserName
	if displayName == "" {
		displayName = rc.UserEmail
	}
	if displayName == "" {
		displayName = rc.L("mgmt_user")
	}

	// Branding from AppConfig
	logoURL := "/static/go-app/management/htmx/m_logo.svg"
	title := "Rage Accounts"
	if rc.AppConfig != nil {
		if rc.AppConfig.BannerBranding.LogoURL != "" {
			logoURL = rc.AppConfig.BannerBranding.LogoURL
			// Config may contain legacy relative paths like "web/m_logo.svg"; resolve to htmx static dir
			if !strings.HasPrefix(logoURL, "/") && !strings.HasPrefix(logoURL, "http") {
				// Strip legacy "web/" prefix if present
				logoURL = strings.TrimPrefix(logoURL, "web/")
				logoURL = "/static/go-app/management/htmx/" + logoURL
			}
		}
		if rc.AppConfig.BannerBranding.Title != "" {
			title = rc.AppConfig.BannerBranding.Title
		}
	}

	// Title container children
	titleChildren := []g.Node{
		Span(Class("dashboard-title"), g.Text(title)),
	}
	if rc.AppConfig != nil && rc.AppConfig.BannerBranding.ShowBannerVersion {
		titleChildren = append(titleChildren,
			Span(Class("dashboard-version"), g.Text("v"+rc.AppVersion)),
		)
	}

	return Header(Class("dashboard-header"),
		Div(Class("dashboard-header-left"),
			// Sidebar toggle for mobile — checkbox+CSS approach
			Label(Class("sidebar-toggle"), g.Attr("for", "sidebar-toggle"),
				g.Raw(HamburgerIconSVG),
			),
			Div(Class("dashboard-logo-group"),
				Img(Src(logoURL+"?v="+rc.CacheBustVersion), Alt(title), Class("dashboard-logo")),
				Div(Class("dashboard-title-container"),
					g.Group(titleChildren),
				),
			),
		),
		Div(Class("dashboard-header-right"),
			Div(Class("user-menu-container"),
				Button(Type("button"), Class("user-menu-button"), ID("user-menu-btn"),
					Div(Class("user-avatar"),
						g.Raw(PersonLoggedInIconSmallSVG),
					),
					Span(Class("user-name"), g.Text(displayName)),
					g.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="dropdown-icon"><path d="M7 10l5 5 5-5z"/></svg>`),
				),
				Div(Class("user-dropdown"), ID("user-dropdown"),
					Div(Class("user-dropdown-header"),
						Div(Class("user-dropdown-avatar"),
							g.Raw(PersonLoggedInIconLargeSVG),
						),
						Div(Class("user-dropdown-info"),
							Div(Class("user-dropdown-name"), g.Text(func() string {
								if rc.UserEmail != "" {
									return rc.UserEmail
								}
								return rc.L("mgmt_user")
							}())),
							Div(Class("user-dropdown-email"), g.Text(rc.UserSubject)),
						),
					),
					Div(Class("user-dropdown-divider")),
					A(Href(rc.Paths.HTMXManagementProfile),
						Class("user-dropdown-item"),
						g.Attr("hx-get", rc.Paths.HTMXManagementProfile),
						g.Attr("hx-target", "#dashboard-main"),
						g.Attr("hx-swap", "innerHTML"),
						g.Attr("hx-push-url", "true"),
						g.Raw(PersonIconSVG),
						Span(g.Text(rc.L("mgmt_my_profile"))),
					),
					A(Href("/logout"), Class("user-dropdown-item"),
						g.Raw(SignOutIconSVG),
						Span(g.Text(rc.L("mgmt_sign_out"))),
					),
				),
			),
		),
		// Dropdown toggle + close-on-outside-click (event delegation for HTMX compatibility)
		Script(g.Raw(`(function(){
  if(window.__dropdownInit) return;
  window.__dropdownInit=true;
  document.addEventListener("click",function(e){
    var btn=document.getElementById("user-menu-btn");
    var dd=document.getElementById("user-dropdown");
    if(!btn||!dd) return;
    if(btn.contains(e.target)){
      dd.classList.toggle("show");
      return;
    }
    dd.classList.remove("show");
  });
  document.body.addEventListener("htmx:beforeSwap",function(){
    var dd=document.getElementById("user-dropdown");
    if(dd) dd.classList.remove("show");
  });
})();`)),
	)
}
