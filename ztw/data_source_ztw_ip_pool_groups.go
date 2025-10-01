package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipgroups"
)

func dataSourceIPPoolGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPPoolGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "ID of the IP address group or IP pool",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the IP address group or IP pool",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the IP group or IP pool",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "IP addresses included in the IP group or IP pool",
			},
			"creator_context": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates that the IP group or IP pool is created in Cloud & Branch Connector (EC) (only applicable value)",
			},
			"is_non_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the group is view-only (true) or editable (false)",
			},
			"extranet_ip_pool": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_predefined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceIPPoolGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *ipgroups.IPGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting ip pool group id: %d\n", id)
		res, err := ipgroups.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting ip pool group : %s\n", name)
		res, err := ipgroups.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("ip_addresses", resp.IPAddresses)
		_ = d.Set("creator_context", resp.CreatorContext)
		_ = d.Set("is_non_editable", resp.IsNonEditable)
		_ = d.Set("extranet_ip_pool", resp.ExtranetIPPool)
		_ = d.Set("is_predefined", resp.IsPredefined)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any ip pool group with name '%s' or id '%d'", name, id))
	}

	return nil
}
