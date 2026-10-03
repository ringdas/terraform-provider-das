---
page_title: "users Data Source - das"
subcategory: "Organization"
description: |-
  Every user, all pages: a person of the Organization's directory, invited by email.
---

# das_users (Data Source)

Every user, all pages: a person of the Organization's directory, invited by email.

API: `GET /api/users`.

## Example Usage

```terraform
data "das_users" "example" {
}
```

## Schema

### Optional

- `q` (String) Only users whose email or name contains this text (case-insensitive).
- `status` (String) Only users of this status. One of `INVITED`, `ACTIVE`, `SUSPENDED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) Invitation instant, ISO-8601 UTC.
- `created_by` (String) Actor who invited the user.
- `display_name` (String) Name shown in the UI.
- `email` (String) Lower-case email; the first sign-in links the identity by it.
- `id` (String) User id, `usr_<ULID>`.
- `linked` (Boolean) Whether the identity provider's identity is linked (the user signed in once).
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `status` (String) `INVITED` until the first sign-in, then `ACTIVE`; `SUSPENDED` has no access. One of `INVITED`, `ACTIVE`, `SUSPENDED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

