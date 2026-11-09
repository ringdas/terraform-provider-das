---
page_title: "tag_manager_policy_target_dimensions Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_target_dimension, all pages: a target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin.
---

# das_tag_manager_policy_target_dimensions (Data Source)

Every tag_manager_policy_target_dimension, all pages: a target dimension of the credential's Network (the feature `policies`): what a policy's target names, with the definition of its values (`input`, the field language of DAS Policies) and its weight (1-10), written by a Network admin.

API: `GET /api/policy-target-dimensions`.

## Example Usage

```terraform
data "das_tag_manager_policy_target_dimensions" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `id` (String) Target dimension id, `dim_<ULID>`.
- `input_json` (String) The definition of its values in the field language of DAS Policies (`{ kind, list?, min?, max?, options?, optionsSource?, operators? }`). JSON-encoded (`jsondecode()` it).
- `key` (String) Key, snake_case, unique in the Network, never changed: what a target names it by.
- `label` (String) Name.
- `weight` (Number, integer) Weight of a condition on it when policies compete, 1-10.

