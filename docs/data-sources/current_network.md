---
page_title: "current_network Data Source - das"
subcategory: "Organization"
description: |-
  The Network of the provider's Network credential.
---

# das_current_network (Data Source)

The Network of the provider's Network credential.

API: `GET /api/network`.

## Example Usage

```terraform
data "das_current_network" "example" {
}
```

## Schema

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits (e.g. `NT-1234567`): unique in the Organization, not the slug and not the id, fixed once given.
- `slug` (String) Host label of the Network, unique in the Organization.
- `timezone` (String) IANA time zone where the Network's days and months begin.
- `updated_at` (String) Last change, ISO-8601 UTC.

