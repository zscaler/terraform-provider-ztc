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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

func TestAccResourceForwardingGateway_Basic(t *testing.T) {
	var gateway zia_forwarding_gateway.ECGateway
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZIAForwardingGateway)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingGatewayConfigure(resourceTypeAndName, initialName, variable.ForwardGWDescription, variable.ForwardGWPrimaryType, variable.ForwardGWSecondaryType, variable.ForwardGWType, variable.ForwardGWFailClose),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ForwardGWDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "primary_type", variable.ForwardGWPrimaryType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "secondary_type", variable.ForwardGWSecondaryType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ForwardGWType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fail_closed", strconv.FormatBool(variable.ForwardGWFailClose)),
				),
			},

			// Update test
			{
				Config: testAccCheckForwardingGatewayConfigure(resourceTypeAndName, updatedName, variable.ForwardGWDescription, variable.ForwardGWPrimaryType, variable.ForwardGWSecondaryType, variable.ForwardGWType, variable.ForwardGWFailClose),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingGatewayExists(resourceTypeAndName, &gateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ForwardGWDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "primary_type", variable.ForwardGWPrimaryType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "secondary_type", variable.ForwardGWSecondaryType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ForwardGWType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fail_closed", strconv.FormatBool(variable.ForwardGWFailClose)),
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

func testAccCheckForwardingGatewayDestroy(s *terraform.State) error {
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

func testAccCheckForwardingGatewayExists(resource string, rule *zia_forwarding_gateway.ECGateway) resource.TestCheckFunc {
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

		receivedRule, _, err := zia_forwarding_gateway.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckForwardingGatewayConfigure(resourceTypeAndName, generatedName, description, primaryType, secondaryType, gatewayType string, failClosed bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
  name           = "%s"
  description    = "%s"
  fail_closed    = "%s"
  primary_type   = "%s"
  secondary_type = "%s"
  type           = "%s"
}

  data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the zia forwarding gateway
		resourcetype.ZIAForwardingGateway,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(failClosed),
		primaryType,
		secondaryType,
		gatewayType,

		// Data source type and name
		resourcetype.ZIAForwardingGateway,
		resourceName,

		// Reference to the resource
		resourcetype.ZIAForwardingGateway,
		resourceName,
	)
}
