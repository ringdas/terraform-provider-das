---
page_title: "current_network_features Data Source - das"
subcategory: "Organization"
description: |-
  The feature switches of the credential's Network.
---

# das_current_network_features (Data Source)

The feature switches of the credential's Network.

API: `GET /api/network/features`.

## Example Usage

```terraform
data "das_current_network_features" "example" {
}
```

## Schema

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

