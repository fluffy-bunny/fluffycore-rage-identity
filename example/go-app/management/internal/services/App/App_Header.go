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
					app.Raw(go_app_common.HamburgerMenuIconSmallSVG),
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
					app.Raw(go_app_common.PersonLoggedInIconSmallSVG),
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
				app.Raw(go_app_common.SignOutIconSmallSVG),
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignIn)),
			),
	)
}

func (s *service) renderUserDropdown() app.UI {
	return app.Div().Class("user-dropdown").Body(
		app.Div().Class("user-dropdown-header").Body(
			app.Div().Class("user-dropdown-avatar").Body(
				app.Raw(go_app_common.PersonLoggedInIconLargeSVG),
			),
			app.Div().Class("user-dropdown-info").Body(
				app.Div().Class("user-dropdown-name").Text(func() string {
					if s.profile.Email != "" {
						return s.profile.Email
					}
					return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUser)
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
				app.Raw(go_app_common.PersonIconSmallSVG),
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyMyProfile)),
			),
		app.Button().
			Class("user-dropdown-item").
			OnClick(s.handleSignOut).
			Body(
				app.Raw(go_app_common.SignOutIconSmallSVG),
				app.Span().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeySignOut)),
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
			go_app_common.HomeIconSmallSVG,
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_Profile,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyProfile),
			go_app_common.PersonIconSmallSVG,
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_PasswordManager,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPassword),
			go_app_common.LockIconSmallSVG,
		),
		app.If(appConfig.EnabledWebAuthN,
			func() app.UI {
				return s.renderSidebarLink(
					contracts_routes.WellknownRoute_PasskeyManager,
					s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasskeys),
					go_app_common.PasskeyIconSmallSVG,
				)
			},
		),
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_LinkedAccounts,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLinkedAccounts),
			go_app_common.LinkIconSmallSVG,
		),
	)
}

func (s *service) renderUnauthenticatedSidebar() app.UI {
	return app.Div().Class("sidebar-section").Body(
		s.renderSidebarLink(
			contracts_routes.WellknownRoute_Home,
			s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyHome),
			go_app_common.HomeIconSmallSVG,
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
