---
page_title: "network_audit_event Data Source - das"
subcategory: "Organization"
description: |-
  One row of the audit trail of the Network named by `network_id`.
---

# das_network_audit_event (Data Source)

One row of the audit trail of the Network named by `network_id`.

API: `GET /api/networks/{networkId}/audit-events/{auditEventId}`.

## Example Usage

```terraform
data "das_network_audit_event" "example" {
  network_id     = "<networkId>"
  audit_event_id = "<auditEventId>"
}
```

## Schema

### Required

- `audit_event_id` (String) The audit event, `aud_<ULID>`.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `action` (String) What was done, `<resource>.<verb>`, e.g. `api-key.create`; `audit-chain.genesis` for a chain's first link.
- `actor` (String) Who acted: `human:<email>`, `service:<apiKeyName>` or `agent:<name>/<session>`.
- `actor_network_id` (String) The Network the principal's credential is bound to; null for an Organization-level one.
- `actor_role` (String) The coarse role of the principal. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `actor_roles` (List of String) The role names the principal held (GET /api/roles of the Organisation API describes each); empty for an API key.
- `chain` (String) The chain the row is a link of: `org`, or `network:<networkId>`.
- `changes` (Attributes List) What the mutation changed, one entry per row written; null before the chain. Empty for an act that changes no row (a test email).
- `channel` (String) The door: the console, the das CLI (its User-Agent, a declaration), an API key, an AI agent, the identity provider over SCIM, or the platform itself. One of `CONSOLE`, `CLI`, `API`, `AGENT`, `SCIM`, `SYSTEM`.
- `credential_id` (String) The credential's id: an API key's `key_<ULID>`, an Organization API key's `okey_<ULID>`, a session's subject or token id, a SCIM token's `scimtok_<ULID>`.
- `credential_kind` (String) The kind of credential, never its secret. One of `API_KEY`, `ORG_API_KEY`, `COGNITO`, `NETWORK_SESSION`, `SCIM_TOKEN`.
- `hash` (String) sha256(prevHash || RFC 8785 JSON of the record), lower-case hex; null before the chain.
- `id` (String) The row, `aud_<ULID>` on a Network's trail, `orgaud_<ULID>` on the org layer's.
- `occurred_at` (String) When the change committed, ISO-8601 UTC with milliseconds.
- `prev_hash` (String) The hash of the link before (64 zeros for the genesis); null before the chain.
- `reason` (String) Why, when the operation records it.
- `request_id` (String) The X-Request-Id of the request that made the change.
- `resource_id` (String) The resource's id.
- `resource_type` (String) The kind of resource, e.g. `api-key`.
- `seq` (Number, integer) Its position in the chain, from 1 (the genesis) without gaps; null for a row written before the chain (sealed by the genesis).

<a id="nestedatt--changes"></a>
### Nested Schema for `changes`

- `fields` (Map of String) The changed fields by name: `{ before, after }` (a CREATE's before and a DELETE's after are null; instants ISO-8601, decimals as strings), or `{ redacted: true, changed }` for a secret, a hash or a secret reference, whose value is never recorded.
- `id` (String) The row's id.
- `model` (String) The model (table), e.g. `ApiKey`.
- `op` (String) What happened to the row. One of `CREATE`, `UPDATE`, `DELETE`.

