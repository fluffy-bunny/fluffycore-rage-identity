package App

import (
	"context"

	"github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/common"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func (s *service) renderHeader(ctx context.Context) app.UI {

	appConfig := s.appConfigAccessor.GetAppConfig(ctx)

	showVersion := func() app.UI {
		if appConfig.BannerBranding.ShowBannerVersion {
			return app.Span().Class("app-version").Text("v" + common.AppVersion)
		}
		return app.Div()
	}
	return app.Div().Class("app-header").Body(
		app.Div().Class("header-content").Body(
			app.Div().Class("logo-title-group").Body(
				app.Img().
					Src(appConfig.BannerBranding.LogoURL).
					Alt("Rage Accounts Logo").
					Class("app-logo"),
				app.Div().Class("title-version-group").Body(
					app.H1().Class("app-title").Text(appConfig.BannerBranding.Title),
					showVersion(),
				),
			),
		),
	)
}
