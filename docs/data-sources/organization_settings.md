---
page_title: "organization_settings Data Source - das"
subcategory: "Organization"
description: |-
  Every organization_setting, all pages: the Organization's value of a setting of the settings catalog (an `org` or `org-default` setting).
---

# das_organization_settings (Data Source)

Every organization_setting, all pages: the Organization's value of a setting of the settings catalog (an `org` or `org-default` setting).

API: `GET /api/settings`.

## Example Usage

```terraform
data "das_organization_settings" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `description` (String) What the setting decides.
- `domain` (String) The owning domain: `delivery`, `audience`, `tag-manager` or `platform`.
- `key` (String) The setting's key, `<domain>.<kebab-name>`.
- `level` (String) `ORG`: one Organization value no Network changes; `ORG_DEFAULT`: the Organization's default, a Network may set its own. One of `ORG`, `ORG_DEFAULT`.
- `updated_at` (String) Last change of the value, ISO-8601 UTC; null while unset.
- `value_json` (String) The Organization's value; null while unset (the setting's built-in default applies). JSON-encoded (`jsondecode()` it).

