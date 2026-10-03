---
page_title: "network_session_token Ephemeral Resource - das"
subcategory: "Organization"
description: |-
  A short-lived Network token for the person of the provider's Organization session, to call a product API in the same run. Never stored in state.
---

# das_network_session_token (Ephemeral Resource)

A short-lived Network token for the person of the provider's Organization session, to call a product API in the same run. Never stored in state. Requires Terraform 1.10 or later; its values never enter state or plan files.

API: `POST /api/networks/{networkId}/session-tokens`.

## Example Usage

```terraform
ephemeral "das_network_session_token" "example" {
  network_id = "<networkId>"
}
```

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `expires_at` (String) When the token expires, ISO-8601 UTC.
- `expires_in` (Number, integer) Its lifetime in seconds.
- `roles` (List of String) The Network roles the token carries.
- `token` (String, Sensitive) The JWT - only in the first response. An idempotent replay returns null: a credential is never stored.
- `token_type` (String, Sensitive) Always `Bearer`.

