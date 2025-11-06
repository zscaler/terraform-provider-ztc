package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

func dataSourceNetworkServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkServicesRead,
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
	}
}

func dataSourceNetworkServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *networkservices.NetworkServices
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network services id: %d\n", id)
		res, err := networkservices.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network services : %s\n", name)
		res, err := networkservices.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	/*
		protocol, _ := d.Get("protocol").(string)
		if resp == nil && protocol != "" {
			log.Printf("[INFO] Getting network services : %s\n", protocol)
			res, err := zClient.networkservices.GetByProtocol(d.Get("protocol").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			resp = res
		}
	*/
	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("tag", resp.Tag)
		_ = d.Set("type", resp.Type)
		_ = d.Set("description", resp.Description)
		_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

		if err := d.Set("src_tcp_ports", flattenNetwordPorts(resp.SrcTCPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_tcp_ports", flattenNetwordPorts(resp.DestTCPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("src_udp_ports", flattenNetwordPorts(resp.SrcUDPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_udp_ports", flattenNetwordPorts(resp.DestUDPPorts)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id))
	}

	return nil
}
