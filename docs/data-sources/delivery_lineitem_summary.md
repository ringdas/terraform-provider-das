---
page_title: "delivery_lineitem_summary Data Source - das"
subcategory: "DAS Delivery"
description: |-
  The summary of the Network's line items (the feature `lineitem-summary`).
---

# das_delivery_lineitem_summary (Data Source)

The summary of the Network's line items (the feature `lineitem-summary`).

API: `GET /api/lineitem-summary`.

## Example Usage

```terraform
data "das_delivery_lineitem_summary" "example" {
}
```

## Schema

### Read-Only

- `booked` (Attributes) Budget of every line item that is not archived, in the Network's currency.
- `by_status` (Attributes List) How many line items are in each lifecycle status, every status listed.
- `total` (Number, integer) Every line item of the Network, archived ones included.

<a id="nestedatt--booked"></a>
### Nested Schema for `booked`

- `amount` (String) Decimal string, e.g. "1250.00".
- `currency` (String) ISO-4217 code, upper-case.

<a id="nestedatt--by_status"></a>
### Nested Schema for `by_status`

- `count` (Number, integer) Line items in it.
- `status` (String) Lifecycle status. One of `DRAFT`, `PENDING_APPROVAL`, `REJECTED`, `SCHEDULED`, `LIVE`, `PAUSED`, `ENDED`, `ARCHIVED`.

