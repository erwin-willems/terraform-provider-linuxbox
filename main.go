package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/erwin-willems/terraform-provider-linuxbox/linuxbox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return linuxbox.Provider()
		},
	})
}


func main2() {
	Yumtest()
}
