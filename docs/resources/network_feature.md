---
page_title: "network_feature Resource - das"
subcategory: "Organization"
description: |-
  A feature switch of a Network. Every feature is OFF for a Network until switched on; destroy switches it off again.
---

# das_network_feature (Resource)

A feature switch of a Network. Every feature is OFF for a Network until switched on; destroy switches it off again.

API: `PUT /api/networks/{networkId}/features/{featureId}` creates, `GET /api/networks/{networkId}/features/{featureId}` reads, `PUT /api/networks/{networkId}/features/{featureId}` updates. Destroy: A feature is off unless switched on; destroy switches it off.

## Example Usage

```terraform
resource "das_network_feature" "example" {
  enabled    = true
  feature_id = "lineitem-summary"
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
}
```

## Schema

### Required

- `enabled` (Boolean) true switches the feature on for the Network, false off.
- `feature_id` (String) A feature id of the feature registry (kebab-case). An id the registry does not know is 404 `NOT_FOUND`. It cannot change: the API offers no update and never deletes it.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network. It cannot change: the API offers no update and never deletes it.

### Read-Only

- `blocked_by` (List of String) Requirements that are off - why a switched-on feature is not effective.
- `description` (String) What the feature gives a Network.
- `effective` (Boolean) Whether the feature is on for the Network: switched on and every requirement effective.
- `id` (String) The identity, networkId/featureId joined by `/` (the import id).
- `requires` (List of String) Features that must be on for this one to count as on.
- `updated_at` (String) Last change of the switch, ISO-8601 UTC; null while never set.
- `updated_by` (String) Actor of the last change of the switch; null while never set.

## Import

```terraform
import {
  to = das_network_feature.example
  id = "<networkId>/<featureId>"
}
```
