package ztw

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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_dns_rules"
)

var (
	trafficForwardingDNSLock          sync.Mutex
	trafficForwardingDNSStartingOrder int
)

func resourceTrafficForwardingDNSRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficForwardingDNSRuleCreate,
		ReadContext:   resourceTrafficForwardingDNSRuleRead,
		UpdateContext: resourceTrafficForwardingDNSRuleUpdate,
		DeleteContext: resourceTrafficForwardingDNSRuleDelete,

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
					resp, err := traffic_dns_rules.GetRulesByName(ctx, service, id)
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
				Optional:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the forwarding rule",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The order of execution for the forwarding rule order",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the Firewall Filtering policy rule is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK",
					"REDIR_ZPA",
					"REDIR_REQ",
				}, false),
			},
			"src_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"dest_addresses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.",
			},
			"locations":       setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations."),
			"location_groups": setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the forwarding rule applies"),
			"ec_groups":       setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies"),
			"src_ip_groups":   setIDsSchemaTypeCustom(nil, "Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group"),
			"dest_ip_groups":  setIDsSchemaTypeCustom(nil, "User-defined destination IP address groups to which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group"),
			"dns_gateway":     setIdNameSchemaCustom(1, "The dns gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method."),
			"zpa_ip_group":    setIdNameSchemaCustom(1, "The zpa ip group for which the rule is applicable. This field is applicable only for action REDIR_ZPA"),
		},
	}
}

func validatePredefinedDNSRules(req traffic_dns_rules.ECDNSRules) error {
	if req.Name == "ZPA Resolver" || req.Name == "Redirect Resolution of Zscaler Domains to WAN CTR" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	return nil
}

func resourceTrafficForwardingDNSRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandForwardingDNSRule(d)
	log.Printf("[INFO] Creating zia forwarding dns rule\n%+v\n", req)

	start := time.Now()

	trafficForwardingDNSLock.Lock()
	if trafficForwardingDNSStartingOrder == 0 {
		list, _ := traffic_dns_rules.GetAll(ctx, service)
		for _, r := range list {
			if r.Order > trafficForwardingDNSStartingOrder {
				trafficForwardingDNSStartingOrder = r.Order
			}
		}
		if trafficForwardingDNSStartingOrder == 0 {
			trafficForwardingDNSStartingOrder = 1
		} else {
			trafficForwardingDNSStartingOrder++
		}
	}
	trafficForwardingDNSLock.Unlock()
	startWithoutLocking := time.Now()

	// Store the intended order from HCL
	intendedOrder := req.Order
	intendedRank := req.Rank
	if intendedRank < 7 {
		// always start rank 7 rules at the next available order after all ranked rules
		req.Rank = 7
	}
	req.Order = trafficForwardingDNSStartingOrder
	resp, err := traffic_dns_rules.Create(ctx, service, &req)

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

	log.Printf("[INFO] Created zia forwarding dns rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
	// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
	resourceType := "traffic_forwarding_dns_rule"

	reorderWithBeforeReorder(
		OrderRule{Order: intendedOrder, Rank: intendedRank},
		resp.ID,
		resourceType,
		func() (int, error) {
			allRules, err := traffic_dns_rules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			// Custom updateOrder that handles predefined rules
			rule, err := traffic_dns_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = traffic_dns_rules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	markOrderRuleAsDone(resp.ID, resourceType)

	return resourceTrafficForwardingDNSRuleRead(ctx, d, meta)
}

func resourceTrafficForwardingDNSRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia forwarding dns rule id is set"))
	}
	resp, err := traffic_dns_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing forwarding dns rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting forwarding dns rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("state", resp.State)
	// _ = d.Set("type", resp.Type)
	_ = d.Set("action", resp.Action)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationsGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ec_groups", flattenIDExtensionsListIDs(resp.ECGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_ip_groups", flattenIDExtensionsListIDs(resp.SrcIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDExtensionsListIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dns_gateway", flattenIDNameSet(resp.DNSGateway)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTrafficForwardingDNSRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding dns rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("forwarding dns rule ID not set"))
	}
	log.Printf("[INFO] Updating forwarding dns rule ID: %v\n", id)
	req := expandForwardingDNSRule(d)

	if _, err := traffic_dns_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	existingRules, err := traffic_dns_rules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all forwarding dns rules: %v", err)
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

	_, err = traffic_dns_rules.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating forwarding dns rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "traffic_forwarding_dns_rule",
		func() (int, error) {
			allRules, err := traffic_dns_rules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := traffic_dns_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = traffic_dns_rules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceTrafficForwardingDNSRuleRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "traffic_forwarding_dns_rule")

	return nil
}

func resourceTrafficForwardingDNSRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("forwarding control rule ID not set: %v", id))
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := traffic_dns_rules.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving forwarding control rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedDNSRules(*rule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting forwarding dns rule ID: %v", id)
	if _, err := traffic_dns_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting forwarding dns rule %d: %v", id, err))
	}

	d.SetId("")
	log.Printf("[INFO] Forwarding control rule deleted")

	return nil
}

func expandForwardingDNSRule(d *schema.ResourceData) traffic_dns_rules.ECDNSRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandForwardingControlRule: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	// Process the DestCountries to add the prefix where needed
	rawDestCountries := SetToStringList(d, "dest_countries")
	processedDestCountries := make([]string, len(rawDestCountries))
	for i, country := range rawDestCountries {
		if country != "ANY" && country != "NONE" && len(country) == 2 { // Assuming the 2 letter code is an ISO Alpha-2 Code
			processedDestCountries[i] = "COUNTRY_" + country
		} else {
			processedDestCountries[i] = country
		}
	}

	result := traffic_dns_rules.ECDNSRules{
		ID:              id,
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Order:           order,
		Rank:            d.Get("rank").(int),
		State:           d.Get("state").(string),
		Action:          d.Get("action").(string),
		SrcIps:          SetToStringList(d, "src_ips"),
		DestAddresses:   SetToStringList(d, "dest_addresses"),
		Locations:       expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups: expandIDNameExtensionsSet(d, "location_groups"),
		SrcIpGroups:     expandIDNameExtensionsSet(d, "src_ip_groups"),
		DestIpGroups:    expandIDNameExtensionsSet(d, "dest_ip_groups"),
		ECGroups:        expandIDNameExtensionsSet(d, "ec_groups"),
		DNSGateway:      expandIDNameSet(d, "dns_gateway"),
	}
	return result
}
