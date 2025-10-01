package ztw

/*
import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/method"
	"github.com/zscaler/terraform-provider-ztw/ztw/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
)

func TestAccResourceLocationTemplateBasic(t *testing.T) {
	var template locationtemplate.LocationTemplate
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZCBC_Location_Template)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationTemplateTest1Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationTemplateExists(resourceTypeAndName, &template),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "desc", variable.LocTemplateResourceDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "template.#", "1"),
				),
			},

			// Update test1
			{
				Config: testAccCheckLocationTemplateTest1Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationTemplateExists(resourceTypeAndName, &template),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "desc", variable.LocTemplateResourceDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "template.#", "1"),
				),
			},
			{
				Config: testAccCheckLocationTemplateTest2Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationTemplateExists(resourceTypeAndName, &template),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "desc", variable.LocTemplateResourceDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "template.#", "1"),
				),
			},
			// Update test2
			{
				Config: testAccCheckLocationTemplateTest2Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationTemplateExists(resourceTypeAndName, &template),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "desc", variable.LocTemplateResourceDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "template.#", "1"),
				),
			},
		},
	})
}

func testAccCheckLocationTemplateDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.locationtemplate

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZCBC_Location_Template {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		template, err := locationtemplate.Get(ctx, service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if template != nil {
			return fmt.Errorf("location templates with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckLocationTemplateExists(resource string, rule *locationtemplate.LocationTemplate) resource.TestCheckFunc {
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
		service := apiClient.locationtemplate

		receivedRule, err := locationtemplate.Get(ctx, service, id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckLocationTemplateTest1Configure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
// location template resource
%s

data "%s" "%s" {
	id = "${%s.id}"
  }
`,
		// resource variables
		LocationTemplateTest1ResourceHCL(generatedName, description),

		// data source variables
		resourcetype.ZCBC_Location_Template,
		generatedName,
		resourceTypeAndName,
	)
}

func LocationTemplateTest1ResourceHCL(generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "tf-acc-test-%s"
    desc = "%s"
    template {
      template_prefix = "testAcc-tf"
      auth_required = true
      ips_control = true
      ofw_enabled = true
      xff_forward_enabled = true
      enforce_bandwidth_control = true
      up_bandwidth = 10
      dn_bandwidth = 10
    }
}
`,
		// resource variables
		resourcetype.ZCBC_Location_Template,
		generatedName,
		generatedName,
		description,
	)
}

func testAccCheckLocationTemplateTest2Configure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`

// location template resource
%s

data "%s" "%s" {
	id = "${%s.id}"
  }
`,
		// resource variables
		LocationTemplateTest2ResourceHCL(generatedName, description),

		// data source variables
		resourcetype.ZCBC_Location_Template,
		generatedName,
		resourceTypeAndName,
	)
}

func LocationTemplateTest2ResourceHCL(generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "tf-acc-test-%s"
    desc = "%s"
    template {
      template_prefix = "testAcc-tf"
      aup_enabled = true
      aup_timeout_in_days = 10
      ips_control = true
      ofw_enabled = true
      xff_forward_enabled = true
      enforce_bandwidth_control = true
      up_bandwidth = 10
      dn_bandwidth = 10
    }
}
`,
		// resource variables
		resourcetype.ZCBC_Location_Template,
		generatedName,
		generatedName,
		description,
	)
}
*/
