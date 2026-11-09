---
page_title: "tag_manager_policy_configurations Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_configuration, all pages: a configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach.
---

# das_tag_manager_policy_configurations (Data Source)

Every tag_manager_policy_configuration, all pages: a configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach.

API: `GET /api/policy-configurations`.

## Example Usage

```terraform
data "das_tag_manager_policy_configurations" "example" {
}
```

## Schema

### Optional

- `configuration_type_id` (String) Only those of this configuration type, `cfgtype_<ULID>`.
- `policy_type_id` (String) Only those of this policy type, `poltype_<ULID>`.
- `q` (String) Only configurations whose name contains this text, regardless of case (`%` and `_` are plain characters), or whose id it is; 1-200 characters.
- `status` (List of String) Only configurations with one of these statuses; repeat the parameter for several (`?status=ACTIVE&status=ARCHIVED`). Left out: every status but `ARCHIVED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `body_json` (String) Its values, shaped by its configuration type's `fields`. JSON-encoded (`jsondecode()` it).
- `configuration_type_id` (String) Its configuration type, `cfgtype_<ULID>`; never changes.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `description` (String) What it is for; null when none.
- `id` (String) Configuration id, `cfg_<ULID>`.
- `name` (String) Name.
- `policy_type_id` (String) Its policy type, `poltype_<ULID>`; never changes.
- `status` (String) `ACTIVE`; `INACTIVE`: its attachments are skipped, the policies still apply; `ARCHIVED`: hidden from the default list, attached by no live policy. One of `ACTIVE`, `INACTIVE`, `ARCHIVED`.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.
- `usage` (Attributes) The policies attaching it.

<a id="nestedatt--items--usage"></a>
### Nested Schema for `items.usage`

- `active_policies` (Number, integer) Of those, the `ACTIVE` ones.
- `policies` (Number, integer) Policies that are not `ARCHIVED` and attach it, enabled or paused.

