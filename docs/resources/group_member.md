---
page_title: "group_member Resource - das"
subcategory: "Organization"
description: |-
  A person's membership of a group.
---

# das_group_member (Resource)

A person's membership of a group.

API: `PUT /api/groups/{groupId}/members/{userId}` creates, `GET /api/groups/{groupId}/members/{userId}` reads, `DELETE /api/groups/{groupId}/members/{userId}` deletes.

## Example Usage

```terraform
resource "das_group_member" "example" {
  group_id = "<groupId>"
  user_id  = "<userId>"
}
```

## Schema

### Required

- `group_id` (String) The group, `grp_<ULID>`. Changing it replaces the resource.
- `user_id` (String) The user, `usr_<ULID>`. Changing it replaces the resource.

### Read-Only

- `created_at` (String) When the member was added, ISO-8601 UTC.
- `created_by` (String) Actor who added the member.
- `display_name` (String) The member's name.
- `email` (String) The member's email.
- `id` (String) The identity, groupId/userId joined by `/` (the import id).
- `status` (String) The member's status. One of `INVITED`, `ACTIVE`, `SUSPENDED`.

## Import

```terraform
import {
  to = das_group_member.example
  id = "<groupId>/<userId>"
}
```
