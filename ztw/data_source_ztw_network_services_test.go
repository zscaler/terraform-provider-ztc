package ztw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/method"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/variable"
)

func TestAccDataSourceNetworkServices_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.NetworkServices)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServicesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServicesConfigure(resourceTypeAndName, generatedName, variable.NetworkServicesDescription, variable.NetworkServicesType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "src_tcp_ports.#", "3"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "dest_tcp_ports.#", "3"),
				),
			},
		},
	})
}
