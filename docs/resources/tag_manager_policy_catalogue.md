---
page_title: "tag_manager_policy_catalogue Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings. Its key never changes. `values` is the list of active values, kept in your order and replaced whole, each a plain text unique in the catalogue (archived values included) with an optional `label` and its `id`, which the API gives it: write a value's `id` (from the state after the first apply) to change its text and keep it, so the mappings' rows follow it; a value without `id` is the one with that text, or a new one. Any change of the catalogue, its label included, plans every `id` you did not write as known after apply. A value left out is archived: `archived_values` lists it, the mappings' rows keep it, and writing its text again restores it. Configurations, policies and defaults hold their own copies of values: a change here changes none of them. Destroy archives the catalogue (P11), refused while a policy type that is not archived reads it in a field, a target dimension reads it in its input or a mapping that is not archived has it as `from` or `to`; its key stays taken, so a new catalogue with it is 409: import the archived one and set `archived = false` to restore it.
---

# das_tag_manager_policy_catalogue (Resource)

A value catalogue of the credential's Network (the feature `policies`), written by an operator: the values fields of policy types and target dimensions offer as `catalogue.<key>`, and the sides of mappings. Its key never changes. `values` is the list of active values, kept in your order and replaced whole, each a plain text unique in the catalogue (archived values included) with an optional `label` and its `id`, which the API gives it: write a value's `id` (from the state after the first apply) to change its text and keep it, so the mappings' rows follow it; a value without `id` is the one with that text, or a new one. Any change of the catalogue, its label included, plans every `id` you did not write as known after apply. A value left out is archived: `archived_values` lists it, the mappings' rows keep it, and writing its text again restores it. Configurations, policies and defaults hold their own copies of values: a change here changes none of them. Destroy archives the catalogue (P11), refused while a policy type that is not archived reads it in a field, a target dimension reads it in its input or a mapping that is not archived has it as `from` or `to`; its key stays taken, so a new catalogue with it is 409: import the archived one and set `archived = false` to restore it.

API: `POST /api/policy-catalogues` creates, `GET /api/policy-catalogues/{catalogueId}` reads, `PATCH /api/policy-catalogues/{catalogueId}` updates. Destroy: A catalogue is archivable (P11): it is never deleted, so destroy archives it (refused while a policy type, a target dimension or a mapping that is not archived reads it).

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

- `key` (String) Key of the catalogue, snake_case, 1-64 characters, unique in the Network and never changed: fields and target dimensions read it as `catalogue.<key>`, mappings name it by it. An archived catalogue keeps its key: 409 `ALREADY_EXISTS` naming its id (restore it instead). It cannot change: the API offers no update and never deletes it.
- `label` (String) Name, 1-200 characters.
- `values` (Attributes List) The active values, 0-1000, kept and answered in the order written.

### Optional

- `archived` (Boolean) `true` archives the catalogue (`archivedAt` set by the server): refused with 409 `IN_USE` while a policy type that is not `ARCHIVED` reads it in a field, a target dimension in its input or a mapping that is not archived has it as `from` or `to`. `false` restores an archived catalogue; an archived one takes no other change (409 `ARCHIVED`), and archiving it again changes nothing. Applied after create, by an update.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null while it is not.
- `archived_values` (Attributes List) The values a write of `values` left out, read-only: mapping rows may still name them, and writing one's text or id in `values` again restores it.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.

<a id="nestedatt--values"></a>
### Nested Schema for `values`

- `id` (String) The value's id, `catval_<ULID>`, as GET answers it: with it, this value takes the text and label written (mapping rows follow it); without it, the value with this text is meant, active or archived (an archived one is restored), or a new value with a new id. An id the catalogue has not, or one listed twice, is 400; on create, leave it out (400).
- `label` (String) How the value is shown, 1-200 characters; null or left out, the value itself (null in the answer, so a read value is written back as it is).
- `value` (String) The value's text, 1-200 characters, unique in the catalogue with its archived values (400 `duplicate`). Configurations, policies and defaults hold their own copy of a text: a change here changes none of them, and none is checked against the catalogue.

<a id="nestedatt--archived_values"></a>
### Nested Schema for `archived_values`

- `archived_at` (String) When a write of `values` left it out, ISO-8601 UTC.
- `id` (String) Value id, `catval_<ULID>`.
- `label` (String) How the value is shown; null: the text itself.
- `value` (String) The value's text.

## Import

```terraform
import {
  to = das_tag_manager_policy_catalogue.example
  id = "<catalogueId>"
}
```
