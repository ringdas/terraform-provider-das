---
page_title: "scim_token Resource - das"
subcategory: "Organization"
description: |-
  A SCIM token the identity provider provisions the directory with. Its secret never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version, for the identity provider's configuration. Destroy revokes it and deletes the stored secret.
---

# das_scim_token (Resource)

A SCIM token the identity provider provisions the directory with. Its secret never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version, for the identity provider's configuration. Destroy revokes it and deletes the stored secret.

API: `POST /api/scim-tokens` creates, `GET /api/scim-tokens/{tokenId}` reads, `DELETE /api/scim-tokens/{tokenId}` deletes.

## Example Usage

```terraform
resource "das_scim_token" "example" {
}
```

## Schema

### Optional

- `retire_previous` (Boolean) The live tokens stop working in 24 hours, enough to update the identity provider (default false). Changing it replaces the resource.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted it.
- `expires_at` (String) When a newer token retires it, ISO-8601 UTC; null = no end.
- `hint` (String) The secret's last four characters, to recognise it.
- `id` (String) The id the API assigned.
- `last_request` (String) The last request, e.g. `PATCH /Users/{id} 200`.
- `last_used_at` (String) The identity provider's last request with it, ISO-8601 UTC.
- `request_count` (Number, integer) Requests made with it.
- `secret_ref` (String) Where the secret was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the token is revoked.

## Import

```terraform
import {
  to = das_scim_token.example
  id = "<tokenId>"
}
```
