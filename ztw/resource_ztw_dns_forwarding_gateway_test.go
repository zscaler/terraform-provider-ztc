package ztw

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/method"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

func TestAccResourceDNSForwardingGateway_Basic(t *testing.T) {
	var gateway dns_forwarding_gateway.DNSGateway
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DNSForwardingGateway)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSForwardingGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDNSForwardingGatewayConfigure(resourceTypeAndName, initialName, variable.ForwardECDNSGatewayPrimary, variable.ForwardECDNSGatewaySecondary, variable.ForwardDNSFailureBehavior),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSForwardingGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_primary", variable.ForwardECDNSGatewayPrimary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_secondary", variable.ForwardECDNSGatewaySecondary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "failure_behavior", variable.ForwardDNSFailureBehavior),
				),
			},

			// Update test
			{
				Config: testAccCheckDNSForwardingGatewayConfigure(resourceTypeAndName, updatedName, variable.ForwardECDNSGatewayPrimary, variable.ForwardECDNSGatewaySecondary, variable.ForwardDNSFailureBehavior),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSForwardingGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_primary", variable.ForwardECDNSGatewayPrimary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_secondary", variable.ForwardECDNSGatewaySecondary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "failure_behavior", variable.ForwardDNSFailureBehavior),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDNSForwardingGatewayDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZIAForwardingGateway {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, _, err := zia_forwarding_gateway.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("forwarding gateway with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDNSForwardingGatewayExists(resource string, rule *dns_forwarding_gateway.DNSGateway) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		receivedRule, _, err := dns_forwarding_gateway.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckDNSForwardingGatewayConfigure(resourceTypeAndName, generatedName, primary, secondary, failBehavior string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
  name                             = "%s"
  ec_dns_gateway_options_primary   = "%s"
  ec_dns_gateway_options_secondary = "%s"
  failure_behavior                 = "%s"
}

  data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the dns forwarding gateway
		resourcetype.DNSForwardingGateway,
		resourceName,
		generatedName,
		primary,
		secondary,
		failBehavior,

		// Data source type and name
		resourcetype.DNSForwardingGateway,
		resourceName,

		// Reference to the resource
		resourcetype.DNSForwardingGateway,
		resourceName,
	)
}
