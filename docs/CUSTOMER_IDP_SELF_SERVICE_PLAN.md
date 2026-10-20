# Customer Self-Service IDP Settings Portal

**Status: design only — not implemented.**

## Goal

A designated user at a customer with a claimed domain (e.g. `evilcorp.com`, backed by an `OIDCProtocol` IDP record) can log into the existing self-service management portal and, from a new tab, update their own IDP's connection settings — authority/issuer, client ID, client secret, scope — without filing a support ticket.

## Reference model

Auth0 has a similar feature ("Self-Service SSO"): the app admin issues a one-time setup link scoped to a single connection, the customer's IT admin fills in a wizard, and the resulting connection is typically **promoted to production by the app admin** rather than going live the instant the customer saves it. That's more staged than what's described above (an always-available tab, self-serve end to end). This doc follows the always-available-tab shape since that's the stated target, but the Auth0 staged-promotion idea is called out below as an option worth considering given the blast radius of this feature (see Security).

## Why this isn't a small feature

Everything needed for this already has a working analog somewhere in the codebase, but none of the actual pieces exist yet. Confirmed by direct inspection:

- **No per-user authorization primitive.** `RageUser` has no roles/claims/groups — the only IDP-related field on a user is `Identity.idp_slug` (which IDP a login came through, not a management relationship). There is nothing today that expresses "user X manages IDP Y."
- **The IDP store is in-memory and read-only.** `SingletonIDPService` (`proto/oidc/idp/idp.proto`) has exactly one RPC — `ListIDP`. The concrete implementation (`pkg/services/idp/inmemory/inmemory-idp-service.go`) loads a `proto_oidc_models.IDPs` once at boot from `cmd/server/config/idps.json` and never writes it back anywhere. No `Create`/`Update`/`Delete` RPC exists, not even stubbed.
- **No reversible encryption exists in the repo.** The only crypto-adjacent service is `IPasswordHasher` (one-way bcrypt/argon2id) — unsuitable, since a `client_secret` must be retrievable in plaintext to actually perform OAuth2 token exchange, not just verified.
- **No database dependency exists at all.** `go.mod` has no SQL driver, ORM, or embedded-DB dependency. Any persistence for customer-edited settings is greenfield.

None of this is a reason not to build it — it's the actual scope. The design below closes each gap using the closest existing pattern in the codebase rather than inventing new conventions where one already fits.

## Proposed design

### 1. Authorization: "this user manages IDP slug X"

Reuse a mechanism that's already live and working end-to-end for a different purpose. `example/services/AuthorizationCodeClaimsAugmentor/AuthorizationCodeClaimsAugmentor.go` already reads a generic `ExampleUser.Metadata` key/value pair (`integrity_id`) at token-mint time and injects it as a claim via `IAuthorizationCodeClaimsAugmentor.AugmentTokenClaims`. Do the same thing for this feature:

- Add a metadata entry on the designated user's `ExampleUser.Metadata`, e.g. `{"managed_idp_slug": "evilcorp"}`.
- Extend the claims augmentor to mint a `managed_idp_slug` claim onto the ID token when that metadata key is present.
- The new management tab and its handler check `s.ClaimsPrincipal().GetClaimsByType("managed_idp_slug")` — same call shape as the existing `isSupportAdmin`/`ACRClaimedDomain` checks in `example/runtime/startup.go` / `profile.go` / `password.go`.

**This assignment is an internal operation, not self-service.** Only your team sets the `managed_idp_slug` metadata on a customer's user record (via an internal/support-portal tool, not something the customer can grant themselves). This mirrors how Auth0 treats the initial setup-link issuance as an admin action.

### 2. Writable IDP store

Add one narrowly-scoped RPC, not a generic `UpdateIDP`:

```proto
rpc UpdateOIDCProtocol (UpdateOIDCProtocolRequest) returns (UpdateOIDCProtocolResponse) {}

message UpdateOIDCProtocolRequest {
    string slug = 1;
    string authority = 2;
    string client_id = 3;
    string client_secret = 4;  // empty = "don't change"
    string scope = 5;
}
```

Deliberately excludes `claimed_domains`, `enabled`, `hidden`, `auto_create`, `multi_factor_required` — those stay admin/config-file-only. A customer should never be able to flip their own `auto_create` or claim a different domain.

**Persistence.** The in-memory store needs to become durable and — this is the part to confirm before building — durable *across replicas*, not just across restarts, if this service ever runs as more than one instance. Two options:
- MVP: embed a pure-Go SQLite (e.g. `modernc.org/sqlite`, no cgo) as a local durable store, with `idps.json` still used as the first-boot seed and DB rows overriding JSON-provided defaults per slug on read.
- Production-correct if multi-replica: a real shared DB (Postgres) — SQLite-on-local-disk will silently diverge between replicas/pods.

**Open question to confirm before building:** does this service run as a single instance today, or multiple replicas behind a load balancer? That answer picks between the two options above.

### 3. Secret encryption at rest

New service, `ISecretCipher` (AES-256-GCM, envelope-encrypted with a master key from env/KMS — not yet existing anywhere in the repo, this is new infrastructure). `client_secret` is encrypted before being written to the store and decrypted only at the point of use — where `pkg/services/codeexchanges/genericoidc/genericoidc-code-exchange.go` currently reads `v.Oidc.ClientSecret` directly for the real token exchange call, that read becomes a decrypt-on-demand.

**Never redisplay the secret.** Once saved, the UI shows only a masked indicator (e.g. "configured, last 4: ••1234") — same pattern as Stripe/Auth0 API-key UIs. The edit form's secret field is write-only; leaving it blank means "keep the existing value."

### 4. New management tab

Follows the exact existing pattern (`example/go-app/management/htmx/handlers/profile/` as the template):
- New path constant in `pkg/wellknown/wellknown_echo/wellknown.go` (e.g. `HTMXManagementSSOPath`).
- New `handlers/sso/sso.go` — `service` embedding `BaseHandler`, `Ctor`, `AddScopedIHandler`.
- New `components/sso.go` render component (Authority, Client ID, Client Secret [masked/write-only], Scope fields).
- Registered in `registration.go` alongside the other management handlers.
- New nav entry in `components/sidebar.go` — **rendered conditionally**, only when `managed_idp_slug` claim is present, same as any other claim-gated UI element in this codebase.

### 5. Change safety (the part that matters most)

A wrong or malicious edit here doesn't just affect one account — it repoints where an entire claimed domain's employees authenticate. Treat it like a credential-rotation operation, not a profile edit:

- **Step-up re-authentication before accepting the save.** The app already has a working "force fresh login" primitive: `pkg/models/wellknown.go`'s `ACRNoSSO = "urn:rage:no-sso"`, already used by the management auth-gate's self-login trigger (`example/runtime/startup.go`'s `fetch('/api/login', {..., acrValues: 'urn:rage:no-sso'})`). Reuse it: require the user to re-authenticate via that same flow immediately before the settings form will submit, and check for a recent-enough auth timestamp/ACR on the session before accepting the POST.
- **Validate before saving**, not after. Fetch the new authority's OIDC discovery document (`/.well-known/openid-configuration`) synchronously on submit and surface a clear inline error if it doesn't resolve or doesn't look like a valid OIDC issuer — don't let a typo silently lock out the whole company on their next login attempt.
- **Audit every change**, following the existing `BaseHandler.SubmitAuditEvent` convention (see `pkg/services/echo/handlers/passwordreset/password-reset.go` for the call shape): `com.fluffybunny.identity.idp.oidc_settings.updated`, subject = the acting user, attributes tagging `idp_slug` and which fields changed — **never include the secret value itself in the audit payload**, only a boolean "secret_changed."
- **Notify a fixed internal contact** (security/support alias) on every change, independent of the audit log, so a compromised customer account doesn't get to quietly redirect trust without anyone noticing.
- **Consider (optional for v1, worth it for v2):** an Auth0-style staged rollout — the customer's edit is saved as a "pending" configuration that doesn't take effect until either a cooldown window passes or an internal admin promotes it. This trades away "instant self-service" for a real safety net against a locked-out customer or a compromised customer account; worth revisiting once this ships and you see how it's actually used.

## Phased order

1. Confirm deployment topology (single instance vs multi-replica) — determines the persistence choice.
2. Build `ISecretCipher` + the writable IDP store (`UpdateOIDCProtocol` RPC + persistence), independent of any UI.
3. Extend the claims augmentor for `managed_idp_slug`; build the internal tool/process for assigning it to a customer user (admin-only).
4. Build the management tab (handler + component + sidebar entry), gated on the claim.
5. Add step-up re-auth requirement + pre-save discovery-document validation.
6. Wire audit events + notification.
7. Decide on staged-rollout (pending/promote) based on real usage, once the always-available version has shipped.
