# SAML Support Plan (Ping)

**Status: not implemented — deferred.** This is a design-only document. We are not building SAML support speculatively; this plan exists so that if a large client requires it, we aren't starting from zero. Do not begin implementation without a confirmed customer requirement.

## Why deferred

Unlike OIDC, SAML implementations vary significantly IdP to IdP (signing choices, attribute naming, encryption, IdP-initiated vs SP-initiated). Building this properly is a real investment — worth doing when a client forces the issue, not before.

## Library recommendation

**`github.com/russellhaering/gosaml2`** (+ `github.com/russellhaering/goxmldsig` for signature verification).

This is what HashiCorp Vault's SAML auth method uses, specifically because Vault has to interoperate with real-world enterprise IdPs (Ping, Okta, ADFS, Azure AD, OneLogin) without controlling the customer's IdP config. It provides:
- IdP metadata parsing, including multiple signing certs (rotation)
- Signature validation against either the Response or the Assertion being signed (Ping can do either)
- Encrypted-assertion support (common in Ping setups)
- Configurable clock-skew tolerance
- AuthnRequest building + redirect-binding URL construction for SP-initiated flow

It's low-level/headless — you call `ServiceProvider.RetrieveAssertionInfo(rawResponse)` and wire the HTTP glue yourself, which fits this codebase's DI/handler conventions better than the alternative.

**Alternative: `github.com/crewjam/saml`** — more popular, more complete (SP + IdP roles, nicer metadata tooling), but its `samlsp` layer wants to own routing/sessions/cookies. Its lower-level `saml.ServiceProvider` type (bypassing `samlsp`) is roughly on par with gosaml2 for this use case. Worth a quick spike of both against a real Ping tenant before committing, but start with gosaml2.

## How this fits the existing external-IDP architecture

The current GitHub/OIDC external-IDP flow has three stages, and SAML should slot into the same shape rather than being bolted on separately.

### 1. IDP model (`proto/oidc/models/idp.proto`)

Add a `SamlProtocol` message and a `Protocol_Saml` oneof variant alongside the existing `OIDCProtocol` / `GithubOAuth2Protocol` / `OAuth2Protocol`. Per-IDP fields needed:

- IdP `EntityID`, SSO URL, and either a pasted signing cert (PEM, possibly multiple for rollover) or an IdP metadata URL to fetch/refresh
- Expected `NameIDFormat`
- **Configurable attribute-name → claim mapping** (email attribute URI, etc.) — non-negotiable, since attribute naming varies per Ping tenant
- Whether assertions are expected to be encrypted
- Whether IdP-initiated SSO is allowed for this IDP (see quirks below)
- Clock-skew tolerance

This reuses the existing `claimed_domains` / `auto_create` / `email_verification_required` / `multi_factor_required` fields already on `IDP` — no changes needed there.

### 2. SP-wide config

A new `SAMLSPConfig` in `contracts_config.Config`, mirroring the existing `SelfIDPConfig` pattern: our SP `EntityID`, ACS base URL, and an optional SP signing/decryption keypair.

### 3. Start-login handler

Add a `Protocol_Saml` case to the existing protocol switch in `pkg/services/echo/handlers/externalidp/externalidp.go` and `externalidp/api/externalidp.go`. Instead of PKCE/nonce, build a SAML `AuthnRequest`, store the AuthnRequest ID + `IdpHint`/`ParentState`/`Directive` in the **same `ExternalOauth2Cookie` mechanism** already used for OAuth2/OIDC (repurposed, no new state-storage plumbing needed), and redirect via HTTP-Redirect binding using the generated `externalState` as `RelayState`.

### 4. ACS handler

A new, dedicated `/saml/acs` POST handler. SAML's binding is a form POST of `SAMLResponse`/`RelayState`, not a GET with `code`/`state`, so it can't reuse `pkg/services/echo/handlers/oauth2/callback` directly. It should:
- Look up the pending state by `RelayState` (when present — see IdP-initiated below)
- Validate the assertion via gosaml2 (signature, conditions, audience, `InResponseTo` if SP-initiated)
- Extract identity via the configured attribute mapping
- Mint the same internal JWT the GitHub/OIDC exchangers mint (`tokenService.MintToken`), and hand off into the link-dance

### 5. Refactor recommendation (do this regardless of SAML)

`callback.go`'s "link dance" (existing-link lookup → auto-link-by-email → auto-create → claimed-domain verification → MFA → cookie-setting, ~470 lines) is currently written as inline closures inside `Do()`. Before adding SAML, extract it into a shared, protocol-agnostic function/service (e.g. `IExternalIdentityLinker.LinkOrCreate(ctx, c, idp, externalIdentity, directive, authorizationRequest)`) that both the OAuth2/OIDC callback and the new SAML ACS handler call.

This also fixes a pre-existing drift worth noting: the start-login handler supports `Protocol_Oauth2`, but the callback's exchange switch only has cases for `Protocol_Github` and `Protocol_Oidc` — a raw generic-OAuth2 IDP would silently fail at callback today.

### 6. SP metadata endpoint

`GET /saml/metadata`, serving our SP's metadata XML so it can be imported into the PingFederate/PingOne admin console.

## Ping-specific quirks to design in from day one

- **Signed Response vs signed Assertion vs both** — varies by Ping deployment; validate whichever is actually signed rather than assuming one.
- **Encrypted assertions** — common in Ping configs; needs an SP decryption keypair as a first-class config option.
- **IdP-initiated SSO** — common with Ping (an admin-configured "tile" POSTs directly to the ACS with no preceding `AuthnRequest`, so no `InResponseTo`/`RelayState` to correlate). The ACS handler needs a distinct path for "no matching stored state" that trusts audience-restriction + signature + a per-IDP default landing directive, with a config toggle to disable this per IDP for SP-initiated-only tenants.
- **Attribute-naming variance** — don't hardcode attribute URIs for email/name; keep the mapping configurable per IDP.
- **Cert rotation** — prefer metadata-URL auto-refresh over a single pasted cert.
- **Clock skew** — Ping assertions are short-lived; make the tolerance configurable (same pattern as the existing `jwxt.WithAcceptableSkew` usage for OIDC).

## Testing strategy

- **Keycloak** (Docker, local/CI) configured as a SAML IdP for fast iteration — exercises real signature validation, real XML, real quirks without needing a live Ping tenant.
- **samltest.id** for a quick manual interoperability sanity check.
- A real **PingOne/PingFederate trial tenant** reserved for final validation before shipping — Ping-specific config surface (attribute names, signing choices) can't be fully validated against a generic IdP.

## Suggested phased order (when triggered)

1. Spike gosaml2 against Keycloak locally — validate the `RetrieveAssertionInfo`/AuthnRequest-building APIs actually fit before touching proto/DI.
2. Proto + config additions (`Protocol_Saml`, `SAMLSPConfig`).
3. Extract the link-dance refactor (independent, de-risks everything after it).
4. SP-initiated flow: start-handler branch + ACS handler + metadata endpoint, validated against Keycloak.
5. IdP-initiated support + attribute-mapping config.
6. Validate against a real Ping trial tenant; adjust for whatever it does differently than Keycloak.

## Trigger condition

Begin implementation only when a specific client contractually requires SAML/Ping SSO. Revisit this document at that point — library versions and Ping's own SAML behavior may have shifted.
