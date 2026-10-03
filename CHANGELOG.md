# Changelog

All notable changes to `terraform-provider-das`. The Release PR writes each section from the `.changeset/*.md` files (docs/engineering/release-process.md); do not edit it by hand.

## [0.1.0] - 2026-10-03

### Added

- **Minted secrets delivered to the secret store**: minting an API key (on
  both surfaces) or a SCIM token takes `secretDelivery`; `SECRET_STORE`
  writes the token to the installation's secret store
  (`<prefix>/api-keys/<keyId>`, `<prefix>/scim-tokens/<tokenId>`) and answers
  only its version-pinned `secretRef`, never the secret (503 `NOT_CONFIGURED`
  without a store); revoking deletes it. Keys and SCIM tokens answer
  `secretRef`, and read back by id (`GET /api/network/api-keys/{keyId}`,
  `GET /api/networks/{networkId}/api-keys/{keyId}`,
  `GET /api/scim-tokens/{tokenId}`), so Terraform manages them without a
  secret in its state. (#120)
- **Organization API keys: the Organization layer's machine credential**.
  `/api/org-api-keys` mints, lists, reads and revokes keys `daso_<secret>`
  that hold one Organization role and act with exactly what a person holding
  it may do - never more than their minter (only an owner mints or revokes an
  owner key), never as a person or in a Network. Only the secret's sha256 is
  stored; the token is answered once or delivered to the secret store
  (`secretDelivery: SECRET_STORE`). The DAS Terraform provider manages them
  as `das_org_api_key` and accepts one as `api_key` for the Organization
  layer; the das CLI takes one in `DAS_API_KEY` for `das org` commands and
  refuses it for a product command before any request. (#119)
- **The DAS Terraform provider, `ringdas/das`: everything the CLI can do,
  as code**. A Go provider (terraform-plugin-framework) derived from the API
  contracts like the `das` CLI: each API declares its Terraform surface
  (`x-das-terraform` in its contract, `terraformSurface.ts`), and the provider
  builds a resource, a data source, an action or an ephemeral resource from
  each entry - Networks and their features and settings, people, groups,
  memberships and role bindings, sign-in, identity providers (the client
  secret is write-only, never in state), email senders and DAS Delivery line
  items, plus every read and the test-email actions. The IaC coverage gate now
  holds every operation of every API to the surface or to an `IAC_DEBT` entry
  with an open issue, replacing the manifest verdicts; the provider's own
  tests hold its embedded contracts, schema, lifecycle and documentation to
  the contracts, and a changed contract needs a provider changeset. Inviting
  a suspended person again reinstates them (the same user, audited as
  `user.reinstate`), so a destroyed and re-applied `das_user` works. (#118)

### Fixed

- **Released to the Terraform and OpenTofu registries**: the provider ships
  under MPL-2.0 with its release configuration (goreleaser, the registry
  manifest, the public signing key); every version the Release PR bumps is
  mirrored to `ringdas/terraform-provider-das`, tagged and built there into
  signed zips the registries read. (#121)
