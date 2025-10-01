package ztw

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/activation"
)

func resourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationStatusCreate,
		ReadContext:   resourceActivationStatusRead,
		DeleteContext: resourceFuncNoOp,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"org_edit_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EDITS_CLEARED",
					"EDITS_PRESENT",
					"EDITS_ACTIVATED_ON_RESTART",
				}, false),
			},
			"org_last_activate_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CAC_ACTV_UNKNOWN",
					"CAC_ACTV_UI",
					"CAC_ACTV_OLD_UI",
					"CAC_ACTV_SUPERADMIN",
					"CAC_ACTV_AUTOSYNC",
					"CAC_ACTV_TIMER",
				}, false),
			},
			"admin_status_map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Admin status",
			},
			"admin_activate_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ADM_LOGGED_IN",
					"ADM_EDITING",
					"ADM_ACTV_QUEUED",
					"ADM_ACTIVATING",
					"ADM_ACTV_DONE",
					"ADM_ACTV_FAIL",
					"ADM_EXPIRED",
				}, false),
			},
		},
	}
}

func resourceActivationStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandActivationStatus(d)
	log.Printf("[INFO] Performing configuration activation\n%+v\n", req)

	resp, err := activation.UpdateActivationStatus(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Configuration activation successfull. %v\n", resp.AdminActivateStatus)
	d.SetId("activation")
	return resourceActivationStatusRead(ctx, d, meta)
}

func resourceActivationStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := activation.GetActivationStatus(ctx, service)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Cannot obtain activation %s from ZTW", d.Id())
			// Activation is not an actual object; hence no ID should be set.
			// d.SetId("")
			// return nil
		}

		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading activation status: %+v\n", resp)
	_ = d.Set("org_edit_status", resp.OrgEditStatus)
	_ = d.Set("org_last_activate_status", resp.OrgLastActivateStatus)
	_ = d.Set("admin_status_map", resp.AdminStatusMap)
	_ = d.Set("admin_activate_status", resp.AdminActivateStatus)

	return nil
}

func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Delete doesn't actually do anything, because an activation can't be deleted.
	return nil
}

func expandActivationStatus(d *schema.ResourceData) activation.ECAdminActivation {
	return activation.ECAdminActivation{
		OrgEditStatus:         d.Get("org_edit_status").(string),
		OrgLastActivateStatus: d.Get("org_last_activate_status").(string),
		AdminStatusMap:        d.Get("admin_status_map").(map[string]interface{}),
		AdminActivateStatus:   d.Get("admin_activate_status").(string),
	}
}
