// Package main is the entry point for the plugin. It sets up the provider and starts the plugin.
package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/josh-silvas/terraform-provider-awx/internal/awx"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// func main() {
// 	plugin.Serve(&plugin.ServeOpts{
// 		ProviderAddr: "registry.terraform.io/josh-silvas/awx",
// 		ProviderFunc: func() *schema.Provider {
// 			return awx.Provider()
// 		},
// 	})
// }

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "registry.terraform.io/josh-silvas/awx",
		ProviderFunc: func() *schema.Provider {
			return awx.Provider()
		},
	}

	plugin.Serve(opts)
}
