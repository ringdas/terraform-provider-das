---
page_title: "group Data Source - das"
subcategory: "Organization"
description: |-
  A group of people of the Organization's directory; role bindings name it. Read by its identity.
---

# das_group (Data Source)

A group of people of the Organization's directory; role bindings name it. Read by its identity.

API: `GET /api/groups/{groupId}`.

## Example Usage

```terraform
data "das_group" "example" {
  id = "<groupId>"
}
```

## Schema

### Required

- `id` (String) The group, `grp_<ULID>`.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who created the group.
- `description` (String) What the group is for.
- `external_id` (String) The identity provider's group (Entra object id or SAML group value): its members are members for their session; null = managed in DAS only.
- `member_count` (Number, integer) Members in DAS, added by an admin or by SCIM (members only through a sign-in's groups claim are not counted).
- `name` (String) Unique name in the Organization.
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `updated_at` (String) Last change, ISO-8601 UTC.

