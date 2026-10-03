# Terraform Provider for DAS

`ringdas/das` configures a DAS Organization as code: everything the `das`
CLI can do. Its resources, data sources, actions and ephemeral resources are
derived from the API contracts of the DAS release it ships with, so a new
DAS feature reaches Terraform with the next provider version.

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

resource "das_group" "traffickers" {
  name = "traffickers"
}

resource "das_network_role_binding" "traffickers" {
  network_id = "net_..."
  group_name = das_group.traffickers.name
  role       = "OPERATOR"
}
```

Documentation of every type: [docs/](docs/index.md). Requires Terraform 1.11
or later (write-only attributes); actions need 1.14.

## Development

The provider is built from `internal/specs/*.openapi.json`, copies of the
monorepo's API contracts. After an API change, `npm run gen:contracts` (from
the repository root) refreshes them and `docs/`. Nothing is written per
resource: see `docs/architecture/terraform-provider.md` in the monorepo.

```sh
go test ./...                 # the gates (needs terraform on PATH for the lifecycle test)
go run ./cmd/docgen           # docs/
go build -o bin/ .            # then a dev_overrides block for "ringdas/das" in ~/.terraformrc
```
