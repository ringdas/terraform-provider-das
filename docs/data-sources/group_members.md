---
page_title: "group_members Data Source - das"
subcategory: "Organization"
description: |-
  Every group_member, all pages: a person's membership of a group.
---

# das_group_members (Data Source)

Every group_member, all pages: a person's membership of a group.

API: `GET /api/groups/{groupId}/members`.

## Example Usage

```terraform
data "das_group_members" "example" {
  group_id = "<groupId>"
}
```

## Schema

### Required

- `group_id` (String) The group, `grp_<ULID>`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) When the member was added, ISO-8601 UTC.
- `created_by` (String) Actor who added the member.
- `display_name` (String) The member's name.
- `email` (String) The member's email.
- `group_id` (String) The group.
- `status` (String) The member's status. One of `INVITED`, `ACTIVE`, `SUSPENDED`.
- `user_id` (String) The member.

