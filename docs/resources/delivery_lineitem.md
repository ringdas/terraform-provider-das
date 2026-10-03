---
page_title: "delivery_lineitem Resource - das"
subcategory: "DAS Delivery"
description: |-
  A line item of the credential's Network: what to serve, when and for which budget. It starts as `DRAFT`. The API has no update and no delete yet (a line item is archived, never deleted): a change is refused at plan, and destroy only forgets it.
---

# das_delivery_lineitem (Resource)

A line item of the credential's Network: what to serve, when and for which budget. It starts as `DRAFT`. The API has no update and no delete yet (a line item is archived, never deleted): a change is refused at plan, and destroy only forgets it.

API: `POST /api/lineitems` creates, `GET /api/lineitems/{lineitemId}` reads. Destroy: A line item is archivable (P11): it is never deleted, and the API has no archive operation yet.

## Example Usage

```terraform
resource "das_delivery_lineitem" "example" {
  budget   = {
    amount = "15000.00"
    currency = "PLN"
  }
  end_at   = "2027-01-31T23:59:59Z"
  name     = "Homepage takeover, January"
  pacing   = "EVEN"
  priority = 8
  start_at = "2027-01-01T00:00:00Z"
  type     = "STANDARD"
}
```

## Schema

### Required

- `budget` (Attributes) Budget, not negative, in the Network's currency (`defaultCurrency` of GET /api/session). It cannot change: the API offers no update and never deletes it.
- `end_at` (String) End of the flight, an ISO-8601 instant after startAt. It cannot change: the API offers no update and never deletes it.
- `name` (String) Name, 1-200 characters. It cannot change: the API offers no update and never deletes it.
- `pacing` (String) How delivery is spread over the flight. One of `EVEN`, `ASAP`, `FRONT_LOADED`. It cannot change: the API offers no update and never deletes it.
- `priority` (Number, integer) Serving priority, 1 (highest) to 16. It cannot change: the API offers no update and never deletes it.
- `start_at` (String) Start of the flight, an ISO-8601 instant. It cannot change: the API offers no update and never deletes it.
- `type` (String) Guaranteed volume, a share of voice, or house inventory. One of `STANDARD`, `SPONSORSHIP`, `HOUSE`. It cannot change: the API offers no update and never deletes it.

### Read-Only

- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `status` (String) Lifecycle status; a new line item is `DRAFT`. One of `DRAFT`, `PENDING_APPROVAL`, `REJECTED`, `SCHEDULED`, `LIVE`, `PAUSED`, `ENDED`, `ARCHIVED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--budget"></a>
### Nested Schema for `budget`

- `amount` (String) Decimal string, e.g. "1250.00".
- `currency` (String) ISO-4217 code, upper-case.

## Import

```terraform
import {
  to = das_delivery_lineitem.example
  id = "<lineitemId>"
}
```
