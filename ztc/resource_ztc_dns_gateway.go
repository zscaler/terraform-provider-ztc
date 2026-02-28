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
	dnsgateway "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/dns_gateway"
)

func resourceDNSGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSGatewayCreate,
		ReadContext:   resourceDNSGatewayRead,
		UpdateContext: resourceDNSGatewayUpdate,
		DeleteContext: resourceDNSGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("gateway_id", idInt)
				} else {
					resp, err := dnsgateway.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("gateway_id", resp.ID)
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
			"gateway_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DNS Gateway",
			},
			"dns_gateway_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of the DNS Gateway",
				ValidateFunc: validation.StringInSlice([]string{
					"EC_DNS_GW",
				}, false),
			},
			"ec_dns_gateway_options_primary": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary DNS gateway option for Edge Connector",
				ValidateFunc: validation.StringInSlice([]string{
					"LAN_PRI_DNS_AS_PRI",
					"LAN_SEC_DNS_AS_SEC",
					"WAN_PRI_DNS_AS_PRI",
					"WAN_SEC_DNS_AS_SEC",
				}, false),
			},
			"ec_dns_gateway_options_secondary": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary DNS gateway option for Edge Connector",
				ValidateFunc: validation.StringInSlice([]string{
					"LAN_PRI_DNS_AS_PRI",
					"LAN_SEC_DNS_AS_SEC",
					"WAN_PRI_DNS_AS_PRI",
					"WAN_SEC_DNS_AS_SEC",
				}, false),
			},
			"failure_behavior": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Choose what happens if the DNS server is unreachable",
				ValidateFunc: validation.StringInSlice([]string{
					"FAIL_RET_ERR",
					"FAIL_ALLOW_IGNORE_DNAT",
				}, false),
			},
			"primary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the primary custom DNS server",
			},
			"secondary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the secondary custom DNS server",
			},
		},
	}
}

func resourceDNSGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	req := expandDNSGateway(d)
	log.Printf("[INFO] Creating ztc dns gateway\n%+v\n", req)

	resp, err := dnsgateway.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ztc dns gateway request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("gateway_id", resp.ID)

	return resourceDNSGatewayRead(ctx, d, meta)
}

func resourceDNSGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no dns gateway id is set"))
	}
	resp, err := dnsgateway.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing ztc_dns_gateway %s from state because it no longer exists in ZTC", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting ztc dns gateway:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("gateway_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("dns_gateway_type", resp.DNSGatewayType)
	_ = d.Set("ec_dns_gateway_options_primary", resp.ECDnsGatewayOptionsPrimary)
	_ = d.Set("ec_dns_gateway_options_secondary", resp.ECDnsGatewayOptionsSecondary)
	_ = d.Set("failure_behavior", resp.FailureBehavior)
	_ = d.Set("primary_ip", resp.PrimaryIP)
	_ = d.Set("secondary_ip", resp.SecondaryIP)

	return nil
}

func resourceDNSGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] dns gateway ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating ztc dns gateway ID: %v\n", id)
	req := expandDNSGateway(d)

	if _, err := dnsgateway.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, _, err := dnsgateway.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSGatewayRead(ctx, d, meta)
}

func resourceDNSGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] dns gateway ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ztc dns gateway ID: %v\n", (d.Id()))

	if _, err := dnsgateway.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ztc dns gateway deleted")

	return nil
}

func expandDNSGateway(d *schema.ResourceData) dnsgateway.DNSGateway {
	id, _ := getIntFromResourceData(d, "gateway_id")
	result := dnsgateway.DNSGateway{
		ID:                           id,
		Name:                         d.Get("name").(string),
		DNSGatewayType:               d.Get("dns_gateway_type").(string),
		ECDnsGatewayOptionsPrimary:   d.Get("ec_dns_gateway_options_primary").(string),
		ECDnsGatewayOptionsSecondary: d.Get("ec_dns_gateway_options_secondary").(string),
		FailureBehavior:              d.Get("failure_behavior").(string),
		PrimaryIP:                    d.Get("primary_ip").(string),
		SecondaryIP:                  d.Get("secondary_ip").(string),
	}
	return result
}
