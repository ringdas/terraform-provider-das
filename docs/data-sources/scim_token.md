---
page_title: "scim_token Data Source - das"
subcategory: "Organization"
description: |-
  A SCIM token the identity provider provisions the directory with. Read by its identity.
---

# das_scim_token (Data Source)

A SCIM token the identity provider provisions the directory with. Read by its identity.

API: `GET /api/scim-tokens/{tokenId}`.

## Example Usage

```terraform
data "das_scim_token" "example" {
  id = "<tokenId>"
}
```

## Schema

### Required

- `id` (String) The SCIM token, `scimtok_<ULID>`.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted it.
- `expires_at` (String) When a newer token retires it, ISO-8601 UTC; null = no end.
- `hint` (String) The secret's last four characters, to recognise it.
- `last_request` (String) The last request, e.g. `PATCH /Users/{id} 200`.
- `last_used_at` (String) The identity provider's last request with it, ISO-8601 UTC.
- `request_count` (Number, integer) Requests made with it.
- `secret_ref` (String) Where the secret was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the token is revoked.

