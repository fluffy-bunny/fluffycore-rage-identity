package App

import (
	go_app_common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/common"
	"github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// renderDashboardHeader is for the main dashboard layout
func (s *service) renderDashboardHeader() app.UI {
	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)

	return app.Header().Class("dashboard-header").Body(
		app.Div().Class("dashboard-header-left").Body(
			app.Button().
				Class("sidebar-toggle").
				OnClick(s.handleToggleSidebar).
				Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="3" y1="12" x2="21" y2="12"></line>
						<line x1="3" y1="6" x2="21" y2="6"></line>
						<line x1="3" y1="18" x2="21" y2="18"></line>
					</svg>`),
				),
			app.Div().Class("dashboard-logo-group").Body(
				app.Img().
					Src(func() string {
						if appConfig.BannerBranding.LogoURL != "" {
							return appConfig.BannerBranding.LogoURL
						}
						return "/web/m_logo.svg"
					}()).
					Alt(func() string {
						if appConfig.BannerBranding.Title != "" {
							return appConfig.BannerBranding.Title
						}
						return "Rage Accounts"
					}()).
					Class("dashboard-logo"),
				app.Div().Class("dashboard-title-container").Body(
					app.Span().Class("dashboard-title").Text(func() string {
						if appConfig.BannerBranding.Title != "" {
							return appConfig.BannerBranding.Title
						}
						return "Rage Accounts"
					}()),
					app.If(appConfig.BannerBranding.ShowBannerVersion, func() app.UI {
						return app.Span().Class("dashboard-version").Text("v" + common.AppVersion)
					}),
				),
			),
		),
		app.Div().Class("dashboard-header-right").Body(
			app.If(s.isAuthenticated, func() app.UI {
				return s.renderAuthenticatedMenu()
			}).Else(func() app.UI {
				return s.renderUnauthenticatedMenu()
			}),
		),
	)
}

func (s *service) renderAuthenticatedMenu() app.UI {
	return app.Div().Class("user-menu-container").Body(
		app.Button().
			Class("user-menu-button").
			OnClick(func(ctx app.Context, e app.Event) {
				e.PreventDefault()
				s.showUserMenu = !s.showUserMenu
				ctx.Update()
			}).
			Body(
				app.Div().Class("user-avatar").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
								<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
							</svg>`),
				),
				app.Span().Class("user-name").Text(func() string {
					if s.profile.Email != "" {
						return s.profile.Email
					}
					return "User"
				}()),
				app.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="dropdown-icon">
							<path d="M7 10l5 5 5-5z"/>
						</svg>`),
			),
		app.If(s.showUserMenu, func() app.UI {
			return s.renderUserDropdown()
		}),
	)
}

func (s *service) renderUnauthenticatedMenu() app.UI {
	return app.Div().Class("login-menu-container").Body(
		app.Button().
			Class("login-menu-button").
			OnClick(s.handleSignin).
			Body(
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLogin)),
				app.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="dropdown-icon">
							<path d="M7 10l5 5 5-5z"/>
						</svg>`),
			),
		app.If(s.showUserMenu, func() app.UI {
			return s.renderLoginDropdown()
		}),
	)
}

func (s *service) renderLoginDropdown() app.UI {
	return app.Div().Class("login-dropdown").Body(
		app.Button().
			Class("login-dropdown-item").
			OnClick(func(ctx app.Context, e app.Event) {
				e.PreventDefault()
				s.showUserMenu = false
				// Simulate login - in real app would navigate to login page
				s.isAuthenticated = true
				ctx.Update()
			}).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
					<polyline points="10 17 15 12 10 7"></polyline>
					<line x1="15" y1="12" x2="3" y2="12"></line>
				</svg>`),
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignIn)),
			),
	)
}

func (s *service) renderUserDropdown() app.UI {
	return app.Div().Class("user-dropdown").Body(
		app.Div().Class("user-dropdown-header").Body(
			app.Div().Class("user-dropdown-avatar").Body(
				app.Raw(`<svg width="40" height="40" viewBox="0 0 24 24" fill="currentColor">
					<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
				</svg>`),
			),
			app.Div().Class("user-dropdown-info").Body(
				app.Div().Class("user-dropdown-name").Text(func() string {
					if s.profile.Email != "" {
						return s.profile.Email
					}
					return "User"
				}()),
				app.Div().Class("user-dropdown-email").Text(func() string {
					if s.profile.Subject != "" {
						return s.profile.Subject
					}
					return ""
				}()),
			),
		),
		app.Div().Class("user-dropdown-divider"),
		app.A().
			Href(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Profile)).
			Class("user-dropdown-item").
			OnClick(func(ctx app.Context, e app.Event) {
				e.PreventDefault()
				s.showUserMenu = false
				ctx.Navigate(contracts_routes.GetFixedRoute(contracts_routes.WellknownRoute_Profile))
			}).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
					<circle cx="12" cy="7" r="4"></circle>
				</svg>`),
				app.Span().Text("My Profile"),
			),
		app.Button().
			Class("user-dropdown-item").
			OnClick(s.handleSignOut).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
					<polyline points="16 17 21 12 16 7"></polyline>
					<line x1="21" y1="12" x2="9" y2="12"></line>
				</svg>`),
				app.Span().Text("Sign Out"),
			),
	)
}

func (s *service) renderSidebar() app.UI {
	sidebarClass := "dashboard-sidebar"
	if s.showSidebar {
		sidebarClass += " show"
	}

	return app.Nav().Class(sidebarClass).Body(
		app.Div().Class("sidebar-content").Body(
			app.If(s.isAuthenticated, func() app.UI {
				return s.renderAuthenticatedSidebar()
			}).Else(func() app.UI {
				return s.renderUnauthenticatedSidebar()
			}),
		),
	)
}

func (s *service) renderAuthenticatedSidebar() app.UI {
	appConfig := s.appConfigAccessor.GetAppConfig(s.AppContext)

	return app.Div().Class("sidebar-section").Body(
		app.Div().Class("sidebar-section-title").Text("Account"),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_Home,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyHome),
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
				<polyline points="9 22 9 12 15 12 15 22"></polyline>
			</svg>`,
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_Profile,
			"Profile",
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>`,
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_PasswordManager,
			"Password",
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
						<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
					</svg>`,
		),
		app.If(appConfig.EnabledWebAuthN,
			func() app.UI {
				return s.renderSidebarLink(
					contracts_routes.WellknownRoute_PasskeyManager,
					"Passkeys",
					go_app_common.PasskeyIconSmallSVG,
				)
			},
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_LinkedAccounts,
			"Linked Accounts",
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
						<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
					</svg>`,
		),
	)
}

func (s *service) renderUnauthenticatedSidebar() app.UI {
	return app.Div().Class("sidebar-section").Body(
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_Home,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyHome),
			`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
				<polyline points="9 22 9 12 15 12 15 22"></polyline>
			</svg>`,
		),
	)
}

func (s *service) renderSidebarLink(route contracts_routes.WellknownRoute, label string, iconSVG string) app.UI {
	isActive := s.currentPage == route
	linkClass := "sidebar-link"
	if isActive {
		linkClass += " active"
	}

	return app.A().
		Href(contracts_routes.GetFixedRoute(route)).
		Class(linkClass).
		OnClick(func(ctx app.Context, e app.Event) {
			e.PreventDefault()
			s.showSidebar = false
			ctx.Navigate(contracts_routes.GetFixedRoute(route))
		}).
		Body(
			app.Raw(iconSVG),
			app.Span().Text(label),
		)
}
