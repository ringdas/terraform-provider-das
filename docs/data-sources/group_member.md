---
page_title: "group_member Data Source - das"
subcategory: "Organization"
description: |-
  A person's membership of a group. Read by its identity.
---

# das_group_member (Data Source)

A person's membership of a group. Read by its identity.

API: `GET /api/groups/{groupId}/members/{userId}`.

## Example Usage

```terraform
data "das_group_member" "example" {
  group_id = "<groupId>"
  user_id  = "<userId>"
}
```

## Schema

### Required

- `group_id` (String) The group, `grp_<ULID>`.
- `user_id` (String) The user, `usr_<ULID>`.

### Read-Only

- `created_at` (String) When the member was added, ISO-8601 UTC.
- `created_by` (String) Actor who added the member.
- `display_name` (String) The member's name.
- `email` (String) The member's email.
- `status` (String) The member's status. One of `INVITED`, `ACTIVE`, `SUSPENDED`.

