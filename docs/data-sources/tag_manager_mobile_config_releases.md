---
page_title: "tag_manager_mobile_config_releases Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  The releases of the Network's mobile configuration (the feature `mobile-config`), newest first, without their sources, filtered by app and environment (every page followed).
---

# das_tag_manager_mobile_config_releases (Data Source)

The releases of the Network's mobile configuration (the feature `mobile-config`), newest first, without their sources, filtered by app and environment (every page followed).

API: `GET /api/mobile-config/releases`.

## Example Usage

```terraform
data "das_tag_manager_mobile_config_releases" "example" {
}
```

## Schema

### Optional

- `environment` (String) Only releases for this environment. One of `PRODUCTION`, `STAGING`, `DEVELOPMENT`.
- `site_code` (String) Only releases of this app, UPPER_SNAKE (e.g. `ONET_IOS`).

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `bundle_sha256` (String) Hex SHA-256 of the signed bundle, as a device receives it.
- `bundle_size` (Number, integer) Length of the signed bundle in bytes.
- `environment` (String) The app's environment the release is for; the bundle names it in lower case (`production`). One of `PRODUCTION`, `STAGING`, `DEVELOPMENT`.
- `id` (String) Release id, `mcrel_<ULID>`.
- `published_at` (String) When it was published, ISO-8601 UTC; signed into the bundle as its issue time.
- `published_by` (String) Who published it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `rolled_back_from_id` (String) The release whose source this one republished (a rollback), `mcrel_<ULID>`; null for a publication.
- `signing_key_fingerprint` (String) Hex SHA-256 of the raw Ed25519 public key that signed the bundle (`fingerprint` of GET /api/mobile-config/signing-key).
- `site_code` (String) The app the release is for, UPPER_SNAKE; the bundle's app id.
- `source_format` (String) The format of the source: `DCFG_RULES_V1`, the rules document of segments and flags. One of `DCFG_RULES_V1`.
- `version` (Number, integer) The release's number among the releases of its `siteCode` and `environment`, from 1; signed into the bundle.

