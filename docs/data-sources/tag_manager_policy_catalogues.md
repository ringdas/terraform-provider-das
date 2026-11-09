---
page_title: "tag_manager_policy_catalogues Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_catalogue, all pages: a value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings.
---

# das_tag_manager_policy_catalogues (Data Source)

Every tag_manager_policy_catalogue, all pages: a value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings.

API: `GET /api/policy-catalogues`.

## Example Usage

```terraform
data "das_tag_manager_policy_catalogues" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) Catalogue id, `cat_<ULID>`.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `catalogue.<key>` and a mapping name it by.
- `label` (String) Name.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `values` (Attributes List) The values, in the order written.

<a id="nestedatt--items--values"></a>
### Nested Schema for `items.values`

- `label` (String) How the value is shown; null: the value itself.
- `value` (String) The value, unique in the catalogue: what bodies, targets and mappings name.

