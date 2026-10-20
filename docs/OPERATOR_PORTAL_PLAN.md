# Operator Portal — Design & Implementation Plan

**Status: planned.** Unlike the SAML doc, this isn't contingent on a client request — it closes a real gap (no admin surface exists at all today; IDPs and OIDC clients are only ever edited by hand-editing JSON files and restarting the process).

## Goals

- A dedicated internal console for operators to manage IDPs, OIDC clients, and user administration tasks (e.g. tagging a user with `managed_idp_slug` per the customer self-service design).
- Reachable **only via VPN** — not merely a gated path behind the same public listener as `/management/*` or `/oidc-login/*`.
- A small, fixed set of named operator accounts, defined in static config, with passwords that are **not** changeable through any in-app flow.

## Non-goals for v1

- Fine-grained per-operator RBAC. A short, named list of operators (e.g. `bob@gmail.com`) is a flat superuser set for v1 — no role hierarchy yet. The config shape below leaves room to add scoped permissions later without changing the identity model.
- SSO/MFA for operator login. Called out as a strong v2 candidate given how privileged this surface is, but not required to ship v1 — VPN-only network access plus a config-injected credential is the v1 bar.

## Network isolation

"Behind an auth check on the same port" isn't sufficient — you were explicit that this can't be mixed in with the account portal. Concretely:

- **Same binary, second listener.** Not a separate `cmd/` binary — one deployable process (`cmd/server`), a second `echo.Echo` instance bound to its own port, started the same way the existing OIDC server already runs: `pkg/runtime/startup.go`'s `OnPreServerStartup` already launches the oidcserver asynchronously via `fluffycore_async.ExecuteWithPromiseAsync(...)` inside the same process — the operator listener follows that exact precedent, just as a second async future. It gets its own route table, its own middleware stack, and its own cookie/session store — none of that is shared with the public listener just because they're in one process. Simpler to build and ship (one image, one deploy), at the cost of one thing worth being explicit about: an unrecovered panic in the operator listener's goroutine could, if not guarded, affect the whole process. Mitigate with the same panic-recovery middleware pattern already used elsewhere in this codebase (e.g. `grpc_recovery.UnaryServerInterceptor`), applied to the operator echo instance too.
- **Its own dedicated port.** The existing `docker-compose.yml` already reserves a block of sequential ports for this service (`30051-30053`, `9044`, `9045` — the last one currently unbound by any app code I could find, likely reserved). Pick the next free port in that convention (e.g. `9046`) for the operator listener. One pod/container can still back two separate Kubernetes `Service` objects on its two ports — the public one gets an `Ingress`, the operator one doesn't (or gets an internal-only one) — so this doesn't weaken the VPN-only story at all, it just collapses the deployment into one unit instead of two.
- **The actual VPN boundary is enforced at the infra layer, not in this repo.** This app's job is to expose a cleanly separable listener; whether it's reachable is a Traefik-entrypoint / security-group / VPC-routing decision made outside the Go code (this repo already uses Traefik per `docs/TRAEFIK.md`) — e.g. an internal-only Traefik entrypoint with no public router attached, or a Kubernetes `Service` with no `Ingress`, reachable only from VPN-routed client IPs. Flagging this explicitly so it doesn't get silently assumed to be "handled by the app."
- **No shared session/cookie surface.** The operator portal uses its own cookie name and session store, entirely distinct from `Auth`/`SSO`/`AccountManagementSession`. Even if the two listeners were ever accidentally reachable from the same origin, an operator session could never be confused with, or escalate into, a customer session, and vice versa.

## Operator identity & authentication

Config-driven, matching the pattern you described:

```json
"operatorConfig": {
  "operators": [
    { "email": "bob@gmail.com", "passwordHashEnvVar": "OPERATOR_BOB_PASSWORD_HASH" }
  ]
}
```

- Loaded the same way `idps.json`/`oidc-clients.json` already are — `fluffycore_utils.ReplaceEnv` substitutes `${OPERATOR_BOB_PASSWORD_HASH}`-style placeholders from the environment before parsing, so the actual hash never lives in a committed file.
- The operator portal has its **own login form**, separate from the OIDC/account login. On submit, it looks up the email in the static operator list and verifies the password against the configured hash using the already-existing `IPasswordHasher.VerifyPassword` — no new hashing primitive needed, same bcrypt/argon2id service used everywhere else in this codebase.
- **There is no "change password" capability anywhere in the operator portal**, by design. Rotating an operator's password means rotating the env var/secret and redeploying — exactly the pattern you described, and it means a compromised operator session (or a bug in the portal) can never be used to lock out or take over the operator account itself.
- Session cookie: short-lived, `Secure`, `HttpOnly`, `SameSite=Strict`, scoped only to the operator listener's own port — a distinct cookie name/session store from the public listener's, enforced in app code regardless of both living in one process.

## Authorization inside the portal

v1: membership in the static operator list *is* the authorization check — every configured operator has full access. No per-operator scoping yet. If/when you need to split "can edit IDPs" from "can view audit logs," add a `Scopes []string` field per operator entry in that same config block rather than redesigning anything.

## Service reuse

Since this is now one process with one root DI container (not a separate binary), the operator portal's handlers reuse the exact same backing services every `/management/*` handler already uses — no parallel service layer to build or keep in sync:

- `RageUserService()` / `IdpServiceServer()` / `ClientServiceServer()` / `AuditStore()` / `PasswordHasher()` — same singletons, same data, resolved off the same root container. The operator listener just needs the second `echo.Echo` instance to run through the same handler-registration/scoped-container framework the public listener already uses (`contracts_handler.AddScopedIHandleWithMetadata` + the per-request scoped-container middleware every existing handler already relies on) — that's what makes `BaseHandler`-style DI resolution work the same way on both listeners, not something to reimplement per listener.
- New operator-only handlers should embed `*services_echo_handlers_base.BaseHandler` exactly like every existing handler in this codebase, rather than inventing a parallel base type.

**One thing that should *not* be reused as-is: the identity/claims layer.** `ClaimsPrincipal()`/`IClaimsPrincipal` represents an OIDC-authenticated end user (claims minted from a JWT via the normal login flow) — that model doesn't fit an operator, who isn't a `RageUser` and isn't authenticating via OIDC at all, just a static config entry checked against an env-injected password hash. Operator handlers need their own, much simpler "is this session's cookie mapped to a configured operator" check rather than going through `ClaimsPrincipal`. The underlying cookie-signing/session infrastructure (`IWellknownCookies`) can still be reused — it's a generic encrypted-cookie service — just with a new, distinct cookie name (e.g. `CookieName_OperatorSession`) so it's never the same cookie as `Auth`/`SSO`/`AccountManagementSession`, consistent with the "no shared session surface" requirement above.

## Capabilities

This is where the portal earns its keep — it needs real backing services, not just a UI, closing gaps already identified this session:

- **IDP management (full CRUD). Done.** `CreateIDP`/`UpdateIDP`/`DeleteIDP` were added to the proto alongside `ListIDP` (see `ExampleIDPStore` in `proto/oidc/idp/idp.proto`), backed by the persistent, sqlite-embedded store in `pkg/services/idp/sqlite` — shared with the customer self-service portal, not built twice. The operator screen touches everything: `claimed_domains`, `enabled`, `hidden`, `auto_create`, `multi_factor_required`, plus the OIDC protocol fields; the customer self-service portal (separate, not yet built) is expected to restrict itself to a customer's own `OIDCProtocol` fields only.
- **OIDC client management. Done, with a correction to this plan's original premise.** The proto looked CRUD-complete but `pkg/services/client/inmemory/inmemory-client-service.go` only ever implemented `ListClient`/`GetClient`/`ValidateClientSecret` — `CreateClient`/`UpdateClient`/`DeleteClient` fell through to `UnimplementedClientServiceServer` and errored on every call. Those three are now implemented (in-memory only — see the persistence note below), and the operator screen manages a working subset of `Client`'s ~30 fields (name/description/namespace/enabled/require_client_secret/allow_offline_access/allowed_scopes/allowed_grant_types/allowed_redirect_uris), not all of them. `AddNewClientSecret` and claims updates are still unimplemented — secret rotation and per-client claims are a follow-up, not part of this screen. The never-mounted grpc-gateway wrapper (`proto/oidc/client/client.pb.gw.go`) stays un-resurrected, as originally planned.
- **Neither screen is a hand-built htmx handler per resource.** Both are one call each into `pkg/services/operatorportal/resource` — a small, generic CRUD scaffold (field metadata in, a full list/create/edit/delete htmx screen out) built for this portal instead of reaching for a JS admin framework. A third resource (users, below) means a third adapter file in that same shape.
- **User administration.** View users, edit metadata (this is exactly where an operator would set `managed_idp_slug` per the customer self-service design), disable/enable accounts, force logout/session revocation. Not yet built — a natural next `resource.Spec`, backed by `RageUserService`.
- **Audit log viewing.** Absorb and retire the existing `/support/*` "audits" handler. That handler currently sits on the **public** listener, gated only by an `isSupportAdmin` role/group claim — the same trust tier as everything moving into the VPN-only operator portal. Running two separate "internal staff" surfaces at two different trust levels is worse than consolidating; migrate audit viewing into the operator portal and retire the public `/support/*` route once it's ported. Not yet built.
- **Operator-attributed audit events for IDP/client mutations.** Not yet wired — `resource.Spec.Create`/`Update`/`Delete` are exactly where a `submitAuditEvent` call belongs (see "Audit logging for operator actions" below), but the two resources built so far don't call it yet.

## Persistence dependency

**Resolved for IDPs**, not yet for clients. IDPs are backed by the persistent, sqlite-embedded store in `pkg/services/idp/sqlite`, shared with the customer self-service portal. OIDC clients are still in-memory only (`pkg/services/client/inmemory`) — creates/edits/deletes made through the operator portal are lost on restart. The open question from before still applies to a future client persistence pass: confirm single-instance vs multi-replica deployment before choosing SQLite-embedded (fine for single-instance, same as IDPs) vs a real shared DB (needed for multi-replica).

## Audit logging for operator actions

Reuse `BaseHandler.SubmitAuditEvent` — same convention as every other mutation in this codebase. Operators aren't `RageUser` subjects, so the audit event's `subject`/actor should be the operator's static config email, with an explicit `actor_type: "operator"` attribute, so operator-performed mutations are unambiguously distinguishable from customer or self-service ones in the audit trail, permanently.

## Phased implementation order

1. ~~Confirm single-instance vs multi-replica deployment~~ — resolved for IDPs (single-instance, sqlite-embedded); still open for a future client-persistence pass (see "Persistence dependency" above).
2. **Done.** Persistent, writable IDP store + `CreateIDP`/`UpdateIDP`/`DeleteIDP` RPCs (`pkg/services/idp/sqlite`, shared with the customer self-service plan).
3. **Done.** The second `echo.Echo` listener (`pkg/services/operatorportal`) is launched from `pkg/runtime/startup.go`'s `OnPreServerStartup`, gated behind `operatorConfig.enabled` (default off), on its own port (`operatorConfig.port`, default `9046`) via a plain `*http.Server` + a second async future — same precedent as the existing oidcserver launch, own `echo/v5/middleware.Recover()`.
4. **Done.** Operator config loading (`operatorConfig.operators[]`, env-substituted password hashes) + dedicated login form + dedicated session cookie. No password-change UI.
5. **Done.** IDP management screen — a `resource.Spec` (see `pkg/services/operatorportal/idp_resource.go`) over the step-2 store.
6. **Done**, with the `inmemory-client-service.go` CRUD gap from the Capabilities section fixed along the way. OIDC client management screen — a `resource.Spec` (see `client_resource.go`) over a managed subset of `Client`'s fields.
7. User administration (metadata tagging, disable/enable, session revocation). Not started.
8. Port audit-log viewing in from `/support/*`; retire the public `/support/*` route. Not started.
9. Coordinate with infra to restrict the new port to VPN-only routing (separate internal-only `Service`/Traefik entrypoint pointed at that port, no public ingress, security-group rule) — this is the step that actually enforces "VPN only," everything above just makes it possible to enforce. Not started — this repo can only expose a cleanly separable listener, not enforce network reachability.
10. (v2 candidate) TOTP MFA for operator login, reusing the existing TOTP service. Not started.
11. (follow-up, surfaced by step 6) Wire `BaseHandler`-style audit events into the IDP/client `resource.Spec` Create/Update/Delete callbacks — not yet done; see "Operator-attributed audit events" above.
12. (follow-up) `AddNewClientSecret` and per-client claims updates on the in-memory client service — still unimplemented; not needed by the current screen, but a real gap if either becomes needed.
