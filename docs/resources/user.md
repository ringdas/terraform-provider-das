---
page_title: "user Resource - das"
subcategory: "Organization"
description: |-
  A person of the Organization's directory, invited by email. They hold no role until one is bound. The API never deletes a person: destroy suspends them.
---

# das_user (Resource)

A person of the Organization's directory, invited by email. They hold no role until one is bound. The API never deletes a person: destroy suspends them.

API: `POST /api/users` creates, `GET /api/users/{userId}` reads, `PATCH /api/users/{userId}` updates. Destroy: People are suspended, never deleted (their audit trail stays).

## Example Usage

```terraform
resource "das_user" "example" {
  display_name = "example-display-name"
  email        = "someone@example.com"
}
```

## Schema

### Required

- `display_name` (String) Name shown in the UI.
- `email` (String) Their email; the first sign-in with it links their identity. It cannot change: the API offers no update and never deletes it.

### Optional

- `status` (String) `SUSPENDED` removes all access; `ACTIVE` restores it (an unlinked user goes back to `INVITED`). One of `ACTIVE`, `SUSPENDED`. Applied after create, by an update.

### Read-Only

- `created_at` (String) Invitation instant, ISO-8601 UTC.
- `created_by` (String) Actor who invited the user.
- `id` (String) The id the API assigned.
- `invitation` (Attributes) The invitation an admin sent and how far its delivery got; null for a person SCIM provisioned or the first owner.
- `linked` (Boolean) Whether the identity provider's identity is linked (the user signed in once).
- `sign_in` (Attributes) How this person signs in under the Organization's sign-in configuration now: an enabled identity provider whose email domains cover their address takes them to it; otherwise a password of the Organization's own sign-in while the sign-in policy allows it.
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--invitation"></a>
### Nested Schema for `invitation`

- `attempted_at` (String) The last attempt; null before the first.
- `completed_at` (String) When the delivery finished; null while it is pending.
- `emailed_at` (String) When the invitation email went out; null before, or when the installation sends no email.
- `error` (String) Why the last attempt did not finish (it is retried); null when it did.
- `identity_provider` (String) The identity provider the email named (`IDENTITY_PROVIDER`), otherwise null.
- `invited_at` (String) When, ISO-8601 UTC.
- `invited_by` (String) Actor who invited them, or invited them again.
- `method` (String) How the delivery let them in (`signIn.method` at the time); null until it acted. One of `IDENTITY_PROVIDER`, `PASSWORD`.
- `pool_invited_at` (String) When their account of the Organization's own sign-in was created and its email with a temporary password sent (`PASSWORD`); null before, or when they already had a password.

<a id="nestedatt--sign_in"></a>
### Nested Schema for `sign_in`

- `identity_provider` (String) The display name of the identity provider (`IDENTITY_PROVIDER`), otherwise null.
- `method` (String) `IDENTITY_PROVIDER`: with their company account through `identityProvider`; `PASSWORD`: with their email and a password (an invitation's is temporary, emailed to them); `NONE`: no provider covers their email domain and the policy allows identity providers only, so they cannot sign in. One of `IDENTITY_PROVIDER`, `PASSWORD`, `NONE`.

## Import

```terraform
import {
  to = das_user.example
  id = "<userId>"
}
```
