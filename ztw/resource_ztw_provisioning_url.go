package ztw

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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/provisioning/provisioning_url"
)

func resourceProvisioningURL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProvisioningURLCreate,
		ReadContext:   resourceProvisioningURLRead,
		UpdateContext: resourceProvisioningURLUpdate,
		DeleteContext: resourceProvisioningURLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("provurl_id", id)
				} else {
					resp, err := provisioning_url.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						d.Set("provurl_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"provurl_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prov_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prov_url_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ONPREM",
					"CLOUD",
					"DISABLED",
					"MON_DELETED",
					"DELETED",
				}, false),
			},
			"prov_url_data": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_template": IdSchema(),
						"form_factor": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SMALL",
								"MEDIUM",
								"LARGE",
							}, false),
						},
						"cloud_provider_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AWS",
								"AZURE",
								"GCP",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceProvisioningURLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandProvisioningURLDetails(d)
	log.Printf("[INFO] Creating zia provisioning url\n%+v\n", req)

	resp, _, err := provisioning_url.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia provisioning url request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("provurl_id", resp.ID)

	return resourceProvisioningURLRead(ctx, d, meta)
}

func resourceProvisioningURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "provurl_id")
	if !ok {
		return diag.Errorf("no provisioning url id is set")
	}
	resp, err := provisioning_url.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia rule labels %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting provisioning url:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("desc", resp.Desc)
	_ = d.Set("prov_url", resp.ProvUrl)
	_ = d.Set("prov_url_type", resp.ProvUrlType)

	if err := d.Set("prov_url_data", flattenProvURLDataSimple(&resp.ProvUrlData)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceProvisioningURLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "provurl_id")
	if !ok {
		log.Printf("[ERROR] provisioning url ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia provisioning url ID: %v\n", id)
	req := expandProvisioningURLDetails(d)
	if _, err := provisioning_url.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := provisioning_url.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceProvisioningURLRead(ctx, d, meta)
}

func resourceProvisioningURLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "provurl_id")
	if !ok {
		log.Printf("[ERROR] provisioning url ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia provisioning url ID: %v\n", (d.Id()))

	if _, err := provisioning_url.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia provisioning url deleted")

	return nil
}

func expandProvisioningURLDetails(d *schema.ResourceData) provisioning_url.ProvisioningURL {
	id, _ := getIntFromResourceData(d, "provurl_id")
	result := provisioning_url.ProvisioningURL{
		ID:          id,
		Name:        d.Get("name").(string),
		Desc:        d.Get("desc").(string),
		ProvUrl:     d.Get("prov_url").(string),
		ProvUrlType: d.Get("prov_url_type").(string),
	}

	provUrlData := expandLocationProvURLData(d)
	if provUrlData != nil {
		result.ProvUrlData = *provUrlData
	}
	return result
}

func expandLocationProvURLData(d *schema.ResourceData) *provisioning_url.ProvUrlData {
	provUrlDataObj, ok := d.GetOk("prov_url_data")
	if !ok {
		return nil
	}
	provUrls, ok := provUrlDataObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(provUrls.List()) == 0 {
		return nil
	}

	provUrlDataMap, ok := provUrls.List()[0].(map[string]interface{})
	if !ok {
		return nil
	}

	locTemplateObj, ok := provUrlDataMap["location_template"]
	if !ok {
		return nil
	}

	locTemplateSet, ok := locTemplateObj.(*schema.Set)
	if !ok || locTemplateSet.Len() == 0 {
		return nil
	}

	locTemplateMap, ok := locTemplateSet.List()[0].(map[string]interface{})
	if !ok {
		return nil
	}

	templateID, ok := locTemplateMap["id"].(int)
	if !ok {
		return nil
	}

	result := &provisioning_url.ProvUrlData{
		LocationTemplate: locationtemplate.LocationTemplate{
			ID: templateID,
		},
	}

	if formFactor, ok := provUrlDataMap["form_factor"].(string); ok && formFactor != "" {
		result.FormFactor = formFactor
	}

	if cloudProviderType, ok := provUrlDataMap["cloud_provider_type"].(string); ok && cloudProviderType != "" {
		result.CloudProviderType = cloudProviderType
	}

	return result
}

func flattenProvURLDataSimple(provUrlData *provisioning_url.ProvUrlData) []map[string]interface{} {
	if provUrlData == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"location_template":   flattenLocationTemplateSimple(&provUrlData.LocationTemplate),
			"form_factor":         provUrlData.FormFactor,
			"cloud_provider_type": provUrlData.CloudProviderType,
		},
	}
}

func flattenLocationTemplateSimple(locTemplate *locationtemplate.LocationTemplate) []map[string]interface{} {
	if locTemplate == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id": locTemplate.ID,
		},
	}
}
