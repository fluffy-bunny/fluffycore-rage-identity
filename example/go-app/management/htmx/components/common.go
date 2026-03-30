package components

import (
	"bytes"
	"net/http"
	"strings"

	common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/common"
	management_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	echo "github.com/labstack/echo/v5"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// RenderContext holds common data needed by all management components.
type RenderContext struct {
	Paths            *wellknown_echo.Paths
	CSRF             string
	Localizer        *i18n.Localizer
	CacheBustVersion string
	ActivePage       string
	DeepLinkPath     string // initial page to load (for deep linking)
	UserEmail        string
	UserName         string
	UserSubject      string
	AppVersion       string
	AppConfig        *management_config.AppConfig
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

// L localizes a key.
func (rc *RenderContext) L(key string) string {
	msg, _ := rc.Localizer.LocalizeMessage(&i18n.Message{ID: key})
	return msg
}

// LF localizes a key with placeholder replacements.
func (rc *RenderContext) LF(key string, replace map[string]string) string {
	template := rc.L(key)
	for k, v := range replace {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template
}

// RenderNode writes a gomponents Node to the echo response.
func RenderNode(c *echo.Context, code int, node g.Node) error {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(code, buf.String())
}

// IsHTMXRequest returns true if the request was made by HTMX (has HX-Request header).
func IsHTMXRequest(c *echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

// --- Shared sub-components ---

// CsrfInput renders a hidden CSRF token input.
func CsrfInput(csrf string) g.Node {
	return Input(Type("hidden"), Name("csrf"), Value(csrf))
}

// ErrorBanner renders an error alert banner.
func ErrorBanner(message string) g.Node {
	if message == "" {
		return nil
	}
	return Div(Class("error-message"), g.Attr("role", "alert"), g.Attr("aria-live", "assertive"),
		Span(g.Text(message)),
	)
}

// SuccessBanner renders a success alert banner.
func SuccessBanner(message string) g.Node {
	if message == "" {
		return nil
	}
	return Div(Class("success-message"), g.Attr("role", "status"), g.Attr("aria-live", "polite"),
		Span(g.Text(message)),
	)
}

// ProfileCard renders a card in the profile/management layout.
func ProfileCard(children ...g.Node) g.Node {
	return Div(Class("profile-card"), g.Group(children))
}

// CardHeader renders a card header with icon and title/subtitle.
func CardHeader(iconSVG, title, subtitle, iconColorClass string) g.Node {
	return Div(Class("card-header"),
		Div(Class("card-header-content"),
			Div(Class("card-icon "+iconColorClass), g.Raw(iconSVG)),
			Div(Class("card-title-group"),
				H2(g.Text(title)),
				P(Class("card-description"), g.Text(subtitle)),
			),
		),
	)
}

// InfoRow renders a label-value info row.
func InfoRow(label, value string) g.Node {
	return Div(Class("info-row"),
		Span(Class("info-label"), g.Text(label)),
		Span(Class("info-value"), g.Text(value)),
	)
}

// FormGroup renders a form-group with label and input.
func FormGroup(labelText, inputType, id, name, value string, extraAttrs ...g.Node) g.Node {
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

// FormGroupWithError renders a form-group with label, input, and an inline error slot.
func FormGroupWithError(labelText, inputType, id, name, value, fieldError string, extraAttrs ...g.Node) g.Node {
	attrs := []g.Node{
		Type(inputType),
		ID(id),
		Name(name),
	}
	if value != "" {
		attrs = append(attrs, Value(value))
	}
	if fieldError != "" {
		attrs = append(attrs, g.Attr("aria-invalid", "true"), g.Attr("aria-describedby", id+"-error"))
	}
	attrs = append(attrs, extraAttrs...)
	nodes := []g.Node{
		Label(g.Attr("for", id), g.Text(labelText)),
		Input(attrs...),
	}
	if fieldError != "" {
		nodes = append(nodes, FieldError(id, fieldError))
	}
	return Div(Class("form-group"), g.Group(nodes))
}

// FieldError renders an inline field-level error message.
func FieldError(fieldID, message string) g.Node {
	if message == "" {
		return nil
	}
	return Div(
		ID(fieldID+"-error"),
		Class("field-error"),
		g.Attr("role", "alert"),
		g.Text(message),
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
		g.Attr("hx-target", "#dashboard-main"),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-push-url", "true"),
		g.Text(text),
	)
}

// HtmxForm renders a form with HTMX post attributes targeting the dashboard main area.
// Includes action and method for progressive enhancement (no-JS fallback).
func HtmxForm(postURL, indicatorID string, children ...g.Node) g.Node {
	return FormEl(
		Action(postURL),
		Method("post"),
		g.Attr("hx-post", postURL),
		g.Attr("hx-target", "#dashboard-main"),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-indicator", "#"+indicatorID),
		g.Group(children),
	)
}

// ToggleSwitch renders a CSS toggle switch.
func ToggleSwitch(id, name string, checked bool, label string) g.Node {
	attrs := []g.Node{
		Type("checkbox"),
		ID(id),
		Name(name),
	}
	if checked {
		attrs = append(attrs, g.Attr("checked"))
	}
	return Div(Class("preference-row"),
		Div(Class("preference-info"),
			Span(Class("info-label"), g.Text(label)),
		),
		Div(Class("preference-control"),
			Label(Class("toggle-switch"),
				Input(attrs...),
				Span(Class("toggle-slider")),
			),
		),
	)
}

// SVG icon re-exports from common package for convenient use in components.
var (
	HomeIconSVG                = common.HomeIconSmallSVG
	PersonIconSVG              = common.PersonIconSmallSVG
	PersonLoggedInIconSmallSVG = common.PersonLoggedInIconSmallSVG
	PersonLoggedInIconLargeSVG = common.PersonLoggedInIconLargeSVG
	LockIconSVG                = common.LockIconSmallSVG
	PasskeyIconSVG             = common.PasskeyIconSmallSVG
	LinkIconSVG                = common.LinkIconSmallSVG
	SettingsIconSVG            = common.SettingsIconSmallSVG
	SignOutIconSVG             = common.SignOutIconSmallSVG
	HamburgerIconSVG           = common.HamburgerMenuIconSmallSVG
	ShieldIconLargeSVG         = common.SheildIconLargeSVG
	PasskeyLargeIconSVG        = common.PasskeyIconLargeSVG
	LinkLargeIconSVG           = common.LinkIconLargeSVG
	LockLargeIconSVG           = common.LockIconLargeSVG
	GoogleIconSVG              = common.GoogleIconSmallSVG
	MicrosoftIconSVG           = common.MicrosoftIconSmallSVG
	GitHubIconSVG              = common.GitHubIconSmallSVG
)
