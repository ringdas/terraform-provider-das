---
page_title: "sign_in_policy Resource - das"
subcategory: "Organization"
description: |-
  How people sign in to the Organization: local sign-in, MFA, session length. One per Organization; destroy leaves the policy as it is.
---

# das_sign_in_policy (Resource)

How people sign in to the Organization: local sign-in, MFA, session length. One per Organization; destroy leaves the policy as it is.

API: `PUT /api/identity/sign-in-policy` creates, `GET /api/identity/sign-in-policy` reads, `PUT /api/identity/sign-in-policy` updates. Destroy: The Organization always has a sign-in policy; destroy stops managing it.

## Example Usage

```terraform
resource "das_sign_in_policy" "example" {
  local_sign_in = "ALLOWED"
  mfa           = "OFF"
  session_days  = 1
}
```

## Schema

### Required

- `local_sign_in` (String) `IDP_ONLY` needs an enabled identity provider (409 `LOCKOUT` otherwise). One of `ALLOWED`, `IDP_ONLY`.
- `mfa` (String) Multi-factor sign-in with an authenticator app. One of `OFF`, `OPTIONAL`, `REQUIRED`.
- `session_days` (Number, integer) How many days a sign-in lasts in the web apps.

### Read-Only

- `applied_at` (String) When the pool last took the policy, ISO-8601 UTC.
- `error` (String) Why the pool refused it.
- `id` (String) A constant: the Organization has one sign_in_policy.
- `state` (String) Whether the user pool carries what was saved: `APPLIED`, `PENDING` (the duty platform.identity-sync applies it within seconds) or `FAILED` (the pool refused it; `error` says why). One of `APPLIED`, `PENDING`, `FAILED`.

## Import

```terraform
import {
  to = das_sign_in_policy.example
  id = "sign_in_policy"
}
```
