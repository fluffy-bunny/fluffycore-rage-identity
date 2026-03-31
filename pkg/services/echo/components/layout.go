package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// PageShell renders a full HTML5 page with bootstrap header/footer wrapping body content.
func PageShell(bodyAttrs []g.Node, children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    "Bare - Start Bootstrap Template",
		Language: "en",
		Head: []g.Node{
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1, shrink-to-fit=no")),
			Meta(Name("description")),
			Meta(Name("author")),
			Link(g.Attr("rel", "icon"), Type("image/x-icon"), Href("static/assets/favicon.ico")),
			Link(Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css"),
				g.Attr("rel", "stylesheet"),
				g.Attr("integrity", "sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH"),
				g.Attr("crossorigin", "anonymous")),
			Link(Href("/static/css/error.css"), g.Attr("rel", "stylesheet"), g.Attr("crossorigin", "anonymous")),
		},
		Body: []g.Node{
			Body(append(bodyAttrs, g.Group(children))...),
			footerScripts(),
		},
	})
}

// PageShellWithNavbar renders a full HTML5 page with navbar, bootstrap header/footer.
func PageShellWithNavbar(rc *RenderContext, children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    "Bare - Start Bootstrap Template",
		Language: "en",
		Head: []g.Node{
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1, shrink-to-fit=no")),
			Meta(Name("description")),
			Meta(Name("author")),
			Link(g.Attr("rel", "icon"), Type("image/x-icon"), Href("static/assets/favicon.ico")),
			Link(Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css"),
				g.Attr("rel", "stylesheet"),
				g.Attr("integrity", "sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH"),
				g.Attr("crossorigin", "anonymous")),
			Link(Href("/static/css/error.css"), g.Attr("rel", "stylesheet"), g.Attr("crossorigin", "anonymous")),
		},
		Body: []g.Node{
			Navbar(rc),
			Body(g.Group(children)),
			footerScripts(),
		},
	})
}

// Navbar renders the Bootstrap responsive navbar.
func Navbar(rc *RenderContext) g.Node {
	var authItem g.Node
	if rc.IsAuthenticated {
		authItem = Li(
			A(Class("dropdown-item"), Href(rc.Paths.Logout), g.Text("Logout")),
		)
	} else {
		authItem = Li(
			A(Class("dropdown-item"), Href(rc.Paths.Login), g.Text("Login")),
		)
	}
	return Nav(Class("navbar navbar-expand-lg navbar-dark bg-dark"),
		Div(Class("container"),
			A(Class("navbar-brand"), Href(rc.Paths.Home), g.Text("Echo Starter")),
			Button(Class("navbar-toggler"), Type("button"),
				g.Attr("data-bs-toggle", "collapse"),
				g.Attr("data-bs-target", "#navbarSupportedContent"),
				g.Attr("aria-controls", "navbarSupportedContent"),
				g.Attr("aria-expanded", "false"),
				g.Attr("aria-label", "Toggle navigation"),
				Span(Class("navbar-toggler-icon")),
			),
			Div(Class("collapse navbar-collapse"), ID("navbarSupportedContent"),
				Ul(Class("navbar-nav ms-auto mb-2 mb-lg-0"),
					Li(Class("nav-item"),
						A(Class("nav-link active"), g.Attr("aria-current", "page"), Href(rc.Paths.About), g.Text("About")),
					),
					Li(Class("nav-item"),
						A(Class("nav-link active"), g.Attr("aria-current", "page"), Href(rc.Paths.Login), g.Text("Login")),
					),
					Li(Class("nav-item dropdown"),
						A(Class("nav-link dropdown-toggle"), ID("navbarDropdown"), Href("#"),
							g.Attr("role", "button"),
							g.Attr("data-bs-toggle", "dropdown"),
							g.Attr("aria-expanded", "false"),
							g.Text(rc.Username),
						),
						Ul(Class("dropdown-menu dropdown-menu-end"), g.Attr("aria-labelledby", "navbarDropdown"),
							authItem,
							Li(
								A(Class("dropdown-item"), Href(rc.Paths.Profile), g.Text("Profile")),
							),
						),
					),
				),
			),
		),
	)
}

func footerScripts() g.Node {
	return g.Group([]g.Node{
		Script(Src("https://ajax.googleapis.com/ajax/libs/jquery/3.7.1/jquery.min.js")),
		Script(Src("https://unpkg.com/@popperjs/core@2")),
		Script(Src("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"),
			g.Attr("integrity", "sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"),
			g.Attr("crossorigin", "anonymous")),
		Script(Src("/static/js/common.js")),
	})
}
