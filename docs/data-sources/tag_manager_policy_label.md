---
page_title: "tag_manager_policy_label Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case. Read by its identity.
---

# das_tag_manager_policy_label (Data Source)

A label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case. Read by its identity.

API: `GET /api/policy-labels/{labelId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_label" "example" {
  id = "<labelId>"
}
```

## Schema

### Required

- `id` (String) The label, `label_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `name` (String) Name, unique in the Network regardless of case.

