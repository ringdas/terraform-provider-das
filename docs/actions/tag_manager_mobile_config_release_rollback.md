---
page_title: "tag_manager_mobile_config_release_rollback Action - das"
subcategory: "DAS Tag Manager"
description: |-
  Rolls an app's mobile configuration back to an earlier release (the feature `mobile-config`), as a Network admin: its source is republished as the next release of the app and environment, signed anew.
---

# das_tag_manager_mobile_config_release_rollback (Action)

Rolls an app's mobile configuration back to an earlier release (the feature `mobile-config`), as a Network admin: its source is republished as the next release of the app and environment, signed anew. Requires Terraform 1.14 or later.

API: `POST /api/mobile-config/releases/{releaseId}/rollback`.

## Example Usage

```terraform
action "das_tag_manager_mobile_config_release_rollback" "example" {
  config {
  release_id = "mcrel_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  }
}
```

Run it with `terraform apply -invoke=action.das_tag_manager_mobile_config_release_rollback.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `release_id` (String) The release, `mcrel_<ULID>`. One of another Network is 404 `NOT_FOUND`.

