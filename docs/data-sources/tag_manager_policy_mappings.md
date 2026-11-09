---
page_title: "tag_manager_policy_mappings Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_mapping, all pages: a mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`).
---

# das_tag_manager_policy_mappings (Data Source)

Every tag_manager_policy_mapping, all pages: a mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`).

API: `GET /api/policy-mappings`.

## Example Usage

```terraform
data "das_tag_manager_policy_mappings" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `cell_json` (String) The type of each cell's value in the field language of DAS Policies (`{ kind, min?, max?, unit?, options?, ... }`); null: a cell holds no value. Never changed. JSON-encoded (`jsondecode()` it).
- `created_at` (String) When it was created, ISO-8601 UTC.
- `from` (String) The catalogue whose values are the rows' keys, by its key; never changed.
- `id` (String) Mapping id, `map_<ULID>`.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `mapping.<key>` names it by.
- `label` (String) Name.
- `rows` (Attributes List) The rows, in the order written.
- `to` (String) The catalogue whose values are the cells' `to`, by its key; never changed.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--items--rows"></a>
### Nested Schema for `items.rows`

- `cells` (Attributes List) The values of `to` the key goes to, in the order written.
- `key` (String) A value of the `from` catalogue.

<a id="nestedatt--items--rows--cells"></a>
### Nested Schema for `items.rows.cells`

- `to` (String) A value of the `to` catalogue.
- `value_json` (String) The cell's value, of the mapping's `cell` type; absent when the mapping has no `cell`. JSON-encoded (`jsondecode()` it).

