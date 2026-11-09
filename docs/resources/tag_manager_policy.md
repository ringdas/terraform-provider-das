---
page_title: "tag_manager_policy Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A policy of the credential's Network (the feature `policies`), written by an operator: where (`inventory`) and for which requests (`target`, groups ORed, rows of a group ANDed) a policy type's configurations apply, saved whole in one mutation. Its policy type never changes (a new one replaces it). `inventory` rows are paths by code (`site`, `site/area`, `site/area/slot`, `*` for a whole part, `*` alone for the whole network), written normalized (`SITE`, not `SITE/*`) so the plan stays empty. `labels`, `inventory` and `target` (and each attachment's `target` and `inventory`) are sets, written in any order and replaced whole: each label with both its `id` and its `name`, each target row with its `op` and its `values` (kept in your order, which a match ignores), each value as one `jsonencode()`. `attachments` is a list in your order, each naming a configuration by id (a configuration is its own resource). An attachment's `enabled`, `target` and `inventory` may be left out for their defaults. Destroy archives it (`ARCHIVED`): a policy is archived, never deleted.
---

# das_tag_manager_policy (Resource)

A policy of the credential's Network (the feature `policies`), written by an operator: where (`inventory`) and for which requests (`target`, groups ORed, rows of a group ANDed) a policy type's configurations apply, saved whole in one mutation. Its policy type never changes (a new one replaces it). `inventory` rows are paths by code (`site`, `site/area`, `site/area/slot`, `*` for a whole part, `*` alone for the whole network), written normalized (`SITE`, not `SITE/*`) so the plan stays empty. `labels`, `inventory` and `target` (and each attachment's `target` and `inventory`) are sets, written in any order and replaced whole: each label with both its `id` and its `name`, each target row with its `op` and its `values` (kept in your order, which a match ignores), each value as one `jsonencode()`. `attachments` is a list in your order, each naming a configuration by id (a configuration is its own resource). An attachment's `enabled`, `target` and `inventory` may be left out for their defaults. Destroy archives it (`ARCHIVED`): a policy is archived, never deleted.

API: `POST /api/policies` creates, `GET /api/policies/{policyId}` reads, `PATCH /api/policies/{policyId}` updates. Destroy: A policy is archivable (P11): it is never deleted, so destroy archives it.

## Example Usage

```terraform
resource "das_tag_manager_policy" "example" {
  attachments    = [{
    configuration_id = "cfg_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  }]
  inventory      = [{
    path = "*"
  }]
  name           = "Fast timeout on the news apps"
  policy_type_id = "poltype_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
}
```

## Schema

### Required

- `inventory` (Attributes Set) Where the policy applies: at least one include (400 `required` at `inventory`), excludes cut what their paths cover out; replaces every row.
- `name` (String) Name, 1-200 characters.
- `policy_type_id` (String) The policy type, `poltype_<ULID>`, not `ARCHIVED` (409 `ARCHIVED`); never changed. It cannot change: the API offers no update and never deletes it.

### Optional

- `attachments` (Attributes List) The configurations it attaches, in order (the first is position 1); replaces them all, keeping the id of each attachment of a configuration still attached.
- `context` (String) Where it is resolved: one of the policy type's contexts (400 `not_in_options`); default the type's first. One of `SLOT`, `PAGE`.
- `description` (String) What it is for, up to 2000 characters. Leave it out for null.
- `end_at` (String) When it stops applying, ISO-8601, after `startAt` (400 `range_order`); null or left out: never. Leave it out for null.
- `labels` (Attributes Set) The labels it carries; replaces them all.
- `start_at` (String) When it starts applying, ISO-8601; null or left out: at once. Leave it out for null.
- `status` (String) `ACTIVE` (default) or `INACTIVE`. One of `ACTIVE`, `INACTIVE`.
- `target` (Attributes Set) The target, in groups (DNF); empty: every request of the inventory; replaces every row.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.

<a id="nestedatt--inventory"></a>
### Nested Schema for `inventory`

- `exclude` (Boolean) true: everything the path covers is left out; false (default): it is included. Defaults to `false`.
- `path` (String) Where, by code: `site`, `site/area` or `site/area/slot` (`NEWS_APP_IOS/home/right`), 1-3 parts joined by `/`; `*` stands for a whole part (`NEWS_APP_IOS/*/right`: the slot right of every area), `*` alone for the whole network; a shorter path covers everything under it. A part is `*` or a code without `/`, spaces and control or invisible characters, never a piece of one (`APP_*`): 400 `invalid_format` at it. Stored normalized: trailing `/*` parts dropped (`SITE/*` is `SITE`, `*/*` is `*`). Only the format is checked, not whether the site, area or slot exists. Up to 200 characters.

<a id="nestedatt--attachments"></a>
### Nested Schema for `attachments`

- `configuration_id` (String) A configuration of the policy's type, `cfg_<ULID>`, not `ARCHIVED` (409 `ARCHIVED` at it; `INACTIVE` is skipped by the resolver), attached once.
- `enabled` (Boolean) false: paused in this policy, the next attachment applies; true (default). Defaults to `true`.
- `inventory` (Attributes Set) Where in the policy's inventory the attachment applies: an include must be covered by an include of the policy and by none of its excludes (400 `outside_policy_inventory` at its `path`); empty (default): all of it. Defaults to `[]`.
- `target` (Attributes Set) The attachment's sub-target, groups ORed as the policy's; empty (default): the policy's target alone. Defaults to `[]`.

<a id="nestedatt--attachments--inventory"></a>
### Nested Schema for `attachments.inventory`

- `exclude` (Boolean) true: everything the path covers is left out; false (default): it is included. Defaults to `false`.
- `path` (String) Where, by code: `site`, `site/area` or `site/area/slot` (`NEWS_APP_IOS/home/right`), 1-3 parts joined by `/`; `*` stands for a whole part (`NEWS_APP_IOS/*/right`: the slot right of every area), `*` alone for the whole network; a shorter path covers everything under it. A part is `*` or a code without `/`, spaces and control or invisible characters, never a piece of one (`APP_*`): 400 `invalid_format` at it. Stored normalized: trailing `/*` parts dropped (`SITE/*` is `SITE`, `*/*` is `*`). Only the format is checked, not whether the site, area or slot exists. Up to 200 characters.

<a id="nestedatt--attachments--target"></a>
### Nested Schema for `attachments.target`

- `dimension` (String) The target dimension's key (`seg_ab`, `app_version`), one of the Network's target dimensions (GET /api/policy-target-dimensions); one of the policy type's dimensions (400 `dimension_not_assigned`).
- `op` (String) The operator, one of the dimension's `operators` (400 `operator_not_allowed`): the request's value is one of `values` (`IN`), none of them (`NOT_IN`), at least (`GTE`) or at most (`LTE`) the one value, or within `[from, to]` (`BETWEEN`); a list dimension's request matches IN when it holds any of `values`. One of `IN`, `NOT_IN`, `GTE`, `LTE`, `BETWEEN`.
- `target_group` (Number, integer) The OR group the row belongs to, 1-100: rows of one group must all match, any group matching is enough.
- `values` (List of String) The values `op` compares with, each shaped by its dimension's `input` (GET /api/policy-target-dimensions; an item of a list dimension): `IN` / `NOT_IN` one or more distinct (one for a `bool`), kept and answered in the order written, which a match ignores; `GTE` / `LTE` one; `BETWEEN` `[from, to]`, from ≤ to for a number. A value of a dimension whose `input` reads `optionsSource: catalogue.<key>` is a value of that catalogue (GET /api/policy-catalogues). 1-500; an invalid one is 400 at `<row>.values.<i>` with the field language's `code`.

<a id="nestedatt--labels"></a>
### Nested Schema for `labels`

- `id` (String) An existing label, `label_<ULID>` (400 `not_found` otherwise).
- `name` (String) A label's name: with `id`, exactly that label's name (400 `name_mismatch` otherwise); alone, the label of that name regardless of case, created with the save when the Network has none.

<a id="nestedatt--target"></a>
### Nested Schema for `target`

- `dimension` (String) The target dimension's key (`seg_ab`, `app_version`), one of the Network's target dimensions (GET /api/policy-target-dimensions); one of the policy type's dimensions (400 `dimension_not_assigned`).
- `op` (String) The operator, one of the dimension's `operators` (400 `operator_not_allowed`): the request's value is one of `values` (`IN`), none of them (`NOT_IN`), at least (`GTE`) or at most (`LTE`) the one value, or within `[from, to]` (`BETWEEN`); a list dimension's request matches IN when it holds any of `values`. One of `IN`, `NOT_IN`, `GTE`, `LTE`, `BETWEEN`.
- `target_group` (Number, integer) The OR group the row belongs to, 1-100: rows of one group must all match, any group matching is enough.
- `values` (List of String) The values `op` compares with, each shaped by its dimension's `input` (GET /api/policy-target-dimensions; an item of a list dimension): `IN` / `NOT_IN` one or more distinct (one for a `bool`), kept and answered in the order written, which a match ignores; `GTE` / `LTE` one; `BETWEEN` `[from, to]`, from ≤ to for a number. A value of a dimension whose `input` reads `optionsSource: catalogue.<key>` is a value of that catalogue (GET /api/policy-catalogues). 1-500; an invalid one is 400 at `<row>.values.<i>` with the field language's `code`.

## Import

```terraform
import {
  to = das_tag_manager_policy.example
  id = "<policyId>"
}
```
