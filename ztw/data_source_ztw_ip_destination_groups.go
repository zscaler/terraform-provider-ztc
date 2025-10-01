package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipdestinationgroups"
)

func dataSourceIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPDestinationGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "ID of the destination IP group",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the destination IP group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the group",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the destination IP group",
			},
			"addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "IP addresses or domain names included in the group",
			},
			"countries": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Countries included in the group",
			},
			"is_non_editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the group is view-only (true) or editable (false)",
			},
		},
	}
}

func dataSourceIPDestinationGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *ipdestinationgroups.IPDestinationGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for ip destination groups id: %d\n", id)
		res, err := ipdestinationgroups.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for ip destination groups : %s\n", name)
		res, err := ipdestinationgroups.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("type", resp.Type)
		_ = d.Set("addresses", resp.Addresses)
		_ = d.Set("description", resp.Description)
		_ = d.Set("countries", resp.Countries)
		_ = d.Set("is_non_editable", resp.IsNonEditable)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any ip destination groups with name '%s' or id '%d'", name, id))
	}

	return nil
}
