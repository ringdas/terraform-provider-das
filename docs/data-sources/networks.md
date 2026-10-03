---
page_title: "networks Data Source - das"
subcategory: "Organization"
description: |-
  Every network, all pages: a Network of the Organization: a tenant with its own host label, time zone and default currency.
---

# das_networks (Data Source)

Every network, all pages: a Network of the Organization: a tenant with its own host label, time zone and default currency.

API: `GET /api/networks`.

## Example Usage

```terraform
data "das_networks" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits (e.g. `NT-1234567`): unique in the Organization, not the slug and not the id, fixed once given.
- `slug` (String) Host label of the Network, unique in the Organization.
- `timezone` (String) IANA time zone where the Network's days and months begin.
- `updated_at` (String) Last change, ISO-8601 UTC.

