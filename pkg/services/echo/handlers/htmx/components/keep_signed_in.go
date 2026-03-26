package components

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// KeepSignedInData holds data for the keep-signed-in partial.
type KeepSignedInData struct {
	*RenderContext
}

// KeepSignedInPartial renders the keep-signed-in preference step.
func KeepSignedInPartial(data KeepSignedInData) g.Node {
	return g.Group([]g.Node{
		H2(g.Text(data.L("keep_signed_in"))),
		P(g.Text("Choose whether to stay signed in on this device")),
		HtmxForm(data.Paths.HTMXKeepSignedIn, "ksi-indicator",
			CsrfInput(data.CSRF),
			Div(Class("form-group"), g.Attr("style", "display:flex;align-items:center;gap:10px;"),
				Input(Type("checkbox"), ID("keepSignedIn"), Name("keepSignedIn"), Value("true"),
					g.Attr("style", "width:20px;height:20px;")),
				Label(g.Attr("for", "keepSignedIn"), g.Attr("style", "margin-bottom:0;"),
					g.Text(data.L("keep_me_signed_in"))),
			),
			ButtonGroup(
				PrimaryButton(data.L("continue"), "ksi-indicator"),
			),
		),
	})
}
