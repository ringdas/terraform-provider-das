---
page_title: "delivery_order_step_completion Action - das"
subcategory: "DAS Delivery"
description: |-
  Marks the current `WORK` step of an order done (the feature `order-workspace`). The credential must be a person's holding the step's role; the order moves to the next step, or is completed.
---

# das_delivery_order_step_completion (Action)

Marks the current `WORK` step of an order done (the feature `order-workspace`). The credential must be a person's holding the step's role; the order moves to the next step, or is completed. Requires Terraform 1.14 or later.

API: `POST /api/orders/{orderId}/step-completions`.

## Example Usage

```terraform
action "das_delivery_order_step_completion" "example" {
  config {
  comment  = "Brief agreed with the advertiser"
  order_id = "ord_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  step_key = "brief"
  }
}
```

Run it with `terraform apply -invoke=action.das_delivery_order_step_completion.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `order_id` (String) The order, `ord_<ULID>`. One of another Network is 404 `NOT_FOUND`.
- `step_key` (String) The key of the order's current step, a `WORK` step (as GET /api/orders/{orderId} shows it).

### Optional

- `comment` (String) A note for whoever takes the next step, 1-2000 characters.
- `step_revision` (Number, integer) The order's `stepRevision` the caller saw (GET /api/orders/{orderId}). When given and the order has moved since - another step taken, a rejection and the work redone - the write is 409 `STEP_NOT_CURRENT` instead of acting on a state the caller never saw.

