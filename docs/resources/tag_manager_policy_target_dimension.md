---
page_title: "tag_manager_policy_target_dimension Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin. Its key never changes. Which policy types may target it is each type's `dimensions`. Destroy deletes it (an archived policy type loses it), refused while a policy type that is not archived has it among its dimensions or a target names it.
---

# das_tag_manager_policy_target_dimension (Resource)

A target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin. Its key never changes. Which policy types may target it is each type's `dimensions`. Destroy deletes it (an archived policy type loses it), refused while a policy type that is not archived has it among its dimensions or a target names it.

API: `POST /api/policy-target-dimensions` creates, `GET /api/policy-target-dimensions/{dimensionId}` reads, `PATCH /api/policy-target-dimensions/{dimensionId}` updates, `DELETE /api/policy-target-dimensions/{dimensionId}` deletes.

## Example Usage

```terraform
resource "das_tag_manager_policy_target_dimension" "example" {
  input  = {
    kind = "int"
    max = 99
    min = 0
  }
  key    = "seg_ab"
  label  = "A/B segment"
  weight = 5
}
```

## Schema

### Required

- `input` (Dynamic) The definition of the dimension's values in the field language of DAS Policies, without id and label (`{ kind, list?, min?, max?, options?, optionsSource?, operators? }`; kinds `string`, `int`, `decimal`, `bool`; `operators` lists the operators its target rows take, by default every one on a single `int` or `decimal`, `IN` and `NOT_IN` otherwise). An invalid definition is 400 at `input.<path>` with the field language's `code`.
- `key` (String) Key of the dimension, snake_case, 1-64 characters, unique in the Network and never changed: the request carries the value under it. Changing it replaces the resource.
- `label` (String) Name, 1-200 characters.

### Optional

- `weight` (Number, integer) Weight of a condition on the dimension when policies compete (decision D31), 1-10; default 5.

### Read-Only

- `id` (String) The id the API assigned.

## Import

```terraform
import {
  to = das_tag_manager_policy_target_dimension.example
  id = "<dimensionId>"
}
```
