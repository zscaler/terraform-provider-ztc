package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/location"
)

func dataSourceLocationManagement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"non_editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enforce_bandwidth_control": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable to specify the maximum bandwidth limits for download (Mbps) and upload (Mbps).",
			},
			"up_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Upload bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.",
			},
			"dn_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Download bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country of the location",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the location",
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Language of the location",
			},
			"tz": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timezone of the location. If not specified, it defaults to GMT",
			},
			"auth_required": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `"Indicates whether to enforce authentication.
				Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled."`,
			},
			"xff_forward_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `"Enable XFF Forwarding for a location.
				When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
				Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location."`,
			},
			"ec_location": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is a Cloud or Branch Connector location (true) or a generic location (false).",
			},
			"ofw_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable firewall for this location",
			},
			"ips_control": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable IPS for this location",
			},
			"aup_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable Acceptable Use Policy (AUP) for this location",
			},
			"caution_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable Caution for this location",
			},
			"exclude_from_dynamic_groups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to exclude this location from dynamic location groups when created",
			},
			"exclude_from_manual_groups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to exclude this location from manual location groups when created",
			},
			"public_cloud_account_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Prefix of Cloud & Branch Connector location template",
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `"Enable XFF Forwarding for a location.
								When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
								Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location."`,
						},
					},
				},
			},
		},
	}
}

func dataSourceLocationManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *location.Locations
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for location id: %d\n", id)
		res, err := location.GetLocation(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for location name: %s\n", name)
		res, err := location.GetLocationByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("parent_id", resp.ParentID)
		_ = d.Set("enforce_bandwidth_control", resp.EnforceBandwidthControl)
		_ = d.Set("up_bandwidth", resp.UpBandwidth)
		_ = d.Set("dn_bandwidth", resp.DnBandwidth)
		_ = d.Set("country", resp.Country)
		_ = d.Set("state", resp.State)
		_ = d.Set("language", resp.Language)
		_ = d.Set("tz", resp.TZ)
		_ = d.Set("auth_required", resp.AuthRequired)
		_ = d.Set("xff_forward_enabled", resp.XFFForwardEnabled)
		_ = d.Set("ec_location", resp.ECLocation)
		_ = d.Set("ofw_enabled", resp.OFWEnabled)
		_ = d.Set("ips_control", resp.IPSControl)
		_ = d.Set("aup_enabled", resp.AUPEnabled)
		_ = d.Set("caution_enabled", resp.CautionEnabled)
		_ = d.Set("exclude_from_dynamic_groups", resp.ExcludeFromDynamicGroups)
		_ = d.Set("exclude_from_manual_groups", resp.ExcludeFromManualGroups)

		if err := d.Set("public_cloud_account_id", flattenIDNameExternalID(resp.PublicCloudAccountID)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("couldn't find any location with name '%s'", name)
	}

	return nil
}

func flattenIDNameExternalID(item *common.CommonIDName) []interface{} {
	r := map[string]interface{}{
		"id":   item.ID,
		"name": item.Name,
	}
	return []interface{}{r}
}
