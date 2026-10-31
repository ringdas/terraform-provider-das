---
page_title: "delivery_sales_channels Data Source - das"
subcategory: "DAS Delivery"
description: |-
  Every delivery_sales_channel, all pages: a sales channel of the credential's Network (the feature `sales-channels`): a self-service tool or a partner the Network sells through, with the ad products it offers (`ad_product_ids`) and the API keys that sell there (`api_key_ids`, each linked to one channel at most).
---

# das_delivery_sales_channels (Data Source)

Every delivery_sales_channel, all pages: a sales channel of the credential's Network (the feature `sales-channels`): a self-service tool or a partner the Network sells through, with the ad products it offers (`ad_product_ids`) and the API keys that sell there (`api_key_ids`, each linked to one channel at most).

API: `GET /api/sales-channels`.

## Example Usage

```terraform
data "das_delivery_sales_channels" "example" {
}
```

## Schema

### Optional

- `status` (String) Only sales channels in this status. One of `ENABLED`, `DISABLED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `ad_product_ids` (List of String) The ad products the channel offers, `adp_<ULID>`; its keys see and order those that are `ENABLED`.
- `api_key_ids` (List of String) The API keys of the Network linked to the channel, `key_<ULID>`; each is linked to one channel at most.
- `archived_at` (String) When it was archived; null while not archived.
- `created_at` (String) Creation instant, ISO-8601 UTC.
- `description` (String) What the channel is; null when not given.
- `id` (String) Sales channel id, `sch_<ULID>`.
- `key` (String) Key, kebab-case, unique in the Network and never changed.
- `name` (String) Name, shown with the channel's orders.
- `status` (String) `ENABLED` while the Network sells through it; `DISABLED` (a new channel): its keys see no ad product and create no order, and still see only its orders. One of `ENABLED`, `DISABLED`.
- `updated_at` (String) Last change, ISO-8601 UTC.

