---
page_title: "delivery_ad_product_versions Data Source - das"
subcategory: "DAS Delivery"
description: |-
  Every delivery_ad_product_version, all pages: a version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps.
---

# das_delivery_ad_product_versions (Data Source)

Every delivery_ad_product_version, all pages: a version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps.

API: `GET /api/ad-products/{adProductId}/versions`.

## Example Usage

```terraform
data "das_delivery_ad_product_versions" "example" {
  ad_product_id = "<adProductId>"
}
```

## Schema

### Required

- `ad_product_id` (String) The ad product, `adp_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `ad_product_id` (String) The ad product, `adp_<ULID>`.
- `created_at` (String) When the version was added, ISO-8601 UTC.
- `definition` (Attributes) The definition, exactly as it was added.
- `id` (String) Version id, `adpv_<ULID>`.
- `version` (Number, integer) The version number: 1 for the first, then one more for each.

