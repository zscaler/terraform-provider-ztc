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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservicegroups"
)

func TestAccResourceNetworkServiceGroups_Basic(t *testing.T) {
	var services networkservicegroups.NetworkServiceGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.NetworkServiceGroups)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServiceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, initialName, variable.NetworkServicesGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServiceGroupsExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.NetworkServicesGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "services.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, updatedName, variable.NetworkServicesGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServiceGroupsExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.NetworkServicesGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "services.#", "1"),
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

func testAccCheckFWNetworkServiceGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.NetworkServiceGroups {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := networkservicegroups.GetNetworkServiceGroups(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network services group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFWNetworkServiceGroupsExists(resource string, rule *networkservicegroups.NetworkServiceGroups) resource.TestCheckFunc {
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

		receivedRule, err := networkservicegroups.GetNetworkServiceGroups(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "ztw_network_services" "example1" {
	name = "ICMP_ANY"
  }

data "ztw_network_services" "example2" {
	name = "TCP_ANY"
  }

resource "%s" "%s" {
    name = "%s"
    description = "%s"
    services {
        id = [
            data.ztw_network_services.example1.id,
            data.ztw_network_services.example2.id,
        ]
    }
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the network services group
		resourcetype.NetworkServiceGroups,
		resourceName,
		generatedName,
		description,

		// Data source type and name
		resourcetype.NetworkServiceGroups,
		resourceName,

		// Reference to the resource
		resourcetype.NetworkServiceGroups,
		resourceName,
	)
}
