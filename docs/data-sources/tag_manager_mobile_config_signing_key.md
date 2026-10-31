---
page_title: "tag_manager_mobile_config_signing_key Data Source - das"
subcategory: "DAS Tag Manager"
description: |-
  The public half of the Organization's Ed25519 key that signs mobile configuration (the feature `mobile-config`): what an app embeds to verify a bundle.
---

# das_tag_manager_mobile_config_signing_key (Data Source)

The public half of the Organization's Ed25519 key that signs mobile configuration (the feature `mobile-config`): what an app embeds to verify a bundle.

API: `GET /api/mobile-config/signing-key`.

## Example Usage

```terraform
data "das_tag_manager_mobile_config_signing_key" "example" {
}
```

## Schema

### Read-Only

- `algorithm` (String) The signature algorithm: Ed25519 (RFC 8032). One of `ED25519`.
- `fingerprint` (String) Hex SHA-256 of the raw public key; a release names the key that signed it by it.
- `public_key` (String) The raw 32-byte Ed25519 public key, base64.
- `public_key_pem` (String) The same key as a PEM SubjectPublicKeyInfo (`-----BEGIN PUBLIC KEY-----`).

