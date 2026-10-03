---
page_title: "delivery_session Data Source - das"
subcategory: "DAS Delivery"
description: |-
  Who the credential is in DAS Delivery for its Network.
---

# das_delivery_session (Data Source)

Who the credential is in DAS Delivery for its Network.

API: `GET /api/session`.

## Example Usage

```terraform
data "das_delivery_session" "example" {
}
```

## Schema

### Read-Only

- `accent` (Attributes) The accent of the Delivery interface for the Network (the setting `delivery.accent-color`).
- `actor` (String) Audit-grammar actor, `human:<email>`.
- `features` (List of String) Ids of the features effectively on for the Network (the UI shows what they own); switched in the Organisation API.
- `network` (Attributes) The credential's Network, as the Organisation API records it.
- `networks` (Attributes List) The Organization's Networks the UI's switcher offers, by display name; switching navigates to that Network's host and signs in there.
- `role` (String) Coarse role of the credential. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `roles` (List of String) The Network roles the credential's Network token carries (GET /api/roles of the Organisation API describes each); empty for other credentials.
- `scopes` (List of String) Scopes `das:<product>:<resource>:<action>`.

<a id="nestedatt--accent"></a>
### Nested Schema for `accent`

- `color` (String) `#rrggbb`; null = the DAS default. The UI derives it for both themes and keeps the default when it misses WCAG AA.
- `source` (String) Where it comes from: the Network's own value, the Organization's, or the built-in default. One of `NETWORK`, `ORGANIZATION`, `DEFAULT`.

<a id="nestedatt--network"></a>
### Nested Schema for `network`

- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits.
- `slug` (String) Host label of the Network.
- `timezone` (String) IANA time zone where the Network's days begin; dates are shown in it.

<a id="nestedatt--networks"></a>
### Nested Schema for `networks`

- `display_name` (String) Human-readable name.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits.
- `slug` (String) Host label of the Network: `<slug>.delivery.<org domain>`.

