package main

import (
       "github.com/hashicorp/terraform/plugin"
	   "github.com/GSLabDev/terraform-provider-ad/ad"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ad.Provider})
}
