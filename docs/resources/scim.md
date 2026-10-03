---
page_title: "scim Resource - das"
subcategory: "Organization"
description: |-
  SCIM provisioning of the directory, on or off. Destroy switches it off.
---

# das_scim (Resource)

SCIM provisioning of the directory, on or off. Destroy switches it off.

API: `PUT /api/identity/scim` creates, `GET /api/identity/scim` reads, `PUT /api/identity/scim` updates. Destroy: SCIM is off unless switched on; destroy switches it off.

## Example Usage

```terraform
resource "das_scim" "example" {
  enabled = true
}
```

## Schema

### Required

- `enabled` (Boolean) Off: the SCIM endpoint answers 403 to every token.

### Read-Only

- `id` (String) A constant: the Organization has one scim.

## Import

```terraform
import {
  to = das_scim.example
  id = "scim"
}
```
