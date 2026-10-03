---
page_title: "organization_setting Data Source - das"
subcategory: "Organization"
description: |-
  The Organization's value of a setting of the settings catalog (an `org` or `org-default` setting). Read by its identity.
---

# das_organization_setting (Data Source)

The Organization's value of a setting of the settings catalog (an `org` or `org-default` setting). Read by its identity.

API: `GET /api/settings/{key}`.

## Example Usage

```terraform
data "das_organization_setting" "example" {
  key = "example-key"
}
```

## Schema

### Required

- `key` (String) A key of the settings catalog, `<domain>.<kebab-name>`, e.g. `delivery.accent-color`. A key the catalog does not know is 404 `NOT_FOUND`.

### Read-Only

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `level` (String) `ORG`: one Organization value no Network changes; `ORG_DEFAULT`: the Organization's default, a Network may set its own. One of `ORG`, `ORG_DEFAULT`.
- `updated_at` (String) Last change of the value, ISO-8601 UTC; null while unset.
- `value` (Dynamic) The Organization's value; null while unset (the setting's built-in default applies).

