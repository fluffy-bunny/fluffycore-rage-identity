// Package operatorportal implements the VPN-only internal operator console.
//
// This is deliberately NOT built on the OIDC/RageUser identity model used by
// the public listener: operators are a small, fixed set of accounts defined in
// static config, authenticated against a password hash injected via an
// environment variable (never stored in config, never changeable through this
// portal — rotating an operator's password means rotating the env var and
// redeploying). It runs as a second echo.Echo instance, on its own port, in
// the same process as the public OIDC server, reusing the same shared
// singleton services (IWellknownCookies, IPasswordHasher, IAuditStore, ...)
// off the same root DI container, but with its own session cookie
// (CookieName_OperatorSession) that is never the same cookie as Auth/SSO/
// AccountManagementSession.
//
// IDP and OIDC client management screens (see idp_resource.go,
// client_resource.go) are built on the generic CRUD scaffold in
// pkg/services/operatorportal/resource rather than a JS admin framework —
// see docs/OPERATOR_PORTAL_PLAN.md for why, and for the full design.
package operatorportal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_events "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/events"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	resource "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/operatorportal/resource"
	proto_events_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/events/types"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	echo "github.com/labstack/echo/v5"
	echo_middleware "github.com/labstack/echo/v5/middleware"
	xid "github.com/rs/xid"
	"github.com/rs/zerolog"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// operatorMainTargetID is the id of the element every htmx interaction in
// this portal swaps into — both hand-written routes below and every
// resource.Mount screen (see shellPage/resource.Shell).
const operatorMainTargetID = "operator-main"

// findOperator returns the configured operator matching email (case-insensitive),
// or nil if none match.
func findOperator(config *contracts_config.OperatorConfig, email string) *contracts_config.OperatorAccount {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, op := range config.Operators {
		if strings.ToLower(op.Email) == email {
			return op
		}
	}
	return nil
}

// authenticate checks email/password against the statically-configured operator
// list. The password hash for a matched operator is read from its configured
// env var at call time (not cached), so rotating the env var takes effect on
// the next login attempt without a restart.
//
// Note: if the email doesn't match any configured operator, this returns
// immediately without calling the password hasher, so response time differs
// slightly between "unknown email" and "wrong password" for a known email.
// Given this portal's tiny, fixed, VPN-only audience, that's an accepted
// simplification, not an oversight.
func authenticate(ctx context.Context, config *contracts_config.OperatorConfig, passwordHasher contracts_identity.IPasswordHasher, email, password string) (*contracts_config.OperatorAccount, bool) {
	op := findOperator(config, email)
	if op == nil {
		return nil, false
	}
	hash := os.Getenv(op.PasswordHashEnvVar)
	if hash == "" {
		zerolog.Ctx(ctx).Error().Str("email", op.Email).Str("envVar", op.PasswordHashEnvVar).
			Msg("operator portal: password hash env var is not set for this operator; denying login")
		return nil, false
	}
	err := passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
		HashedPassword: hash,
		Password:       password,
	})
	if err != nil {
		return nil, false
	}
	return op, true
}

// submitAuditEvent records an operator-portal action, tagged so it is always
// distinguishable from customer/self-service audit events — operators aren't
// RageUser subjects, so the "subject" here is the operator's static email.
func submitAuditEvent(ctx context.Context, auditStore contracts_events.IAuditStore, eventType, operatorEmail string, data any, extraAttributes map[string]string) error {
	attributes := map[string]*proto_events_types.CloudEvent_CloudEventAttributeValue{
		"time": {
			Attr: &proto_events_types.CloudEvent_CloudEventAttributeValue_CeTimestamp{CeTimestamp: timestamppb.Now()},
		},
		"datacontenttype": {
			Attr: &proto_events_types.CloudEvent_CloudEventAttributeValue_CeString{CeString: "application/json"},
		},
		"actor_type": {
			Attr: &proto_events_types.CloudEvent_CloudEventAttributeValue_CeString{CeString: "operator"},
		},
	}
	if operatorEmail != "" {
		attributes["subject"] = &proto_events_types.CloudEvent_CloudEventAttributeValue{
			Attr: &proto_events_types.CloudEvent_CloudEventAttributeValue_CeString{CeString: operatorEmail},
		}
	}
	for k, v := range extraAttributes {
		attributes[k] = &proto_events_types.CloudEvent_CloudEventAttributeValue{
			Attr: &proto_events_types.CloudEvent_CloudEventAttributeValue_CeString{CeString: v},
		}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = auditStore.Submit(ctx, &contracts_events.SubmitRequest{
		CloudEvent: &proto_events_types.CloudEvent{
			SpecVersion: "1.0",
			Id:          xid.New().String(),
			Source:      "/rage/operator",
			Type:        eventType,
			Attributes:  attributes,
			Data:        &proto_events_types.CloudEvent_TextData{TextData: string(payload)},
		},
	})
	return err
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

func loginPage(errorMsg string) g.Node {
	children := []g.Node{
		H1(g.Text("Operator Login")),
	}
	if errorMsg != "" {
		children = append(children, P(Class("error"), g.Text(errorMsg)))
	}
	children = append(children,
		FormEl(Method("post"), Action("/operator/login"),
			Div(Label(g.Text("Email")), Input(Type("email"), Name("email"), Required())),
			Div(Label(g.Text("Password")), Input(Type("password"), Name("password"), Required())),
			Button(Type("submit"), g.Text("Sign in")),
		),
	)
	return HTML(
		Head(TitleEl(g.Text("Operator Portal")),
			StyleEl(g.Raw(`body{font-family:sans-serif;max-width:420px;margin:80px auto;}
.error{color:#b00020;} div{margin-bottom:12px;} label{display:block;margin-bottom:4px;}`)),
		),
		Body(g.Group(children)),
	)
}

// navLink renders a sidebar entry. Highlighting is prefix-based (except for
// the dashboard root, which only matches itself) so it stays "active" while
// an operator is on a resource's new/edit sub-page, not just its list page.
func navLink(label, href, activePath string) g.Node {
	active := activePath == href || (href != "/operator" && strings.HasPrefix(activePath, href+"/")) || activePath == href+"/"
	class := "op-nav-link"
	if active {
		class += " op-nav-link-active"
	}
	return A(Href(href), Class(class),
		g.Attr("hx-get", href),
		g.Attr("hx-target", "#"+operatorMainTargetID),
		g.Attr("hx-swap", "innerHTML"),
		g.Attr("hx-push-url", "true"),
		g.Text(label),
	)
}

// shellPage is the resource.Shell implementation for every screen in this
// portal (dashboard + both resource.Mount screens) — nav, branding, and the
// htmx bootstrap script live here exactly once.
func shellPage(c *echo.Context, title string, body g.Node) g.Node {
	email, _ := c.Get("operatorEmail").(string)
	activePath := c.Request().URL.Path

	return HTML(
		Head(
			Meta(g.Attr("charset", "utf-8")),
			Meta(Name("viewport"), g.Attr("content", "width=device-width, initial-scale=1")),
			TitleEl(g.Text("Operator Portal — "+title)),
			Script(Src("https://unpkg.com/htmx.org@2.0.4"),
				g.Attr("integrity", "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"),
				g.Attr("crossorigin", "anonymous")),
			StyleEl(g.Raw(operatorPortalCSS)),
		),
		Body(
			Header(Class("op-topbar"),
				Div(Class("op-brand"), g.Text("Operator Portal")),
				Div(Class("op-topbar-right"),
					Span(g.Text(email)),
					FormEl(Method("post"), Action("/operator/logout"),
						Button(Type("submit"), Class("op-btn"), g.Text("Sign out")),
					),
				),
			),
			Div(Class("op-layout"),
				Nav(Class("op-sidebar"),
					navLink("Dashboard", "/operator", activePath),
					navLink("Identity Providers", "/operator/idps", activePath),
					navLink("OIDC Clients", "/operator/clients", activePath),
				),
				Main(ID(operatorMainTargetID), Class("op-main"), body),
			),
		),
	)
}

func dashboardFragment(email string) g.Node {
	return Div(
		H1(g.Text("Operator Portal")),
		P(g.Text("Signed in as "+email+".")),
		P(g.Text("Manage Identity Providers and OIDC Clients from the nav on the left. User administration and audit-log viewing land here in later phases — see docs/OPERATOR_PORTAL_PLAN.md.")),
	)
}

// requireOperatorSession resolves the current operator session cookie, if any.
// Returns ("", false) when there is no valid session.
func requireOperatorSession(c *echo.Context, wellknownCookies contracts_cookies.IWellknownCookies) (string, bool) {
	resp, err := wellknownCookies.GetOperatorSessionCookie(c)
	if err != nil || resp == nil || resp.OperatorSessionCookie == nil || resp.OperatorSessionCookie.Email == "" {
		return "", false
	}
	return resp.OperatorSessionCookie.Email, true
}

// requireOperatorAuthMiddleware gates every route in the authenticated group
// (dashboard + both resource screens): no valid operator session redirects
// to the login form instead of reaching the handler. On success it stashes
// the operator's email on the echo context under "operatorEmail" for
// shellPage/handlers to read back.
func requireOperatorAuthMiddleware(wellknownCookies contracts_cookies.IWellknownCookies) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			email, ok := requireOperatorSession(c, wellknownCookies)
			if !ok {
				return c.Redirect(http.StatusFound, "/operator/login")
			}
			c.Set("operatorEmail", email)
			return next(c)
		}
	}
}

// NewEcho builds the standalone operator-portal echo.Echo instance. The
// caller is responsible for starting it (e.Start(":<port>")) on its own port,
// separate from the public listener, and for making sure that port is only
// reachable via VPN at the infra layer (see docs/OPERATOR_PORTAL_PLAN.md).
func NewEcho(rootContainer di.Container, config *contracts_config.Config) *echo.Echo {
	wellknownCookies := di.Get[contracts_cookies.IWellknownCookies](rootContainer)
	passwordHasher := di.Get[contracts_identity.IPasswordHasher](rootContainer)
	auditStore := di.Get[contracts_events.IAuditStore](rootContainer)
	idpStore := di.Get[proto_oidc_idp.IFluffyCoreExampleIDPStoreServer](rootContainer)
	clientService := di.Get[proto_oidc_client.IFluffyCoreClientServiceServer](rootContainer)
	operatorConfig := config.OperatorConfig

	e := echo.New()
	e.Use(echo_middleware.Recover())

	e.GET("/operator/login", func(c *echo.Context) error {
		if _, ok := requireOperatorSession(c, wellknownCookies); ok {
			return c.Redirect(http.StatusFound, "/operator")
		}
		return renderNode(c, http.StatusOK, loginPage(""))
	})

	e.POST("/operator/login", func(c *echo.Context) error {
		ctx := c.Request().Context()
		email := c.FormValue("email")
		password := c.FormValue("password")
		op, ok := authenticate(ctx, operatorConfig, passwordHasher, email, password)
		if !ok {
			return renderNode(c, http.StatusUnauthorized, loginPage("Invalid email or password."))
		}
		if err := wellknownCookies.SetOperatorSessionCookie(c, &contracts_cookies.SetOperatorSessionCookieRequest{
			OperatorSessionCookie: &contracts_cookies.OperatorSessionCookie{Email: op.Email},
		}); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("operator portal: SetOperatorSessionCookie")
			return renderNode(c, http.StatusInternalServerError, loginPage("Something went wrong. Please try again."))
		}
		if err := submitAuditEvent(ctx, auditStore, "com.fluffybunny.identity.operator.login", op.Email, map[string]string{}, map[string]string{"mutation": "operator_login"}); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("operator portal: submitAuditEvent(login)")
		}
		return c.Redirect(http.StatusFound, "/operator")
	})

	e.POST("/operator/logout", func(c *echo.Context) error {
		ctx := c.Request().Context()
		email, _ := requireOperatorSession(c, wellknownCookies)
		wellknownCookies.DeleteOperatorSessionCookie(c)
		if email != "" {
			if err := submitAuditEvent(ctx, auditStore, "com.fluffybunny.identity.operator.logout", email, map[string]string{}, map[string]string{"mutation": "operator_logout"}); err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("operator portal: submitAuditEvent(logout)")
			}
		}
		return c.Redirect(http.StatusFound, "/operator/login")
	})

	authed := e.Group("/operator", requireOperatorAuthMiddleware(wellknownCookies))

	authed.GET("", func(c *echo.Context) error {
		email, _ := c.Get("operatorEmail").(string)
		body := dashboardFragment(email)
		if isHTMXRequest(c) {
			return renderNode(c, http.StatusOK, body)
		}
		return renderNode(c, http.StatusOK, shellPage(c, "Dashboard", body))
	})

	// IDP and OIDC client management: the generic CRUD scaffold in
	// pkg/services/operatorportal/resource, given field metadata + the
	// backing service's CRUD calls — see idp_resource.go/client_resource.go.
	// Auditing operator mutations to these resources is a follow-up (the
	// scaffold's Create/Update/Delete callbacks are exactly where a
	// SubmitAuditEvent call would go, mirroring the public listener's
	// BaseHandler.SubmitAuditEvent convention) — not yet wired.
	resource.Mount(authed, "", newIDPResourceSpec(idpStore), shellPage, operatorMainTargetID)
	resource.Mount(authed, "", newClientResourceSpec(clientService), shellPage, operatorMainTargetID)

	return e
}
