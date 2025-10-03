package ztw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/method"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/variable"
)

func TestAccDataSourceDNSForwardingGateway_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DNSForwardingGateway)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSForwardingGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDNSForwardingGatewayConfigure(resourceTypeAndName, generatedName, variable.ForwardECDNSGatewayPrimary, variable.ForwardECDNSGatewaySecondary, variable.ForwardDNSFailureBehavior),
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
