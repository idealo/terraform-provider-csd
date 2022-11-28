package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/idealo/terraform-provider-idealo-tools/idealo_tools"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug: debugMode,
		// TODO: update this string with the full name of your provider as used in your configs
		ProviderAddr: "idealo.com/transport/idealo-tools",
		ProviderFunc: func() *schema.Provider {
			return idealo_tools.Provider()
		},
	})
}
