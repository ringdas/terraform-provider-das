---
page_title: "tag_manager_policy_mapping Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`). Its key, `from`, `to` and `cell` never change: a change replaces it (reference `das_tag_manager_policy_catalogue.<name>.key` so the catalogues are created first). `cell` is the type of a cell's value in the field language of DAS Policies; without it a cell holds no value, with it each cell's `value_json` is one `jsonencode()` of its value. `rows` and their `cells` are lists kept in your order and replaced whole. A mapping is the template a `map` field of a configuration starts from: its rows change no configuration. Destroy deletes it, refused while a policy type that is not archived reads it in a field; to change `cell`, create a new mapping with a new key, point the type's field at it, then remove the old one. A mapping without a cell type is such a new mapping too: removing `cell` from the Terraform configuration (`.tf`) never turns this mapping into one.
---

# das_tag_manager_policy_mapping (Resource)

A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`). Its key, `from`, `to` and `cell` never change: a change replaces it (reference `das_tag_manager_policy_catalogue.<name>.key` so the catalogues are created first). `cell` is the type of a cell's value in the field language of DAS Policies; without it a cell holds no value, with it each cell's `value_json` is one `jsonencode()` of its value. `rows` and their `cells` are lists kept in your order and replaced whole. A mapping is the template a `map` field of a configuration starts from: its rows change no configuration. Destroy deletes it, refused while a policy type that is not archived reads it in a field; to change `cell`, create a new mapping with a new key, point the type's field at it, then remove the old one. A mapping without a cell type is such a new mapping too: removing `cell` from the Terraform configuration (`.tf`) never turns this mapping into one.

API: `POST /api/policy-mappings` creates, `GET /api/policy-mappings/{mappingId}` reads, `PATCH /api/policy-mappings/{mappingId}` updates, `DELETE /api/policy-mappings/{mappingId}` deletes.

## Example Usage

```terraform
resource "das_tag_manager_policy_mapping" "example" {
  from  = "sizes_app"
  key   = "size_map_app"
  label = "Size map APP"
  rows  = [{
    cells = [{
      to = "GAM"
    }, {
      to = "PREBID"
    }]
    key = "320x50"
  }, {
    cells = [{
      to = "PREBID"
    }]
    key = "native"
  }]
  to    = "size_consumers"
}
```

## Schema

### Required

- `from` (String) The catalogue whose values are the rows' keys, by its key (GET /api/policy-catalogues; 400 when the Network has none). Never changed. Changing it replaces the resource.
- `key` (String) Key of the mapping, snake_case, 1-64 characters, unique in the Network and never changed: a `map` field reads it as `mapping.<key>`. Changing it replaces the resource.
- `label` (String) Name, 1-200 characters.
- `rows` (Attributes List) The rows, 0-1000, kept and answered in the order written. Validated against the catalogues' values as they are when the write runs; an invalid row is 400 at `rows.<i>...` with the field language's `code`.
- `to` (String) The catalogue whose values are the cells' `to`, by its key (400 when the Network has none). Never changed. Changing it replaces the resource.

### Optional

- `cell` (Dynamic) The type of each cell's value in the field language of DAS Policies, a definition of one value (`{ kind, min?, max?, unit?, options?, optionLabels?, maxLength? }`; kinds `string`, `int`, `decimal`, `bool`; no `list`, no `optionsSource`). Null or left out, a cell holds no value: the row's `to` values are what its key goes to. Never changed. An invalid definition is 400 at `cell.<path>` with the field language's `code`. Changing it replaces the resource.

### Read-Only

- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--rows"></a>
### Nested Schema for `rows`

- `cells` (Attributes List) The values of `to` the key goes to, one or more (400 `empty`), kept and answered in the order written.
- `key` (String) A value of the `from` catalogue (400 `not_in_options` otherwise), once in the mapping (400 `duplicate`).

<a id="nestedatt--rows--cells"></a>
### Nested Schema for `rows.cells`

- `to` (String) A value of the `to` catalogue (400 `not_in_options` otherwise), once per row (400 `duplicate`).
- `value_json` (String) The cell's value, of the mapping's `cell` type: required with a `cell` (400 `required`), refused without one (400 `unknown_key`); an invalid one is 400 at it with the field language's `code`. JSON-encoded (`jsondecode()` it).

## Import

```terraform
import {
  to = das_tag_manager_policy_mapping.example
  id = "<mappingId>"
}
```
