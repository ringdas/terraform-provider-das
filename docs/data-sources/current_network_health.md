---
page_title: "current_network_health Data Source - das"
subcategory: "Organization"
description: |-
  The health checks of the credential's Network.
---

# das_current_network_health (Data Source)

The health checks of the credential's Network.

API: `GET /api/network/health`.

## Example Usage

```terraform
data "das_current_network_health" "example" {
}
```

## Schema

### Read-Only

- `checked_at` (String) When the checks were evaluated, ISO-8601 UTC.
- `checks` (Attributes List) Every registered check, passing or failing.
- `network_id` (String) The Network.
- `status` (String) The worst failing severity: `CRITICAL` = evidence something is broken, `WARNING` = a prediction; `OK` when every check passes. One of `OK`, `WARNING`, `CRITICAL`.

<a id="nestedatt--checks"></a>
### Nested Schema for `checks`

- `domain` (String) The domain that owns the check. One of `delivery`, `audience`, `tag-manager`, `platform`.
- `evidence` (List of String) Ids or facts that show the failure; empty when passing.
- `failing` (Boolean) Whether the check fails now.
- `id` (String) Check id, `<domain>.<name>`.
- `remediation` (String) Where the fix is described.
- `severity` (String) `CRITICAL` = evidence, `WARNING` = prediction. One of `CRITICAL`, `WARNING`.
- `summary` (String) The outcome in one line.
- `title` (String) What the check watches.

