---
page_title: "session Data Source - das"
subcategory: "Organization"
description: |-
  The Organization session of the provider's credential and the Networks it may enter.
---

# das_session (Data Source)

The Organization session of the provider's credential and the Networks it may enter.

API: `GET /api/session`.

## Example Usage

```terraform
data "das_session" "example" {
}
```

## Schema

### Read-Only

- `actor` (String) The actor audit events record, `human:<email>`.
- `groups` (Attributes List) The user's groups: memberships and groups mapped to the identity provider's groups of this session.
- `networks` (Attributes List) The Networks the user may enter, by name.
- `role` (String) The highest coarse rank of those roles; null with none. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `roles` (List of String) The Organization roles held, directly or through a group.
- `scopes` (List of String) The scopes those roles grant on the Organization layer.
- `user` (Attributes) The user.

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

- `id` (String) Group id.
- `name` (String) Group name.

<a id="nestedatt--networks"></a>
### Nested Schema for `networks`

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `default_currency` (String) ISO-4217 code of the Network's money.
- `display_name` (String) Human-readable name.
- `id` (String) Network id, `net_<ULID>`.
- `network_code` (String) The number people use for the Network, `NT-` and seven digits.
- `roles` (List of String) The Network roles the user holds there, directly or through a group.
- `slug` (String) Host label of the Network.
- `timezone` (String) IANA time zone of the Network.
- `updated_at` (String) Last change, ISO-8601 UTC.

<a id="nestedatt--user"></a>
### Nested Schema for `user`

- `display_name` (String) Name.
- `email` (String) Email.
- `id` (String) User id.
- `status` (String) Status. One of `INVITED`, `ACTIVE`, `SUSPENDED`.

