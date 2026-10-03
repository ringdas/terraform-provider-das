---
page_title: "tag_manager_session Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  Who the credential is in DAS Tag Manager for its Network.
---

# das_tag_manager_session (Data Source)

Who the credential is in DAS Tag Manager for its Network.

API: `GET /api/session`.

## Example Usage

```terraform
data "das_tag_manager_session" "example" {
}
```

## Schema

### Read-Only

- `actor` (String) Audit-grammar actor, `human:<email>`.
- `features` (List of String) Ids of the features effectively on for the Network (the UI shows what they own); switched in the Organisation API.
- `network` (Attributes) The credential's Network, as the Organisation API records it.
- `role` (String) Coarse role of the credential. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `scopes` (List of String) Scopes `das:<product>:<resource>:<action>`.

<a id="nestedatt--network"></a>
### Nested Schema for `network`

- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `slug` (String) Host label of the Network.
- `timezone` (String) IANA time zone where the Network's days begin; dates are shown in it.

