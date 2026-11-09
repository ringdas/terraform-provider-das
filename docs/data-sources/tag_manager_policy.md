---
page_title: "tag_manager_policy Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A policy of the credential's Network (the feature `policies`), written by an operator: where (`inventory`) and for which requests (`target`, groups ORed, rows of a group ANDed) a policy type's configurations apply, saved whole in one mutation. Read by its identity.
---

# das_tag_manager_policy (Data Source)

A policy of the credential's Network (the feature `policies`), written by an operator: where (`inventory`) and for which requests (`target`, groups ORed, rows of a group ANDed) a policy type's configurations apply, saved whole in one mutation. Read by its identity.

API: `GET /api/policies/{policyId}`.

## Example Usage

```terraform
data "das_tag_manager_policy" "example" {
  id = "<policyId>"
}
```

## Schema

### Required

- `id` (String) The policy, `pol_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `attachments` (Attributes List) The configurations it attaches, by position.
- `context` (String) Where it is resolved: `SLOT` or `PAGE`, one of its type's contexts. One of `SLOT`, `PAGE`.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `description` (String) What it is for; null when none.
- `end_at` (String) When it stops applying, ISO-8601 UTC; null: never.
- `inventory` (Attributes List) Where it applies: at least one include.
- `labels` (Attributes List) The labels it carries, by name.
- `name` (String) Name.
- `policy_type_id` (String) Its policy type, `poltype_<ULID>`; never changes.
- `start_at` (String) When it starts applying, ISO-8601 UTC; null: from the start.
- `status` (String) `ACTIVE`; `INACTIVE`: it does not apply; `ARCHIVED`: hidden from the default list. One of `ACTIVE`, `INACTIVE`, `ARCHIVED`.
- `target` (Attributes List) Its target in groups (DNF); empty: every request of its inventory.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.

<a id="nestedatt--attachments"></a>
### Nested Schema for `attachments`

- `configuration_id` (String) The configuration, `cfg_<ULID>`.
- `enabled` (Boolean) false: paused in this policy, the next attachment applies.
- `id` (String) Attachment id, `polcfg_<ULID>`; kept while the policy attaches the same configuration.
- `inventory` (Attributes List) Where in the policy's inventory it applies; empty: all of it.
- `position` (Number, integer) Its place among the policy's attachments, from 1.
- `target` (Attributes List) The sub-target, in the policy's canonical order; empty: the policy's target alone.

<a id="nestedatt--attachments--inventory"></a>
### Nested Schema for `attachments.inventory`

- `exclude` (Boolean) true: everything the path covers is left out.
- `path` (String) `site`, `site/area` or `site/area/slot` by code (`NEWS_APP_IOS/home/right`), `*` for a whole part, `*` alone for the whole network; normalized (no trailing `/*`).

<a id="nestedatt--attachments--target"></a>
### Nested Schema for `attachments.target`

- `dimension` (String) The target dimension's key.
- `op` (String) The operator: `IN`, `NOT_IN`, `GTE`, `LTE` or `BETWEEN`. One of `IN`, `NOT_IN`, `GTE`, `LTE`, `BETWEEN`.
- `target_group` (Number, integer) The OR group of the row.
- `values` (List of String) The values `op` compares with, each shaped by the dimension's `input`: `IN` / `NOT_IN` in the order written (a match ignores it), `GTE` / `LTE` one, `BETWEEN` `[from, to]`.

<a id="nestedatt--inventory"></a>
### Nested Schema for `inventory`

- `exclude` (Boolean) true: everything the path covers is left out.
- `path` (String) `site`, `site/area` or `site/area/slot` by code (`NEWS_APP_IOS/home/right`), `*` for a whole part, `*` alone for the whole network; normalized (no trailing `/*`).

<a id="nestedatt--labels"></a>
### Nested Schema for `labels`

- `id` (String) Label id, `label_<ULID>`.
- `name` (String) Name, unique in the Network regardless of case.

<a id="nestedatt--target"></a>
### Nested Schema for `target`

- `dimension` (String) The target dimension's key.
- `op` (String) The operator: `IN`, `NOT_IN`, `GTE`, `LTE` or `BETWEEN`. One of `IN`, `NOT_IN`, `GTE`, `LTE`, `BETWEEN`.
- `target_group` (Number, integer) The OR group of the row.
- `values` (List of String) The values `op` compares with, each shaped by the dimension's `input`: `IN` / `NOT_IN` in the order written (a match ignores it), `GTE` / `LTE` one, `BETWEEN` `[from, to]`.

