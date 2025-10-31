package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/account_groups"
)

func dataSourceAccountGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The ID of the AWS account group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the AWS account group. Must be non-null, non-empty, unique, and 128 characters or fewer in length.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the IP group or IP pool",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud type. The default and manadatory value is AWS. Supported values are AWS, AZURE, GCP",
			},
			"cloud_connector_groups": UIDNameSchema(),
			"public_cloud_accounts":  UIDNameSchema(),
		},
	}
}

func dataSourceAccountGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *account_groups.AccountGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting account group id: %d\n", id)
		res, err := account_groups.GetAccountGroup(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(res) > 0 {
			resp = &res[0]
		}
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting account group : %s\n", name)
		res, err := account_groups.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("cloud_type", resp.CloudType)

		if err := d.Set("public_cloud_accounts", flattenIDExtensionsListIDs(resp.PublicCloudAccounts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("cloud_connector_groups", flattenIDExtensionsListIDs(resp.CloudConnectorGroups)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any account group with name '%s' or id '%d'", name, id))
	}

	return nil
}
