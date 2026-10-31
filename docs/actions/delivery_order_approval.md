---
page_title: "delivery_order_approval Action - das"
subcategory: "DAS Delivery"
description: |-
  Approves or rejects the current `APPROVAL` step of an order (the feature `order-workspace`), as an approver who did not do the work before it. Never with an API key (P9): the credential must be an approver's session.
---

# das_delivery_order_approval (Action)

Approves or rejects the current `APPROVAL` step of an order (the feature `order-workspace`), as an approver who did not do the work before it. Never with an API key (P9): the credential must be an approver's session. Requires Terraform 1.14 or later.

API: `POST /api/orders/{orderId}/approvals`.

## Example Usage

```terraform
action "das_delivery_order_approval" "example" {
  config {
  order_id = "ord_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  outcome  = "APPROVED"
  step_key = "approve"
  }
}
```

Run it with `terraform apply -invoke=action.das_delivery_order_approval.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `order_id` (String) The order, `ord_<ULID>`. One of another Network is 404 `NOT_FOUND`.
- `outcome` (String) `APPROVED` moves the order to the next step; `REJECTED` sends it back to the nearest earlier `WORK` step. One of `APPROVED`, `REJECTED`.
- `step_key` (String) The key of the order's current step, an `APPROVAL` step (as GET /api/orders/{orderId} shows it).

### Optional

- `comment` (String) Why, 1-2000 characters; required with `REJECTED`.
- `step_revision` (Number, integer) The order's `stepRevision` the caller saw (GET /api/orders/{orderId}). When given and the order has moved since - another step taken, a rejection and the work redone - the write is 409 `STEP_NOT_CURRENT` instead of acting on a state the caller never saw.

