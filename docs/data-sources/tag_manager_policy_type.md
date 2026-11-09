---
page_title: "tag_manager_policy_type Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  A policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin. Read by its identity.
---

# das_tag_manager_policy_type (Data Source)

A policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin. Read by its identity.

API: `GET /api/policy-types/{policyTypeId}`.

## Example Usage

```terraform
data "das_tag_manager_policy_type" "example" {
  id = "<policyTypeId>"
}
```

## Schema

### Required

- `id` (String) The policy type, `poltype_<ULID>`. One of another Network is 404 `NOT_FOUND`.

### Read-Only

- `additive` (Boolean) `false`: one matching policy wins; `true`: every matching policy applies.
- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `configuration_types` (Attributes List) Its configuration types, sorted by key.
- `contexts` (List of String) The contexts its policies may take; the first is a new policy's default.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `description` (String) What the type configures; null when none.
- `dimensions` (List of String) The keys of the target dimensions its policies may target, directly and in their attachments, sorted by key.
- `key` (String) Key, snake_case, unique in the Network, never changed: the type's name in a resolution.
- `label` (String) Name.
- `platforms` (List of String) The platforms it is on, sorted by key.
- `status` (String) `ACTIVE`; `INACTIVE`: not resolved, greyed out in the panel; `ARCHIVED`: hidden from the default list. One of `ACTIVE`, `INACTIVE`, `ARCHIVED`.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.

<a id="nestedatt--configuration_types"></a>
### Nested Schema for `configuration_types`

- `fields` (List of String) The field definitions of its configurations' bodies, in the field language of DAS Policies.
- `id` (String) Configuration type id, `cfgtype_<ULID>`.
- `key` (String) Key, snake_case, unique in the policy type.
- `label` (String) Name; null when none (the panel shows the key).

