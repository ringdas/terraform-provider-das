---
page_title: "tag_manager_policy_catalogue Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings. Its key never changes. `values` is a list kept in your order and replaced whole, each value a plain text unique in the catalogue with an optional `label`. A value removed from `values` leaves the mappings on the catalogue with it (the rows it keys, the cells it is the `to` of, a row left without cells): their next plan shows that drift, to settle in their `rows`. A value a configuration, a policy's target or a field's `default` of a policy type names (any of them not archived) is not removed (409 IN_USE): apply them without it first, then the catalogue (two applies). Destroy deletes it, refused while a policy type that is not archived reads it in a field, a target dimension reads it in its input or a mapping has it as `from` or `to`.
---

# das_tag_manager_policy_catalogue (Resource)

A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions read as `catalogue.<key>`, and the sides of mappings. Its key never changes. `values` is a list kept in your order and replaced whole, each value a plain text unique in the catalogue with an optional `label`. A value removed from `values` leaves the mappings on the catalogue with it (the rows it keys, the cells it is the `to` of, a row left without cells): their next plan shows that drift, to settle in their `rows`. A value a configuration, a policy's target or a field's `default` of a policy type names (any of them not archived) is not removed (409 IN_USE): apply them without it first, then the catalogue (two applies). Destroy deletes it, refused while a policy type that is not archived reads it in a field, a target dimension reads it in its input or a mapping has it as `from` or `to`.

API: `POST /api/policy-catalogues` creates, `GET /api/policy-catalogues/{catalogueId}` reads, `PATCH /api/policy-catalogues/{catalogueId}` updates, `DELETE /api/policy-catalogues/{catalogueId}` deletes.

## Example Usage

```terraform
resource "das_tag_manager_policy_catalogue" "example" {
  key    = "size_consumers"
  label  = "Size consumers"
  values = [{
    value = "GAM"
  }, {
    label = "Prebid"
    value = "PREBID"
  }]
}
```

## Schema

### Required

- `key` (String) Key of the catalogue, snake_case, 1-64 characters, unique in the Network and never changed: fields and target dimensions read it as `catalogue.<key>`, mappings name it by it. Changing it replaces the resource.
- `label` (String) Name, 1-200 characters.
- `values` (Attributes List) The values, 0-1000, kept and answered in the order written.

### Read-Only

- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--values"></a>
### Nested Schema for `values`

- `label` (String) How the value is shown, 1-200 characters; null or left out, the value itself (null in the answer, so a read value is written back as it is).
- `value` (String) The value, 1-200 characters, unique in the catalogue (400 `duplicate`): bodies, targets and mappings name it by this text, so a new text is the removal of one value and the addition of another.

## Import

```terraform
import {
  to = das_tag_manager_policy_catalogue.example
  id = "<catalogueId>"
}
```
