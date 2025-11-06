package ztc

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
)

func TestAccDataSourceTrafficForwardingRule_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingRule)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.IPSourceGroup)
	sourceIPGroupHCL := testAccCheckIPSourceGroupsConfigure(sourceIPGroupTypeAndName, "tf-acc-test-"+sourceIPGroupGeneratedName, variable.IPSRCGroupDescription)

	// Generate Destination IP Group HCL Resource
	dstIPGroupTypeAndName, _, dstIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.IPDestinationGroup)
	dstIPGroupHCL := testAccCheckIPDestinationGroupsConfigure(dstIPGroupTypeAndName, "tf-acc-test-"+dstIPGroupGeneratedName, variable.IPDSTGroupDescription, variable.IPDSTGroupTypeDSTNFQDN)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.ForwardControlRuleDescription, variable.ForwardControlRuleType, variable.ForwardControlMethod, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "forward_method", resourceTypeAndName, "forward_method"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "src_ip_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "dest_ip_groups.#", "1"),
				),
			},
		},
	})
}
*/
