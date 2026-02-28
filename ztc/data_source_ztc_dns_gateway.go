package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnsgateway "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/dns_gateway"
)

func dataSourceDNSGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the DNS Gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the DNS Gateway",
			},
			"dns_gateway_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the DNS Gateway",
			},
			"ec_dns_gateway_options_primary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary DNS gateway option for Edge Connector",
			},
			"ec_dns_gateway_options_secondary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary DNS gateway option for Edge Connector",
			},
			"failure_behavior": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Choose what happens if the DNS server is unreachable",
			},
			"primary_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the primary custom DNS server",
			},
			"secondary_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the secondary custom DNS server",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when it was last modified",
			},
			"last_modified_by": UIDNameSchemaLite(),
		},
	}
}

func dataSourceDNSGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dnsgateway.DNSGateway
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for ztc dns gateway id: %d\n", id)
		res, err := dnsgateway.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for ztc dns gateway name: %s\n", name)
		res, err := dnsgateway.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("dns_gateway_type", resp.DNSGatewayType)
		_ = d.Set("ec_dns_gateway_options_primary", resp.ECDnsGatewayOptionsPrimary)
		_ = d.Set("ec_dns_gateway_options_secondary", resp.ECDnsGatewayOptionsSecondary)
		_ = d.Set("failure_behavior", resp.FailureBehavior)
		_ = d.Set("primary_ip", resp.PrimaryIP)
		_ = d.Set("secondary_ip", resp.SecondaryIP)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if resp.LastModifiedBy != nil {
			_ = d.Set("last_modified_by", flattenCustomIDNameSet(resp.LastModifiedBy))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any dns gateway with name '%s' or id '%d'", name, id))
	}

	return nil
}
