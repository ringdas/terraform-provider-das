---
page_title: "tag_manager_policy_site_platform Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning. Written by a Network admin; a position is taken once in the Network.
---

# das_tag_manager_policy_site_platform (Resource)

A site-platform rule of the credential's Network (the feature `policies`): the sites whose code matches `pattern` (`*` for any characters) are on `platform`, the first matching rule by `position` winning. Written by a Network admin; a position is taken once in the Network.

API: `POST /api/policy-site-platforms` creates, `GET /api/policy-site-platforms/{mappingId}` reads, `PATCH /api/policy-site-platforms/{mappingId}` updates, `DELETE /api/policy-site-platforms/{mappingId}` deletes.

## Example Usage

```terraform
resource "das_tag_manager_policy_site_platform" "example" {
  pattern  = "*_APP_*"
  platform = "APP"
  position = 1
}
```

## Schema

### Required

- `pattern` (String) A site code (an exception, placed before the patterns) or a pattern where `*` stands for any characters (`*_APP_*`, `*`), 1-200 characters.
- `platform` (String) The platform of the sites it matches: a platform key (GET /api/policy-platforms).
- `position` (Number, integer) Where the rule is evaluated, 1-10000: the first rule by position whose pattern matches a site's code gives its platform. Unique in the Network (409 `ALREADY_EXISTS`).

### Read-Only

- `id` (String) The id the API assigned.

## Import

```terraform
import {
  to = das_tag_manager_policy_site_platform.example
  id = "<mappingId>"
}
```
