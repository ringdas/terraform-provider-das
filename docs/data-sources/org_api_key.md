---
page_title: "org_api_key Data Source - das"
subcategory: "Organization"
description: |-
  An Organization API key (`daso_<secret>`): the Organization layer's machine credential, holding one Organization role, for Terraform and CI. Read by its identity.
---

# das_org_api_key (Data Source)

An Organization API key (`daso_<secret>`): the Organization layer's machine credential, holding one Organization role, for Terraform and CI. Read by its identity.

API: `GET /api/org-api-keys/{keyId}`.

## Example Usage

```terraform
data "das_org_api_key" "example" {
  id = "<keyId>"
}
```

## Schema

### Required

- `id` (String) The Organization API key, `okey_<ULID>`.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted the key.
- `expires_at` (String) When the key stops authenticating, ISO-8601 UTC; null = no expiry.
- `name` (String) Name of the key; audit events show it as `service:<name>`.
- `revoked_at` (String) When the key was revoked, ISO-8601 UTC; null while it is not.
- `role` (String) The Organization role the key holds: it acts with exactly what a person holding this role may do. One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`, `ORG_AUDITOR`.
- `secret_ref` (String) Where the token was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the key is revoked.

