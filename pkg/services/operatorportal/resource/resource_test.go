package resource

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// fakeStore is a trivial in-memory backing service, standing in for a real
// resource's own service (sqlite-backed IDPs, in-memory clients, ...) so
// this file can exercise Mount's routing/rendering without any of that.
type fakeStore struct {
	rows map[string]Row
}

func newFakeSpec(store *fakeStore) *Spec {
	return &Spec{
		Key:          "widgets",
		Title:        "Widgets",
		SingularName: "Widget",
		IDField:      "id",
		Columns: []Field{
			{Key: "id", Label: "ID"},
			{Key: "name", Label: "Name"},
			{Key: "active", Label: "Active", Kind: Checkbox},
		},
		Fields: []Field{
			{Key: "id", Label: "ID", Immutable: true},
			{Key: "name", Label: "Name"},
			{Key: "active", Label: "Active", Kind: Checkbox},
			{Key: "tags", Label: "Tags", Kind: StringList},
		},
		List: func(ctx context.Context) ([]Row, error) {
			var rows []Row
			for _, r := range store.rows {
				rows = append(rows, r)
			}
			return rows, nil
		},
		Get: func(ctx context.Context, id string) (Row, error) {
			row, ok := store.rows[id]
			if !ok {
				return nil, errNotFoundForTest(id)
			}
			return row, nil
		},
		Create: func(ctx context.Context, row Row) (Row, error) {
			if _, exists := store.rows[row["id"]]; exists {
				return nil, errNotFoundForTest("duplicate " + row["id"])
			}
			store.rows[row["id"]] = row
			return row, nil
		},
		Update: func(ctx context.Context, id string, row Row) (Row, error) {
			if _, ok := store.rows[id]; !ok {
				return nil, errNotFoundForTest(id)
			}
			row["id"] = id
			store.rows[id] = row
			return row, nil
		},
		Delete: func(ctx context.Context, id string) error {
			if _, ok := store.rows[id]; !ok {
				return errNotFoundForTest(id)
			}
			delete(store.rows, id)
			return nil
		},
	}
}

type testError struct{ msg string }

func (e testError) Error() string { return e.msg }
func errNotFoundForTest(id string) error {
	return testError{msg: "not found: " + id}
}

func testShell(c *echo.Context, title string, body g.Node) g.Node {
	return Div(Class("shell"), H1(g.Text(title)), Div(ID("operator-main"), body))
}

func newTestEcho(spec *Spec) *echo.Echo {
	e := echo.New()
	Mount(e, "/operator", spec, testShell, "operator-main")
	return e
}

func doRequest(t *testing.T, e *echo.Echo, method, path string, form url.Values, htmx bool) *httptest.ResponseRecorder {
	t.Helper()
	var body *strings.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	} else {
		body = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestListEmpty(t *testing.T) {
	e := newTestEcho(newFakeSpec(&fakeStore{rows: map[string]Row{}}))

	rec := doRequest(t, e, http.MethodGet, "/operator/widgets", nil, false)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "No Widgets yet.")
	// Direct (non-htmx) nav gets the full shell, not just the fragment.
	require.Contains(t, rec.Body.String(), "<h1>Widgets</h1>")
}

func TestCreateShowsUpInList(t *testing.T) {
	e := newTestEcho(newFakeSpec(&fakeStore{rows: map[string]Row{}}))

	form := url.Values{"id": {"w1"}, "name": {"First Widget"}, "active": {"on"}, "tags": {"a, b, c"}}
	rec := doRequest(t, e, http.MethodPost, "/operator/widgets/new", form, true)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "w1")
	require.Contains(t, body, "First Widget")
	// htmx request gets the bare fragment (its own heading + table), not the
	// shell wrapper this test's testShell would add around it.
	require.NotContains(t, body, `class="shell"`)
	require.Contains(t, body, "op-table")
}

func TestCreateRequiresID(t *testing.T) {
	e := newTestEcho(newFakeSpec(&fakeStore{rows: map[string]Row{}}))

	form := url.Values{"id": {""}, "name": {"No ID"}}
	rec := doRequest(t, e, http.MethodPost, "/operator/widgets/new", form, true)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "id is required")
}

func TestEditFormPrefillsExistingValues(t *testing.T) {
	store := &fakeStore{rows: map[string]Row{
		"w1": {"id": "w1", "name": "Existing", "active": "true", "tags": "x, y"},
	}}
	e := newTestEcho(newFakeSpec(store))

	rec := doRequest(t, e, http.MethodGet, "/operator/widgets/w1/edit", nil, true)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `value="Existing"`)
	require.Contains(t, body, "checked")
	require.Contains(t, body, `value="x, y"`)
}

func TestUpdateRoundTrip(t *testing.T) {
	store := &fakeStore{rows: map[string]Row{
		"w1": {"id": "w1", "name": "Old Name", "active": "false"},
	}}
	e := newTestEcho(newFakeSpec(store))

	form := url.Values{"name": {"New Name"}, "active": {"on"}}
	rec := doRequest(t, e, http.MethodPost, "/operator/widgets/w1/edit", form, true)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "New Name")
	require.Equal(t, "New Name", store.rows["w1"]["name"])
	require.Equal(t, "true", store.rows["w1"]["active"])
}

func TestDeleteRemovesRow(t *testing.T) {
	store := &fakeStore{rows: map[string]Row{
		"w1": {"id": "w1", "name": "Gone Soon"},
	}}
	e := newTestEcho(newFakeSpec(store))

	rec := doRequest(t, e, http.MethodPost, "/operator/widgets/w1/delete", url.Values{}, true)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotContains(t, rec.Body.String(), "Gone Soon")
	_, stillThere := store.rows["w1"]
	require.False(t, stillThere)
}

func TestSplitListAndJoinListRoundTrip(t *testing.T) {
	require.Equal(t, []string{"a", "b", "c"}, SplitList(" a, b ,c "))
	require.Nil(t, SplitList("  "))
	require.Equal(t, "a, b, c", JoinList([]string{"a", "b", "c"}))
}

func TestParseBoolFormatBoolRoundTrip(t *testing.T) {
	require.True(t, ParseBool("true"))
	require.True(t, ParseBool("on"))
	require.False(t, ParseBool("false"))
	require.False(t, ParseBool(""))
	require.Equal(t, "true", FormatBool(true))
	require.Equal(t, "false", FormatBool(false))
}
