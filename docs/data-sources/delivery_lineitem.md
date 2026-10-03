---
page_title: "delivery_lineitem Data Source - das"
subcategory: "DAS Delivery"
description: |-
  A line item of the credential's Network: what to serve, when and for which budget. Read by its identity.
---

# das_delivery_lineitem (Data Source)

A line item of the credential's Network: what to serve, when and for which budget. Read by its identity.

API: `GET /api/lineitems/{lineitemId}`.

## Example Usage

```terraform
data "das_delivery_lineitem" "example" {
  id = "<lineitemId>"
}
```

## Schema

### Required

- `id` (String) The line item, `li_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived; null while not archived.
- `budget` (Attributes) Budget of the flight, in the Network's currency.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `end_at` (String) End of the flight, ISO-8601 UTC.
- `name` (String) Name, shown in lists and reports.
- `pacing` (String) How delivery is spread over the flight. One of `EVEN`, `ASAP`, `FRONT_LOADED`.
- `priority` (Number, integer) Serving priority, 1 (highest) to 16.
- `start_at` (String) Start of the flight, ISO-8601 UTC.
- `status` (String) Lifecycle status; a new line item is `DRAFT`. One of `DRAFT`, `PENDING_APPROVAL`, `REJECTED`, `SCHEDULED`, `LIVE`, `PAUSED`, `ENDED`, `ARCHIVED`.
- `type` (String) Guaranteed volume, a share of voice, or house inventory. One of `STANDARD`, `SPONSORSHIP`, `HOUSE`.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--budget"></a>
### Nested Schema for `budget`

- `amount` (String) Decimal string, e.g. "1250.00".
- `currency` (String) ISO-4217 code, upper-case.

