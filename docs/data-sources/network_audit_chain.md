---
page_title: "network_audit_chain Data Source - das"
subcategory: "Organization"
description: |-
  The hash chain of the audit trail of the Network named by `network_id`.
---

# das_network_audit_chain (Data Source)

The hash chain of the audit trail of the Network named by `network_id`.

API: `GET /api/networks/{networkId}/audit-chain`.

## Example Usage

```terraform
data "das_network_audit_chain" "example" {
  network_id = "<networkId>"
}
```

## Schema

### Required

- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network.

### Read-Only

- `chain` (String) `org`, or `network:<networkId>`.
- `head_hash` (String) The newest link's hash; null while the chain is empty.
- `head_seq` (Number, integer) The newest link's seq; 0 while the chain is empty.
- `last_check` (Attributes) The newest verification; null while none ran.

<a id="nestedatt--last_check"></a>
### Nested Schema for `last_check`

- `break_reason` (String) Why it did not hold (AuditChainBreak.reason); null when intact.
- `break_seq` (Number, integer) The first link that did not hold; null when intact.
- `checked_at` (String) When it ran, ISO-8601 UTC.
- `from_seq` (Number, integer) The first link it verified (1 = the whole chain).
- `head_hash` (String) The head's hash when it ran (also logged by the duty, a witness outside the database).
- `head_seq` (Number, integer) The chain's head when it ran.
- `status` (String) The verdict. One of `INTACT`, `BROKEN`.
- `through_seq` (Number, integer) The last link it verified; 0 for an empty chain.

