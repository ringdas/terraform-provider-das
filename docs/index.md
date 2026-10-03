---
page_title: "Provider: DAS"
subcategory: ""
description: |-
  The DAS provider configures a DAS Organization as code: everything the `das` CLI can do. Its resources, data sources, actions and ephemeral resources are derived from the API contracts of the DAS version it was released with.
---

# DAS Provider

The DAS provider configures a DAS Organization as code: everything the `das` CLI can do. Its resources, data sources, actions and ephemeral resources are derived from the API contracts of the DAS version it was released with.

Every resource, data source, action and ephemeral resource below is derived from the API contracts of the DAS release this provider version ships with; a newer DAS adds them to a newer provider version. Requires Terraform 1.11 or later (write-only attributes); actions need Terraform 1.14.

Secrets never enter Terraform state: a client secret is a write-only attribute, a minted credential is an ephemeral resource.

## Example Usage

```terraform
terraform {
  required_providers {
    das = {
      source = "ringdas/das"
    }
  }
}

# The credential comes from DAS_API_KEY or DAS_TOKEN.
provider "das" {
  org_domain = "das.example.com"
}
```

## Contracts

| API | Module | Version |
| --- | --- | --- |
| DAS Audience API | `audience` | 0.0.0 |
| DAS Delivery API | `delivery` | 0.1.0 |
| DAS Organisation API | `org` | 0.1.0 |
| DAS Tag Manager API | `tag-manager` | 0.0.0 |

## Schema

### Optional

- `org_domain` (String) The Organization's domain, e.g. `das.example.com`. The Organisation API answers on it, a product API on `<network>.<product>.<org_domain>`. Environment variable `DAS_ORG_DOMAIN`.
- `scheme` (String) `https` (default), or `http` for a local stack. Environment variable `DAS_ORG_SCHEME`.
- `network` (String) The slug of the Network a product API (Delivery, Audience, Tag Manager) is called for. Omitted, it is the API key's own Network. Environment variable `DAS_NETWORK`.
- `api_key` (String, Sensitive) An API key: a Network's, `das_<networkId>_<secret>`, or the Organization's, `daso_<secret>` (`das_org_api_key`), which administers the Organization layer only - a product API then takes `token`. The tenant of a call is always the credential's. Environment variable `DAS_API_KEY`.
- `token` (String, Sensitive) A bearer token instead of an API key: an Organization session token, or a Network token. Environment variable `DAS_TOKEN`.
- `endpoints` (Map of String) Base URLs by API module (`org`, `delivery`, `audience`, `tag-manager`) instead of the hosts derived from `org_domain`: local development and tests.
