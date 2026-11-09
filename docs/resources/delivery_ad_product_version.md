---
page_title: "delivery_ad_product_version Resource - das"
subcategory: "DAS Delivery"
description: |-
  A version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps. A version is immutable and becomes the ad product's current version when it is created: a changed definition is refused at plan - add a new resource instead - and destroy only forgets it.
---

# das_delivery_ad_product_version (Resource)

A version of an ad product's definition (the feature `ad-product-catalog`): goals, billing models, creative source and process steps. A version is immutable and becomes the ad product's current version when it is created: a changed definition is refused at plan - add a new resource instead - and destroy only forgets it.

API: `POST /api/ad-products/{adProductId}/versions` creates, `GET /api/ad-products/{adProductId}/versions/{version}` reads. Destroy: A version is immutable and is never deleted: the ad product keeps its history.

## Example Usage

```terraform
resource "das_delivery_ad_product_version" "example" {
  ad_product_id = "adp_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  definition    = {
    billing_models = ["CPM", "FLAT_FEE"]
    creative_source = "SUPPLIED"
    goals = ["AWARENESS", "REACH"]
    steps = [{
      key = "brief"
      kind = "WORK"
      name = "Take the brief"
      role = "OPERATOR"
    }, {
      key = "approve"
      kind = "APPROVAL"
      name = "Approve the plan"
      role = "APPROVER"
    }]
  }
}
```

## Schema

### Required

- `ad_product_id` (String) The ad product, `adp_<ULID>`. One of another Network is 404 `NOT_FOUND`. It cannot change: the API offers no update and never deletes it.
- `definition` (Attributes) The definition of the new version: goals, billing models, creative source and process steps. It cannot change: the API offers no update and never deletes it.

### Read-Only

- `created_at` (String) When the version was added, ISO-8601 UTC.
- `id` (String) The identity, adProductId/version joined by `/` (the import id).
- `version` (Number, integer) The version number: 1 for the first, then one more for each.

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

## Import

```terraform
import {
  to = das_delivery_ad_product_version.example
  id = "<adProductId>/<version>"
}
```
