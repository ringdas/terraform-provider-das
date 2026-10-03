---
page_title: "network_role_binding Data Source - das"
subcategory: "Organization"
description: |-
  A role in the Network named by `network_id`, held by a person or a group, administered by the Organization. Read by its identity.
---

# das_network_role_binding (Data Source)

A role in the Network named by `network_id`, held by a person or a group, administered by the Organization. Read by its identity.

API: `GET /api/networks/{networkId}/role-bindings/{bindingId}`.

## Example Usage

```terraform
data "das_network_role_binding" "example" {
  network_id = "<networkId>"
  id         = "<bindingId>"
}
```

## Schema

### Required

- `id` (String) The Network role binding, `rlb_<ULID>`.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `created_at` (String) When the role was granted, ISO-8601 UTC.
- `created_by` (String) Actor who granted the role.
- `display_name` (String) The user's name; null for a group.
- `email` (String) The user's email; null for a group.
- `group_id` (String) The group holding it, `grp_<ULID>`; null for a user.
- `group_name` (String) The group's name; null for a user.
- `role` (String) The role (GET /api/roles describes each). One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`, `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`.
- `subject_type` (String) Whether a user or a group holds the role. One of `USER`, `GROUP`.
- `user_id` (String) The user holding it, `usr_<ULID>`; null for a group.

