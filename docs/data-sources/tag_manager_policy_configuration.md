---
page_title: "tag_manager_policy_configuration Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach. Read by its identity.
---

# das_tag_manager_policy_configuration (Data Source)

A configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach. Read by its identity.

API: `GET /api/policy-configurations/{configurationId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_configuration" "example" {
  id = "<configurationId>"
}
```

## Schema

### Required

- `id` (String) The configuration, `cfg_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `body_json` (String) Its values, shaped by its configuration type's `fields`. JSON-encoded (`jsondecode()` it).
- `configuration_type_id` (String) Its configuration type, `cfgtype_<ULID>`; never changes.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `description` (String) What it is for; null when none.
- `name` (String) Name.
- `policy_type_id` (String) Its policy type, `poltype_<ULID>`; never changes.
- `status` (String) `ACTIVE`; `INACTIVE`: its attachments are skipped, the policies still apply; `ARCHIVED`: hidden from the default list, attached by no live policy. One of `ACTIVE`, `INACTIVE`, `ARCHIVED`.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.
- `usage` (Attributes) The policies attaching it.

<a id="nestedatt--usage"></a>
### Nested Schema for `usage`

- `active_policies` (Number, integer) Of those, the `ACTIVE` ones.
- `policies` (Number, integer) Policies that are not `ARCHIVED` and attach it, enabled or paused.

