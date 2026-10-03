---
page_title: "group Resource - das"
subcategory: "Organization"
description: |-
  A group of people of the Organization's directory; role bindings name it.
---

# das_group (Resource)

A group of people of the Organization's directory; role bindings name it.

API: `POST /api/groups` creates, `GET /api/groups/{groupId}` reads, `PATCH /api/groups/{groupId}` updates, `DELETE /api/groups/{groupId}` deletes.

## Example Usage

```terraform
resource "das_group" "example" {
  name = "example-name"
}
```

## Schema

### Required

- `name` (String) Unique name: 1-100 letters, digits, space, dot, underscore or hyphen.

### Optional

- `description` (String) What the group is for (default empty).
- `external_id` (String) The identity provider's group to map, or null.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who created the group.
- `id` (String) The id the API assigned.
- `member_count` (Number, integer) Members in DAS, added by an admin or by SCIM (members only through a sign-in's groups claim are not counted).
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `updated_at` (String) Last change, ISO-8601 UTC.

## Import

```terraform
import {
  to = das_group.example
  id = "<groupId>"
}
```
