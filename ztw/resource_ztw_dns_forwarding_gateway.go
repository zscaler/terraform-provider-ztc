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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
)

func resourceDNSForwardingGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSForwardingGatewayCreate,
		ReadContext:   resourceDNSForwardingGatewayRead,
		UpdateContext: resourceDNSForwardingGatewayUpdate,
		DeleteContext: resourceDNSForwardingGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("gateway_id", idInt)
				} else {
					resp, err := dns_forwarding_gateway.GetByName(ctx, service, id)
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
				Description: "The name of the DNS Forwarding Gateway",
			},
			"failure_behavior": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Choose what happens if the DNS server is unreachable.",
				ValidateFunc: validation.StringInSlice([]string{
					"FAIL_RET_ERR",
					"FAIL_ALLOW_IGNORE_DNAT",
				}, false),
			},
			"primary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the primary custom DNS server.",
			},
			"secondary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the secondary custom DNS server.",
			},
			"ec_dns_gateway_options_primary": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the primary LAN DNS Server",
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
				Description: "IP address of the secondary LAN DNS Server.",
				ValidateFunc: validation.StringInSlice([]string{
					"LAN_PRI_DNS_AS_PRI",
					"LAN_SEC_DNS_AS_SEC",
					"WAN_PRI_DNS_AS_PRI",
					"WAN_SEC_DNS_AS_SEC",
				}, false),
			},
		},
	}
}

func resourceDNSForwardingGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandDNSForwardingGateway(d)
	log.Printf("[INFO] Creating ZTW DNS forwarding gateway \n%+v\n", req)

	resp, _, err := dns_forwarding_gateway.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZTW DNS forwarding gateway request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("gateway_id", resp.ID)

	return resourceDNSForwardingGatewayRead(ctx, d, meta)
}

func resourceDNSForwardingGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ZTW forwarding gateway id is set"))
	}
	resp, _, err := dns_forwarding_gateway.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing ZTW forwarding gateway %s from state because it no longer exists in ZTW", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting ZTW forwarding gateway:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("gateway_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("failure_behavior", resp.FailureBehavior)
	_ = d.Set("primary_ip", resp.PrimaryIP)
	_ = d.Set("secondary_ip", resp.SecondaryIP)
	_ = d.Set("ec_dns_gateway_options_primary", resp.ECDNSGatewayOptionsPrimary)
	_ = d.Set("ec_dns_gateway_options_secondary", resp.ECDNSGatewayOptionsSecondary)

	return nil
}

func resourceDNSForwardingGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] gateway ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating ZTW DNS forwarding gateway ID: %v\n", id)
	req := expandDNSForwardingGateway(d)
	if _, _, err := dns_forwarding_gateway.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := dns_forwarding_gateway.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSForwardingGatewayRead(ctx, d, meta)
}

func resourceDNSForwardingGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] ZTW DNS forwarding gateway not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ZTW DNS forwarding gateway ID: %v\n", (d.Id()))

	if _, err := dns_forwarding_gateway.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ztw forwarding gateway deleted")

	return nil
}

func expandDNSForwardingGateway(d *schema.ResourceData) dns_forwarding_gateway.DNSGateway {
	id, _ := getIntFromResourceData(d, "gateway_id")
	result := dns_forwarding_gateway.DNSGateway{
		ID:                           id,
		Name:                         d.Get("name").(string),
		FailureBehavior:              d.Get("failure_behavior").(string),
		PrimaryIP:                    d.Get("primary_ip").(string),
		SecondaryIP:                  d.Get("secondary_ip").(string),
		ECDNSGatewayOptionsPrimary:   d.Get("ec_dns_gateway_options_primary").(string),
		ECDNSGatewayOptionsSecondary: d.Get("ec_dns_gateway_options_secondary").(string),
	}
	return result
}
