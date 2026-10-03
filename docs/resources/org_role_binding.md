---
page_title: "org_role_binding Resource - das"
subcategory: "Organization"
description: |-
  An Organization-layer role held by a person (by email) or a group (by name).
---

# das_org_role_binding (Resource)

An Organization-layer role held by a person (by email) or a group (by name).

API: `POST /api/org-role-bindings` creates, `GET /api/org-role-bindings/{bindingId}` reads, `DELETE /api/org-role-bindings/{bindingId}` deletes.

## Example Usage

```terraform
resource "das_org_role_binding" "example" {
  email = "someone@example.com"
  role  = "ORG_ADMIN"
}
```

## Schema

### Required

- `role` (String) An Organization role. One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`. Changing it replaces the resource.

### Optional

- `email` (String) The user's email (a user of the directory). Changing it replaces the resource.
- `group_name` (String) The group's name. Changing it replaces the resource.

### Read-Only

- `created_at` (String) When the role was granted, ISO-8601 UTC.
- `created_by` (String) Actor who granted the role.
- `display_name` (String) The user's name; null for a group.
- `group_id` (String) The group holding it, `grp_<ULID>`; null for a user.
- `id` (String) The id the API assigned.
- `network_id` (String) The Network the role is held in; null on the Organization layer.
- `subject_type` (String) Whether a user or a group holds the role. One of `USER`, `GROUP`.
- `user_id` (String) The user holding it, `usr_<ULID>`; null for a group.

## Import

```terraform
import {
  to = das_org_role_binding.example
  id = "<bindingId>"
}
```
