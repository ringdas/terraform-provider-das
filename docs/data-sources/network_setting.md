---
page_title: "network_setting Data Source - das"
subcategory: "Organization"
description: |-
  A Network's own value of a setting (an override of the Organization's default). Read by its identity.
---

# das_network_setting (Data Source)

A Network's own value of a setting (an override of the Organization's default). Read by its identity.

API: `GET /api/networks/{networkId}/settings/{key}`.

## Example Usage

```terraform
data "das_network_setting" "example" {
  network_id = "<networkId>"
  key        = "example-key"
}
```

## Schema

### Required

- `key` (String) A key of the settings catalog, `<domain>.<kebab-name>`, e.g. `delivery.accent-color`. A key the catalog does not know is 404 `NOT_FOUND`.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `level` (String) `ORG`: the Organization's alone (a Network cannot set it); `ORG_DEFAULT`: a Network may override it; `NETWORK`: the Network's only. One of `ORG`, `ORG_DEFAULT`, `NETWORK`.
- `network_value` (Dynamic) The Network's own value; null when it has none.
- `organization_value` (Dynamic) The Organization's value; null when unset.
- `source` (String) Where the value that applies comes from: the Network's own, the Organization's, or the setting's built-in default. One of `NETWORK`, `ORGANIZATION`, `DEFAULT`.
- `updated_at` (String) Last change of the Network's own value, ISO-8601 UTC; null when it has none.
- `updated_by` (String) Actor of the last change of the Network's own value; null when it has none.
- `value` (Dynamic) What applies for the Network; null = the setting's built-in default.

