---
page_title: "tag_manager_policy_catalogues Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_catalogue, all pages: a value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings.
---

# das_tag_manager_policy_catalogues (Data Source)

Every tag_manager_policy_catalogue, all pages: a value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings.

API: `GET /api/policy-catalogues`.

## Example Usage

```terraform
data "das_tag_manager_policy_catalogues" "example" {
}
```

## Schema

### Optional

- `include_archived` (Boolean) `true`: the archived ones too, each with its `archivedAt`. `false` or left out: only those not archived.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `archived_at` (String) When it was archived, ISO-8601 UTC; null while it is not.
- `archived_values` (Attributes List) The values a write of `values` left out, read-only: mapping rows may still name them, and writing one's text or id in `values` again restores it.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) Catalogue id, `cat_<ULID>`.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `catalogue.<key>` and a mapping name it by.
- `label` (String) Name.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `values` (Attributes List) The active values, in the order written.

<a id="nestedatt--items--archived_values"></a>
### Nested Schema for `items.archived_values`

- `archived_at` (String) When a write of `values` left it out, ISO-8601 UTC.
- `id` (String) Value id, `catval_<ULID>`.
- `label` (String) How the value is shown; null: the text itself.
- `value` (String) The value's text.

<a id="nestedatt--items--values"></a>
### Nested Schema for `items.values`

- `id` (String) Value id, `catval_<ULID>`: never changes; mapping rows point at it, so they follow a change of its text.
- `label` (String) How the value is shown; null: the text itself.
- `value` (String) The value's text, unique in the catalogue with its archived values.

