---
page_title: "network_features Data Source - das"
subcategory: "Organization"
description: |-
  Every network_feature, all pages: a feature switch of a Network.
---

# das_network_features (Data Source)

Every network_feature, all pages: a feature switch of a Network.

API: `GET /api/networks/{networkId}/features`.

## Example Usage

```terraform
data "das_network_features" "example" {
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

- `blocked_by` (List of String) Requirements that are off - why a switched-on feature is not effective.
- `description` (String) What the feature gives a Network.
- `effective` (Boolean) Whether the feature is on for the Network: switched on and every requirement effective.
- `enabled` (Boolean) The Network's own switch (false while never set).
- `feature_id` (String) The feature id.
- `network_id` (String) The Network.
- `requires` (List of String) Features that must be on for this one to count as on.
- `updated_at` (String) Last change of the switch, ISO-8601 UTC; null while never set.
- `updated_by` (String) Actor of the last change of the switch; null while never set.

