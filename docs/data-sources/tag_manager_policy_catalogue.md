---
page_title: "tag_manager_policy_catalogue Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings. Read by its identity.
---

# das_tag_manager_policy_catalogue (Data Source)

A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings. Read by its identity.

API: `GET /api/policy-catalogues/{catalogueId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_catalogue" "example" {
  id = "<catalogueId>"
}
```

## Schema

### Required

- `id` (String) The catalogue, `cat_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null while it is not.
- `archived_values` (Attributes List) The values a write of `values` left out, read-only: mapping rows may still name them, and writing one's text or id in `values` again restores it.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `catalogue.<key>` and a mapping name it by.
- `label` (String) Name.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `values` (Attributes List) The active values, in the order written.

<a id="nestedatt--archived_values"></a>
### Nested Schema for `archived_values`

- `archived_at` (String) When a write of `values` left it out, ISO-8601 UTC.
- `id` (String) Value id, `catval_<ULID>`.
- `label` (String) How the value is shown; null: the text itself.
- `value` (String) The value's text.

<a id="nestedatt--values"></a>
### Nested Schema for `values`

- `id` (String) Value id, `catval_<ULID>`: never changes; mapping rows point at it, so they follow a change of its text.
- `label` (String) How the value is shown; null: the text itself.
- `value` (String) The value's text, unique in the catalogue with its archived values.

