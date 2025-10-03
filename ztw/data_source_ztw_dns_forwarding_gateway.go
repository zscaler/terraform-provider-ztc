package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
)

func dataSourceDNSForwardingGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSForwardingGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the DNS Forwarding Gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the DNS Forwarding Gateway",
			},
			"failure_behavior": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Choose what happens if the DNS server is unreachable.",
			},
			"dns_gateway_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the DNS Forwarding Gateway",
			},
			"primary_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the primary custom DNS server.",
			},
			"secondary_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the secondary custom DNS server.",
			},
			"ec_dns_gateway_options_primary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the primary LAN DNS Server",
			},
			"ec_dns_gateway_options_secondary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the secondary LAN DNS Server.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the it was last modified",
			},
			"last_modified_by": UIDNameSchema(),
		},
	}
}

func dataSourceDNSForwardingGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dns_forwarding_gateway.DNSGateway
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for ztw dns forwarding gateway  id: %d\n", id)
		res, _, err := dns_forwarding_gateway.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for ztw dns  forwarding gateway name: %s\n", name)
		res, err := dns_forwarding_gateway.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("failure_behavior", resp.FailureBehavior)
		_ = d.Set("dns_gateway_type", resp.DNSGatewayType)
		_ = d.Set("primary_ip", resp.PrimaryIP)
		_ = d.Set("secondary_ip", resp.SecondaryIP)
		_ = d.Set("ec_dns_gateway_options_primary", resp.ECDNSGatewayOptionsPrimary)
		_ = d.Set("ec_dns_gateway_options_secondary", resp.ECDNSGatewayOptionsSecondary)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any dns forwarding gateway with name '%s'", name))
	}

	return nil
}
