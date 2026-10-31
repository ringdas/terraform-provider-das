---
page_title: "org_api_key Resource - das"
subcategory: "Organization"
description: |-
  An Organization API key (`daso_<secret>`): the Organization layer's machine credential, holding one Organization role, for Terraform and CI. Its token never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version. A key never holds more than its minter. Destroy revokes it and deletes the stored token.
---

# das_org_api_key (Resource)

An Organization API key (`daso_<secret>`): the Organization layer's machine credential, holding one Organization role, for Terraform and CI. Its token never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version. A key never holds more than its minter. Destroy revokes it and deletes the stored token.

API: `POST /api/org-api-keys` creates, `GET /api/org-api-keys/{keyId}` reads, `DELETE /api/org-api-keys/{keyId}` deletes.

## Example Usage

```terraform
resource "das_org_api_key" "example" {
  name = "terraform"
  role = "ORG_ADMIN"
}
```

## Schema

### Required

- `name` (String) Name of the key; audit events show it as `service:<name>`. Changing it replaces the resource.
- `role` (String) The Organization role the key holds: it acts with exactly what a person holding this role may do. One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`, `ORG_AUDITOR`. Changing it replaces the resource.

### Optional

- `expires_at` (String) Optional: after this instant (ISO-8601 UTC, in the future) the key authenticates nothing (401 `CREDENTIAL_EXPIRED`). Omit for no expiry. Changing it replaces the resource.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted the key.
- `id` (String) The id the API assigned.
- `revoked_at` (String) When the key was revoked, ISO-8601 UTC; null while it is not.
- `secret_ref` (String) Where the token was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the key is revoked.

## Import

```terraform
import {
  to = das_org_api_key.example
  id = "<keyId>"
}
```
