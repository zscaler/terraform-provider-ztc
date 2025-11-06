package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zscaler/terraform-provider-ztc/ztc"
	"github.com/zscaler/terraform-provider-ztc/ztc/common"
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
	log.Printf(`ZTC Terraform Provider

Version %s

https://registry.terraform.io/providers/zscaler/ztc/latest/docs

`, common.Version())
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ztc.ZTCProvider,
		ProviderAddr: "registry.terraform.io/zscaler/ztc",
		Debug:        debug,
	})
}
