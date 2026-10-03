---
page_title: "user Data Source - das"
subcategory: "Organization"
description: |-
  A person of the Organization's directory, invited by email. Read by its identity.
---

# das_user (Data Source)

A person of the Organization's directory, invited by email. Read by its identity.

API: `GET /api/users/{userId}`.

## Example Usage

```terraform
data "das_user" "example" {
  id = "<userId>"
}
```

## Schema

### Required

- `id` (String) The user, `usr_<ULID>`.

### Read-Only

- `created_at` (String) Invitation instant, ISO-8601 UTC.
- `created_by` (String) Actor who invited the user.
- `display_name` (String) Name shown in the UI.
- `email` (String) Lower-case email; the first sign-in links the identity by it.
- `linked` (Boolean) Whether the identity provider's identity is linked (the user signed in once).
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `status` (String) `INVITED` until the first sign-in, then `ACTIVE`; `SUSPENDED` has no access. One of `INVITED`, `ACTIVE`, `SUSPENDED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

