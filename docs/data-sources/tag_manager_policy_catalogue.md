---
page_title: "tag_manager_policy_catalogue Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings. Read by its identity.
---

# das_tag_manager_policy_catalogue (Data Source)

A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings. Read by its identity.

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

- `created_at` (String) When it was created, ISO-8601 UTC.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `catalogue.<key>` and a mapping name it by.
- `label` (String) Name.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `values` (Attributes List) The values, in the order written.

<a id="nestedatt--values"></a>
### Nested Schema for `values`

- `label` (String) How the value is shown; null: the value itself.
- `value` (String) The value, unique in the catalogue: what bodies, targets and mappings name.

