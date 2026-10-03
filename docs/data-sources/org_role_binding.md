---
page_title: "org_role_binding Data Source - das"
subcategory: "Organization"
description: |-
  An Organization-layer role held by a person (by email) or a group (by name). Read by its identity.
---

# das_org_role_binding (Data Source)

An Organization-layer role held by a person (by email) or a group (by name). Read by its identity.

API: `GET /api/org-role-bindings/{bindingId}`.

## Example Usage

```terraform
data "das_org_role_binding" "example" {
  id = "<bindingId>"
}
```

## Schema

### Required

- `id` (String) The Organization role binding, `orb_<ULID>`.

### Read-Only

- `created_at` (String) When the role was granted, ISO-8601 UTC.
- `created_by` (String) Actor who granted the role.
- `display_name` (String) The user's name; null for a group.
- `email` (String) The user's email; null for a group.
- `group_id` (String) The group holding it, `grp_<ULID>`; null for a user.
- `group_name` (String) The group's name; null for a user.
- `network_id` (String) The Network the role is held in; null on the Organization layer.
- `role` (String) The role (GET /api/roles describes each). One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`, `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`.
- `subject_type` (String) Whether a user or a group holds the role. One of `USER`, `GROUP`.
- `user_id` (String) The user holding it, `usr_<ULID>`; null for a group.

