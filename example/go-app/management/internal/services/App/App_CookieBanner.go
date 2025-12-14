package App

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func (s *service) renderCookieBanner() app.UI {
	return app.Div().Class("cookie-banner").Body(
		app.Div().Class("cookie-content").Body(
			app.Span().Class("cookie-text").Text("We use cookies to enhance your experience."),
			app.Button().
				Class("cookie-accept-btn").
				OnClick(s.handleAcceptCookies).
				Text("Accept"),
		),
	)
}
