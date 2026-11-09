---
page_title: "delivery_ad_product_version Data Source - das"
subcategory: "DAS Delivery"
description: |-
  A version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps. Read by its identity.
---

# das_delivery_ad_product_version (Data Source)

A version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps. Read by its identity.

API: `GET /api/ad-products/{adProductId}/versions/{version}`.

## Example Usage

```terraform
data "das_delivery_ad_product_version" "example" {
  ad_product_id = "<adProductId>"
  version       = 1
}
```

## Schema

### Required

- `ad_product_id` (String) The ad product, `adp_<ULID>`. One of another Network is 404 `NOT_FOUND`.
- `version` (Number, integer) The version number, an integer from 1.

### Read-Only

- `created_at` (String) When the version was added, ISO-8601 UTC.
- `definition` (Attributes) The definition, exactly as it was added.
- `id` (String) Version id, `adpv_<ULID>`.

<a id="nestedatt--definition"></a>
### Nested Schema for `definition`

- `billing_models` (List of String) How a campaign of the ad product may be billed: 1-6 distinct billing models.
- `creative_source` (String) Where the creatives come from: supplied by the buyer, built from a template, or generated from a product catalog. One of `SUPPLIED`, `TEMPLATE`, `CATALOG`.
- `goals` (List of String) What a campaign of the ad product may aim for: 1-7 distinct goals.
- `steps` (Attributes List) The process that delivers a campaign, in order: 1-20 steps with distinct keys.

<a id="nestedatt--definition--steps"></a>
### Nested Schema for `definition.steps`

- `key` (String) Key of the step, kebab-case, 2-60 characters, unique in the definition.
- `kind` (String) `WORK` does the work; `APPROVAL` approves the work before it. One of `WORK`, `APPROVAL`.
- `name` (String) Name of the step, 1-200 characters.
- `role` (String) The Network role that performs the step (GET /api/roles of the Organisation API). An `APPROVAL` step takes only a role that may approve: APPROVER. A `WORK` step takes only a role that may write: NETWORK_ADMIN, OPERATOR, checked when a version is added; a version added before that check may name another role there. One of `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`.

