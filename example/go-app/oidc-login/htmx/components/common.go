package components

import (
	"bytes"
	"net/http"
	"strings"

	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	echo "github.com/labstack/echo/v5"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// RenderContext holds the common data needed by all components.
type RenderContext struct {
	Paths            *wellknown_echo.Paths
	CSRF             string
	Localizer        *i18n.Localizer
	CacheBustVersion string
}

// Localize returns the localized message for the given key.
func Localize(l *i18n.Localizer, key string) string {
	message, _ := l.LocalizeMessage(&i18n.Message{ID: key})
	return message
}

// RenderNode writes a gomponents Node to the echo response.
func RenderNode(c *echo.Context, code int, node g.Node) error {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(code, buf.String())
}

// NewRenderContext creates a RenderContext from the echo context.
func NewRenderContext(c *echo.Context, localizer *i18n.Localizer) *RenderContext {
	csrfValue := c.Get("csrf")
	csrfStr := ""
	if csrfValue != nil {
		if str, ok := csrfValue.(string); ok {
			csrfStr = str
		}
	}
	return &RenderContext{
		Paths:     wellknown_echo.NewPaths(),
		CSRF:      csrfStr,
		Localizer: localizer,
	}
}

// L is a shorthand for Localize using the RenderContext's localizer.
func (rc *RenderContext) L(key string) string {
	return Localize(rc.Localizer, key)
}

// LF localizes a key and replaces {placeholders} with values from the map.
func (rc *RenderContext) LF(key string, replace map[string]string) string {
	template := Localize(rc.Localizer, key)
	for k, v := range replace {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template
}

// --- Shared sub-components ---

// ErrorMessages renders error messages if any exist.
func ErrorMessages(errors []string) g.Node {
	if len(errors) == 0 {
		return nil
	}
	children := make([]g.Node, len(errors))
	for i, e := range errors {
		children[i] = Span(g.Text(e))
	}
	return Div(Class("error-message"), g.Group(children))
}

// CsrfInput renders a hidden CSRF token input.
func CsrfInput(csrf string) g.Node {
	return Input(Type("hidden"), Name("csrf"), Value(csrf))
}

// FormGroupField renders a form-group with label and input.
func FormGroupField(labelText, inputType, id, name, value string, extraAttrs ...g.Node) g.Node {
	attrs := []g.Node{
		Type(inputType),
		ID(id),
		Name(name),
	}
	if value != "" {
		attrs = append(attrs, Value(value))
	}
	attrs = append(attrs, extraAttrs...)
	return Div(Class("form-group"),
		Label(g.Attr("for", id), g.Text(labelText)),
		Input(attrs...),
	)
}

// ButtonGroup renders a flex button group container.
func ButtonGroup(children ...g.Node) g.Node {
	return Div(Class("button-group"), g.Group(children))
}

// PrimaryButton renders a primary submit button with optional HTMX indicator.
func PrimaryButton(text, indicatorID string) g.Node {
	return Button(Type("submit"), Class("btn-primary"),
		g.Text(text),
		Span(ID(indicatorID), Class("htmx-indicator"), g.Attr("role", "status"), g.Text(" ...")),
	)
}

// SecondaryButton renders a secondary button that navigates via HTMX GET.
func SecondaryButton(text, hxGetURL string) g.Node {
	return Button(Type("button"), Class("btn-secondary"),
		g.Attr("hx-get", hxGetURL),
		g.Attr("hx-target", "#main-content"),
		g.Attr("hx-swap", "innerHTML"),
		g.Text(text),
	)
}

// HtmxForm renders a form with HTMX post attributes.
func HtmxForm(postURL, indicatorID string, children ...g.Node) g.Node {
	return FormEl(
		g.Attr("hx-post", postURL),
		g.Attr("hx-target", "#main-content"),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-indicator", "#"+indicatorID),
		g.Group(children),
	)
}

// SocialIdpButtons renders branded social IDP login buttons matching the WASM layout.
func SocialIdpButtons(idps []*proto_oidc_models.IDP, csrf, postURL, orText string) g.Node {
	if len(idps) == 0 {
		return nil
	}
	buttons := make([]g.Node, 0, len(idps))
	for _, idp := range idps {
		icon, cssClass := socialIdpBranding(idp.Slug)
		buttons = append(buttons, FormEl(
			g.Attr("style", "display:contents"),
			g.Attr("hx-post", postURL),
			g.Attr("hx-target", "#main-content"),
			g.Attr("hx-swap", "innerHTML"),
			CsrfInput(csrf),
			Input(Type("hidden"), Name("idp_hint"), Value(idp.Slug)),
			Button(Type("submit"), Class("social-btn "+cssClass),
				g.Raw(icon),
				Span(g.Text(socialIdpLabel(idp.Slug))),
			),
		))
	}
	return Div(Class("social-login-section"),
		Div(Class("divider"), Span(g.Text(orText))),
		Div(Class("social-buttons"), g.Group(buttons)),
	)
}

// socialIdpBranding returns the SVG icon and CSS class for a known IDP slug.
func socialIdpBranding(slug string) (string, string) {
	switch slug {
	case "google-social":
		return GoogleIconSVG, "google-btn"
	case "microsoft-social":
		return MicrosoftIconSVG, "microsoft-btn"
	case "github-social":
		return GitHubIconSVG, "github-btn"
	default:
		return "", ""
	}
}

// socialIdpLabel returns the display label for a known IDP slug.
func socialIdpLabel(slug string) string {
	switch slug {
	case "google-social":
		return "Google"
	case "microsoft-social":
		return "Microsoft"
	case "github-social":
		return "GitHub"
	default:
		return slug
	}
}

// PasskeyLoginSection renders the passkey login option with divider.
// Uses client-side JavaScript to call the WebAuthn LoginUser() flow.
func PasskeyLoginSection(csrf, keepSignedInURL, text string) g.Node {
	return Div(Class("passkey-login-section"),
		Div(Class("divider"), Span(g.Text("OR"))),
		Button(Type("button"), Class("passkey-btn"),
			ID("passkey-login-btn"),
			g.Raw(PasskeyIconSVG),
			Span(g.Text(text)),
		),
		Script(g.Raw(`document.getElementById("passkey-login-btn").addEventListener("click",function(){
  this.disabled=true;
  LoginUser("",false,
    function(errMsg){
      document.getElementById("passkey-login-btn").disabled=false;
      alert("Passkey authentication failed: "+errMsg);
    },
    function(data){
      if(data.directive==="displayKeepSignedInPage"){
        htmx.ajax("GET","`+keepSignedInURL+`",{target:"#main-content",swap:"innerHTML"});
      } else if(data.directive==="redirect"&&data.directiveRedirect&&data.directiveRedirect.redirectUri){
        window.location.href=data.directiveRedirect.redirectUri;
      }
    }
  );
});`)),
	)
}

// --- SVG Icon Constants (matching WASM branding) ---

const PasskeyIconSVG = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path></svg>`

const GoogleIconSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>`

const MicrosoftIconSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M11.4 11.4H2V2h9.4v9.4z" fill="#F25022"/><path d="M22 11.4h-9.4V2H22v9.4z" fill="#7FBA00"/><path d="M11.4 22H2v-9.4h9.4V22z" fill="#00A4EF"/><path d="M22 22h-9.4v-9.4H22V22z" fill="#FFB900"/></svg>`

const GitHubIconSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.137 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z"/></svg>`
