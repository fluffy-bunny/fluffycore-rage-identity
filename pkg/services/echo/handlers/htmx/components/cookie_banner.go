package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// CookieBanner renders a GDPR cookie consent banner.
// Uses localStorage to track acceptance, matching the WASM implementation.
func CookieBanner(rc *RenderContext) g.Node {
	cookieText := rc.L("cookie_consent_statement")
	acceptText := rc.L("accept")

	return g.Group([]g.Node{
		Div(ID("cookie-banner"), Class("cookie-consent-banner"),
			g.Attr("style", "display:none"),
			Div(Class("cookie-consent-content"),
				Div(Class("cookie-consent-text"),
					Span(Class("cookie-consent-icon"), g.Raw("&#x1F36A;")),
					Span(g.Text(cookieText)),
				),
				Button(Class("cookie-consent-button"),
					ID("cookie-accept-btn"),
					g.Text(acceptText),
				),
			),
		),
		Script(g.Raw(`(function(){
  var b=document.getElementById("cookie-banner");
  if(localStorage.getItem("cookiesAccepted")==="true"){return;}
  b.style.display="";
  document.getElementById("cookie-accept-btn").addEventListener("click",function(){
    localStorage.setItem("cookiesAccepted","true");
    b.style.display="none";
  });
})();`)),
	})
}
