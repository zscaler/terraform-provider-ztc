package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservicegroups"
)

func dataSourceNetworkServiceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkServiceGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "ID of network service",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Name of network service",
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_tcp_ports":  dataNetworkPortsSchema("Source TCP ports"),
						"dest_tcp_ports": dataNetworkPortsSchema("Destination TCP ports"),
						"src_udp_ports":  dataNetworkPortsSchema("Source UDP ports"),
						"dest_udp_ports": dataNetworkPortsSchema("Destination UDP ports"),
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of network service: standard, predefined, or custom",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of network service",
						},
						"is_name_l10n_tag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkServiceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *networkservicegroups.NetworkServiceGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network service group id: %d\n", id)
		res, err := networkservicegroups.GetNetworkServiceGroups(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network service group : %s\n", name)
		res, err := networkservicegroups.GetNetworkServiceGroupsByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)

		if err := d.Set("services", flattenServices(resp.Services)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenServices(service []networkservicegroups.Services) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"id":               val.ID,
			"name":             val.Name,
			"description":      val.Description,
			"is_name_l10n_tag": val.IsNameL10nTag,
		}
	}

	return services
}
