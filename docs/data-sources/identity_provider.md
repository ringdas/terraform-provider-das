---
page_title: "identity_provider Data Source - das"
subcategory: "Organization"
description: |-
  An identity provider (OIDC or SAML) federated into the Organization's sign-in. Read by its identity.
---

# das_identity_provider (Data Source)

An identity provider (OIDC or SAML) federated into the Organization's sign-in. Read by its identity.

API: `GET /api/identity-providers/{providerId}`.

## Example Usage

```terraform
data "das_identity_provider" "example" {
  id = "<providerId>"
}
```

## Schema

### Required

- `id` (String) The identity provider, `idp_<ULID>`.

### Read-Only

- `applied_at` (String) When the pool last took it, ISO-8601 UTC.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `display_name` (String) Shown on the sign-in page.
- `email_attribute` (String) The provider's attribute carrying the email.
- `email_domains` (List of String) People with these email domains go straight to this provider; its emails link accounts only within them.
- `enabled` (Boolean) Offered on the sign-in page.
- `error` (String) Why the pool refused it.
- `groups_attribute` (String) The provider's attribute carrying the groups (mapped to DAS groups by `externalId`).
- `metadata_url` (String) SAML: the federation metadata URL.
- `name` (String) Its name in the user pool: unique and immutable.
- `oidc_client_id` (String) OIDC: the client id at the provider.
- `oidc_client_secret_set` (Boolean) OIDC: whether a client secret is set (its value is never returned).
- `oidc_client_secret_updated_at` (String) OIDC: when the client secret was last set, ISO-8601 UTC.
- `oidc_issuer` (String) OIDC: the issuer.
- `protocol` (String) How it federates. One of `SAML`, `OIDC`.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.
- `updated_at` (String) Last change, ISO-8601 UTC.
- `updated_by` (String) Actor of the last change.

