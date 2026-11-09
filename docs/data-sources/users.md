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
- `invitation` (Attributes) The invitation an admin sent and how far its delivery got; null for a person SCIM provisioned or the first owner.
- `linked` (Boolean) Whether the identity provider's identity is linked (the user signed in once).
- `sign_in` (Attributes) How this person signs in under the Organization's sign-in configuration now: an enabled identity provider whose email domains cover their address takes them to it; otherwise a password of the Organization's own sign-in while the sign-in policy allows it.
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `status` (String) `INVITED` until the first sign-in, then `ACTIVE`; `SUSPENDED` has no access. One of `INVITED`, `ACTIVE`, `SUSPENDED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--items--invitation"></a>
### Nested Schema for `items.invitation`

- `attempted_at` (String) The last attempt; null before the first.
- `completed_at` (String) When the delivery finished; null while it is pending.
- `emailed_at` (String) When the invitation email went out; null before, or when the installation sends no email.
- `error` (String) Why the last attempt did not finish (it is retried); null when it did.
- `identity_provider` (String) The identity provider the email named (`IDENTITY_PROVIDER`), otherwise null.
- `invited_at` (String) When, ISO-8601 UTC.
- `invited_by` (String) Actor who invited them, or invited them again.
- `method` (String) How the delivery let them in (`signIn.method` at the time); null until it acted. One of `IDENTITY_PROVIDER`, `PASSWORD`.
- `pool_invited_at` (String) When their account of the Organization's own sign-in was created and its email with a temporary password sent (`PASSWORD`); null before, or when they already had a password.

<a id="nestedatt--items--sign_in"></a>
### Nested Schema for `items.sign_in`

- `identity_provider` (String) The display name of the identity provider (`IDENTITY_PROVIDER`), otherwise null.
- `method` (String) `IDENTITY_PROVIDER`: with their company account through `identityProvider`; `PASSWORD`: with their email and a password (an invitation's is temporary, emailed to them); `NONE`: no provider covers their email domain and the policy allows identity providers only, so they cannot sign in. One of `IDENTITY_PROVIDER`, `PASSWORD`, `NONE`.

