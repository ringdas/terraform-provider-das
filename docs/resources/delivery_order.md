---
page_title: "delivery_order Resource - das"
subcategory: "DAS Delivery"
description: |-
  An order of the credential's Network (the feature `order-workspace`), delivered with an `ENABLED` ad product: it pins the product's current version and starts on its first step. Its name and flight change in place; its ad product, goal and billing model are refused at plan once created. Steps move with the actions `das_delivery_order_step_completion` and `das_delivery_order_approval`. Destroy cancels it (`CANCELLED`): an order is archived, never deleted; one already `COMPLETED` or `CANCELLED` is refused (409 `STEP_NOT_CURRENT`).
---

# das_delivery_order (Resource)

An order of the credential's Network (the feature `order-workspace`), delivered with an `ENABLED` ad product: it pins the product's current version and starts on its first step. Its name and flight change in place; its ad product, goal and billing model are refused at plan once created. Steps move with the actions `das_delivery_order_step_completion` and `das_delivery_order_approval`. Destroy cancels it (`CANCELLED`): an order is archived, never deleted; one already `COMPLETED` or `CANCELLED` is refused (409 `STEP_NOT_CURRENT`).

API: `POST /api/orders` creates, `GET /api/orders/{orderId}` reads, `PATCH /api/orders/{orderId}` updates. Destroy: An order is archivable (P11): it is never deleted, so destroy cancels it.

## Example Usage

```terraform
resource "das_delivery_order" "example" {
  ad_product_id = "adp_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  billing_model = "FLAT_FEE"
  end_at        = "2027-01-31T23:59:59Z"
  goal          = "AWARENESS"
  name          = "Homepage takeover, January"
  start_at      = "2027-01-01T00:00:00Z"
}
```

## Schema

### Required

- `ad_product_id` (String) The ad product the order is delivered with, `adp_<ULID>`: `ENABLED` and with a version. It cannot change: the API offers no update and never deletes it.
- `billing_model` (String) How the order is billed: one of the billing models of the ad product's current version. One of `CPM`, `VIEWABLE_CPM`, `CPC`, `CPV`, `CPA`, `FLAT_FEE`. It cannot change: the API offers no update and never deletes it.
- `end_at` (String) End of the order's flight, an ISO-8601 instant after startAt.
- `goal` (String) What the order aims for: one of the goals of the ad product's current version. One of `AWARENESS`, `REACH`, `TRAFFIC`, `ENGAGEMENT`, `VIDEO_VIEWS`, `LEADS`, `SALES`. It cannot change: the API offers no update and never deletes it.
- `name` (String) Name, 1-200 characters.
- `start_at` (String) Start of the order's flight, an ISO-8601 instant.

### Optional

- `status` (String) `CANCELLED` stops the order: it leaves every queue and takes no further step. One of `CANCELLED`. Applied after create, by an update.

### Read-Only

- `ad_product_version` (Number, integer) The version of the ad product the order follows, pinned when it was created.
- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `current_step` (Attributes) The step the order waits on; null once it is completed or cancelled.
- `id` (String) The id the API assigned.
- `sales_channel_id` (String) The sales channel the order came through, `sch_<ULID>`, from the credential that created it (a key of the channel); null when the Network's own staff created it.
- `step_revision` (Number, integer) How many step writes the order has had (0 when created). Send it back as `stepRevision` of a step write so it is refused (409 `STEP_NOT_CURRENT`) if the order moved meanwhile.
- `steps` (Attributes List) Every step of the pinned version, in order, with where each stands.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--current_step"></a>
### Nested Schema for `current_step`

- `key` (String) Key of the step in the pinned version.
- `kind` (String) `WORK`: mark it done (POST .../step-completions); `APPROVAL`: approve or reject (POST .../approvals). One of `WORK`, `APPROVAL`.
- `name` (String) Name of the step.
- `role` (String) The Network role that takes the step. One of `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`.

<a id="nestedatt--steps"></a>
### Nested Schema for `steps`

- `key` (String) Key of the step in the pinned version.
- `kind` (String) `WORK` does the work; `APPROVAL` approves the work before it. One of `WORK`, `APPROVAL`.
- `name` (String) Name of the step.
- `role` (String) The Network role that takes the step. One of `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`.
- `state` (String) `DONE` behind the order, `CURRENT` where it waits, `PENDING` ahead of it (or never reached, when cancelled). One of `DONE`, `CURRENT`, `PENDING`.

## Import

```terraform
import {
  to = das_delivery_order.example
  id = "<orderId>"
}
```
