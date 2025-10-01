package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-ztw/ztw"
	"github.com/zscaler/terraform-provider-ztw/ztw/common"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(common.Version())
		return
	}
	var debug bool
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		debug = true
	}
	log.Printf(`ZTW Terraform Provider

Version %s

https://registry.terraform.io/providers/zscaler/ztw/latest/docs

`, common.Version())
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ztw.ZTWProvider,
		ProviderAddr: "registry.terraform.io/zscaler/ztw",
		Debug:        debug,
	})
}
