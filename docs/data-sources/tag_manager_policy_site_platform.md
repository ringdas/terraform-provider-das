---
page_title: "tag_manager_policy_site_platform Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning. Read by its identity.
---

# das_tag_manager_policy_site_platform (Data Source)

A site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning. Read by its identity.

API: `GET /api/policy-site-platforms/{mappingId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_site_platform" "example" {
  id = "<mappingId>"
}
```

## Schema

### Required

- `id` (String) The site-platform rule, `spmap_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `pattern` (String) A site code, or a pattern where `*` stands for any characters.
- `platform` (String) The platform key of the sites it matches.
- `position` (Number, integer) Where it is evaluated: the first matching rule by position wins.

