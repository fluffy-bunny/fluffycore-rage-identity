package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ErrorData holds data for the error page partial.
type ErrorData struct {
	*RenderContext
	ErrorCode    string
	ErrorMessage string
	// ReturnURL is an optional URL to return the user to the client app.
	// When set, shows a "Return to App" link instead of "Start Over".
	ReturnURL string
}

// ErrorPartial renders the error page.
func ErrorPartial(data ErrorData) g.Node {
	var actionButton g.Node
	if data.ReturnURL != "" {
		// Fatal error — OIDC flow is unrecoverable. Send user back to client app.
		actionButton = ButtonGroup(
			A(Href(data.ReturnURL), Class("btn-primary"),
				g.Text(data.L("return_to_app")),
			),
		)
	} else {
		// Recoverable error — restart the OIDC login flow.
		actionButton = ButtonGroup(
			Button(Type("button"), Class("btn-primary"),
				g.Attr("hx-post", data.Paths.HTMXStartOver),
				g.Attr("hx-target", "#main-content"),
				g.Attr("hx-swap", "innerHTML"),
				g.Text(data.L("start_over")),
			),
		)
	}

	return g.Group([]g.Node{
		Div(Class("error-section"),
			H2(Class("error-heading"),
				g.Text(data.L("error")),
			),
			g.If(data.ErrorCode != "",
				P(Class("error-code"),
					g.Text(data.ErrorCode),
				),
			),
			P(Class("error-message"),
				g.Text(data.ErrorMessage),
			),
			actionButton,
		),
	})
}
