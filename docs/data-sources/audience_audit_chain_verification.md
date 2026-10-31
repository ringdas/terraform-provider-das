---
page_title: "audience_audit_chain_verification Data Source - das"
subcategory: "DAS Audience"
description: |-
  A verification of a segment of the Network's audit chain in DAS Audience, computed at read.
---

# das_audience_audit_chain_verification (Data Source)

A verification of a segment of the Network's audit chain in DAS Audience, computed at read.

API: `GET /api/audit-chain/verification`.

## Example Usage

```terraform
data "das_audience_audit_chain_verification" "example" {
}
```

## Schema

### Optional

- `from_seq` (Number, integer) The first link to verify (default 1, the genesis).

### Read-Only

- `chain` (String) `org`, or `network:<networkId>`.
- `checked_from_seq` (Number, integer) The first seq verified.
- `checked_through_seq` (Number, integer) The last seq verified without a break.
- `first_break` (Attributes) The first link that does not hold; null when intact.
- `head_hash` (String) The head's hash now.
- `head_seq` (Number, integer) The chain's head now.
- `intact` (Boolean) Whether every link of the segment holds (and, from seq 1, the genesis seal).
- `next_from_seq` (Number, integer) Where the next segment starts; null at the head or at a break.

<a id="nestedatt--first_break"></a>
### Nested Schema for `first_break`

- `reason` (String) HASH_MISMATCH: a field changed behind the application; PREV_HASH_MISMATCH: a link does not name the one before; SEQ_GAP: a link was deleted; GENESIS_MISSING: the first link is no genesis; LEGACY_DIGEST_MISMATCH: a row before the chain changed; TRUNCATED: the newest links were deleted. One of `HASH_MISMATCH`, `PREV_HASH_MISMATCH`, `SEQ_GAP`, `GENESIS_MISSING`, `LEGACY_DIGEST_MISMATCH`, `TRUNCATED`.
- `seq` (Number, integer) The seq where the chain breaks.

