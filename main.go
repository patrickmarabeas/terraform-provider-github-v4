package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/patrickmarabeas/terraform-provider-githubv4/github"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: github.Provider,
	})
}
