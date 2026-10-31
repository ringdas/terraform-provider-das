---
page_title: "delivery_order_step_completions Data Source - das"
subcategory: "DAS Delivery"
description: |-
  The step history of an order (the feature `order-workspace`): who did or decided which step, when, and how, oldest first (every page followed).
---

# das_delivery_order_step_completions (Data Source)

The step history of an order (the feature `order-workspace`): who did or decided which step, when, and how, oldest first (every page followed).

API: `GET /api/orders/{orderId}/step-completions`.

## Example Usage

```terraform
data "das_delivery_order_step_completions" "example" {
  order_id = "<orderId>"
}
```

## Schema

### Required

- `order_id` (String) The order, `ord_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `actor` (String) Audit-grammar actor of who did or decided the step, `human:<email>`.
- `comment` (String) The note or the reason given; null when none.
- `created_at` (String) When, ISO-8601 UTC.
- `id` (String) History entry id, `ordsc_<ULID>`.
- `order_id` (String) The order, `ord_<ULID>`.
- `outcome` (String) `DONE` for the work of a `WORK` step; `APPROVED` or `REJECTED` for an `APPROVAL` step. One of `DONE`, `APPROVED`, `REJECTED`.
- `role` (String) The role the step belongs to. One of `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`.
- `step_key` (String) Key of the step.
- `step_kind` (String) Kind of the step. One of `WORK`, `APPROVAL`.

