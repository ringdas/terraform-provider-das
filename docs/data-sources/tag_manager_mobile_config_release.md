---
page_title: "tag_manager_mobile_config_release Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  One release of the Network's mobile configuration (the feature `mobile-config`), with its source as published.
---

# das_tag_manager_mobile_config_release (Data Source)

One release of the Network's mobile configuration (the feature `mobile-config`), with its source as published.

API: `GET /api/mobile-config/releases/{releaseId}`.

## Example Usage

```terraform
data "das_tag_manager_mobile_config_release" "example" {
  release_id = "<releaseId>"
}
```

## Schema

### Required

- `release_id` (String) The release, `mcrel_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `bundle_sha256` (String) Hex SHA-256 of the signed bundle, as a device receives it.
- `bundle_size` (Number, integer) Length of the signed bundle in bytes.
- `environment` (String) The app's environment the release is for; the bundle names it in lower case (`production`). One of `PRODUCTION`, `STAGING`, `DEVELOPMENT`.
- `id` (String) Release id, `mcrel_<ULID>`.
- `published_at` (String) When it was published, ISO-8601 UTC; signed into the bundle as its issue time.
- `published_by` (String) Who published it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `rolled_back_from_id` (String) The release whose source this one republished (a rollback), `mcrel_<ULID>`; null for a publication.
- `signing_key_fingerprint` (String) Hex SHA-256 of the raw Ed25519 public key that signed the bundle (`fingerprint` of GET /api/mobile-config/signing-key).
- `site_code` (String) The app the release is for, UPPER_SNAKE; the bundle's app id.
- `source` (Dynamic) The source as published, in `sourceFormat` (`DCFG_RULES_V1`: `{ segments?, flags }`).
- `source_format` (String) The format of the source: `DCFG_RULES_V1`, the rules document of segments and flags. One of `DCFG_RULES_V1`.
- `version` (Number, integer) The release's number among the releases of its `siteCode` and `environment`, from 1; signed into the bundle.

