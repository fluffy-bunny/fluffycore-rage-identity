package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// CookieIconSVG is an inline SVG cookie icon for the consent banner.
const CookieIconSVG = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M18.12,9.78a3,3,0,0,1-3.9-3.9A3,3,0,0,1,12,3a9,9,0,1,0,9,9A3,3,0,0,1,18.12,9.78Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"/></svg>`

// CookieBanner renders a GDPR cookie consent banner.
func CookieBanner(rc *RenderContext) g.Node {
	cookieText := rc.L("cookie_consent_statement")
	acceptText := rc.L("accept")

	return g.Group([]g.Node{
		Div(ID("cookie-banner"), Class("cookie-consent-banner"),
			g.Attr("style", "display:none"),
			Div(Class("cookie-consent-content"),
				Div(Class("cookie-consent-text"),
					Span(Class("cookie-consent-icon"), g.Raw(CookieIconSVG)),
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
