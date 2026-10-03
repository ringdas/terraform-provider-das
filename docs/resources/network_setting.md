---
page_title: "network_setting Resource - das"
subcategory: "Organization"
description: |-
  A Network's own value of a setting (an override of the Organization's default). Destroy clears it, and the Organization's value applies again.
---

# das_network_setting (Resource)

A Network's own value of a setting (an override of the Organization's default). Destroy clears it, and the Organization's value applies again.

API: `PUT /api/networks/{networkId}/settings/{key}` creates, `GET /api/networks/{networkId}/settings/{key}` reads, `PUT /api/networks/{networkId}/settings/{key}` updates, `DELETE /api/networks/{networkId}/settings/{key}` deletes.

## Example Usage

```terraform
resource "das_network_setting" "example" {
  key        = "delivery.accent-color"
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
  value      = "#aa5500"
}
```

## Schema

### Required

- `key` (String) A key of the settings catalog, `<domain>.<kebab-name>`, e.g. `delivery.accent-color`. A key the catalog does not know is 404 `NOT_FOUND`. Changing it replaces the resource.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network. Changing it replaces the resource.
- `value` (Dynamic) The value, in the setting's format (its description says which): a string such as `#4fb0e8` for an accent color, or an object such as a Network's email branding.

### Read-Only

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `id` (String) The identity, networkId/key joined by `/` (the import id).
- `level` (String) `ORG`: the Organization's alone (a Network cannot set it); `ORG_DEFAULT`: a Network may override it; `NETWORK`: the Network's only. One of `ORG`, `ORG_DEFAULT`, `NETWORK`.
- `network_value` (Dynamic) The Network's own value; null when it has none.
- `organization_value` (Dynamic) The Organization's value; null when unset.
- `source` (String) Where the value that applies comes from: the Network's own, the Organization's, or the setting's built-in default. One of `NETWORK`, `ORGANIZATION`, `DEFAULT`.
- `updated_at` (String) Last change of the Network's own value, ISO-8601 UTC; null when it has none.
- `updated_by` (String) Actor of the last change of the Network's own value; null when it has none.

## Import

```terraform
import {
  to = das_network_setting.example
  id = "<networkId>/<key>"
}
```
