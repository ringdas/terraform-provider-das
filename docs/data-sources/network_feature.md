---
page_title: "network_feature Data Source - das"
subcategory: "Organization"
description: |-
  A feature switch of a Network. Read by its identity.
---

# das_network_feature (Data Source)

A feature switch of a Network. Read by its identity.

API: `GET /api/networks/{networkId}/features/{featureId}`.

## Example Usage

```terraform
data "das_network_feature" "example" {
  network_id = "<networkId>"
  feature_id = "<featureId>"
}
```

## Schema

### Required

- `feature_id` (String) A feature id of the feature registry (kebab-case). An id the registry does not know is 404 `NOT_FOUND`.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `blocked_by` (List of String) Requirements that are off - why a switched-on feature is not effective.
- `description` (String) What the feature gives a Network.
- `effective` (Boolean) Whether the feature is on for the Network: switched on and every requirement effective.
- `enabled` (Boolean) The Network's own switch (false while never set).
- `requires` (List of String) Features that must be on for this one to count as on.
- `updated_at` (String) Last change of the switch, ISO-8601 UTC; null while never set.
- `updated_by` (String) Actor of the last change of the switch; null while never set.

