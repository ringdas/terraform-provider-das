---
page_title: "identity Data Source - das"
subcategory: "Organization"
description: |-
  The Organization's sign-in: user pool, identity providers, SCIM and their sync state.
---

# das_identity (Data Source)

The Organization's sign-in: user pool, identity providers, SCIM and their sync state.

API: `GET /api/identity`.

## Example Usage

```terraform
data "das_identity" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Each thing the duty applies, with its state.
- `mode` (String) The user pool the API administers: Cognito, the dev identity provider (local development), or none (settings are stored, not applied). One of `COGNITO`, `DEV`, `NONE`.
- `pool` (Attributes) The pool as last described; null until the duty reached it.
- `redirects` (Attributes) The web app client's redirect URLs: one per host people sign in on (Cognito has no wildcards).
- `scim` (Attributes) SCIM provisioning at a glance.
- `synced_at` (String) When the duty last reached the pool, ISO-8601 UTC.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `applied_at` (String) When the pool last took it, ISO-8601 UTC.
- `attempted_at` (String) The duty's last attempt, ISO-8601 UTC.
- `error` (String) Why the pool refused the last attempt.
- `key` (String) `sign-in-policy`, `web-client`, `cli-client` or `provider:<idp id>`.
- `label` (String) What it is, for people.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.

<a id="nestedatt--pool"></a>
### Nested Schema for `pool`

- `clients` (Attributes List) The app clients the duty keeps in step.
- `oidc_redirect_uri` (String) An OIDC provider's redirect URI.
- `region` (String) The pool's region.
- `saml_acs_url` (String) A SAML provider's reply URL (Entra: Reply URL / Assertion Consumer Service URL).
- `saml_entity_id` (String) A SAML provider's identifier of the pool (Entra: Identifier / Entity ID).
- `sign_in_url` (String) The sign-in page (Hosted UI) people are sent to.
- `user_pool_id` (String) The pool id, e.g. `eu-central-1_AbC123`.

<a id="nestedatt--pool--clients"></a>
### Nested Schema for `pool.clients`

- `client_id` (String) The app client id.
- `kind` (String) The web apps (redirect URLs follow the Networks) or the das CLI (loopback redirects). One of `WEB`, `CLI`.
- `name` (String) Its name in the pool.

<a id="nestedatt--redirects"></a>
### Nested Schema for `redirects`

- `desired` (Number, integer) Hosts people sign in on: the Organization, every module, every Network's modules.
- `limit` (Number, integer) Cognito's limit per app client (100).
- `missing_hosts` (List of String) Hosts without a registered redirect URL: nobody can sign in on them.
- `registered` (Number, integer) Redirect URLs the client carries; null until the duty applied them.

<a id="nestedatt--scim"></a>
### Nested Schema for `scim`

- `active_tokens` (Number, integer) Tokens that authenticate now.
- `base_url` (String) The SCIM endpoint (Entra: Tenant URL).
- `enabled` (Boolean) Whether the identity provider may provision.
- `last_request` (String) The last request, e.g. `PATCH /Users/{id} 200`.
- `last_request_at` (String) The identity provider's last request, ISO-8601 UTC.

