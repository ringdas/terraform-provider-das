---
page_title: "tag_manager_policy_labels Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_label, all pages: a label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case.
---

# das_tag_manager_policy_labels (Data Source)

Every tag_manager_policy_label, all pages: a label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case.

API: `GET /api/policy-labels`.

## Example Usage

```terraform
data "das_tag_manager_policy_labels" "example" {
}
```

## Schema

### Optional

- `q` (String) Only labels whose name contains this text, regardless of case (`%` and `_` are plain characters), 1-64 characters.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `id` (String) Label id, `label_<ULID>`.
- `name` (String) Name, unique in the Network regardless of case.

