package ztc

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_log_rules"
)

var (
	trafficForwardingLogLock          sync.Mutex
	trafficForwardingLogStartingOrder int
)

func resourceTrafficForwardingLogRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficForwardingLogRuleRuleCreate,
		ReadContext:   resourceTrafficForwardingLogRuleRuleRead,
		UpdateContext: resourceTrafficForwardingLogRuleRuleUpdate,
		DeleteContext: resourceTrafficForwardingLogRuleRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := traffic_log_rules.GetRulesByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("rule_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the forwarding rule",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates whether the forwarding rule is enabled or disabled",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The order of execution for the forwarding rule order",
			},

			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"forward_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of traffic forwarding method selected from the available options",
				ValidateFunc: validation.StringInSlice([]string{
					"ECSELF",
				}, false),
			},
			"ec_groups":     setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies"),
			"locations":     setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations."),
			"proxy_gateway": setIdNameSchemaCustom(1, "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method."),
		},
	}
}

func validatePredefinedLogRules(req traffic_log_rules.ECTrafficLogRules) error {
	if req.Name == "ZPA Resolver" || req.Name == "Redirect Resolution of Zscaler Domains to WAN CTR" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	return nil
}

func resourceTrafficForwardingLogRuleRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandForwardingLogRule(d)
	log.Printf("[INFO] Creating ztc traffic log forwarding rule\n%+v\n", req)

	start := time.Now()

	for {
		trafficForwardingLogLock.Lock()
		if trafficForwardingLogStartingOrder == 0 {
			list, _ := traffic_log_rules.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > trafficForwardingLogStartingOrder {
					trafficForwardingLogStartingOrder = r.Order
				}
			}
			if trafficForwardingLogStartingOrder == 0 {
				trafficForwardingLogStartingOrder = 1
			}
		}
		trafficForwardingLogLock.Unlock()
		startWithoutLocking := time.Now()

		intendedOrder := req.Order
		intendedRank := req.Rank
		if intendedRank < 7 {
			// always start rank 7 rules at the next available order after all ranked rules
			req.Rank = 7
		}
		req.Order = trafficForwardingLogStartingOrder
		resp, err := traffic_log_rules.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, intendedOrder, req.Rank, currentOrderVsRankWording(ctx, zClient), err))
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created ztc traffic log forwarding rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
		resourceType := "traffic_forwarding_log_rule"

		reorderWithBeforeReorder(
			OrderRule{Order: intendedOrder, Rank: intendedRank},
			resp.ID,
			resourceType,
			func() (int, error) {
				allRules, err := traffic_log_rules.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := traffic_log_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}

				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = traffic_log_rules.Update(ctx, service, id, rule)
				return err
			},
			nil, // Remove beforeReorder function to avoid adding too many rules to the map
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		markOrderRuleAsDone(resp.ID, resourceType)

		return resourceTrafficForwardingLogRuleRuleRead(ctx, d, meta)
	}
}

func resourceTrafficForwardingLogRuleRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ztc traffic log forwarding rule id is set"))
	}
	resp, err := traffic_log_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing traffic log forwarding rule %s from state because it no longer exists in ZTC", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting traffic log forwarding rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("state", resp.State)
	_ = d.Set("forward_method", resp.ForwardMethod)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ec_groups", flattenIDExtensionsListIDs(resp.ECGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTrafficForwardingLogRuleRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] traffic log forwarding rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("traffic log forwarding rule ID not set"))
	}
	log.Printf("[INFO] Updating traffic log forwarding rule ID: %v\n", id)
	req := expandForwardingLogRule(d)

	if _, err := traffic_log_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	existingRules, err := traffic_log_rules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all traffic log forwarding rules: %v", err)
	}
	sort.Slice(existingRules, func(i, j int) bool {
		return existingRules[i].Rank < existingRules[j].Rank || (existingRules[i].Rank == existingRules[j].Rank && existingRules[i].Order < existingRules[j].Order)
	})
	intendedOrder := req.Order
	intendedRank := req.Rank
	nextAvailableOrder := existingRules[len(existingRules)-1].Order
	// always start rank 7 rules at the next available order after all ranked rules
	req.Rank = 7

	req.Order = nextAvailableOrder

	_, err = traffic_log_rules.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating traffic log forwarding rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "traffic_forwarding_log_rule",
		func() (int, error) {
			allRules, err := traffic_log_rules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := traffic_log_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = traffic_log_rules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceTrafficForwardingLogRuleRuleRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "traffic_forwarding_log_rule")

	return nil
}

func resourceTrafficForwardingLogRuleRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("traffic log forwarding rule ID not set: %v", id))
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := traffic_log_rules.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving traffic log forwarding rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedLogRules(*rule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting traffic log forwarding rule ID: %v", id)
	if _, err := traffic_log_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting traffic log forwarding rule %d: %v", id, err))
	}

	d.SetId("")
	log.Printf("[INFO] Traffic log forwarding rule deleted")

	return nil
}

func expandForwardingLogRule(d *schema.ResourceData) traffic_log_rules.ECTrafficLogRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandForwardingControlRule: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := traffic_log_rules.ECTrafficLogRules{
		ID:            id,
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Order:         order,
		Rank:          d.Get("rank").(int),
		State:         d.Get("state").(string),
		ForwardMethod: d.Get("forward_method").(string),
		Locations:     expandIDNameExtensionsSet(d, "locations"),
		ECGroups:      expandIDNameExtensionsSet(d, "ec_groups"),
		ProxyGateway:  expandIDNameSet(d, "proxy_gateway"),
	}
	return result
}
