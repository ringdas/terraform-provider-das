---
page_title: "identity_providers Data Source - das"
subcategory: "Organization"
description: |-
  Every identity_provider, all pages: an identity provider (OIDC or SAML) federated into the Organization's sign-in.
---

# das_identity_providers (Data Source)

Every identity_provider, all pages: an identity provider (OIDC or SAML) federated into the Organization's sign-in.

API: `GET /api/identity-providers`.

## Example Usage

```terraform
data "das_identity_providers" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `applied_at` (String) When the pool last took it, ISO-8601 UTC.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `display_name` (String) Shown on the sign-in page.
- `email_attribute` (String) The provider's attribute carrying the email.
- `email_domains` (List of String) People with these email domains go straight to this provider; its emails link accounts only within them.
- `enabled` (Boolean) Offered on the sign-in page.
- `error` (String) Why the pool refused it.
- `groups_attribute` (String) The provider's attribute carrying the groups (mapped to DAS groups by `externalId`).
- `id` (String) Identity provider id, `idp_<ULID>`.
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

