---
page_title: "current_network_session Data Source - das"
subcategory: "Organization"
description: |-
  Who the credential is in its Network: principal, roles and scopes.
---

# das_current_network_session (Data Source)

Who the credential is in its Network: principal, roles and scopes.

API: `GET /api/network/session`.

## Example Usage

```terraform
data "das_current_network_session" "example" {
}
```

## Schema

### Read-Only

- `actor` (String) The actor audit events record.
- `credential` (String) The kind of credential. One of `API_KEY`, `COGNITO`, `NETWORK_SESSION`.
- `key_id` (String) The API key, `key_<ULID>`, when the credential is one (a product introspects keys here); null otherwise.
- `network` (Attributes) A Network: one tenant of the Organization.
- `role` (String) The coarse rank; null without a role. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `roles` (List of String) The Network roles of a user; empty for an API key.
- `scopes` (List of String) The scopes the credential holds in the Network.
- `user_id` (String) The user of a Network token; null for an API key.

<a id="nestedatt--network"></a>
### Nested Schema for `network`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits (e.g. `NT-1234567`): unique in the Organization, not the slug and not the id, fixed once given.
- `slug` (String) Host label of the Network, unique in the Organization.
- `timezone` (String) IANA time zone where the Network's days and months begin.
- `updated_at` (String) Last change, ISO-8601 UTC.

