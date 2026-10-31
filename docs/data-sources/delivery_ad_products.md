---
page_title: "delivery_ad_products Data Source - das"
subcategory: "DAS Delivery"
description: |-
  Every delivery_ad_product, all pages: an ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells.
---

# das_delivery_ad_products (Data Source)

Every delivery_ad_product, all pages: an ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells.

API: `GET /api/ad-products`.

## Example Usage

```terraform
data "das_delivery_ad_products" "example" {
}
```

## Schema

### Optional

- `status` (String) Only ad products in this status. One of `ENABLED`, `DISABLED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `current_version` (Number, integer) The version in force (the newest one); null until the first version is added.
- `description` (String) What the ad product is for; null when not given.
- `id` (String) Ad product id, `adp_<ULID>`.
- `key` (String) Key, kebab-case, unique in the Network and never changed.
- `name` (String) Name, shown in lists.
- `status` (String) `ENABLED` while the Network sells it; a new ad product is `DISABLED`. One of `ENABLED`, `DISABLED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

