---
page_title: "organization_setting Resource - das"
subcategory: "Organization"
description: |-
  The Organization's value of a setting of the settings catalog (an `org` or `org-default` setting). Destroy clears it.
---

# das_organization_setting (Resource)

The Organization's value of a setting of the settings catalog (an `org` or `org-default` setting). Destroy clears it.

API: `PUT /api/settings/{key}` creates, `GET /api/settings/{key}` reads, `PUT /api/settings/{key}` updates, `DELETE /api/settings/{key}` deletes.

## Example Usage

```terraform
resource "das_organization_setting" "example" {
  key   = "delivery.accent-color"
  value = "#0055aa"
}
```

## Schema

### Required

- `key` (String) A key of the settings catalog, `<domain>.<kebab-name>`, e.g. `delivery.accent-color`. A key the catalog does not know is 404 `NOT_FOUND`. Changing it replaces the resource.
- `value` (Dynamic) The value, in the setting's format (its description says which): a string such as `#4fb0e8` for an accent color, or an object such as a Network's email branding.

### Read-Only

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `id` (String) The identity, key joined by `/` (the import id).
- `level` (String) `ORG`: one Organization value no Network changes; `ORG_DEFAULT`: the Organization's default, a Network may set its own. One of `ORG`, `ORG_DEFAULT`.
- `updated_at` (String) Last change of the value, ISO-8601 UTC; null while unset.

## Import

```terraform
import {
  to = das_organization_setting.example
  id = "<key>"
}
```
