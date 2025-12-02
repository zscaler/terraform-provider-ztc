package ztc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
)

func resourceLocationTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationTemplateCreate,
		ReadContext:   resourceLocationTemplateRead,
		UpdateContext: resourceLocationTemplateUpdate,
		DeleteContext: resourceLocationTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("template_id", int(idInt))
				} else {
					resp, err := locationtemplate.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("template_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location Name.",
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"xff_forward_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.",
						},
						"auth_required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.",
						},
						"aup_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable AUP. When set to true, AUP is enabled for the location.",
						},
						"caution_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable Caution. When set to true, a caution notifcation is enabled for the location.",
						},
						"aup_timeout_in_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.",
						},
						"ofw_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable Firewall. When set to true, Firewall is enabled for the location.",
						},
						"ips_control": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.",
						},
						"enforce_bandwidth_control": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"up_bandwidth": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 99999999),
							Description:  "Upload bandwidth in bytes. The value 0 implies no Bandwidth Control enforcement.",
						},
						"dn_bandwidth": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 99999999),
							Description:  "Upload bandwidth in bytes. The value 0 implies no Bandwidth Control enforcement.",
						},
					},
				},
			},
		},
	}
}

func resourceLocationTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandLocationTemplate(d)
	log.Printf("[INFO] Creating cloud connector location template\n%+v\n", req)
	if err := checkLocationTemplateDependencies(req); err != nil {
		return diag.FromErr(err)
	}

	resp, err := locationtemplate.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created cloud connector location template request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("template_id", resp.ID)

	return resourceLocationTemplateRead(ctx, d, meta)
}

func checkLocationTemplateDependencies(template locationtemplate.LocationTemplate) error {
	if template.LocationTemplateDetails.AuthRequired && template.LocationTemplateDetails.CautionEnabled {
		return fmt.Errorf("authentication required must be disabled, when enabling caution")
	}
	if template.LocationTemplateDetails.AupEnabled && template.LocationTemplateDetails.CautionEnabled {
		return fmt.Errorf("enabling AUP and Caution together is not allowed")
	}
	if template.LocationTemplateDetails.AupEnabled && template.LocationTemplateDetails.AupTimeoutInDays == 0 {
		return fmt.Errorf("AUP timeout in days is required, when AUP is enabled")
	}
	if template.LocationTemplateDetails.EnforceBandwidthControl && template.LocationTemplateDetails.UpBandwidth == 0 && template.LocationTemplateDetails.DnBandwidth == 0 {
		return fmt.Errorf("upload and download bandwidth is mandatory when enforce bandwidth setting is on")
	}
	return nil
}

func resourceLocationTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		return diag.Errorf("no location template id is set")
	}
	resp, err := locationtemplate.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.Response.StatusCode == 404 {
			log.Printf("[WARN] Removing location template %s from state because it no longer exists in Cloud Connector", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting location template:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("desc", resp.Description)

	if err := d.Set("template", flattenTemplateDetails(resp.LocationTemplateDetails)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLocationTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] location template ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating location template ID: %v\n", id)
	req := expandLocationTemplate(d)
	if err := checkLocationTemplateDependencies(req); err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := locationtemplate.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceLocationTemplateRead(ctx, d, meta)
}

func resourceLocationTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] location template ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting location template ID: %v\n", (d.Id()))

	if _, err := locationtemplate.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] location template deleted")
	return nil
}

func expandLocationTemplate(d *schema.ResourceData) locationtemplate.LocationTemplate {
	id, _ := getIntFromResourceData(d, "template_id")
	result := locationtemplate.LocationTemplate{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("desc").(string),

		LocationTemplateDetails: expandLocationTemplateDetails(d),
	}
	templateDetails := expandLocationTemplateDetails(d)
	if templateDetails != nil {
		result.LocationTemplateDetails = templateDetails
	}
	return result
}

func expandLocationTemplateDetails(d *schema.ResourceData) *locationtemplate.LocationTemplateDetails {
	templateObj, ok := d.GetOk("template")
	if !ok {
		return nil
	}
	templates, ok := templateObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(templates.List()) > 0 {
		templateObj := templates.List()[0]
		template, ok := templateObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &locationtemplate.LocationTemplateDetails{
			TemplatePrefix:          template["template_prefix"].(string),
			XFFForwardEnabled:       template["xff_forward_enabled"].(bool),
			AuthRequired:            template["auth_required"].(bool),
			CautionEnabled:          template["caution_enabled"].(bool),
			AupEnabled:              template["aup_enabled"].(bool),
			AupTimeoutInDays:        template["aup_timeout_in_days"].(int),
			OFWEnabled:              template["ofw_enabled"].(bool),
			IPSControl:              template["ips_control"].(bool),
			EnforceBandwidthControl: template["enforce_bandwidth_control"].(bool),
			UpBandwidth:             template["up_bandwidth"].(int),
			DnBandwidth:             template["dn_bandwidth"].(int),
		}
	}
	return nil
}
