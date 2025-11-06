package ztc

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
)

func TestAccDataSourceForwardingGateway_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZIAForwardingGateway)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingGatewayConfigure(resourceTypeAndName, generatedName, variable.ForwardGWDescription, variable.ForwardGWPrimaryType, variable.ForwardGWSecondaryType, variable.ForwardGWType, variable.ForwardGWFailClose),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "primary_type", resourceTypeAndName, "primary_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "secondary_type", resourceTypeAndName, "secondary_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fail_closed", strconv.FormatBool(variable.ForwardGWFailClose)),
				),
			},
		},
	})
}
