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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipsourcegroups"
)

func TestAccResourceIPPoolGroupsBasic(t *testing.T) {
	var groups ipgroups.IPGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.IPPoolGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPPoolGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPPoolGroupsConfigure(resourceTypeAndName, initialName, variable.IPPoolGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPPoolGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.IPPoolGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckIPPoolGroupsConfigure(resourceTypeAndName, updatedName, variable.IPPoolGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPPoolGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.IPPoolGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "1"),
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

func testAccCheckIPPoolGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.IPPoolGroup {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := ipsourcegroups.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("ip source group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckIPPoolGroupsExists(resource string, rule *ipgroups.IPGroups) resource.TestCheckFunc {
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

		receivedRule, err := ipgroups.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckIPPoolGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
    ip_addresses = ["192.168.1.0/24"]
  }

data "%s" "%s" {
id = "${%s.%s.id}"
}
`,
		// Resource type and name for the ip group
		resourcetype.IPPoolGroup,
		resourceName,
		generatedName,
		description,

		// Data source type and name
		resourcetype.IPPoolGroup,
		resourceName,

		// Reference to the resource
		resourcetype.IPPoolGroup,
		resourceName,
	)
}
