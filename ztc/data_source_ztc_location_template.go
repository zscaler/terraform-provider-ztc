package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
)

func dataSourceLocationTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of Cloud & Branch Connector location template",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of Cloud & Branch Connector location template",
			},
			"desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of Cloud & Branch Connector location template",
			},
			"editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Cloud & Branch Connector location template is editable",
			},
			"last_mod_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last time Cloud & Branch Connector location template was modified",
			},
			"template": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prefix of Cloud & Branch Connector location template",
						},
						"xff_forward_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: `"Enable XFF Forwarding for a location.
								When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
								Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location."`,
						},
						"auth_required": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: `"Indicates whether to enforce authentication.
								Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled."`,
						},
						"caution_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable Caution for this location",
						},
						"aup_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable Acceptable Use Policy (AUP) for this location",
						},
						"aup_timeout_in_days": {
							Type:     schema.TypeInt,
							Computed: true,
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
					},
				},
			},
			"last_mod_uid": UIDNameSchema(),
		},
	}
}

func dataSourceLocationTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *locationtemplate.LocationTemplate
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for location template id: %d\n", id)
		res, err := locationtemplate.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for location template name: %s\n", name)
		res, err := locationtemplate.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("desc", resp.Description)
		_ = d.Set("editable", resp.Editable)
		_ = d.Set("last_mod_time", resp.LastModTime)

		if err := d.Set("template", flattenTemplateDetails(resp.LocationTemplateDetails)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("last_mod_uid", flattenGeneralPurpose(resp.LastModUid)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.Errorf("couldn't find any location template with name '%s'", name)
	}

	return nil
}

func flattenTemplateDetails(template *locationtemplate.LocationTemplateDetails) []map[string]interface{} {
	if template == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"template_prefix":           template.TemplatePrefix,
			"xff_forward_enabled":       template.XFFForwardEnabled,
			"auth_required":             template.AuthRequired,
			"caution_enabled":           template.CautionEnabled,
			"aup_enabled":               template.AupEnabled,
			"aup_timeout_in_days":       template.AupTimeoutInDays,
			"ofw_enabled":               template.OFWEnabled,
			"ips_control":               template.IPSControl,
			"enforce_bandwidth_control": template.EnforceBandwidthControl,
			"up_bandwidth":              template.UpBandwidth,
			"dn_bandwidth":              template.DnBandwidth,
		},
	}
}
