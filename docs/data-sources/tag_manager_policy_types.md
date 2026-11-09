---
page_title: "tag_manager_policy_types Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Every tag_manager_policy_type, all pages: a policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin.
---

# das_tag_manager_policy_types (Data Source)

Every tag_manager_policy_type, all pages: a policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin.

API: `GET /api/policy-types`.

## Example Usage

```terraform
data "das_tag_manager_policy_types" "example" {
}
```

## Schema

### Optional

- `platform` (String) Only those on this platform, a platform key (`APP`, `WEB`).
- `status` (List of String) Only policy types with one of these statuses; repeat the parameter for several (`?status=ACTIVE&status=ARCHIVED`). Left out: every status but `ARCHIVED`.

### Read-Only

- `items` (Attributes List) Every item, all pages followed.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `additive` (Boolean) `false`: one matching policy wins; `true`: every matching policy applies.
- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `configuration_types` (Attributes List) Its configuration types, sorted by key.
- `contexts` (List of String) The contexts its policies may take; the first is a new policy's default.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `description` (String) What the type configures; null when none.
- `dimensions` (List of String) The keys of the target dimensions its policies may target, directly and in their attachments, sorted by key.
- `id` (String) Policy type id, `poltype_<ULID>`.
- `key` (String) Key, snake_case, unique in the Network, never changed: the type's name in a resolution.
- `label` (String) Name.
- `platforms` (List of String) The platforms it is on, sorted by key.
- `status` (String) `ACTIVE`; `INACTIVE`: not resolved, greyed out in the panel; `ARCHIVED`: hidden from the default list. One of `ACTIVE`, `INACTIVE`, `ARCHIVED`.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.

<a id="nestedatt--items--configuration_types"></a>
### Nested Schema for `items.configuration_types`

- `fields` (List of String) The field definitions of its configurations' bodies, in the field language of DAS Policies.
- `id` (String) Configuration type id, `cfgtype_<ULID>`.
- `key` (String) Key, snake_case, unique in the policy type.
- `label` (String) Name; null when none (the panel shows the key).

