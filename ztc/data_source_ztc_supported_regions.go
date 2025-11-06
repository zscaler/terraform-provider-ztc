package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations"
)

func dataSourceSupportedRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSupportedRegionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique ID of the supported region. When specified, returns a single region.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the supported region. When specified, returns a single region.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud type when retrieving a single region. The default value is AWS. Supported values are AWS, AZURE, GCP",
			},
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all supported regions. Populated when neither id nor name is specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique ID of the supported region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the supported region.",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud type. Supported values are AWS, AZURE, GCP",
						},
					},
				},
			},
		},
	}
}

func dataSourceSupportedRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *common.SupportedRegions
	id, idOk := getIntFromResourceData(d, "id")
	name, _ := d.Get("name").(string)

	// If ID is provided, search by ID
	if idOk {
		log.Printf("[INFO] Getting supported region by id: %d\n", id)
		regions, err := partner_integrations.GetSupportedRegions(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		// Find the region with matching ID
		for _, region := range regions {
			if region.ID == id {
				resp = &region
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("no supported region found with id: %d", id))
		}
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("cloud_type", resp.CloudType)
		log.Printf("[INFO] Retrieved supported region id: %d, name: %s\n", resp.ID, resp.Name)
	} else if name != "" {
		// If name is provided, search by name
		log.Printf("[INFO] Getting supported region by name: %s\n", name)
		res, err := partner_integrations.GetSupportedRegionsByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("cloud_type", resp.CloudType)
		log.Printf("[INFO] Retrieved supported region id: %d, name: %s\n", resp.ID, resp.Name)
	} else {
		// If neither ID nor name is provided, return all regions
		log.Printf("[INFO] Getting all supported regions\n")
		regions, err := partner_integrations.GetSupportedRegions(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		// Flatten all regions into the regions list
		regionsList := make([]map[string]interface{}, 0, len(regions))
		for _, region := range regions {
			regionsList = append(regionsList, map[string]interface{}{
				"id":         region.ID,
				"name":       region.Name,
				"cloud_type": region.CloudType,
			})
		}

		if err := d.Set("regions", regionsList); err != nil {
			return diag.FromErr(err)
		}

		// Set ID to "0" to indicate all regions (d.SetId always takes a string)
		d.SetId("0")
		log.Printf("[INFO] Retrieved %d supported regions\n", len(regions))
	}

	return nil
}
