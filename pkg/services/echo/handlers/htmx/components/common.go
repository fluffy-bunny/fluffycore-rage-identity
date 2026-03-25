package components

import (
	"bytes"
	"net/http"

	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	echo "github.com/labstack/echo/v5"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// RenderContext holds the common data needed by all components.
type RenderContext struct {
	Paths     *wellknown_echo.Paths
	CSRF      string
	Localizer *i18n.Localizer
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

// SocialIdpButtons renders social IDP login buttons.
func SocialIdpButtons(idps []*proto_oidc_models.IDP, csrf, postURL, orText string) g.Node {
	if len(idps) == 0 {
		return nil
	}
	forms := make([]g.Node, len(idps))
	for i, idp := range idps {
		forms[i] = FormEl(
			g.Attr("hx-post", postURL),
			g.Attr("hx-target", "#main-content"),
			g.Attr("hx-swap", "innerHTML"),
			CsrfInput(csrf),
			Input(Type("hidden"), Name("idp_hint"), Value(idp.Slug)),
			Button(Type("submit"), Class("btn-social"), g.Text(idp.Slug)),
		)
	}
	return Div(Class("social-login-section"),
		Div(Class("divider"), Span(g.Text(orText))),
		Div(Class("social-logins"), g.Group(forms)),
	)
}
