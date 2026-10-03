---
page_title: "sign_in_policy Data Source - das"
subcategory: "Organization"
description: |-
  How people sign in to the Organization: local sign-in, MFA, session length.
---

# das_sign_in_policy (Data Source)

How people sign in to the Organization: local sign-in, MFA, session length.

API: `GET /api/identity/sign-in-policy`.

## Example Usage

```terraform
data "das_sign_in_policy" "example" {
}
```

## Schema

### Read-Only

- `applied_at` (String) When the pool last took the policy, ISO-8601 UTC.
- `error` (String) Why the pool refused it.
- `local_sign_in` (String) Whether the pool's own email-and-password sign-in is offered, or only the identity providers. One of `ALLOWED`, `IDP_ONLY`.
- `mfa` (String) Multi-factor sign-in with an authenticator app. One of `OFF`, `OPTIONAL`, `REQUIRED`.
- `session_days` (Number, integer) How many days a sign-in lasts in the web apps.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.

