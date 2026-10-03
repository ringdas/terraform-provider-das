---
page_title: "network Resource - das"
subcategory: "Organization"
description: |-
  A Network of the Organization: a tenant with its own host label, time zone and default currency. The API never deletes a Network, so destroy only forgets it, and a change to an attribute is refused at plan.
---

# das_network (Resource)

A Network of the Organization: a tenant with its own host label, time zone and default currency. The API never deletes a Network, so destroy only forgets it, and a change to an attribute is refused at plan.

API: `POST /api/networks` creates, `GET /api/networks/{networkId}` reads. Destroy: A Network is never deleted through the API; destroy leaves it and its data in place.

## Example Usage

```terraform
resource "das_network" "example" {
  default_currency = "PLN"
  display_name     = "example-display-name"
  slug             = "example-slug"
  timezone         = "Europe/Warsaw"
}
```

## Schema

### Required

- `default_currency` (String) ISO-4217 code, upper-case, e.g. "PLN". It cannot change: the API offers no update and never deletes it.
- `display_name` (String) Human-readable name. It cannot change: the API offers no update and never deletes it.
- `slug` (String) Host label of the Network (`<slug>.<module>.<org domain>`): 1-63 lower-case letters, digits or hyphens, not a reserved platform label (help, developer, api, www, module slugs, ...). It cannot change: the API offers no update and never deletes it.
- `timezone` (String) IANA time zone name, e.g. "Europe/Warsaw". It cannot change: the API offers no update and never deletes it.

### Optional

- `network_code` (String) The Network's code when it already has one (an import from another system), `NT-` and seven digits; omitted, a free code is assigned. It cannot change: the API offers no update and never deletes it.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) Last change, ISO-8601 UTC.

## Import

```terraform
import {
  to = das_network.example
  id = "<networkId>"
}
```
