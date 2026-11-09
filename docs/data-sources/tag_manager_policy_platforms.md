---
page_title: "tag_manager_policy_platforms Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  The platforms of DAS Policies (the feature `policies`): `APP` and `WEB`, by key, what policy types and site-platform rules name.
---

# das_tag_manager_policy_platforms (Data Source)

The platforms of DAS Policies (the feature `policies`): `APP` and `WEB`, by key, what policy types and site-platform rules name.

API: `GET /api/policy-platforms`.

## Example Usage

```terraform
data "das_tag_manager_policy_platforms" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) The platforms.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `key` (String) Its key, UPPER_SNAKE (`APP`, `WEB`): what bodies name it by.
- `label` (String) Human-readable name.

