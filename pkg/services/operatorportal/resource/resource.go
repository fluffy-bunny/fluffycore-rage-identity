// Package resource is a small, self-rolled CRUD scaffold for the operator
// portal: hand it a resource's field metadata plus four callbacks (List/Get/
// Create/Update/Delete) and it produces a full list + create + edit + delete
// htmx screen. It's the thing you reach for instead of pulling a JS admin
// framework (react-admin and friends) into a Go+htmx app — a resource
// definition here is the equivalent of a react-admin <Resource>, minus the
// SPA, the JS bundle, and the second data layer.
//
// This package knows nothing about IDPs or OIDC clients specifically — see
// idp_resource.go and client_resource.go in the parent package for the two
// concrete resources wired up today. Adding a third resource (users, once
// the operator portal grows into user administration) means writing a new
// adapter file in that same shape, not touching this one.
package resource

import (
	"bytes"
	"context"
	"net/http"
	"sort"
	"strings"

	echo "github.com/labstack/echo/v5"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Kind is the input widget a Field renders as, on both the list table
// (informationally) and the create/edit form (as the actual input).
type Kind int

const (
	Text Kind = iota
	TextArea
	Password
	Checkbox
	// StringList is a comma-separated list, rendered as a single text input
	// and parsed/joined to and from []string at the Row boundary.
	StringList
)

// Field describes one attribute of a Row: its form widget, its label, and
// whether it's immutable once created (e.g. the id field itself).
type Field struct {
	Key       string // Row map key and <input name="...">
	Label     string
	Kind      Kind
	Immutable bool // read-only on the edit form; still shown, never submitted for edits
	Help      string
}

// Row is a flat, string-keyed view of one record — the boundary between this
// package's generic rendering and a resource's own backing service, which
// owns the Row<->proto/struct mapping (see idp_resource.go/client_resource.go).
// Checkbox fields round-trip as "true"/"false"; StringList fields round-trip
// as a comma-separated string.
type Row map[string]string

// Spec wires the scaffold to one resource.
type Spec struct {
	Key          string  // URL path segment this resource mounts under, e.g. "idps"
	Title        string  // e.g. "Identity Providers"
	SingularName string  // e.g. "IDP" — used in button/heading/confirm text
	IDField      string  // Key of the Field that identifies a Row; always Immutable
	Columns      []Field // shown, in order, in the list table
	Fields       []Field // shown, in order, in the create/edit form

	List   func(ctx context.Context) ([]Row, error)
	Get    func(ctx context.Context, id string) (Row, error)
	Create func(ctx context.Context, row Row) (Row, error)
	Update func(ctx context.Context, id string, row Row) (Row, error)
	Delete func(ctx context.Context, id string) error
}

// Router is the subset of echo's route-registration API the scaffold needs.
// Both *echo.Echo and *echo.Group satisfy it, so Mount can be pointed at an
// authenticated sub-group without this package needing to know how (or
// whether) the caller gates access.
type Router interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
}

// Shell wraps one fragment of body content in a full HTML page — nav,
// branding, htmx bootstrap script, all the caller's job, none of this
// package's. Mount calls it only for non-htmx requests (a direct browser
// nav or a page refresh); htmx-driven interactions get the bare fragment
// back, swapped into whatever element mainTargetID names.
type Shell func(c *echo.Context, title string, body g.Node) g.Node

// Mount registers list/new/edit/delete routes for spec under
// basePath+"/"+spec.Key. mainTargetID must match the id of the element the
// caller's Shell uses as its htmx swap target, so fragment responses land in
// the right place and action buttons/forms can hx-target it directly.
func Mount(r Router, basePath string, spec *Spec, shell Shell, mainTargetID string) {
	base := basePath + "/" + spec.Key

	respond := func(c *echo.Context, code int, title string, body g.Node) error {
		if isHTMXRequest(c) {
			return renderNode(c, code, body)
		}
		return renderNode(c, code, shell(c, title, body))
	}

	// List
	r.GET(base, func(c *echo.Context) error {
		ctx := c.Request().Context()
		rows, err := spec.List(ctx)
		if err != nil {
			return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, nil, err))
		}
		return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, rows, nil))
	})

	// New
	r.GET(base+"/new", func(c *echo.Context) error {
		return respond(c, http.StatusOK, "New "+spec.SingularName, formPage(spec, base, mainTargetID, Row{}, nil, nil))
	})
	r.POST(base+"/new", func(c *echo.Context) error {
		ctx := c.Request().Context()
		row := parseForm(c, spec.Fields)
		if err := requireID(spec, row); err != nil {
			return respond(c, http.StatusOK, "New "+spec.SingularName, formPage(spec, base, mainTargetID, row, nil, err))
		}
		if _, err := spec.Create(ctx, row); err != nil {
			return respond(c, http.StatusOK, "New "+spec.SingularName, formPage(spec, base, mainTargetID, row, nil, err))
		}
		rows, err := spec.List(ctx)
		return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, rows, err))
	})

	// Edit
	r.GET(base+"/:id/edit", func(c *echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		row, err := spec.Get(ctx, id)
		if err != nil {
			return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, nil, err))
		}
		return respond(c, http.StatusOK, "Edit "+spec.SingularName, formPage(spec, base, mainTargetID, row, &id, nil))
	})
	r.POST(base+"/:id/edit", func(c *echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		row := parseForm(c, spec.Fields)
		if _, err := spec.Update(ctx, id, row); err != nil {
			return respond(c, http.StatusOK, "Edit "+spec.SingularName, formPage(spec, base, mainTargetID, row, &id, err))
		}
		rows, err := spec.List(ctx)
		return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, rows, err))
	})

	// Delete
	r.POST(base+"/:id/delete", func(c *echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		err := spec.Delete(ctx, id)
		rows, listErr := spec.List(ctx)
		if err == nil {
			err = listErr
		}
		return respond(c, http.StatusOK, spec.Title, listPage(spec, base, mainTargetID, rows, err))
	})
}

func isHTMXRequest(c *echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func renderNode(c *echo.Context, code int, node g.Node) error {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(code, buf.String())
}

// requireID checks that the form submitted a non-empty value for the id
// field — the one validation rule generic enough to belong in the scaffold
// itself; anything resource-specific belongs in that resource's Create/
// Update callback, whose error is rendered the same way.
func requireID(spec *Spec, row Row) error {
	if strings.TrimSpace(row[spec.IDField]) == "" {
		return errRequiredID{field: spec.IDField}
	}
	return nil
}

type errRequiredID struct{ field string }

func (e errRequiredID) Error() string { return e.field + " is required" }

// parseForm reads every non-immutable field out of the posted form, plus the
// id field itself (so edits can still see it even though it's immutable/not
// re-submitted as editable), producing a Row ready for Create/Update.
func parseForm(c *echo.Context, fields []Field) Row {
	row := Row{}
	for _, f := range fields {
		switch f.Kind {
		case Checkbox:
			row[f.Key] = boolString(c.FormValue(f.Key) == "on" || c.FormValue(f.Key) == "true")
		default:
			row[f.Key] = c.FormValue(f.Key)
		}
	}
	return row
}

func boolString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func isTrue(v string) bool {
	return v == "true" || v == "on"
}

// ParseBool mirrors the Row boolean convention (Checkbox fields round-trip
// as "true"/"false"/"on") for resource adapters converting a Row back to
// their own backing type.
func ParseBool(v string) bool {
	return isTrue(v)
}

// FormatBool is ParseBool's inverse, for resource adapters converting their
// backing type into a Row.
func FormatBool(b bool) string {
	return boolString(b)
}

// splitList/joinList are exported so resource adapters (idp_resource.go,
// client_resource.go) share the exact same comma-separated-list convention
// this package's StringList form field uses.
func SplitList(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	sort.Strings(result)
	return result
}

func JoinList(v []string) string {
	return strings.Join(v, ", ")
}

// --- rendering ---

func listPage(spec *Spec, base, mainTargetID string, rows []Row, err error) g.Node {
	children := []g.Node{
		Div(Class("op-page-header"),
			H1(g.Text(spec.Title)),
			hxLink("+ New "+spec.SingularName, base+"/new", mainTargetID, "op-btn op-btn-primary"),
		),
		errorBanner(err),
	}
	if rows != nil || err == nil {
		children = append(children, table(spec, base, mainTargetID, rows))
	}
	return Div(g.Group(children))
}

func table(spec *Spec, base, mainTargetID string, rows []Row) g.Node {
	headerCells := make([]g.Node, 0, len(spec.Columns)+1)
	for _, col := range spec.Columns {
		headerCells = append(headerCells, Th(g.Text(col.Label)))
	}
	headerCells = append(headerCells, Th(g.Text("")))

	bodyRows := make([]g.Node, 0, len(rows))
	for _, row := range rows {
		id := row[spec.IDField]
		cells := make([]g.Node, 0, len(spec.Columns)+1)
		for _, col := range spec.Columns {
			val := row[col.Key]
			if col.Kind == Checkbox {
				val = map[bool]string{true: "yes", false: "no"}[isTrue(val)]
			}
			cells = append(cells, Td(g.Text(val)))
		}
		cells = append(cells, Td(Class("op-actions"),
			hxLink("Edit", base+"/"+id+"/edit", mainTargetID, "op-btn"),
			hxDeleteForm(base+"/"+id+"/delete", mainTargetID, spec.SingularName, id),
		))
		bodyRows = append(bodyRows, Tr(g.Group(cells)))
	}
	if len(rows) == 0 {
		bodyRows = append(bodyRows, Tr(Td(g.Attr("colspan", "99"), Class("op-empty"), g.Text("No "+spec.Title+" yet."))))
	}

	return Table(Class("op-table"),
		THead(Tr(g.Group(headerCells))),
		TBody(g.Group(bodyRows)),
	)
}

func formPage(spec *Spec, base, mainTargetID string, row Row, editingID *string, err error) g.Node {
	action := base + "/new"
	isEdit := editingID != nil
	heading := "New " + spec.SingularName
	if isEdit {
		action = base + "/" + *editingID + "/edit"
		heading = "Edit " + spec.SingularName
	}

	formChildren := []g.Node{
		errorBanner(err),
	}
	for _, f := range spec.Fields {
		formChildren = append(formChildren, formField(f, row[f.Key], isEdit))
	}
	formChildren = append(formChildren,
		Div(Class("op-actions"),
			Button(Type("submit"), Class("op-btn op-btn-primary"), g.Text("Save")),
			hxLink("Cancel", base, mainTargetID, "op-btn"),
		),
	)

	return Div(
		H1(g.Text(heading)),
		FormEl(Action(action), Method("post"),
			g.Attr("hx-post", action),
			g.Attr("hx-target", "#"+mainTargetID),
			g.Attr("hx-swap", "innerHTML"),
			g.Group(formChildren),
		),
	)
}

func formField(f Field, value string, isEdit bool) g.Node {
	label := Label(g.Attr("for", f.Key), g.Text(f.Label))
	help := g.Node(nil)
	if f.Help != "" {
		help = P(Class("op-help"), g.Text(f.Help))
	}

	// Immutable fields (the id field, on edit) are shown for context but not
	// posted as editable — rendered readonly, name intentionally omitted so a
	// tampered client-side value can't smuggle an id change through Update.
	readOnly := f.Immutable && isEdit

	switch f.Kind {
	case TextArea:
		attrs := []g.Node{ID(f.Key), Name(f.Key), Rows("3"), g.Text(value)}
		if readOnly {
			attrs = append(attrs, g.Attr("readonly"))
		}
		return Div(Class("op-field"), label, Textarea(attrs...), help)
	case Checkbox:
		attrs := []g.Node{Type("checkbox"), ID(f.Key), Name(f.Key)}
		if isTrue(value) {
			attrs = append(attrs, Checked())
		}
		if readOnly {
			attrs = append(attrs, Disabled())
		}
		return Div(Class("op-field op-field-checkbox"), Input(attrs...), label, help)
	default:
		inputType := "text"
		if f.Kind == Password {
			inputType = "password"
		}
		attrs := []g.Node{Type(inputType), ID(f.Key), Name(f.Key), Value(value)}
		if readOnly {
			attrs = append(attrs, g.Attr("readonly"))
		}
		return Div(Class("op-field"), label, Input(attrs...), help)
	}
}

func errorBanner(err error) g.Node {
	if err == nil {
		return nil
	}
	return Div(Class("op-error"), g.Attr("role", "alert"), g.Text(err.Error()))
}

// hxLink renders a same-page navigation link that swaps mainTargetID via
// htmx, with a plain href so it still works with JS disabled (a full page
// nav to the same route, handled by the same GET handler — see Mount).
func hxLink(text, href, mainTargetID, class string) g.Node {
	return A(Href(href), Class(class),
		g.Attr("hx-get", href),
		g.Attr("hx-target", "#"+mainTargetID),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-push-url", "true"),
		g.Text(text),
	)
}

// hxDeleteForm renders a tiny inline form (not a link) so delete is a POST,
// with hx-confirm giving the operator an "are you sure" prompt.
func hxDeleteForm(action, mainTargetID, singularName, id string) g.Node {
	return FormEl(Action(action), Method("post"), Class("op-inline-form"),
		g.Attr("hx-post", action),
		g.Attr("hx-target", "#"+mainTargetID),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-confirm", "Delete "+singularName+" \""+id+"\"? This cannot be undone."),
		Button(Type("submit"), Class("op-btn op-btn-danger"), g.Text("Delete")),
	)
}
