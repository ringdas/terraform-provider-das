---
page_title: "tag_manager_policy_duplicate Action - das"
subcategory: "DAS Tag Manager"
description: |-
  Duplicates a policy (the feature `policies`), as an operator (D5): a new policy with the same type, context, period, description, labels, target and inventory, attaching the same configurations with their sub-targets and inventory; `<name> (copy)` and `INACTIVE` unless given. Every invocation is a new policy, so it is an action, not a resource: import the copy to manage it.
---

# das_tag_manager_policy_duplicate (Action)

Duplicates a policy (the feature `policies`), as an operator (D5): a new policy with the same type, context, period, description, labels, target and inventory, attaching the same configurations with their sub-targets and inventory; `<name> (copy)` and `INACTIVE` unless given. Every invocation is a new policy, so it is an action, not a resource: import the copy to manage it. Requires Terraform 1.14 or later.

API: `POST /api/policies/{policyId}/duplicates`.

## Example Usage

```terraform
action "das_tag_manager_policy_duplicate" "example" {
  config {
  name      = "Fast timeout, autumn test"
  policy_id = "pol_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  }
}
```

Run it with `terraform apply -invoke=action.das_tag_manager_policy_duplicate.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `policy_id` (String) The policy, `pol_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Optional

- `name` (String) The copy's name; default `<name> (copy)`.
- `status` (String) The copy's status; default `INACTIVE`. One of `ACTIVE`, `INACTIVE`.

