---
page_title: "network_settings Data Source - das"
subcategory: "Organization"
description: |-
  Every network_setting, all pages: a Network's own value of a setting (an override of the Organization's default).
---

# das_network_settings (Data Source)

Every network_setting, all pages: a Network's own value of a setting (an override of the Organization's default).

API: `GET /api/networks/{networkId}/settings`.

## Example Usage

```terraform
data "das_network_settings" "example" {
  network_id = "<networkId>"
}
```

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `key` (String) The setting's key.
- `level` (String) `ORG`: the Organization's alone (a Network cannot set it); `ORG_DEFAULT`: a Network may override it; `NETWORK`: the Network's only. One of `ORG`, `ORG_DEFAULT`, `NETWORK`.
- `network_id` (String) The Network.
- `network_value_json` (String) The Network's own value; null when it has none. JSON-encoded (`jsondecode()` it).
- `organization_value_json` (String) The Organization's value; null when unset. JSON-encoded (`jsondecode()` it).
- `source` (String) Where the value that applies comes from: the Network's own, the Organization's, or the setting's built-in default. One of `NETWORK`, `ORGANIZATION`, `DEFAULT`.
- `updated_at` (String) Last change of the Network's own value, ISO-8601 UTC; null when it has none.
- `updated_by` (String) Actor of the last change of the Network's own value; null when it has none.
- `value_json` (String) What applies for the Network; null = the setting's built-in default. JSON-encoded (`jsondecode()` it).

