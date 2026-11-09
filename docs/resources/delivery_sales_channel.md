---
page_title: "delivery_sales_channel Resource - das"
subcategory: "DAS Delivery"
description: |-
  A sales channel of the credential's Network (the feature `sales-channels`): a self-service tool or a partner the Network sells through, with the ad products it offers (`ad_product_ids`) and the API keys that sell there (`api_key_ids`, each linked to one channel at most). It starts `DISABLED`; set `status = "ENABLED"` to sell through it. Its key never changes, and destroy switches it off (`DISABLED`) and keeps its links, so its keys stay narrowed to it: a sales channel is archived, never deleted.
---

# das_delivery_sales_channel (Resource)

A sales channel of the credential's Network (the feature `sales-channels`): a self-service tool or a partner the Network sells through, with the ad products it offers (`ad_product_ids`) and the API keys that sell there (`api_key_ids`, each linked to one channel at most). It starts `DISABLED`; set `status = "ENABLED"` to sell through it. Its key never changes, and destroy switches it off (`DISABLED`) and keeps its links, so its keys stay narrowed to it: a sales channel is archived, never deleted.

API: `POST /api/sales-channels` creates, `GET /api/sales-channels/{salesChannelId}` reads, `PATCH /api/sales-channels/{salesChannelId}` updates. Destroy: A sales channel is archivable (P11): it is never deleted, so destroy stops selling through it; its keys stay linked and narrowed (fail closed).

## Example Usage

```terraform
resource "das_delivery_sales_channel" "example" {
  ad_product_ids = ["adp_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"]
  api_key_ids    = ["key_01K6P6ZQ4N2V3W4X5Y6Z7A8B9D"]
  description    = "The self-service tool of the Network's advertisers"
  key            = "onetads"
  name           = "OnetAds"
}
```

## Schema

### Required

- `key` (String) Key of the sales channel, kebab-case, 2-60 characters, unique in the Network and never changed. It cannot change: the API offers no update and never deletes it.
- `name` (String) Name, 1-200 characters, shown with the channel's orders.

### Optional

- `ad_product_ids` (List of String) The ad products the channel offers, `adp_<ULID>` of this Network: 0-200 distinct ids. A key of the channel sees and orders only those that are `ENABLED`.
- `api_key_ids` (List of String) The API keys of this Network that sell through the channel, `key_<ULID>` (GET /api/network/api-keys of the Organisation API): 0-20 distinct ids, each linked to one channel at most. The Organisation API is not asked: an id of no key links nothing.
- `description` (String) What the channel is, up to 2000 characters. Leave it out for null.
- `status` (String) `ENABLED` to sell through the channel, `DISABLED` to stop: its keys then see no ad product and create no order. One of `ENABLED`, `DISABLED`. Applied after create, by an update.

### Read-Only

- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `id` (String) The id the API assigned.
- `updated_at` (String) Last change, ISO-8601 UTC.

## Import

```terraform
import {
  to = das_delivery_sales_channel.example
  id = "<salesChannelId>"
}
```
