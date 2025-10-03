package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_log_rules"
)

func dataSourceTrafficForwardingLogRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingLogRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the forwarding rule",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the forwarding rule is enabled or disabled",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The order of execution for the forwarding rule order",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule type selected from the available options",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"forward_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of traffic forwarding method selected from the available options",
			},
			"default_rule": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the forwarding rule is a default rule",
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"proxy_gateway": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"ec_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
		},
	}
}

func dataSourceTrafficForwardingLogRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *traffic_log_rules.ECTrafficLogRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for forwarding log control rule id: %d\n", id)
		res, err := traffic_log_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for forwarding log control rule : %s\n", name)
		res, err := traffic_log_rules.GetRulesByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("state", resp.State)
		_ = d.Set("type", resp.Type)
		_ = d.Set("default_rule", resp.DefaultRule)
		_ = d.Set("forward_method", resp.ForwardMethod)

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("ec_groups", flattenIDNameExtensions(resp.ECGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any forwarding log control rule with name '%s' or id '%d'", name, id))
	}

	return nil
}
