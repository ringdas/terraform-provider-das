---
page_title: "scim Data Source - das"
subcategory: "Organization"
description: |-
  SCIM provisioning of the directory, on or off.
---

# das_scim (Data Source)

SCIM provisioning of the directory, on or off.

API: `GET /api/identity/scim`.

## Example Usage

```terraform
data "das_scim" "example" {
}
```

## Schema

### Read-Only

- `enabled` (Boolean) Off: the SCIM endpoint answers 403 to every token.

