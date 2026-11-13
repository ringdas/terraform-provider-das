---
page_title: "tag_manager_policy_type Resource - das"
subcategory: "DAS Tag Manager"
description: |-
  A policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin. Its key never changes. `platforms` and `configuration_types` are lists replaced whole and answered sorted by key (write them in that order); `dimensions` is a set of target dimension keys (reference `das_tag_manager_policy_target_dimension.<name>.key` so the dimension is created first), and a dimension a live policy of the type targets cannot leave it; each configuration type's `fields` are field definitions in the field language of DAS Policies, one `jsonencode()` per field. Destroy archives it (`ARCHIVED`), refused while a policy or configuration of the type is not archived: a policy type is archived, never deleted.
---

# das_tag_manager_policy_type (Resource)

A policy type of the credential's Network (the feature `policies`), with its platforms, target dimensions and configuration types, written by a Network admin. Its key never changes. `platforms` and `configuration_types` are lists replaced whole and answered sorted by key (write them in that order); `dimensions` is a set of target dimension keys (reference `das_tag_manager_policy_target_dimension.<name>.key` so the dimension is created first), and a dimension a live policy of the type targets cannot leave it; each configuration type's `fields` are field definitions in the field language of DAS Policies, one `jsonencode()` per field. Destroy archives it (`ARCHIVED`), refused while a policy or configuration of the type is not archived: a policy type is archived, never deleted.

API: `POST /api/policy-types` creates, `GET /api/policy-types/{policyTypeId}` reads, `PATCH /api/policy-types/{policyTypeId}` updates. Destroy: A policy type is archivable (P11): it is never deleted, so destroy archives it (refused while its policies or configurations are not archived).

## Example Usage

```terraform
resource "das_tag_manager_policy_type" "example" {
  configuration_types = [{
    fields = ["{\"id\":\"timeoutMs\",\"kind\":\"int\",\"max\":3000,\"min\":100,\"required\":true,\"unit\":\"ms\"}", "{\"default\":true,\"id\":\"sendTmax\",\"kind\":\"bool\"}"]
    key = "prebid_timeout_app"
  }]
  contexts            = ["PAGE"]
  dimensions          = ["app_version", "seg_ab"]
  key                 = "prebid_timeout_app"
  label               = "Prebid timeout"
  platforms           = ["APP"]
}
```

## Schema

### Required

- `configuration_types` (Attributes List) Its configuration types, 1-50, keys distinct (a type without variants has one, keyed like the type). Answered sorted by key.
- `contexts` (List of String) The contexts the type's policies may take, distinct, at least one; the first is a new policy's default. `SLOT`: decided for each slot; `PAGE`: once per page. One still taken by a policy that is not `ARCHIVED` cannot be dropped (409 `IN_USE`).
- `key` (String) Key of the policy type, snake_case, 1-64 characters, unique in the Network and never changed: the type's name in a resolution and between environments. It cannot change: the API offers no update and never deletes it.
- `label` (String) Name, 1-200 characters.
- `platforms` (List of String) The platforms the type is on: platform keys (GET /api/policy-platforms), distinct, at least one. Answered sorted by key.

### Optional

- `additive` (Boolean) `false` (default): one matching policy wins; `true`: every matching policy applies, in the order of decision D31.
- `description` (String) What the type configures, up to 2000 characters. Leave it out for null.
- `dimensions` (Set of String) The target dimensions the type's policies may target, directly and in their attachments; left out, none: keys of the Network's target dimensions (GET /api/policy-target-dimensions), distinct; a key that names none is 400 `unknown_dimension` at `dimensions.<i>`. Answered sorted by key.
- `status` (String) `ACTIVE` (default) or `INACTIVE`: the type is not resolved and the panel greys it out. One of `ACTIVE`, `INACTIVE`.

### Read-Only

- `archived_at` (String) When it was archived, ISO-8601 UTC; null unless `ARCHIVED`.
- `created_at` (String) When it was created, ISO-8601 UTC.
- `created_by` (String) Who created it, audit-grammar actor (`human:<email>`, `service:<name>`).
- `id` (String) The id the API assigned.
- `updated_at` (String) When it last changed, ISO-8601 UTC.
- `updated_by` (String) Who changed it last, audit-grammar actor.

<a id="nestedatt--configuration_types"></a>
### Nested Schema for `configuration_types`

- `fields` (List of String) The schema of the body of its configurations: 0-200 field definitions in the field language of DAS Policies (`{ id, kind, label?, required?, default?, list?, options?, optionsSource?, min?, max?, ... }`; kinds `string`, `int`, `decimal`, `bool`, `section`, `typed`, `map`; `list: true` makes a value a list, of records for a `section`). `optionsSource` names the values a field reads (GET /api/policy-option-sources): a `string` reads a catalogue, `catalogue.<key>`; a `map` reads a mapping, `mapping.<key>`, and holds rows `{ key, cells: [{ to, value? }] }` (e.g. `[{ "key": "320x50", "cells": [{ "to": "GAM" }, { "to": "PREBID" }] }]`). A source the Network has not, or one archived, is 400 `unknown_source`; a `default` is not checked against its catalogue. Stored by the source's id and answered by its key. An invalid definition is 400 at `configurationTypes.<i>.fields.<path>` with the field language's `code`.
- `key` (String) Key of the configuration type, snake_case, 1-64 characters, unique in the policy type; a write matches the stored configuration types by it.
- `label` (String) Name shown in the panel, 1-200 characters; left out, none (the panel shows the key).

## Import

```terraform
import {
  to = das_tag_manager_policy_type.example
  id = "<policyTypeId>"
}
```
