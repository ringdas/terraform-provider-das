---
page_title: "api_key Data Source - das"
subcategory: "Organization"
description: |-
  An API key of the credential's own Network, minted with a Network credential. Read by its identity.
---

# das_api_key (Data Source)

An API key of the credential's own Network, minted with a Network credential. Read by its identity.

API: `GET /api/network/api-keys/{keyId}`.

## Example Usage

```terraform
data "das_api_key" "example" {
  id = "<keyId>"
}
```

## Schema

### Required

- `id` (String) The API key, `key_<ULID>`, of the Network (the credential's, or the one in the path).

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted the key; null for keys minted before it was recorded.
- `expires_at` (String) When the key stops authenticating, ISO-8601 UTC; null = no expiry.
- `name` (String) Name of the key; audit events show it as `service:<name>`.
- `network_id` (String) The Network the key is bound to.
- `revoked_at` (String) When the key was revoked, ISO-8601 UTC; null while it is not.
- `role` (String) Coarse role of the key. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `scopes` (List of String) Scopes of the key.
- `secret_ref` (String) Where the token was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the key is revoked.

