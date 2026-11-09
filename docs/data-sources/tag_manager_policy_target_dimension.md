---
page_title: "tag_manager_policy_target_dimension Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin. Read by its identity.
---

# das_tag_manager_policy_target_dimension (Data Source)

A target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin. Read by its identity.

API: `GET /api/policy-target-dimensions/{dimensionId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_target_dimension" "example" {
  id = "<dimensionId>"
}
```

## Schema

### Required

- `id` (String) The target dimension, `dim_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `input` (Dynamic) The definition of its values in the field language of DAS Policies (`{ kind, list?, min?, max?, options?, optionsSource?, operators? }`).
- `key` (String) Key, snake_case, unique in the Network, never changed: what a target names it by.
- `label` (String) Name.
- `weight` (Number, integer) Weight of a condition on it when policies compete, 1-10.

