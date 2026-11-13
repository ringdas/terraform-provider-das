---
page_title: "tag_manager_policy_mapping Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`), each value by its text (stored by the value's id, so a value whose text changes stays in the rows). Read by its identity.
---

# das_tag_manager_policy_mapping (Data Source)

A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`), each value by its text (stored by the value's id, so a value whose text changes stays in the rows). Read by its identity.

API: `GET /api/policy-mappings/{mappingId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_mapping" "example" {
  id = "<mappingId>"
}
```

## Schema

### Required

- `id` (String) The mapping, `map_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null while it is not.
- `cell` (Dynamic) The type of each cell's value in the field language of DAS Policies (`{ kind, min?, max?, unit?, options?, ... }`); null: a cell holds no value. Never changed.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `from` (String) The catalogue whose values are the rows' keys, by its key; never changed.
- `key` (String) Key, snake_case, unique in the Network, never changed: what `mapping.<key>` names it by.
- `label` (String) Name.
- `rows` (Attributes List) The rows, in the order written.
- `to` (String) The catalogue whose values are the cells' `to`, by its key; never changed.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--rows"></a>
### Nested Schema for `rows`

- `cells` (Attributes List) The values of `to` the key goes to, in the order written.
- `key` (String) A value of the `from` catalogue, by its text as it is now (it may be archived).

<a id="nestedatt--rows--cells"></a>
### Nested Schema for `rows.cells`

- `to` (String) A value of the `to` catalogue, by its text as it is now (it may be archived).
- `value_json` (String) The cell's value, of the mapping's `cell` type; absent when the mapping has no `cell`. JSON-encoded (`jsondecode()` it).

