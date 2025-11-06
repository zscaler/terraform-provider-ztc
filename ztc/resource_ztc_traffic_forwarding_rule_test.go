package ztc

/*
import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
)

func TestAccResourceTrafficForwardingRule_Basic(t *testing.T) {
	var rules forwarding_rules.ForwardingRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingRule)

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
					testAccCheckTrafficForwardingRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ForwardControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ForwardControlRuleType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "forward_method", variable.ForwardControlMethod),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.ForwardControlState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckTrafficForwardingRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.ForwardControlRuleDescription, variable.ForwardControlRuleType, variable.ForwardControlMethod, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ForwardControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ForwardControlRuleType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "forward_method", variable.ForwardControlMethod),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.ForwardControlState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
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

func testAccCheckTrafficForwardingRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficForwardingRule {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := forwarding_rules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("forwarding control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckTrafficForwardingRuleExists(resource string, rule *forwarding_rules.ForwardingRules) resource.TestCheckFunc {
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

		var receivedRule *forwarding_rules.ForwardingRules

		// Integrate retry here
		retryErr := RetryOnError(func() error {
			var innerErr error
			receivedRule, innerErr = forwarding_rules.Get(context.Background(), service, id)
			if innerErr != nil {
				return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, innerErr)
			}
			return nil
		})

		if retryErr != nil {
			return retryErr
		}

		*rule = *receivedRule
		return nil
	}
}

func testAccCheckTrafficForwardingRuleConfigure(resourceTypeAndName, generatedName, name, description, ruleType, forwardMethod string, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL string) string {
	return fmt.Sprintf(`

// source ip group resource
%s

// destination ip group resource
%s

// forwarding control rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		sourceIPGroupHCL,
		dstIPGroupHCL,
		getTrafficForwardingRuleResourceHCL(generatedName, name, description, ruleType, forwardMethod, variable.ForwardControlState, sourceIPGroupTypeAndName, dstIPGroupTypeAndName),

		// data source variables
		resourcetype.TrafficForwardingRule,
		generatedName,
		resourceTypeAndName,
	)
}

func getTrafficForwardingRuleResourceHCL(generatedName, name, description, ruleType, forwardMethod, state string, sourceIPGroupTypeAndName, dstIPGroupTypeAndName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	state = "%s"
	order = 1
	rank = 7
	type = "%s"
    forward_method = "%s"
	src_ip_groups {
		id = ["${%s.id}"]
	}
	dest_ip_groups {
		id = ["${%s.id}"]
	}
	depends_on = [ %s, %s ]
}
		`,
		// resource variables
		resourcetype.TrafficForwardingRule,
		generatedName,
		name,
		description,
		state,
		ruleType,
		forwardMethod,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
	)
}
*/
