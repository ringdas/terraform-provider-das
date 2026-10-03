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
- `linked` (Boolean) Whether the identity provider's identity is linked (the user signed in once).
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `updated_at` (String) Last change, ISO-8601 UTC.

## Import

```terraform
import {
  to = das_user.example
  id = "<userId>"
}
```
