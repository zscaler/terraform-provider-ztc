package ztc

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/method"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/testing/variable"
)

func TestAccDataSourceLocationTemplate_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZCBC_Location_Template)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationTemplateTest1Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "desc", resourceTypeAndName, "desc"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "template.#", "1"),
				),
			},
			{
				Config: testAccCheckLocationTemplateTest2Configure(resourceTypeAndName, generatedName, variable.LocTemplateResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "desc", resourceTypeAndName, "desc"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "template.#", "1"),
				),
			},
		},
	})
}
*/
