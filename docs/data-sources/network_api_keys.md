---
page_title: "network_api_keys Data Source - das"
subcategory: "Organization"
description: |-
  Every network_api_key, all pages: an API key of the Network named by `network_id`, minted by the Organization.
---

# das_network_api_keys (Data Source)

Every network_api_key, all pages: an API key of the Network named by `network_id`, minted by the Organization.

API: `GET /api/networks/{networkId}/api-keys`.

## Example Usage

```terraform
data "das_network_api_keys" "example" {
  network_id = "<networkId>"
}
```

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted the key; null for keys minted before it was recorded.
- `expires_at` (String) When the key stops authenticating, ISO-8601 UTC; null = no expiry.
- `id` (String) API key id, `key_<ULID>`.
- `name` (String) Name of the key; audit events show it as `service:<name>`.
- `network_id` (String) The Network the key is bound to.
- `revoked_at` (String) When the key was revoked, ISO-8601 UTC; null while it is not.
- `role` (String) Coarse role of the key. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `scopes` (List of String) Scopes of the key.
- `secret_ref` (String) Where the token was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the key is revoked.

