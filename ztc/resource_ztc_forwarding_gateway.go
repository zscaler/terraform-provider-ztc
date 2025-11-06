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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

func resourceForwardingGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingGatewayCreate,
		ReadContext:   resourceForwardingGatewayRead,
		UpdateContext: resourceForwardingGatewayUpdate,
		DeleteContext: resourceForwardingGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("gateway_id", idInt)
				} else {
					resp, err := zia_forwarding_gateway.GetByName(ctx, service, id)
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
				Computed:    true,
				Optional:    true,
				Description: "The name of the Forwarding Gateway",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional details about the Forwarding Gateway",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway type",
				ValidateFunc: validation.StringInSlice([]string{
					"PROXYCHAIN",
					"ZIA",
					"ECSELF",
				}, false),
			},
			"fail_closed": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `"A true value indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable.
				 A false value indicates that traffic must be allowed."`,
			},
			"manual_primary": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `"Specifies the primary proxy through which traffic must be forwarded
				Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name."`,
			},
			"manual_secondary": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `"Specifies the secondary proxy through which traffic must be forwarded
				Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name."`,
			},
			"primary_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"AUTO",
					"MANUAL_OVERRIDE",
					"SUBCLOUD",
					"VZEN",
					"PZEN",
					"DC",
				}, false),
				Description: `"Type of the primary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards
				traffic through a selected data center (DC), or override (MANUAL_OVERRIDE) that forwards
				traffic through a specified IP address or domain."`,
			},
			"secondary_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"AUTO",
					"MANUAL_OVERRIDE",
					"SUBCLOUD",
					"VZEN",
					"PZEN",
					"DC",
				}, false),
				Description: `"Type of the secondary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards
				traffic through a selected data center (DC), or override (MANUAL_OVERRIDE) that forwards
				traffic through a specified IP address or domain."`,
			},
			"subcloud_primary":   setIdNameSchemaCustom(1, "If a manual (DC) primary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified DC"),
			"subcloud_secondary": setIdNameSchemaCustom(1, "If a manual (DC) secondary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified DC"),
		},
	}
}

func resourceForwardingGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandForwardingGateway(d)
	log.Printf("[INFO] Creating ZTW forwarding gateway \n%+v\n", req)

	resp, _, err := zia_forwarding_gateway.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZTW forwarding gateway request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("gateway_id", resp.ID)

	return resourceForwardingGatewayRead(ctx, d, meta)
}

func resourceForwardingGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ZTW forwarding gateway id is set"))
	}
	resp, _, err := zia_forwarding_gateway.Get(ctx, service, id)
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
	_ = d.Set("description", resp.Description)
	_ = d.Set("fail_closed", resp.FailClosed)
	_ = d.Set("type", resp.Type)
	_ = d.Set("manual_primary", resp.ManualPrimary)
	_ = d.Set("manual_secondary", resp.ManualSecondary)
	_ = d.Set("primary_type", resp.PrimaryType)
	_ = d.Set("secondary_type", resp.SecondaryType)

	if err := d.Set("subcloud_primary", flattenCommonIDNameExternalID(resp.SubCloudPrimary)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subcloud_secondary", flattenCommonIDNameExternalID(resp.SubCloudSecondary)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceForwardingGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] gateway ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating ZTW forwarding gateway ID: %v\n", id)
	req := expandForwardingGateway(d)
	if _, _, err := zia_forwarding_gateway.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zia_forwarding_gateway.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceForwardingGatewayRead(ctx, d, meta)
}

func resourceForwardingGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] ZTW forwarding gateway not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ZTW forwarding gateway ID: %v\n", (d.Id()))

	if _, err := zia_forwarding_gateway.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ztw forwarding gateway deleted")

	return nil
}

func expandForwardingGateway(d *schema.ResourceData) zia_forwarding_gateway.ECGateway {
	id, _ := getIntFromResourceData(d, "gateway_id")
	result := zia_forwarding_gateway.ECGateway{
		ID:                id,
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Type:              d.Get("type").(string),
		FailClosed:        d.Get("fail_closed").(bool),
		ManualPrimary:     d.Get("manual_primary").(string),
		ManualSecondary:   d.Get("manual_secondary").(string),
		PrimaryType:       d.Get("primary_type").(string),
		SecondaryType:     d.Get("secondary_type").(string),
		SubCloudPrimary:   expandCommonIDNameExternalID(d, "subcloud_primary"),
		SubCloudSecondary: expandCommonIDNameExternalID(d, "subcloud_secondary"),
	}
	return result
}

func expandCommonIDNameExternalID(d *schema.ResourceData, key string) *common.CommonIDNameExternalID {
	idNameList, ok := d.Get(key).(*schema.Set)
	if !ok || idNameList.Len() == 0 {
		return nil
	}

	for _, v := range idNameList.List() {
		item := v.(map[string]interface{})
		return &common.CommonIDNameExternalID{
			ID:   item["id"].(int),
			Name: item["name"].(string),
		}
	}

	return nil
}
