---
page_title: "tag_manager_policy_option_sources Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  The sources a field of a policy type may read in `optionsSource` (the feature `policies`), for one `kind` of field: `string` the Network's catalogues (`catalogue.<key>`), `map` its mappings (`mapping.<key>`), each with its count and its first three values. Set `kind`: the API requires it (400 at read otherwise).
---

# das_tag_manager_policy_option_sources (Data Source)

The sources a field of a policy type may read in `optionsSource` (the feature `policies`), for one `kind` of field: `string` the Network's catalogues (`catalogue.<key>`), `map` its mappings (`mapping.<key>`), each with its count and its first three values. Set `kind`: the API requires it (400 at read otherwise).

API: `GET /api/policy-option-sources`.

## Example Usage

```terraform
data "das_tag_manager_policy_option_sources" "example" {
  kind = "string"
}
```

## Schema

### Optional

- `kind` (String) The kind of field the sources are for: `string` (the catalogues as `catalogue.<key>`) or `map` (the mappings as `mapping.<key>`). Required: left out or another value is 400 `VALIDATION_FAILED` at `kind`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `count` (Number, integer) How many values it holds: a catalogue's values, a mapping's rows.
- `key` (String) What a field names it by in `optionsSource`: `catalogue.<key>` or `mapping.<key>`.
- `label` (String) Its name: the catalogue's or mapping's label.
- `sample` (List of String) Its first three values (a mapping's: its rows' keys), in the source's own order.

