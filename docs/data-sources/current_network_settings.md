---
page_title: "current_network_settings Data Source - das"
subcategory: "Organization"
description: |-
  The setting values that apply to the credential's Network.
---

# das_current_network_settings (Data Source)

The setting values that apply to the credential's Network.

API: `GET /api/network/settings`.

## Example Usage

```terraform
data "das_current_network_settings" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `domain` (String) The owning domain.
- `key` (String) The setting's key.
- `source` (String) Where the value that applies comes from: the Network's own, the Organization's, or the setting's built-in default. One of `NETWORK`, `ORGANIZATION`, `DEFAULT`.
- `value_json` (String) What applies; null = the setting's built-in default. JSON-encoded (`jsondecode()` it).

