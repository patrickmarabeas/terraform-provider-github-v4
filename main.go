package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/patrickmarabeas/terraform-provider-github-v4/github-v4"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: github.Provider,
	})
}
