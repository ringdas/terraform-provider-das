---
page_title: "scim_tokens Data Source - das"
subcategory: "Organization"
description: |-
  Every scim_token, all pages: a SCIM token the identity provider provisions the directory with.
---

# das_scim_tokens (Data Source)

Every scim_token, all pages: a SCIM token the identity provider provisions the directory with.

API: `GET /api/scim-tokens`.

## Example Usage

```terraform
data "das_scim_tokens" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted it.
- `expires_at` (String) When a newer token retires it, ISO-8601 UTC; null = no end.
- `hint` (String) The secret's last four characters, to recognise it.
- `id` (String) SCIM token id, `scimtok_<ULID>`.
- `last_request` (String) The last request, e.g. `PATCH /Users/{id} 200`.
- `last_used_at` (String) The identity provider's last request with it, ISO-8601 UTC.
- `request_count` (Number, integer) Requests made with it.
- `secret_ref` (String) Where the secret was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the token is revoked.

