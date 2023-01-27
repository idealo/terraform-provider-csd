package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/idealo/terraform-provider-idealo-tools/csd"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool

	// TODO: figure out how to use this
	flag.BoolVar(&debugMode, "debug", false, "Set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        debugMode,
		ProviderAddr: "idealo.com/transport/csd",
		ProviderFunc: func() *schema.Provider {
			return csd.Provider()
		},
	})
}
