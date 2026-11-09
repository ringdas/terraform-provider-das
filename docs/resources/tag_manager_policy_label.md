---
page_title: "tag_manager_policy_label Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case. Destroy deletes it and removes it from every policy carrying it.
---

# das_tag_manager_policy_label (Resource)

A label of the credential's Network (the feature `policies`), written by an operator: a free-form tag its policies carry, its name unique in the Network regardless of case. Destroy deletes it and removes it from every policy carrying it.

API: `POST /api/policy-labels` creates, `GET /api/policy-labels/{labelId}` reads, `PATCH /api/policy-labels/{labelId}` updates, `DELETE /api/policy-labels/{labelId}` deletes.

## Example Usage

```terraform
resource "das_tag_manager_policy_label" "example" {
  name = "floorian"
}
```

## Schema

### Required

- `name` (String) The label's name, trimmed, 1-64 characters, unique in the Network regardless of case (409 `ALREADY_EXISTS`).

### Read-Only

- `id` (String) The id the API assigned.

## Import

```terraform
import {
  to = das_tag_manager_policy_label.example
  id = "<labelId>"
}
```
