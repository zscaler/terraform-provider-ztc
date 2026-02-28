package ztc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
	dnsgateway "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/dns_gateway"
)

func TestAccResourceDNSGateway_Basic(t *testing.T) {
	var gateway dnsgateway.DNSGateway
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DNSGateway)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-acc-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDNSGatewayConfigure(resourceTypeAndName, initialName, variable.DNSGatewayECOptionsPrimary, variable.DNSGatewayECOptionsSecondary, variable.DNSGatewayFailureBehavior),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_primary", variable.DNSGatewayECOptionsPrimary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_secondary", variable.DNSGatewayECOptionsSecondary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "failure_behavior", variable.DNSGatewayFailureBehavior),
				),
			},

			// Update test
			{
				Config: testAccCheckDNSGatewayConfigure(resourceTypeAndName, updatedName, variable.DNSGatewayECOptionsPrimary, variable.DNSGatewayECOptionsSecondary, variable.DNSGatewayFailureBehavior),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_primary", variable.DNSGatewayECOptionsPrimary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ec_dns_gateway_options_secondary", variable.DNSGatewayECOptionsSecondary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "failure_behavior", variable.DNSGatewayFailureBehavior),
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

func testAccCheckDNSGatewayDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DNSGateway {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := dnsgateway.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dns gateway with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDNSGatewayExists(resource string, rule *dnsgateway.DNSGateway) resource.TestCheckFunc {
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

		receivedRule, err := dnsgateway.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckDNSGatewayConfigure(resourceTypeAndName, generatedName, primary, secondary, failBehavior string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]

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
		resourcetype.DNSGateway,
		resourceName,
		generatedName,
		primary,
		secondary,
		failBehavior,

		resourcetype.DNSGateway,
		resourceName,

		resourcetype.DNSGateway,
		resourceName,
	)
}
