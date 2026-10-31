---
page_title: "delivery_order Data Source - das"
subcategory: "DAS Delivery"
description: |-
  An order of the credential's Network (the feature `order-workspace`), delivered with an `ENABLED` ad product: it pins the product's current version and starts on its first step. Read by its identity.
---

# das_delivery_order (Data Source)

An order of the credential's Network (the feature `order-workspace`), delivered with an `ENABLED` ad product: it pins the product's current version and starts on its first step. Read by its identity.

API: `GET /api/orders/{orderId}`.

## Example Usage

```terraform
data "das_delivery_order" "example" {
  id = "<orderId>"
}
```

## Schema

### Required

- `id` (String) The order, `ord_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `ad_product_id` (String) The ad product the order is delivered with, `adp_<ULID>`.
- `ad_product_version` (Number, integer) The version of the ad product the order follows, pinned when it was created.
- `archived_at` (String) When it was archived; null while not archived.
- `billing_model` (String) How the order is billed, one of the pinned version's billing models. One of `CPM`, `VIEWABLE_CPM`, `CPC`, `CPV`, `CPA`, `FLAT_FEE`.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `current_step` (Attributes) The step the order waits on; null once it is completed or cancelled.
- `end_at` (String) End of the flight, ISO-8601 UTC.
- `goal` (String) What the order aims for, one of the pinned version's goals. One of `AWARENESS`, `REACH`, `TRAFFIC`, `ENGAGEMENT`, `VIDEO_VIEWS`, `LEADS`, `SALES`.
- `name` (String) Name, shown in lists and queues.
- `sales_channel_id` (String) The sales channel the order came through, `sch_<ULID>`, from the credential that created it (a key of the channel); null when the Network's own staff created it.
- `start_at` (String) Start of the flight, ISO-8601 UTC.
- `status` (String) `IN_PROGRESS` while its steps run, `COMPLETED` after the last, `CANCELLED` when stopped. One of `IN_PROGRESS`, `COMPLETED`, `CANCELLED`.
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

