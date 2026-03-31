package components

import (
	"bytes"
	"fmt"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// EmailData holds common data for email rendering.
type EmailData struct {
	LocalizeMessage func(string) string
	AccountURL      string
	HeadLinks       []EmailHeadLink
}

// EmailHeadLink represents a <link> element in the email head.
type EmailHeadLink struct {
	HREF string
	REL  string
}

// emailBaseLayout renders the full HTML email wrapper (html_begin + header + base_begin/base_end + footer + html_end).
func emailBaseLayout(data EmailData, content g.Node) g.Node {
	orgName := data.LocalizeMessage("organization_name")
	myAccount := data.LocalizeMessage("my_account")
	admin := data.LocalizeMessage("admin")

	var headLinkNodes []g.Node
	for _, link := range data.HeadLinks {
		headLinkNodes = append(headLinkNodes, Link(Href(link.HREF), Rel(link.REL)))
	}

	return g.Group([]g.Node{
		g.Raw("<!DOCTYPE html>"),
		HTML(
			Lang("en"),
			Head(
				Meta(g.Attr("name", "viewport"), Content("width=device-width, initial-scale=1, maximum-scale=1")),
				g.Group(headLinkNodes),
			),
			Body(
				g.Attr("style", "font-family:Arial,sans-serif;font-size:14px;padding:10px;margin-bottom: 20px;"),
				// Organization banner
				Div(
					g.Attr("style", "margin-bottom:10px; text-align: center; padding: 10px; background-color: #80adbe;"),
					P(
						g.Attr("style", "font-size: 24px; color: white;margin-top: 0;margin-bottom: 0;"),
						g.Text(orgName),
					),
				),
				// Body content
				Div(
					Class("body"),
					content,
				),
				// Footer
				Div(
					Class("footer"),
					g.Attr("style", "color:#444;font-size:12px;margin-top:20px;"),
					Hr(),
					g.Textf("You received this email because you have a %s account. Manage your account at ", orgName),
					A(Href(data.AccountURL), g.Text(myAccount)),
					g.Text(". You can also reach out to "),
					A(Href("mailto:rage@test.com"), g.Text(admin)),
					g.Text(" for no help at all."),
				),
			),
		),
	})
}

// GenericEmailHTML renders the generic email HTML template (emails/generic/index).
func GenericEmailHTML(data EmailData, body string) g.Node {
	return emailBaseLayout(data, Div(
		Class("container"),
		g.Raw(body),
	))
}

// GenericEmailText renders the generic email text template (emails/generic/txt).
func GenericEmailText(body string) string {
	return fmt.Sprintf("Hello\n\n%s", body)
}

// TestEmailRouteRow represents a route row in the test email.
type TestEmailRouteRow struct {
	Verbs string
	Path  string
}

// TestEmailHTML renders the test email HTML template (emails/test/index).
func TestEmailHTML(data EmailData, routes []TestEmailRouteRow) g.Node {
	var rows []g.Node
	for idx, r := range routes {
		rows = append(rows, Tr(
			Th(Class("text-start"), g.Attr("scope", "row"), g.Textf("%d", idx)),
			Td(Class("text-start"), g.Text(r.Verbs)),
			Td(Class("text-start"), g.Text(r.Path)),
		))
	}

	return emailBaseLayout(data, Div(
		Class("container"),
		Div(
			Class("text-center mt-5"),
			g.Attr("role", "alert"),
			H1(g.Text("Perfect Corp.")),
			P(Class("lead"), g.Text("Everything good, nothing bad")),
			Div(
				Class("mt-5 alert alert-success"),
				g.Attr("role", "alert"),
				Table(
					Class("table table-striped"),
					THead(
						Tr(
							Th(Class("text-start"), g.Attr("scope", "col"), g.Text("#")),
							Th(Class("text-start"), g.Attr("scope", "col"), g.Text("Verbs")),
							Th(Class("text-start"), g.Attr("scope", "col"), g.Text("Path")),
						),
					),
					TBody(g.Group(rows)),
				),
			),
		),
	))
}

// TestEmailText renders the test email text template (emails/test/txt).
func TestEmailText(user string) string {
	return fmt.Sprintf("\nHello %s\n", user)
}

// RenderEmailNode renders a gomponents Node to an HTML string.
func RenderEmailNode(node g.Node) (string, error) {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
