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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservicegroups"
)

func resourceNetworkServiceGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkServiceGroupsCreate,
		ReadContext:   resourceNetworkServiceGroupsRead,
		UpdateContext: resourceNetworkServiceGroupsUpdate,
		DeleteContext: resourceNetworkServiceGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := networkservicegroups.GetNetworkServiceGroupsByName(ctx, service, id)
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
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"services": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "list of services IDs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func resourceNetworkServiceGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandNetworkServiceGroups(d)
	log.Printf("[INFO] Creating network service groups\n%+v\n", req)

	resp, err := networkservicegroups.CreateNetworkServiceGroups(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia network service groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	return resourceNetworkServiceGroupsRead(ctx, d, meta)
}

func resourceNetworkServiceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no network service groups id is set"))
	}
	resp, err := networkservicegroups.GetNetworkServiceGroups(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network service groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting network service groups :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)

	if err := d.Set("services", flattenServicesSimple(resp.Services)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenServicesSimple(list []networkservicegroups.Services) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]int, len(list))
	for i, item := range list {
		ids[i] = item.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceNetworkServiceGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service groups ID: %v\n", id)
	req := expandNetworkServiceGroups(d)
	if _, err := networkservicegroups.GetNetworkServiceGroups(ctx, service, req.ID); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := networkservicegroups.UpdateNetworkServiceGroups(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkServiceGroupsRead(ctx, d, meta)
}

func resourceNetworkServiceGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network service groups ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		ctx,
		zClient,
		id,
		"NwApplicationGroups",
		func(r *forwarding_rules.ForwardingRules) []common.IDNameExtensions {
			return r.NwApplicationGroups
		},
		func(r *forwarding_rules.ForwardingRules, ids []common.IDNameExtensions) {
			r.NwApplicationGroups = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := networkservicegroups.DeleteNetworkServiceGroups(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] network service groups deleted")

	return nil
}

func expandNetworkServiceGroups(d *schema.ResourceData) networkservicegroups.NetworkServiceGroups {
	id, _ := getIntFromResourceData(d, "group_id")
	result := networkservicegroups.NetworkServiceGroups{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Services:    expandServicesSet(d),
	}

	return result
}

func expandServicesSet(d *schema.ResourceData) []networkservicegroups.Services {
	setInterface, ok := d.GetOk("services")
	if ok {
		set := setInterface.(*schema.Set)
		var result []networkservicegroups.Services
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				idSet, ok := itemMap["id"].(*schema.Set)
				if ok {
					for _, id := range idSet.List() {
						result = append(result, networkservicegroups.Services{
							ID: id.(int),
						})
					}
				}
			}
		}
		return result
	}
	return []networkservicegroups.Services{}
}
