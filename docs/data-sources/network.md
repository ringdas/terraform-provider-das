---
page_title: "network Data Source - das"
subcategory: "Organization"
description: |-
  A Network of the Organization: a tenant with its own host label, time zone and default currency. Read by its identity.
---

# das_network (Data Source)

A Network of the Organization: a tenant with its own host label, time zone and default currency. Read by its identity.

API: `GET /api/networks/{networkId}`.

## Example Usage

```terraform
data "das_network" "example" {
  id = "<networkId>"
}
```

## Schema

### Required

- `id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits (e.g. `NT-1234567`): unique in the Organization, not the slug and not the id, fixed once given.
- `slug` (String) Host label of the Network, unique in the Organization.
- `timezone` (String) IANA time zone where the Network's days and months begin.
- `updated_at` (String) Last change, ISO-8601 UTC.

