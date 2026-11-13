---
page_title: "tag_manager_policy_mapping Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`), each value by its text (stored by the value's id, so a value whose text changes stays in the rows). Its key, `from`, `to` and `cell` never change, and a mapping is never deleted, so a change of one is an error at plan (reference `das_tag_manager_policy_catalogue.<name>.key` so the catalogues are created first). `cell` is the type of a cell's value in the field language of DAS Policies; without it a cell holds no value, with it each cell's `value_json` is one `jsonencode()` of its value. `rows` and their `cells` are lists kept in your order and replaced whole. A mapping is the template a `map` field of a configuration starts from: its rows change no configuration. Destroy archives it (P11), refused while a policy type that is not archived reads it in a field; its key stays taken, so a new mapping with it is 409: import the archived one and set `archived = false` to restore it. A mapping without a cell type is a new mapping too: removing `cell` from the Terraform configuration (`.tf`) never turns this mapping into one.
---

# das_tag_manager_policy_mapping (Resource)

A mapping between two value catalogues of the credential's Network (the feature `policies`), written by an operator: for each value of `from` (a row's `key`) the values of `to` it goes to (its `cells`), each value by its text (stored by the value's id, so a value whose text changes stays in the rows). Its key, `from`, `to` and `cell` never change, and a mapping is never deleted, so a change of one is an error at plan (reference `das_tag_manager_policy_catalogue.<name>.key` so the catalogues are created first). `cell` is the type of a cell's value in the field language of DAS Policies; without it a cell holds no value, with it each cell's `value_json` is one `jsonencode()` of its value. `rows` and their `cells` are lists kept in your order and replaced whole. A mapping is the template a `map` field of a configuration starts from: its rows change no configuration. Destroy archives it (P11), refused while a policy type that is not archived reads it in a field; its key stays taken, so a new mapping with it is 409: import the archived one and set `archived = false` to restore it. A mapping without a cell type is a new mapping too: removing `cell` from the Terraform configuration (`.tf`) never turns this mapping into one.

API: `POST /api/policy-mappings` creates, `GET /api/policy-mappings/{mappingId}` reads, `PATCH /api/policy-mappings/{mappingId}` updates. Destroy: A mapping is archivable (P11): it is never deleted, so destroy archives it (refused while a policy type that is not archived reads it). To change its key, from, to or cell, create a new mapping with a new key, point the policy types' fields at it, then archive the old one by removing its resource.

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

- `from` (String) The catalogue whose values are the rows' keys, by its key (GET /api/policy-catalogues; 400 when the Network has none or it is archived). Never changed. It cannot change: the API offers no update and never deletes it.
- `key` (String) Key of the mapping, snake_case, 1-64 characters, unique in the Network and never changed: a `map` field reads it as `mapping.<key>`. An archived mapping keeps its key: 409 `ALREADY_EXISTS` naming its id (restore it instead). It cannot change: the API offers no update and never deletes it.
- `label` (String) Name, 1-200 characters.
- `rows` (Attributes List) The rows, 0-1000, kept and answered in the order written. Validated against the catalogues' values, archived ones included, as they are when the write runs; an invalid row is 400 at `rows.<i>...` with the field language's `code`.
- `to` (String) The catalogue whose values are the cells' `to`, by its key (400 when the Network has none or it is archived). Never changed. It cannot change: the API offers no update and never deletes it.

### Optional

- `archived` (Boolean) `true` archives the mapping (`archivedAt` set by the server): refused with 409 `IN_USE` while a policy type that is not `ARCHIVED` reads it in a field. `false` restores an archived mapping (400 while one of its catalogues is archived: restore that first); an archived one takes no other change (409 `ARCHIVED`), and archiving it again changes nothing. Applied after create, by an update.
- `cell` (Dynamic) The type of each cell's value in the field language of DAS Policies, a definition of one value (`{ kind, min?, max?, unit?, options?, optionLabels?, maxLength? }`; kinds `string`, `int`, `decimal`, `bool`; no `list`, no `optionsSource`). Null or left out, a cell holds no value: the row's `to` values are what its key goes to. Never changed. An invalid definition is 400 at `cell.<path>` with the field language's `code`. It cannot change: the API offers no update and never deletes it.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null while it is not.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--rows"></a>
### Nested Schema for `rows`

- `cells` (Attributes List) The values of `to` the key goes to, one or more (400 `empty`), kept and answered in the order written.
- `key` (String) A value of the `from` catalogue by its text, active or archived (400 `not_in_options` otherwise), once in the mapping (400 `duplicate`); stored by the value's id, so it follows a change of the value's text.

<a id="nestedatt--rows--cells"></a>
### Nested Schema for `rows.cells`

- `to` (String) A value of the `to` catalogue by its text, active or archived (400 `not_in_options` otherwise), once per row (400 `duplicate`); stored by the value's id, so it follows a change of the value's text.
- `value_json` (String) The cell's value, of the mapping's `cell` type: required with a `cell` (400 `required`), refused without one (400 `unknown_key`); an invalid one is 400 at it with the field language's `code`. JSON-encoded (`jsondecode()` it).

## Import

```terraform
import {
  to = das_tag_manager_policy_mapping.example
  id = "<mappingId>"
}
```
