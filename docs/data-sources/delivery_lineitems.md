---
page_title: "delivery_lineitems Data Source - das"
subcategory: "DAS Delivery"
description: |-
  Every delivery_lineitem, all pages: a line item of the credential's Network: what to serve, when and for which budget.
---

# das_delivery_lineitems (Data Source)

Every delivery_lineitem, all pages: a line item of the credential's Network: what to serve, when and for which budget.

API: `GET /api/lineitems`.

## Example Usage

```terraform
data "das_delivery_lineitems" "example" {
}
```

## Schema

### Optional

- `status` (String) Only line items in this lifecycle status. One of `DRAFT`, `PENDING_APPROVAL`, `REJECTED`, `SCHEDULED`, `LIVE`, `PAUSED`, `ENDED`, `ARCHIVED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `archived_at` (String) When it was archived; null while not archived.
- `budget` (Attributes) Budget of the flight, in the Network's currency.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `end_at` (String) End of the flight, ISO-8601 UTC.
- `id` (String) Line item id, `li_<ULID>`.
- `name` (String) Name, shown in lists and reports.
- `pacing` (String) How delivery is spread over the flight. One of `EVEN`, `ASAP`, `FRONT_LOADED`.
- `priority` (Number, integer) Serving priority, 1 (highest) to 16.
- `start_at` (String) Start of the flight, ISO-8601 UTC.
- `status` (String) Lifecycle status; a new line item is `DRAFT`. One of `DRAFT`, `PENDING_APPROVAL`, `REJECTED`, `SCHEDULED`, `LIVE`, `PAUSED`, `ENDED`, `ARCHIVED`.
- `type` (String) Guaranteed volume, a share of voice, or house inventory. One of `STANDARD`, `SPONSORSHIP`, `HOUSE`.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--items--budget"></a>
### Nested Schema for `items.budget`

- `amount` (String) Decimal string, e.g. "1250.00".
- `currency` (String) ISO-4217 code, upper-case.

