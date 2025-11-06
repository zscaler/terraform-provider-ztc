package ztc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/public_cloud_info"
)

func resourcePublicCloudInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicCloudInfoCreate,
		ReadContext:   resourcePublicCloudInfoRead,
		UpdateContext: resourcePublicCloudInfoUpdate,
		DeleteContext: resourcePublicCloudInfoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("cloud_id", idInt)
				} else {
					resp, err := public_cloud_info.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("cloud_id", resp.ID)
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
			"cloud_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the AWS account. Must be non-null, non-empty, unique, and 128 characters or fewer in length.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cloud type. The default and mandatory value is AWS. Supported values are AWS, AZURE, GCP",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique external ID for the AWS account.",
			},
			"account_details": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.",
						},
						"aws_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The AWS trusting role in your account. The name is non-null, non-empty, and 64 characters or fewer in length.",
						},
						"cloud_watch_group_arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource name (ARN) of the AWS CloudWatch log group.",
						},
						"event_bus_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the event bus that sends notifications to the Zscaler service using EventBridge.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique external ID for the AWS account. If provided, it must match the externalId specified outside of accountDetails.",
						},
						"log_info_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of log information. Supported types are INFO and ERROR.",
						},
						"trouble_shooting_logging": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether logging is enabled for troubleshooting purposes.",
						},
						"trusted_account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the Zscaler AWS account.",
						},
						"trusted_role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the trusted role in the Zscaler AWS account.",
						},
					},
				},
			},
			"account_groups":    setIDsSchemaTypeCustom(nil, "An immutable reference to an entity, which consists of ID, name, etc."),
			"supported_regions": setIDsSchemaTypeCustom(nil, "Regions supported by Zscaler’s Tag Discovery Service."),
		},
	}
}

func resourcePublicCloudInfoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandPublicCloudInfo(d)
	log.Printf("[INFO] Creating zia public cloud info\n%+v\n", req)

	resp, err := public_cloud_info.CreatePublicCloudInfo(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia public cloud info request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("cloud_id", resp.ID)

	return resourcePublicCloudInfoRead(ctx, d, meta)
}

func resourcePublicCloudInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no public cloud info id is set"))
	}
	resp, err := public_cloud_info.GetPublicCloudInfo(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia public cloud info %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia ip source groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("cloud_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("cloud_type", resp.CloudType)
	_ = d.Set("external_id", resp.ExternalID)

	if err := d.Set("account_groups", flattenIDExtensionsListIDs(resp.AccountGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("supported_regions", flattenIDSupportedRegions(resp.SupportedRegions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("account_details", flattenAccountDetails(resp.AccountDetails)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePublicCloudInfoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		log.Printf("[ERROR] public cloud info ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia public cloud info ID: %v\n", id)
	req := expandPublicCloudInfo(d)
	if _, err := public_cloud_info.GetPublicCloudInfo(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := public_cloud_info.UpdatePublicCloudInfo(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePublicCloudInfoRead(ctx, d, meta)
}

func resourcePublicCloudInfoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		log.Printf("[ERROR] public cloud info ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia public cloud info ID: %v\n", (d.Id()))

	if err := public_cloud_info.DeletePublicCloudInfo(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia public cloud info deleted")

	return nil
}

func expandPublicCloudInfo(d *schema.ResourceData) public_cloud_info.PublicCloudInfo {
	id, _ := getIntFromResourceData(d, "cloud_id")
	return public_cloud_info.PublicCloudInfo{
		ID:               id,
		Name:             d.Get("name").(string),
		CloudType:        d.Get("cloud_type").(string),
		ExternalID:       d.Get("external_id").(string),
		AccountGroups:    expandIDNameExtensionsSet(d, "account_groups"),
		SupportedRegions: expandIDSupportedRegionsSet(d, "supported_regions"),
		AccountDetails:   expandAccountDetails(d, "account_details"),
	}
}

func flattenAccountDetails(gp *public_cloud_info.AccountDetails) []map[string]interface{} {
	if gp == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"aws_account_id":           gp.AwsAccountID,
			"aws_role_name":            gp.AwsRoleName,
			"cloud_watch_group_arn":    gp.CloudWatchGroupArn,
			"event_bus_name":           gp.EventBusName,
			"external_id":              gp.ExternalID,
			"trouble_shooting_logging": gp.TroubleShootingLogging,
			"trusted_account_id":       gp.TrustedAccountID,
			"trusted_role":             gp.TrustedRole,
		},
	}
}

func flattenIDSupportedRegions(list []common.SupportedRegions) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

func expandAccountDetails(d *schema.ResourceData, key string) *public_cloud_info.AccountDetails {
	accountDetailsList, ok := d.GetOk(key)
	if !ok {
		return nil
	}

	list := accountDetailsList.([]interface{})
	if len(list) == 0 {
		return nil
	}

	item := list[0].(map[string]interface{})
	return &public_cloud_info.AccountDetails{
		AwsAccountID:           item["aws_account_id"].(string),
		AwsRoleName:            item["aws_role_name"].(string),
		CloudWatchGroupArn:     item["cloud_watch_group_arn"].(string),
		EventBusName:           item["event_bus_name"].(string),
		ExternalID:             item["external_id"].(string),
		TroubleShootingLogging: item["trouble_shooting_logging"].(bool),
		TrustedAccountID:       item["trusted_account_id"].(string),
		TrustedRole:            item["trusted_role"].(string),
	}
}

func expandIDSupportedRegionsSet(d *schema.ResourceData, key string) []common.SupportedRegions {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.SupportedRegions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
					result = append(result, common.SupportedRegions{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.SupportedRegions{}
}
