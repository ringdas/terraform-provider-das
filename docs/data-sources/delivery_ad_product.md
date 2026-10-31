---
page_title: "delivery_ad_product Data Source - das"
subcategory: "DAS Delivery"
description: |-
  An ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells. Read by its identity.
---

# das_delivery_ad_product (Data Source)

An ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells. Read by its identity.

API: `GET /api/ad-products/{adProductId}`.

## Example Usage

```terraform
data "das_delivery_ad_product" "example" {
  id = "<adProductId>"
}
```

## Schema

### Required

- `id` (String) The ad product, `adp_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `current_version` (Number, integer) The version in force (the newest one); null until the first version is added.
- `description` (String) What the ad product is for; null when not given.
- `key` (String) Key, kebab-case, unique in the Network and never changed.
- `name` (String) Name, shown in lists.
- `status` (String) `ENABLED` while the Network sells it; a new ad product is `DISABLED`. One of `ENABLED`, `DISABLED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

