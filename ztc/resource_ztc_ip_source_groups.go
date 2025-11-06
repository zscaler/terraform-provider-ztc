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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipsourcegroups"
)

func resourceIPSourceGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPSourceGroupsGroupsCreate,
		ReadContext:   resourceIPSourceGroupsGroupsRead,
		UpdateContext: resourceIPSourceGroupsGroupsUpdate,
		DeleteContext: resourceIPSourceGroupsGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := ipsourcegroups.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("group_id", resp.ID)
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
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceIPSourceGroupsGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFWIPSourceGroups(d)
	log.Printf("[INFO] Creating zia ip source groups\n%+v\n", req)

	resp, _, err := ipsourcegroups.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	return resourceIPSourceGroupsGroupsRead(ctx, d, meta)
}

func resourceIPSourceGroupsGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ip source groups id is set"))
	}
	resp, err := ipsourcegroups.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip source groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia ip source groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("description", resp.Description)

	return nil
}

func resourceIPSourceGroupsGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia ip source groups ID: %v\n", id)
	req := expandFWIPSourceGroups(d)
	if _, err := ipsourcegroups.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := ipsourcegroups.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceIPSourceGroupsGroupsRead(ctx, d, meta)
}

func resourceIPSourceGroupsGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip source groups ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		ctx,
		zClient,
		id,
		"IPSourceGroups",
		func(r *forwarding_rules.ForwardingRules) []common.IDNameExtensions {
			return r.SrcIpGroups
		},
		func(r *forwarding_rules.ForwardingRules, ids []common.IDNameExtensions) {
			r.SrcIpGroups = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := ipsourcegroups.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia ip source groups deleted")

	return nil
}

func expandFWIPSourceGroups(d *schema.ResourceData) ipsourcegroups.IPSourceGroups {
	id, _ := getIntFromResourceData(d, "group_id")
	return ipsourcegroups.IPSourceGroups{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IPAddresses: SetToStringList(d, "ip_addresses"),
	}
}
