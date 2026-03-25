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
}

// ErrorPartial renders the error page.
func ErrorPartial(data ErrorData) g.Node {
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
			ButtonGroup(
				PrimaryButton(data.L("start_over"), ""),
			),
		),
	})
}
