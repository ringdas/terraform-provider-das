---
page_title: "network_directory Data Source - das"
subcategory: "Organization"
description: |-
  The Organization's Networks a Network session may switch to (slug, name and code).
---

# das_network_directory (Data Source)

The Organization's Networks a Network session may switch to (slug, name and code).

API: `GET /api/network/directory`.

## Example Usage

```terraform
data "das_network_directory" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `display_name` (String) Human-readable name.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits (e.g. `NT-1234567`): unique in the Organization, not the slug and not the id, fixed once given.
- `slug` (String) Host label of the Network: a switch navigates to `<slug>.<module>.<org domain>`.

