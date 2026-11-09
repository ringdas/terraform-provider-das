---
page_title: "delivery_ad_product Resource - das"
subcategory: "DAS Delivery"
description: |-
  An ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells. It starts `DISABLED`; set `status = "ENABLED"` once a `das_delivery_ad_product_version` exists. Its key never changes, and destroy switches it off (`DISABLED`): an ad product is archived, never deleted.
---

# das_delivery_ad_product (Resource)

An ad product of the credential's Network (the feature `ad-product-catalog`): a kind of advertising it sells. It starts `DISABLED`; set `status = "ENABLED"` once a `das_delivery_ad_product_version` exists. Its key never changes, and destroy switches it off (`DISABLED`): an ad product is archived, never deleted.

API: `POST /api/ad-products` creates, `GET /api/ad-products/{adProductId}` reads, `PATCH /api/ad-products/{adProductId}` updates. Destroy: An ad product is archivable (P11): it is never deleted, so destroy stops selling it.

## Example Usage

```terraform
resource "das_delivery_ad_product" "example" {
  description = "A full-page takeover of the homepage for one day"
  key         = "homepage-takeover"
  name        = "Homepage takeover"
}
```

## Schema

### Required

- `key` (String) Key of the ad product, kebab-case, 2-60 characters, unique in the Network and never changed. It cannot change: the API offers no update and never deletes it.
- `name` (String) Name, 1-200 characters.

### Optional

- `description` (String) What the ad product is for, up to 2000 characters. Leave it out for null.
- `status` (String) `ENABLED` to sell the ad product (it needs a version), `DISABLED` to stop. One of `ENABLED`, `DISABLED`. Applied after create, by an update.

### Read-Only

- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `current_version` (Number, integer) The version in force (the newest one); null until the first version is added.
- `id` (String) The id the API assigned.
- `updated_at` (String) Last change, ISO-8601 UTC.

## Import

```terraform
import {
  to = das_delivery_ad_product.example
  id = "<adProductId>"
}
```
