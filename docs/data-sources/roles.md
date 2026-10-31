---
page_title: "roles Data Source - das"
subcategory: "Organization"
description: |-
  The role catalog of both layers.
---

# das_roles (Data Source)

The role catalog of both layers.

API: `GET /api/roles`.

## Example Usage

```terraform
data "das_roles" "example" {
}
```

## Schema

### Read-Only

- `items` (Attributes List) The roles.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

- `description` (String) What the role is for.
- `layer` (String) `ORG`: held on the Organization layer; `NETWORK`: held in one Network. One of `ORG`, `NETWORK`.
- `name` (String) The role. One of `ORG_OWNER`, `ORG_ADMIN`, `ORG_VIEWER`, `ORG_AUDITOR`, `NETWORK_ADMIN`, `OPERATOR`, `VIEWER`, `APPROVER`, `AUDITOR`.
- `permissions` (Attributes List) Every permission of the registry the role gives, with its actions.
- `rank` (String) The coarse rank the access classes compare. One of `VIEWER`, `OPERATOR`, `ADMIN`.
- `scopes` (List of String) The scopes it grants (a `*` resource covers every resource of the product).

