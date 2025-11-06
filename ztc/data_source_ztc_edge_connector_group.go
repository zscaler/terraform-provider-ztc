package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/ecgroup"
)

func dataSourceEdgeConnectorGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEdgeConnectorGroupRead,
		Schema: MergeSchema(
			ecGroupSchemaData(), map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
			},
		),
	}
}

func dataSourceEdgeConnectorGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *ecgroup.EcGroup
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for edge connector group id: %d\n", id)
		res, err := ecgroup.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for edge connector group name: %s\n", name)
		res, err := ecgroup.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("desc", resp.Description)
		_ = d.Set("deploy_type", resp.DeployType)
		_ = d.Set("platform", resp.Platform)
		_ = d.Set("aws_availability_zone", resp.AWSAvailabilityZone)
		_ = d.Set("azure_availability_zone", resp.AzureAvailabilityZone)
		_ = d.Set("max_ec_count", resp.MaxEcCount)
		_ = d.Set("tunnel_mode", resp.TunnelMode)

		// Convert Status from []string to string (take the first element)
		var statusStr string
		if len(resp.Status) > 0 {
			statusStr = resp.Status[0]
		}
		_ = d.Set("status", statusStr)

		if err := d.Set("location", flattenGeneralPurpose(resp.Location)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("prov_template", flattenGeneralPurpose(resp.ProvTemplate)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ec_vms", flattenECVms(resp.ECVMs)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.Errorf("couldn't find any edge connector group with name '%s'", name)
	}

	return nil
}
