package ztc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/activation"
)

func dataSourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceActivationStatusRead,
		Schema: map[string]*schema.Schema{
			"org_edit_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization policy edit status",
			},
			"org_last_activate_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization policy last activation status",
			},
			"admin_activate_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Admin activation status",
			},
			"admin_status_map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Admin status",
			},
		},
	}
}

func dataSourceActivationStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := activation.GetActivationStatus(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		d.SetId("activation")
		_ = d.Set("org_edit_status", resp.OrgEditStatus)
		_ = d.Set("org_last_activate_status", resp.OrgLastActivateStatus)
		_ = d.Set("admin_status_map", resp.AdminStatusMap)
		_ = d.Set("admin_activate_status", resp.AdminActivateStatus)

	} else {
		return diag.Errorf("couldn't find the activation status")
	}

	return nil
}
