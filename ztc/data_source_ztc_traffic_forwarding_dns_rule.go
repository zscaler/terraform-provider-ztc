package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_dns_rules"
)

func dataSourceTrafficForwardingDNSRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingDNSRuleRead,
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
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access permission available for the current user to the rule",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule type selected from the available options",
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
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank assigned to the forwarding rule",
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
			"location_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the location groups to which the forwarding rule applies",
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
			"src_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"src_ip_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.",
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
			"dest_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.",
			},
			"dest_ip_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined destination IP address groups to which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.",
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
			"dns_gateway": {
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
			"zpa_ip_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zpa ip group for which the rule is applicable.",
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

func dataSourceTrafficForwardingDNSRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *traffic_dns_rules.ECDNSRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for forwarding control rule id: %d\n", id)
		res, err := traffic_dns_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for forwarding control rule : %s\n", name)
		res, err := traffic_dns_rules.GetRulesByName(ctx, service, name)
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
		_ = d.Set("action", resp.Action)
		_ = d.Set("type", resp.Type)
		_ = d.Set("src_ips", resp.SrcIps)
		_ = d.Set("dest_addresses", resp.DestAddresses)

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationsGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("ec_groups", flattenIDNameExtensions(resp.ECGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("src_ip_groups", flattenIDNameExtensions(resp.SrcIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dns_gateway", flattenIDNameSet(resp.DNSGateway)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("zpa_ip_group", flattenIDNameSet(resp.ZPAIPGroup)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any forwarding dns rule with name '%s' or id '%d'", name, id))
	}

	return nil
}
