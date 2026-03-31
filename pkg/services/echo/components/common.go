package components

import (
	"bytes"
	"net/http"

	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	echo "github.com/labstack/echo/v5"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// RenderContext holds common data needed by all web components.
type RenderContext struct {
	Paths           *wellknown_echo.Paths
	CSRF            string
	Localizer       *i18n.Localizer
	IsAuthenticated bool
	Username        string
}

// L localizes a message key.
func (rc *RenderContext) L(key string) string {
	if rc.Localizer == nil {
		return key
	}
	msg, _ := rc.Localizer.LocalizeMessage(&i18n.Message{ID: key})
	return msg
}

// RenderNode writes a gomponents Node to the echo response.
func RenderNode(c *echo.Context, code int, node g.Node) error {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(code, buf.String())
}

// CsrfInput renders a hidden CSRF token input.
func CsrfInput(csrf string) g.Node {
	return Input(Type("hidden"), Name("csrf"), Value(csrf))
}

// ErrorList renders Bootstrap-styled error list if errors are present.
func ErrorList(errors []string) g.Node {
	if len(errors) == 0 {
		return nil
	}
	items := make([]g.Node, 0, len(errors))
	for _, e := range errors {
		items = append(items, Li(Class("error-list-item"), g.Text(e)))
	}
	return Div(Class("error-container"),
		Ul(Class("error-list"), g.Group(items)),
	)
}
