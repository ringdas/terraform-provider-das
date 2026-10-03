---
page_title: "identity_provider Resource - das"
subcategory: "Organization"
description: |-
  An identity provider (OIDC or SAML) federated into the Organization's sign-in. A client secret is a write-only attribute: it never enters Terraform state.
---

# das_identity_provider (Resource)

An identity provider (OIDC or SAML) federated into the Organization's sign-in. A client secret is a write-only attribute: it never enters Terraform state.

API: `POST /api/identity-providers` creates, `GET /api/identity-providers/{providerId}` reads, `PATCH /api/identity-providers/{providerId}` updates, `DELETE /api/identity-providers/{providerId}` deletes.

## Example Usage

```terraform
resource "das_identity_provider" "example" {
  display_name                  = "Microsoft Entra ID"
  email_domains                 = ["example.com"]
  groups_attribute              = "groups"
  name                          = "entra"
  oidc_client_id                = "00000000-0000-0000-0000-000000000000"
  oidc_client_secret_wo         = "client-secret"
  oidc_client_secret_wo_version = 1
  oidc_issuer                   = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/v2.0"
  protocol                      = "OIDC"
}
```

## Schema

### Required

- `display_name` (String) Shown on the sign-in page.
- `name` (String) Its name in the user pool, 3-32 letters, digits or hyphens; immutable. Changing it replaces the resource.
- `protocol` (String) How it federates; immutable. One of `SAML`, `OIDC`. Changing it replaces the resource.

### Optional

- `email_attribute` (String) The provider's attribute carrying the email (default: SAML's emailaddress claim, OIDC's `email`).
- `email_domains` (List of String) Email domains that go straight to this provider (default none); a domain belongs to one provider.
- `enabled` (Boolean) Offered on the sign-in page (default true).
- `groups_attribute` (String) The provider's attribute carrying the groups, e.g. Entra's groups claim.
- `metadata_url` (String) SAML: the federation metadata URL (https).
- `oidc_client_id` (String) OIDC: the client id at the provider.
- `oidc_client_secret_wo` (String, Sensitive, Write-only) OIDC: the client secret at the provider; DAS keeps it in the installation's secret store and never returns it. Write-only: never stored in Terraform state (Terraform >= 1.11); change `oidc_client_secret_wo_version` to send a new value.
- `oidc_client_secret_wo_version` (Number, integer) Any number; change it to send `oidc_client_secret_wo` again (a rotation).
- `oidc_issuer` (String) OIDC: the issuer (https).

### Read-Only

- `applied_at` (String) When the pool last took it, ISO-8601 UTC.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `error` (String) Why the pool refused it.
- `id` (String) The id the API assigned.
- `oidc_client_secret_set` (Boolean) OIDC: whether a client secret is set (its value is never returned).
- `oidc_client_secret_updated_at` (String) OIDC: when the client secret was last set, ISO-8601 UTC.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.
- `updated_at` (String) Last change, ISO-8601 UTC.
- `updated_by` (String) Actor of the last change.

## Import

```terraform
import {
  to = das_identity_provider.example
  id = "<providerId>"
}
```
