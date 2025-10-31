package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/public_cloud_info"
)

func dataSourcePublicCloudInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicCloudInfoRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique ID of the AWS account.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the AWS account. Must be non-null, non-empty, unique, and 128 characters or fewer in length.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud type. The default and mandatory value is AWS. Supported values are AWS, AZURE, GCP",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique external ID for the AWS account.",
			},
			"last_mod_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The date and time when the AWS account was last modified.",
			},
			"last_sync_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The last time the AWS account was synced.",
			},
			// "permission_status": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "Indicates whether the provided credentials (external ID and AWS role name) have permission to access the AWS account.",
			// },
			"account_groups": UIDNameSchema(),
			"last_mod_user":  UIDNameSchema(),
			"region_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique ID of the region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the region.",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud type. The default and mandatory value is AWS. Supported Values: AWS, AZURE, GCP",
						},
						"status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the operational status of the region.",
						},
					},
				},
			},
			"supported_regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique ID of the supported region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the supported region.",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud type. The default and mandatory value is AWS. Supported Values: AWS, AZURE, GCP",
						},
					},
				},
			},
			"account_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.",
						},
						"aws_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.",
						},
						"aws_role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The AWS trusting role in your account. The name is non-null, non-empty, and 64 characters or fewer in length.",
						},
						"cloud_watch_group_arn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource name (ARN) of the AWS CloudWatch log group.",
						},
						"event_bus_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the event bus that sends notifications to the Zscaler service using EventBridge.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique external ID for the AWS account. If provided, it must match the externalId specified outside of accountDetails.",
						},
						"log_info_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of log information. Supported types are INFO and ERROR.",
						},
						"trouble_shooting_logging": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether logging is enabled for troubleshooting purposes.",
						},
						"trusted_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Zscaler AWS account.",
						},
						"trusted_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the trusted role in the Zscaler AWS account.",
						},
					},
				},
			},
		},
	}
}

func dataSourcePublicCloudInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *public_cloud_info.PublicCloudInfo
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting public cloud info id: %d\n", id)
		res, err := public_cloud_info.GetPublicCloudInfo(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting public cloud info by name: %s\n", name)
		res, err := public_cloud_info.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("cloud_type", resp.CloudType)
		_ = d.Set("external_id", resp.ExternalID)
		_ = d.Set("last_mod_time", resp.LastModTime)
		_ = d.Set("last_sync_time", resp.LastSyncTime)
		// _ = d.Set("permission_status", resp.PermissionStatus)

		if err := d.Set("account_groups", flattenIDExtensionsListIDs(resp.AccountGroups)); err != nil {
			return diag.FromErr(err)
		}

		if resp.LastModUser != nil {
			if err := d.Set("last_mod_user", flattenCommonIDNameExternalID(resp.LastModUser)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err := d.Set("region_status", flattenRegionStatusList(resp.RegionStatus)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("supported_regions", flattenSupportedRegionsList(resp.SupportedRegions)); err != nil {
			return diag.FromErr(err)
		}

		if resp.AccountDetails != nil {
			accountDetails := []map[string]interface{}{
				{
					"aws_account_id":           resp.AccountDetails.AwsAccountID,
					"aws_role_name":            resp.AccountDetails.AwsRoleName,
					"cloud_watch_group_arn":    resp.AccountDetails.CloudWatchGroupArn,
					"event_bus_name":           resp.AccountDetails.EventBusName,
					"external_id":              resp.AccountDetails.ExternalID,
					"log_info_type":            resp.AccountDetails.LogInfoType,
					"trouble_shooting_logging": resp.AccountDetails.TroubleShootingLogging,
					"trusted_account_id":       resp.AccountDetails.TrustedAccountID,
					"trusted_role":             resp.AccountDetails.TrustedRole,
				},
			}
			if err := d.Set("account_details", accountDetails); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any public cloud info with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenRegionStatusList(regionStatus []common.RegionStatus) []map[string]interface{} {
	if regionStatus == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(regionStatus))
	for _, item := range regionStatus {
		result = append(result, map[string]interface{}{
			"id":         item.ID,
			"name":       item.Name,
			"cloud_type": item.CloudType,
			"status":     item.Status,
		})
	}
	return result
}

func flattenSupportedRegionsList(supportedRegions []common.SupportedRegions) []map[string]interface{} {
	if supportedRegions == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(supportedRegions))
	for _, item := range supportedRegions {
		result = append(result, map[string]interface{}{
			"id":         item.ID,
			"name":       item.Name,
			"cloud_type": item.CloudType,
		})
	}
	return result
}
