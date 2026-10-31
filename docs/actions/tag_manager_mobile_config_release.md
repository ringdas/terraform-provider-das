---
page_title: "tag_manager_mobile_config_release Action - das"
subcategory: "DAS Tag Manager"
description: |-
  Publishes a release of an app's mobile configuration for one environment (the feature `mobile-config`), as a Network admin: compiled, signed and numbered after the last release of the app and environment. A release is immutable, so publishing is an action, not a resource: every invocation is a new release.
---

# das_tag_manager_mobile_config_release (Action)

Publishes a release of an app's mobile configuration for one environment (the feature `mobile-config`), as a Network admin: compiled, signed and numbered after the last release of the app and environment. A release is immutable, so publishing is an action, not a resource: every invocation is a new release. Requires Terraform 1.14 or later.

API: `POST /api/mobile-config/releases`.

## Example Usage

```terraform
action "das_tag_manager_mobile_config_release" "example" {
  config {
  environment   = "PRODUCTION"
  site_code     = "ONET_IOS"
  source        = {
    flags = {
      ads_enabled = {
        default = true
        description = "Master switch for ad rendering in the app."
      }
    }
  }
  source_format = "DCFG_RULES_V1"
  }
}
```

Run it with `terraform apply -invoke=action.das_tag_manager_mobile_config_release.example`, or from a resource's `lifecycle { action_trigger { ... } }`.

## Schema

### Required

- `environment` (String) The app's environment the release is for; signed into the bundle in lower case (`production`). One of `PRODUCTION`, `STAGING`, `DEVELOPMENT`.
- `site_code` (String) The app the release is for, UPPER_SNAKE, 1-64 characters (e.g. `ONET_IOS`); signed into the bundle as its app id.
- `source` (Dynamic) The rules in `sourceFormat`. `DCFG_RULES_V1`: `{ segments?, flags }` - named reusable conditions and the flags by key, each with its `default`, optional `prerequisites` and targeting `rules` (docs/architecture/tag-manager-mobile-config.md "Authoring rules"). Nested at most 100 levels; no key `__proto__`; no U+0000 and no unpaired surrogate in a key or a string.
- `source_format` (String) The format of `source`: `DCFG_RULES_V1`, the rules document of segments and flags. One of `DCFG_RULES_V1`.

