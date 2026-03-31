package components

import (
	"bytes"
	"net/http"

	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	echo "github.com/labstack/echo/v5"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// AutoPostData holds the data needed to render the auto-posting redirect page.
type AutoPostData struct {
	Action     string
	FormParams []models.FormParam
	CSRF       string
}

// AutoPostPage renders a full HTML page that automatically submits a hidden POST form.
// Styled with the same dark theme as the management and OIDC login UIs.
func AutoPostPage(data AutoPostData) g.Node {
	hiddenFields := make([]g.Node, 0, len(data.FormParams)+1)
	for _, fp := range data.FormParams {
		hiddenFields = append(hiddenFields,
			Input(Type("hidden"), Name(fp.Name), Value(fp.Value)),
		)
	}
	hiddenFields = append(hiddenFields,
		Input(Type("hidden"), Name("csrf"), Value(data.CSRF)),
	)

	return c.HTML5(c.HTML5Props{
		Title:    "Redirecting...",
		Language: "en",
		Head: []g.Node{
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1, shrink-to-fit=no")),
			StyleEl(g.Raw(`
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap");
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:"Inter",-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#111827;color:#f9fafb;min-height:100vh;display:flex;align-items:center;justify-content:center}
.redirect-container{text-align:center}
.spinner{width:40px;height:40px;border:3px solid #374151;border-top-color:#10b981;border-radius:50%;animation:spin .8s linear infinite;margin:0 auto 16px}
@keyframes spin{to{transform:rotate(360deg)}}
.redirect-text{font-size:14px;color:#9ca3af;font-weight:500}
`)),
		},
		Body: []g.Node{
			Div(Class("redirect-container"),
				Div(Class("spinner")),
				Div(Class("redirect-text"), g.Text("Redirecting…")),
			),
			FormEl(
				ID("autoForm"),
				Action(data.Action),
				Method("post"),
				g.Attr("style", "display:none"),
				g.Group(hiddenFields),
			),
			Script(g.Raw(`document.getElementById("autoForm").submit();`)),
		},
	})
}

// RenderAutoPost renders the auto-posting redirect page to the echo response.
func RenderAutoPost(c *echo.Context, code int, data AutoPostData) error {
	var buf bytes.Buffer
	if err := AutoPostPage(data).Render(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(code, buf.String())
}
