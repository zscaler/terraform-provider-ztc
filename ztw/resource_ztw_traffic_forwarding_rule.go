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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policymanagement/forwardingrules"
)

var (
	forwardingControlLock          sync.Mutex
	forwardingControlStartingOrder int
)

func resourceForwardingControlRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingControlRuleCreate,
		ReadContext:   resourceForwardingControlRuleRead,
		UpdateContext: resourceForwardingControlRuleUpdate,
		DeleteContext: resourceForwardingControlRuleDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			forwardMethod := d.Get("forward_method").(string)

			// Function to check if an attribute is set and has values
			isSet := func(attr string) bool {
				val, ok := d.GetOk(attr)
				if !ok {
					return false
				}
				// Check if it's a set and has elements
				if set, ok := val.(*schema.Set); ok {
					return set.Len() > 0
				}
				return true
			}

			// If forward_method is ECZPA, certain attributes cannot be set
			if forwardMethod == "ECZPA" {
				prohibitedAttrs := []string{
					"dest_addresses",
					"dest_countries",
					"dest_ip_groups",
					"dest_ip_categories",
					"proxy_gateway",
					"nw_services",
					"nw_service_groups",
					"app_service_groups",
				}
				for _, attr := range prohibitedAttrs {
					if isSet(attr) {
						return fmt.Errorf("%s attribute cannot be set when forward_method is 'ECZPA'", attr)
					}
				}
			}

			// If forward_method is ZIA, DIRECT, LOCAL_SWITCH, or DROP, ZPA-related attributes cannot be set
			if forwardMethod == "ZIA" || forwardMethod == "DIRECT" || forwardMethod == "LOCAL_SWITCH" || forwardMethod == "DROP" {
				prohibitedAttrs := []string{
					"zpa_application_segments",
					"zpa_application_segment_groups",
				}
				for _, attr := range prohibitedAttrs {
					if isSet(attr) {
						return fmt.Errorf("%s attribute cannot be set when forward_method is '%s'", attr, forwardMethod)
					}
				}
			}

			return nil
		},
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
					resp, err := forwardingrules.GetRulesByName(ctx, service, id)
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
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule type selected from the available options",
				ValidateFunc: validation.StringInSlice([]string{
					"FIREWALL",
					"DNS",
					"DNAT",
					"SNAT",
					"FORWARDING",
					"INTRUSION_PREVENTION",
					"EC_DNS",
					"EC_RDR",
					"EC_SELF",
					"DNS_RESPONSE",
				}, false),
			},
			"forward_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of traffic forwarding method selected from the available options",
				ValidateFunc: validation.StringInSlice([]string{
					"DIRECT",
					"LOCAL_SWITCH",
					"ZIA",
					"ECZPA",
					"DROP",
				}, false),
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
			"src_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"source_ip_group_exclusion": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Source IP groups that must be excluded from the rule application",
			},
			"dest_addresses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.",
			},
			"dest_ip_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.",
			},
			"res_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination domain categories to which the rule applies",
			},
			"wan_selection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "WAN selection is only applicable when configuring a hardware device deployed in gateway mode.",
				ValidateFunc: validation.StringInSlice([]string{
					"SMRULEF_ZPA_BROKERS_RULE",
					"SMRULEF_APPC_DYNAMIC_SRC_IPGROUP",
					"SMRULEF_EXCL_SRC_IP",
					"BALANCED_RULE",
					"BESTLINK_RULE",
				}, false),
			},
			"locations":                      setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations."),
			"location_groups":                setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the forwarding rule applies"),
			"ec_groups":                      setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies"),
			"src_ip_groups":                  setIDsSchemaTypeCustom(nil, "Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group"),
			"dest_ip_groups":                 setIDsSchemaTypeCustom(nil, "User-defined destination IP address groups to which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group"),
			"nw_services":                    setIDsSchemaTypeCustom(intPtr(1024), "User-defined network services to which the rule applies. If not set, the rule is not restricted to a specific network service."),
			"nw_service_groups":              setIDsSchemaTypeCustom(nil, "User-defined network service group to which the rule applies. If not set, the rule is not restricted to a specific network service group."),
			"app_service_groups":             setIDsSchemaTypeCustom(nil, "list of application service groups"),
			"proxy_gateway":                  setIdNameSchemaCustom(1, "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method."),
			"zpa_application_segments":       setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"zpa_application_segment_groups": setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"dest_countries":                 getISOCountryCodes(),
			"src_workload_groups":            setIDsSchemaTypeCustom(nil, "The list of preconfigured workload groups to which the policy must be applied"),
		},
	}
}

func validatePredefinedRules(req forwardingrules.ForwardingRules) error {
	if req.Name == "ZPA Forwarding Rule" || req.Name == "Direct rule for Zscaler Cloud Endpoints" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	if req.Name == "Direct rule for WAN Destinations Group" || req.Name == "Direct rule for LAN Destinations Group" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	if req.Name == "Client Connector to ZPA" || req.Name == "ZPA Pool For Stray Traffic" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	return nil
}

func resourceForwardingControlRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandForwardingControlRule(d)
	log.Printf("[INFO] Creating zia forwarding control rule\n%+v\n", req)

	start := time.Now()

	forwardingControlLock.Lock()
	if forwardingControlStartingOrder == 0 {
		list, _ := forwardingrules.GetAll(ctx, service)
		for _, r := range list {
			if r.Order > forwardingControlStartingOrder {
				forwardingControlStartingOrder = r.Order
			}
		}
		if forwardingControlStartingOrder == 0 {
			forwardingControlStartingOrder = 1
		} else {
			forwardingControlStartingOrder++
		}
	}
	forwardingControlLock.Unlock()
	startWithoutLocking := time.Now()

	// Store the intended order from HCL
	intendedOrder := req.Order
	intendedRank := req.Rank
	if intendedRank < 7 {
		// always start rank 7 rules at the next available order after all ranked rules
		req.Rank = 7
	}
	req.Order = forwardingControlStartingOrder
	resp, err := forwardingrules.Create(ctx, service, &req)

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

	log.Printf("[INFO] Created zia forwarding control rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
	// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
	resourceType := "forwarding_control_rule"

	reorderWithBeforeReorder(
		OrderRule{Order: intendedOrder, Rank: intendedRank},
		resp.ID,
		resourceType,
		func() (int, error) {
			allRules, err := forwardingrules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			// Custom updateOrder that handles predefined rules
			rule, err := forwardingrules.Get(ctx, service, id)
			if err != nil {
				return err
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = forwardingrules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	markOrderRuleAsDone(resp.ID, resourceType)

	return resourceForwardingControlRuleRead(ctx, d, meta)
}

func resourceForwardingControlRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia forwarding control rule id is set"))
	}
	resp, err := forwardingrules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing forwarding control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	log.Printf("[INFO] Getting forwarding control rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("forward_method", resp.ForwardMethod)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("state", resp.State)
	_ = d.Set("type", resp.Type)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("res_categories", resp.ResCategories)
	_ = d.Set("wan_selection", resp.WanSelection)
	_ = d.Set("source_ip_group_exclusion", resp.SourceIpGroupExclusion)
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

	if err := d.Set("nw_services", flattenIDExtensionsListIDs(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_service_groups", flattenIDExtensionsListIDs(resp.NwServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_service_groups", flattenIDExtensionsListIDs(resp.AppServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_application_segments", flattenZPAApplicationSegmentsSimple(resp.ZPAApplicationSegments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_application_segment_groups", flattenZPAApplicationSegmentGroupsSimple(resp.ZPAApplicationSegmentGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_workload_groups", flattenIDExtensionsListIDs(resp.SrcWorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting src_workload_groups: %s", err))
	}
	return nil
}

func resourceForwardingControlRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding control rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("forwarding control rule ID not set"))
	}
	log.Printf("[INFO] Updating forwarding control rule ID: %v\n", id)
	req := expandForwardingControlRule(d)

	if _, err := forwardingrules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	existingRules, err := forwardingrules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all forwarding rules: %v", err)
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

	_, err = forwardingrules.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating forwarding control rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "forwarding_control_rule",
		func() (int, error) {
			allRules, err := forwardingrules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := forwardingrules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = forwardingrules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceForwardingControlRuleRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "forwarding_control_rule")

	return nil
}

func resourceForwardingControlRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("forwarding control rule ID not set: %v", id))
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := forwardingrules.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving forwarding control rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedRules(*rule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting forwarding control rule ID: %v", id)
	if _, err := forwardingrules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting forwarding control rule %d: %v", id, err))
	}

	d.SetId("")
	log.Printf("[INFO] Forwarding control rule deleted")

	return nil
}

func expandForwardingControlRule(d *schema.ResourceData) forwardingrules.ForwardingRules {
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

	result := forwardingrules.ForwardingRules{
		ID:                          id,
		Name:                        d.Get("name").(string),
		Description:                 d.Get("description").(string),
		Order:                       order,
		Rank:                        d.Get("rank").(int),
		Type:                        d.Get("type").(string),
		State:                       d.Get("state").(string),
		ForwardMethod:               d.Get("forward_method").(string),
		WanSelection:                d.Get("wan_selection").(string),
		SourceIpGroupExclusion:      d.Get("source_ip_group_exclusion").(bool),
		ResCategories:               SetToStringList(d, "res_categories"),
		SrcIps:                      SetToStringList(d, "src_ips"),
		DestAddresses:               SetToStringList(d, "dest_addresses"),
		DestIpCategories:            SetToStringList(d, "dest_ip_categories"),
		DestCountries:               processedDestCountries,
		Locations:                   expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:             expandIDNameExtensionsSet(d, "location_groups"),
		SrcIpGroups:                 expandIDNameExtensionsSet(d, "src_ip_groups"),
		DestIpGroups:                expandIDNameExtensionsSet(d, "dest_ip_groups"),
		NwServices:                  expandIDNameExtensionsSet(d, "nw_services"),
		AppServiceGroups:            expandIDNameExtensionsSet(d, "app_service_groups"),
		NwServiceGroups:             expandIDNameExtensionsSet(d, "nw_service_groups"),
		NwApplicationGroups:         expandIDNameExtensionsSet(d, "nw_application_groups"),
		ZPAApplicationSegments:      expandZPAApplicationSegmentSet(d, "zpa_application_segments"),
		ZPAApplicationSegmentGroups: expandZPAApplicationSegmentGroupSet(d, "zpa_application_segment_groups"),
		ECGroups:                    expandIDNameExtensionsSet(d, "ec_groups"),
		SrcWorkloadGroups:           expandIDNameExtensionsSet(d, "src_workload_groups"),
		ProxyGateway:                expandIDNameSet(d, "proxy_gateway"),
	}
	return result
}

func currentRuleOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := forwardingrules.GetAll(ctx, service)
	if err != nil {
		return ""
	}
	result := ""
	for i, r := range list {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("Rank %d VS Order %d", r.Rank, r.Order)

	}
	return result
}

func flattenZPAApplicationSegmentsSimple(list []common.ZPAApplicationSegments) []interface{} {
	if len(list) == 0 {
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

func flattenZPAApplicationSegmentGroupsSimple(list []common.ZPAApplicationSegmentGroups) []interface{} {
	if len(list) == 0 {
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}
