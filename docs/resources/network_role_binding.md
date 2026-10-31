---
page_title: "network_role_binding Resource - das"
subcategory: "Organization"
description: |-
  A role in the Network named by `network_id`, held by a person or a group, administered by the Organization.
---

# das_network_role_binding (Resource)

A role in the Network named by `network_id`, held by a person or a group, administered by the Organization.

API: `POST /api/networks/{networkId}/role-bindings` creates, `GET /api/networks/{networkId}/role-bindings/{bindingId}` reads, `DELETE /api/networks/{networkId}/role-bindings/{bindingId}` deletes.

## Example Usage

```terraform
resource "das_network_role_binding" "example" {
  group_name = "traffickers"
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
  role       = "OPERATOR"
}
```

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network. Changing it replaces the resource.
- `role` (String) A Network role. One of `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`. Changing it replaces the resource.

### Optional

- `email` (String) The user's email (a user of the directory). Changing it replaces the resource.
- `group_name` (String) The group's name; never for a role that owns the Organization or reads the audit trail (ORG_OWNER, ORG_AUDITOR, AUDITOR), which only a person holds. Changing it replaces the resource.

### Read-Only

- `created_at` (String) When the role was granted, ISO-8601 UTC.
- `created_by` (String) Actor who granted the role.
- `display_name` (String) The user's name; null for a group.
- `group_id` (String) The group holding it, `grp_<ULID>`; null for a user.
- `id` (String) The id the API assigned.
- `subject_type` (String) Whether a user or a group holds the role. One of `USER`, `GROUP`.
- `user_id` (String) The user holding it, `usr_<ULID>`; null for a group.

## Import

```terraform
import {
  to = das_network_role_binding.example
  id = "<networkId>/<bindingId>"
}
```
