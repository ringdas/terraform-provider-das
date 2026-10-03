---
page_title: "network_api_key Resource - das"
subcategory: "Organization"
description: |-
  An API key of the Network named by `network_id`, minted by the Organization. Its token never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version. Destroy revokes it and deletes the stored token.
---

# das_network_api_key (Resource)

An API key of the Network named by `network_id`, minted by the Organization. Its token never enters Terraform state: the API writes it to the installation's secret store and `secret_ref` names the version. Destroy revokes it and deletes the stored token.

API: `POST /api/networks/{networkId}/api-keys` creates, `GET /api/networks/{networkId}/api-keys/{keyId}` reads, `DELETE /api/networks/{networkId}/api-keys/{keyId}` deletes.

## Example Usage

```terraform
resource "das_network_api_key" "example" {
  name       = "deploy"
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
  role       = "OPERATOR"
  scopes     = ["das:delivery:lineitems:write"]
}
```

## Schema

### Required

- `name` (String) Name of the key; audit events show it as `service:<name>`. Changing it replaces the resource.
- `network_id` (String) The Network administered, `net_<ULID>`. Organization layer only: an Organization admin managing that Network. Changing it replaces the resource.
- `role` (String) Coarse role of the key. One of `VIEWER`, `OPERATOR`, `ADMIN`. Changing it replaces the resource.
- `scopes` (List of String) Scopes `das:<product>:<resource>:<action>`. Changing it replaces the resource.

### Optional

- `expires_at` (String) Optional: after this instant (ISO-8601 UTC, in the future) the key authenticates nothing (401 `CREDENTIAL_EXPIRED`). Omit for no expiry. Changing it replaces the resource.

### Read-Only

- `created_at` (String) Creation instant, ISO-8601 UTC.
- `created_by` (String) Actor who minted the key; null for keys minted before it was recorded.
- `id` (String) The id the API assigned.
- `revoked_at` (String) When the key was revoked, ISO-8601 UTC; null while it is not.
- `secret_ref` (String) Where the token was delivered (minted with `secretDelivery: SECRET_STORE`): a Secrets Manager reference `<ARN>#<version id>` to exactly that version; null when it was answered once. Deleted from the store when the key is revoked.

## Import

```terraform
import {
  to = das_network_api_key.example
  id = "<networkId>/<keyId>"
}
```
