---
page_title: "groups Data Source - das"
subcategory: "Organization"
description: |-
  Every group, all pages: a group of people of the Organization's directory; role bindings name it.
---

# das_groups (Data Source)

Every group, all pages: a group of people of the Organization's directory; role bindings name it.

API: `GET /api/groups`.

## Example Usage

```terraform
data "das_groups" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who created the group.
- `description` (String) What the group is for.
- `external_id` (String) The identity provider's group (Entra object id or SAML group value): its members are members for their session; null = managed in DAS only.
- `id` (String) Group id, `grp_<ULID>`.
- `member_count` (Number, integer) Members in DAS, added by an admin or by SCIM (members only through a sign-in's groups claim are not counted).
- `name` (String) Unique name in the Organization.
- `source` (String) Who writes it: `DAS` (an admin) or `SCIM` (the identity provider's provisioning, which created it or last wrote it). One of `DAS`, `SCIM`.
- `updated_at` (String) Last change, ISO-8601 UTC.

