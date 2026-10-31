# Changelog

All notable changes to `terraform-provider-das`. The Release PR writes each section from the `.changeset/*.md` files (docs/engineering/release-process.md); do not edit it by hand.

## [0.2.0] - 2026-10-31

### Added

- **Ad product catalog in DAS Delivery**: a Network records the kinds of
  advertising it sells and how each is delivered. The feature
  `ad-product-catalog` (off by default) adds `/api/ad-products`: an ad product
  (key unique in the Network and never changed, name, description, status
  `ENABLED` / `DISABLED`, the current version) and its immutable, numbered
  versions, each a validated definition of goals, billing models, creative
  source and process steps with roles and approvals; a new version becomes the
  current one at once. The roles holding a Delivery read scope (`VIEWER`,
  `OPERATOR`, `APPROVER`, `NETWORK_ADMIN`) read the catalog, `NETWORK_ADMIN`
  writes it (the new permission `das:delivery:ad-products`), and every change
  is audited.
  The `das` CLI gains `das delivery ad-products ...` and the Terraform provider
  the resources `das_delivery_ad_product` and `das_delivery_ad_product_version`
  with their data sources; the provider now reads a numeric identity (a version
  number) back from a create's answer. (#188)
- **Mobile configuration releases in DAS Tag Manager: publish, read and roll
  back signed bundles per app and environment.** A Network's admin publishes
  the rules of an app (`siteCode`) for an environment (`PRODUCTION`, `STAGING`,
  `DEVELOPMENT`) with `POST /api/mobile-config/releases` (feature
  `mobile-config`, off by default; permission `das:tag-manager:mobile-config`):
  the rules are compiled and signed with the Organization's Ed25519 key at
  once, numbered 1, 2, ... per app and environment under an advisory lock, and
  kept as an immutable release - the source, the signed bundle with its
  SHA-256 and size, and the key's fingerprint. Every role reads the releases
  (newest first, without sources; one with its source), and a rollback
  republishes an older release's source as the next release, signed anew.
  `GET /api/mobile-config/signing-key` serves the public key an app embeds.
  The raw source passes a depth (100 levels), `__proto__` and storability
  check before anything recursive reads it; a source that does not compile is
  400 and writes nothing. The local stack generates its signing key; an
  Organization answers 503 `NOT_CONFIGURED` for publishing until its key is
  provisioned (#260). The `das`
  CLI (`das tag-manager mobile-config ...`) and the Terraform provider (actions
  `das_tag_manager_mobile_config_release` and `..._release_rollback`, data
  sources `..._releases`, `..._release`, `..._signing_key`) follow the
  contract. (#255)
- **Orders walked through their ad product's steps in DAS Delivery, and the
  queue "Waiting for me"**: the feature `order-workspace` (off by default,
  requires `ad-product-catalog`) adds `/api/orders`. An order is created on an
  `ENABLED` ad product, pins its current version and stands on its first step;
  a `WORK` step is marked done by a person holding the step's role, an
  `APPROVAL` step is approved or rejected (with a comment, back to the nearest
  earlier `WORK` step) by an approver who did not do that work (403
  `SELF_APPROVAL`); `awaitingMe=true` lists the orders whose current step
  belongs to one of the caller's roles. Orders change or are cancelled only
  while in progress, are never deleted, keep a history of every step, and
  every write is audited.
  The platform adds the access class `approval` (a mutation needs the scope
  action `approve`, never held by an API key), the permission
  `das:delivery:orders` (read, write, approve) and the problem codes
  `STEP_NOT_CURRENT` and `SELF_APPROVAL` (every API's Problem enum lists them).
  The Delivery web app builds the Orders module; the kit gains the order
  statuses, step states and outcomes and their glossary terms. The `das` CLI
  gains `das delivery orders ...` and the Terraform provider the resource
  `das_delivery_order`, the actions `das_delivery_order_step_completion` and
  `das_delivery_order_approval` and the data source
  `das_delivery_order_step_completions`. (#202)
- **Sales channels in DAS Delivery: ad products offered to self-service tools
  and partners, and orders that come through them.** A Network's admin creates
  a sales channel (`/api/sales-channels`, feature `sales-channels`, off by
  default) - a self-service tool such as OnetAds, or a partner - with the ad
  products it offers and the API keys that sell through it. A key of a channel
  is narrowed on the server from its credential, never from a parameter: the
  catalog shows it only the `ENABLED` ad products of its `ENABLED` channel's
  offer, `POST /api/orders` takes only those and records the channel on the
  order (`Order.salesChannelId`, null for the Network's own staff), and every
  order operation shows it only its channel's orders (404 otherwise). With
  the feature off or the channel `DISABLED` a key of the channel sells nothing
  and stays narrowed to its channel's orders (fail closed), and it never writes
  the catalog or the channels.

  - `delivery`: the `SalesChannel` model (migration `0006_sales_channels`,
    RLS), the four operations, the narrowing in the catalog and order routes,
    the `salesChannelId` field and list filter of orders, the Sales channel
    column and fact in the Orders UI.
  - `platform`: the feature `sales-channels`, the permission
    `das:delivery:sales-channels`; the secret-name gates accept `apiKeyIds`
    (ids of API keys hold no secret material).
  - `ui`: the glossary term "Sales channel" and the regenerated contract types.
  - `das-cli`: `das delivery sales-channels create|list|get|update`.
  - `terraform-provider-das`: the resource `das_delivery_sales_channel` and
    its data sources; destroy switches the channel off and keeps its links.
- **Strict audit: what every mutation changed, a tamper-evident hash chain,
  and a trail auditors read through the API, the CLI, Terraform and the
  console**. Every audit row in the four databases now records `changes`
  (every changed field before and after, captured automatically by
  `auditedTx` in the mutation's transaction; a secret, its hash or its
  reference only as `{ redacted, changed }`), the channel (`CONSOLE`, `CLI`,
  `API`, `AGENT`, `SCIM`, `SYSTEM`), the principal's role, Network and
  credential id. Rows form a SHA-256 chain per Network and for the org layer
  (`seq`, `prevHash`, `hash` over the RFC 8785 record), opened by a genesis
  that seals the older rows; a NOT VALID CHECK refuses unchained rows. The
  leader's duty `platform.audit-verify` checks the chains every 15 minutes,
  records `AuditChainCheck` / `OrgAuditChainCheck` and logs each head; the
  health checks `platform.audit-chain-broken` (critical) and
  `platform.audit-chain-unverified` (warning) read the verdicts. Each API
  serves the trail (`/api/audit-events` with filters, one row, a JSON Lines
  export that verifies offline, `/api/audit-chain` and a read-only segment
  verification; the Organisation API also on the Network layer and the
  org-layer twin `/api/networks/{networkId}/...`), readable only with
  `das:platform:audit:read`: the new roles `ORG_AUDITOR` and `AUDITOR` and
  the owner, never an admin. The Organization console gets an Audit module,
  an Audit tab per Network and a Network auditor's view, with the kit's new
  `ChangeDiff`; the CLI and the Terraform provider derive the read commands
  and data sources from the contracts. Rollback note: the CHECK constraints
  stay after a code revert, so returning to older code needs a forward fix
  that drops them. (#134)

### Fixed

- **The ad product catalog refuses a `WORK` step whose role may not write.** A
  new version whose `WORK` step names a Network role without
  `das:delivery:*:write` (today `VIEWER`, `APPROVER`, `AUDITOR`) is 400
  `VALIDATION_FAILED` at `definition.steps.<i>.role`, naming the roles that may
  (today `NETWORK_ADMIN`, `OPERATOR`); before, an order of such a version stuck
  on that step unless the person holding the role also held a writing one. The roles come from `ROLE_DEFS`, like the rule of
  `APPROVAL` steps. Versions already stored stay as they are. The contract,
  the CLI reference and the provider's embedded contract describe the rule.
  (#205)
- **DAS Tag Manager has a web app and runs in the local stack.** A new web app
  (`apps/tag-manager/web`, React on `@das/ui`) serves the product's shell on
  `<network>.tag-manager.<org domain>`: sign-in, the product and Network
  switchers, the Network's accent, and every Tag Manager module (Containers,
  Tags, Triggers, Variables, Versions, Preview) as an empty state until it is
  built. GET `/api/session` now answers what the shell needs, as DAS Delivery's
  does: the roles a Network token carries (`roles`), the Network's
  `networkCode`, the Networks the session may switch to (`networks`) and the
  Network's Tag Manager accent (`accent`, the setting
  `tag-manager.accent-color`). `npm run dev` migrates the `tag_manager`
  database and starts the Tag Manager API on base+7 and its web app on base+9
  behind the router.
- **The web apps and the CLI sign in on the user pool's real Hosted UI domain and app client.**
  A new open operation, `GET /api/sign-in`, answers the web apps how to sign a person in: the
  Hosted UI's authorize, token and logout endpoints and the web app client's id, as the duty
  `platform.identity-sync` last recorded them (null without a user pool or before the duty ran).
  The Hosted UI's domain is the one the duty records with the pool (`hostedUiDomain`: a custom
  domain or the Cognito prefix domain `<prefix>.auth.<region>.amazoncognito.com`), so an
  Organization whose pool has no `auth.<org domain>` custom domain signs in; `GET /api/cli` and the
  console's sign-in page address use it too, falling back to `auth.<org domain>` until the duty
  described the pool. The kit's session module (`@das/ui`) reads `GET /api/sign-in` through the
  typed client instead of assuming `auth.<org domain>` and the client id `das-web` (Cognito's
  client ids are random); it keeps them only while the Organisation API answers no sign-in, as the
  local stack's dev identity provider serves them. The kit's typed client sends no credential when
  it is given none, for the open operations read before sign-in. (#69)

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
