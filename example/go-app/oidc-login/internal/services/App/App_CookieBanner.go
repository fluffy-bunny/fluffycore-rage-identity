package App

import (
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func (s *service) renderCookieBanner() app.UI {
	cookieStatement := s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyCookieConsentStatement)
	acceptText := s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyAccept)
	return app.Div().Class("cookie-banner").Body(
		app.Div().Class("cookie-content").Body(
			app.Span().Class("cookie-text").Text(cookieStatement),
			app.Button().
				Class("cookie-accept-btn").
				OnClick(s.handleAcceptCookies).
				Text(acceptText),
		),
	)
}
