# Support Portal Boundary

The `/support/` portal is intentionally scoped to identity-engine operations owned by this service.

## In Scope

- Viewing and filtering identity audit events (`tmp/auditstore.jsonl` in development).
- Operational diagnostics for auth/identity engine behavior.
- Support-admin-only access controls.

## Out of Scope

- Source-of-truth user directory management.
- Global user listing/search across external user stores.
- Business profile CRUD, entitlement administration, and domain-specific account workflows.

## Why

This service is an OIDC engine and identity orchestrator. Integrator systems own canonical user records and associated admin portals.

## Integration Pattern

- Keep identity/auth operations in this service.
- Build full user-management portals against the integrator's canonical user APIs/databases.
- Optionally link those systems with this service through audited API calls and event streams.
