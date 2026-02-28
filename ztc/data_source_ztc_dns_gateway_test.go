package ztc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
)

func TestAccDataSourceDNSGateway_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DNSGateway)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDNSGatewayConfigure(resourceTypeAndName, generatedName, variable.DNSGatewayECOptionsPrimary, variable.DNSGatewayECOptionsSecondary, variable.DNSGatewayFailureBehavior),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ec_dns_gateway_options_primary", resourceTypeAndName, "ec_dns_gateway_options_primary"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ec_dns_gateway_options_secondary", resourceTypeAndName, "ec_dns_gateway_options_secondary"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "failure_behavior", resourceTypeAndName, "failure_behavior"),
				),
			},
		},
	})
}
