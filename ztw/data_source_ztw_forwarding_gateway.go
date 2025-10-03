package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

func dataSourceZIAForwardingGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZIAForwardingGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the forwarding gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the Forwarding Gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional details about the Forwarding Gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway type",
			},
			"fail_closed": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `"A true value indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable.
				 A false value indicates that traffic must be allowed."`,
			},
			"manual_primary": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `"Specifies the primary proxy through which traffic must be forwarded
				Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name."`,
			},
			"manual_secondary": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `"Specifies the secondary proxy through which traffic must be forwarded
				Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name."`,
			},
			"primary_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `"Type of the primary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards
				traffic through a selected data center (DC), or override (MANUAL_OVERRIDE) that forwards
				traffic through a specified IP address or domain."`,
			},
			"secondary_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `"Type of the secondary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards
				traffic through a selected data center (DC), or override (MANUAL_OVERRIDE) that forwards
				traffic through a specified IP address or domain."`,
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the ZPA gateway was last modified",
			},
			"subcloud_primary":   UIDNameSchema(),
			"subcloud_secondary": UIDNameSchema(),
			"last_modified_by":   UIDNameSchema(),
		},
	}
}

func dataSourceZIAForwardingGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *zia_forwarding_gateway.ECGateway
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for zia forwarding gateway  id: %d\n", id)
		res, _, err := zia_forwarding_gateway.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for zia forwarding gateway name: %s\n", name)
		res, err := zia_forwarding_gateway.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("fail_closed", resp.FailClosed)
		_ = d.Set("type", resp.Type)
		_ = d.Set("manual_primary", resp.ManualPrimary)
		_ = d.Set("manual_secondary", resp.ManualSecondary)
		_ = d.Set("primary_type", resp.PrimaryType)
		_ = d.Set("secondary_type", resp.SecondaryType)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("subcloud_primary", flattenCommonIDNameExternalID(resp.SubCloudPrimary)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("subcloud_secondary", flattenCommonIDNameExternalID(resp.SubCloudSecondary)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any proxy with name '%s'", name))
	}

	return nil
}
