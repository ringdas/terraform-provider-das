// The DAS Terraform provider, `ringdas/das`
// (docs/architecture/terraform-provider.md). Everything the `das` CLI can
// do, this provider can do: both derive what they offer from the API
// contracts of the DAS release they ship with.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/ringdas/terraform-provider-das/internal/provider"
)

// version is set at release build time (-ldflags "-X main.version=...");
// it stays "dev" for local builds.
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: provider.Address,
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
