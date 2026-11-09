---
page_title: "tag_manager_policy_site_platforms Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_site_platform, all pages: a site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning.
---

# das_tag_manager_policy_site_platforms (Data Source)

Every tag_manager_policy_site_platform, all pages: a site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning.

API: `GET /api/policy-site-platforms`.

## Example Usage

```terraform
data "das_tag_manager_policy_site_platforms" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `id` (String) Rule id, `spmap_<ULID>`.
- `pattern` (String) A site code, or a pattern where `*` stands for any characters.
- `platform` (String) The platform key of the sites it matches.
- `position` (Number, integer) Where it is evaluated: the first matching rule by position wins.

