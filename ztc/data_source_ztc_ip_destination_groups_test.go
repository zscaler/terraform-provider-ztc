package ztc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
)

func TestAccDataSourceIPDestinationGroups_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.IPDestinationGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, variable.IPDSTGroupDescription, variable.IPDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "addresses.#", "3"),
				),
			},
		},
	})
}
