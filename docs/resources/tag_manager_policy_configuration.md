---
page_title: "tag_manager_policy_configuration Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach. Its policy type and configuration type never change (a new one replaces it). `body_json` is the body, an object shaped by the configuration type's `fields`, as one `jsonencode()`. Destroy archives it (`ARCHIVED`), refused while a policy that is not archived attaches it: a configuration is archived, never deleted.
---

# das_tag_manager_policy_configuration (Resource)

A configuration of the credential's Network (the feature `policies`), written by an operator: values of one of a policy type's configuration types, which policies attach. Its policy type and configuration type never change (a new one replaces it). `body_json` is the body, an object shaped by the configuration type's `fields`, as one `jsonencode()`. Destroy archives it (`ARCHIVED`), refused while a policy that is not archived attaches it: a configuration is archived, never deleted.

API: `POST /api/policy-configurations` creates, `GET /api/policy-configurations/{configurationId}` reads, `PATCH /api/policy-configurations/{configurationId}` updates. Destroy: A configuration is archivable (P11): it is never deleted, so destroy archives it (refused while a policy that is not archived attaches it).

## Example Usage

```terraform
resource "das_tag_manager_policy_configuration" "example" {
  body_json             = "{\"sendTmax\":true,\"timeoutMs\":2000}"
  configuration_type_id = "cfgtype_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
  name                  = "Standard 2000 ms"
  policy_type_id        = "poltype_01K6P6ZQ4N2V3W4X5Y6Z7A8B9C"
}
```

## Schema

### Required

- `body_json` (String) The configuration's values: a JSON object shaped by its configuration type's `fields` (GET /api/policy-types: the field language of DAS Policies; a `section` is a nested object, a field hidden by `showIf` is left out, a field with a `default` may be). A value of an `optionsSource: catalogue.<key>` field is the configuration's own copy of a text, not checked against that catalogue (GET /api/policy-catalogues offers its values); a `map` field's rows are `{ key, cells: [{ to, value? }] }`, each cell with a value of the mapping's cell type when it has one (neither the catalogues nor the mapping's rows are checked). An invalid value is 400 at `body.<path>` with the field language's `code`. JSON-encoded (`jsondecode()` it).
- `configuration_type_id` (String) One of the policy type's configuration types, `cfgtype_<ULID>` (400 otherwise); never changed. It cannot change: the API offers no update and never deletes it.
- `name` (String) Name, 1-200 characters.
- `policy_type_id` (String) The policy type, `poltype_<ULID>`, not `ARCHIVED` (409 `ARCHIVED`); never changed. It cannot change: the API offers no update and never deletes it.

### Optional

- `description` (String) What it is for, up to 2000 characters. Leave it out for null.
- `status` (String) `ACTIVE` (default) or `INACTIVE`: its attachments are skipped, the policies still apply. One of `ACTIVE`, `INACTIVE`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.
- `usage` (Attributes) The policies attaching it.

<a id="nestedatt--usage"></a>
### Nested Schema for `usage`

- `active_policies` (Number, integer) Of those, the `ACTIVE` ones.
- `policies` (Number, integer) Policies that are not `ARCHIVED` and attach it, enabled or paused.

## Import

```terraform
import {
  to = das_tag_manager_policy_configuration.example
  id = "<configurationId>"
}
```
