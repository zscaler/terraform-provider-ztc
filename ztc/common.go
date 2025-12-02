package ztc

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

func sortOrders(ruleOrderMap map[int]orderWithState) RuleIDOrderPairList {
	pl := make(RuleIDOrderPairList, len(ruleOrderMap))
	i := 0
	for k, v := range ruleOrderMap {
		pl[i] = RuleIDOrderPair{k, v.order}
		i++
	}
	sort.Sort(pl)
	return pl
}

type orderWithState struct {
	order OrderRule
	done  bool
}

type listrules struct {
	orders  map[string]map[int]orderWithState
	orderer map[string]int
	sync.Mutex
}

var rules = listrules{
	orders: make(map[string]map[int]orderWithState),
}

type RuleIDOrderPair struct {
	ID    int
	Order OrderRule
}

type RuleIDOrderPairList []RuleIDOrderPair

func (p RuleIDOrderPairList) Len() int { return len(p) }
func (p RuleIDOrderPairList) Less(i, j int) bool {
	if p[i].Order == p[j].Order {
		return p[i].ID < p[j].ID
	}
	return p[i].Order.Rank < p[j].Order.Rank || p[i].Order.Rank == p[j].Order.Rank && p[i].Order.Order < p[j].Order.Order
}
func (p RuleIDOrderPairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func reorderAll(resourceType string, getCount func() (int, error), updateOrder func(id int, order OrderRule) error, beforeReorder func()) {
	ticker := time.NewTicker(time.Second * 30) // create a ticker that ticks every 30 seconds
	defer ticker.Stop()                        // stop the ticker when the loop ends
	numResources := []int{0, 0, 0}
	for {
		select {
		case <-ticker.C:
			rules.Lock()
			size := len(rules.orders[resourceType])
			done := true
			// first check if all rules creation is done
			for _, v := range rules.orders[resourceType] {
				if !v.done {
					done = false
				}
			}
			numResources[0], numResources[1], numResources[2] = numResources[1], numResources[2], size
			if done && numResources[0] == numResources[1] && numResources[1] == numResources[2] {
				// No changes after a while (4 runs), so reorder, and return
				count, _ := getCount()
				// sort by order (ascending)
				sorted := sortOrders(rules.orders[resourceType])
				log.Printf("[INFO] sorting filtering rule after tick; sorted:%v", sorted)
				if beforeReorder != nil {
					beforeReorder()
				}
				for _, v := range sorted {
					if v.Order.Order <= count {
						if err := updateOrder(v.ID, v.Order); err != nil {
							log.Printf("[ERROR] couldn't reorder the rule after tick, the order may not have taken place: %v\n", err)
						}
					}
				}
				rules.Unlock()
				return
			}
			rules.Unlock()
		default:
			time.Sleep(time.Second * 15)
		}
	}
}

func markOrderRuleAsDone(id int, resourceType string) {
	rules.Lock()
	r := rules.orders[resourceType][id]
	r.done = true
	rules.orders[resourceType][id] = r
	rules.Unlock()
}

type OrderRule struct {
	Order int
	Rank  int
}

func reorderWithBeforeReorder(order OrderRule, id int, resourceType string, getCount func() (int, error), updateOrder func(id int, order OrderRule) error, beforeReorder func()) {
	rules.Lock()
	shouldCallReorder := false
	if len(rules.orders) == 0 {
		rules.orders = map[string]map[int]orderWithState{}
		rules.orderer = map[string]int{}
	}
	if _, ok := rules.orderer[resourceType]; ok {
		shouldCallReorder = false
	} else {
		rules.orderer[resourceType] = id
		shouldCallReorder = true
	}
	if len(rules.orders[resourceType]) == 0 {
		rules.orders[resourceType] = map[int]orderWithState{}
	}
	rules.orders[resourceType][id] = orderWithState{order, shouldCallReorder}
	rules.Unlock()
	if shouldCallReorder {
		log.Printf("[INFO] starting to reorder the rules, delegating to rule:%d, order:%d", id, order)
		// one resource will wait until all resources are done and reorder then return
		reorderAll(resourceType, getCount, updateOrder, beforeReorder)
	}
}

func reorder(order OrderRule, id int, resourceType string, getCount func() (int, error), updateOrder func(id int, order OrderRule) error) {
	reorderWithBeforeReorder(order, id, resourceType, getCount, updateOrder, nil)
}

func flattenCommonIDNameExternalID(gp *common.CommonIDNameExternalID) []map[string]interface{} {
	if gp == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":               gp.ID,
			"name":             gp.Name,
			"is_name_l10n_tag": gp.IsNameL10nTag,
			"extensions":       gp.Extensions,
			"deleted":          gp.Deleted,
			"external_id":      gp.ExternalID,
			"association_time": gp.AssociationTime,
		},
	}
}

func flattenGeneralPurpose(gp *common.CommonIDNameExternalID) []map[string]interface{} {
	return flattenCommonIDNameExternalID(gp)
}

func flattenIDNameExtensions(list []common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":   val.ID,
			"name": val.Name,
		}
		if val.Extensions != nil {
			r["extensions"] = val.Extensions
		}
		flattenedList[i] = r
	}
	return flattenedList
}

func flattenListCommonIDNameExternalID(gp []common.CommonIDNameExternalID) []map[string]interface{} {
	if gp == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(gp))
	for _, item := range gp {
		result = append(result, map[string]interface{}{
			"id":               item.ID,
			"name":             item.Name,
			"is_name_l10n_tag": item.IsNameL10nTag,
			"extensions":       item.Extensions,
			"deleted":          item.Deleted,
			"external_id":      item.ExternalID,
			"association_time": item.AssociationTime,
		})
	}
	return result
}

func IdSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed: true,
		// ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					// Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					// Computed: true,
				},
			},
		},
	}
}

func ListIdsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

func UIDNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
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
				"is_name_l10n_tag": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
				},
				"extensions": {
					Type:        schema.TypeMap,
					Computed:    true,
					Description: "General purpose",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"deleted": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "General purpose",
				},
				"external_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "General purpose",
				},
				"association_time": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "General purpose",
				},
			},
		},
	}
}

func UIDNameSchemaLite() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenDNS(dns *common.DNS) []map[string]interface{} {
	if dns == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":       dns.ID,
			"ips":      dns.IPs,
			"dns_type": dns.DNSType,
		},
	}
}

func flattenNW(nw *common.ManagementNw) []map[string]interface{} {
	if nw == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":              nw.ID,
			"ip_start":        nw.IPStart,
			"ip_end":          nw.IPEnd,
			"netmask":         nw.Netmask,
			"default_gateway": nw.DefaultGateway,
			"nw_type":         nw.NWType,
			"dns":             flattenDNS(nw.DNS),
		},
	}
}

func NWSchemaResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ip_start": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"ip_end": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"netmask": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"default_gateway": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"nw_type": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"dns": {
					Type:     schema.TypeSet,
					Computed: true,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"ips": {
								Type:     schema.TypeSet,
								Computed: true,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"dns_type": {
								Type:     schema.TypeString,
								Computed: true,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func NWSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ip_start": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ip_end": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"netmask": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"default_gateway": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"nw_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"dns": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"ips": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"dns_type": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func VPNCredentialsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"common_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"fqdn": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ip_address": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"pre_shared_key": {
					Type:      schema.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"xauth_password": {
					Type:      schema.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"comments": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"disabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"psk": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"location":   UIDNameSchema(),
				"managed_by": UIDNameSchema(),
			},
		},
	}
}

func ecGroupSchemaData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deploy_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"platform": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aws_availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"azure_availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"location": UIDNameSchema(),
		"max_ec_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"prov_template": UIDNameSchema(),
		"tunnel_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ec_vms": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"form_factor": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"management_nw": NWSchema(),
					"ec_instances": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ec_instance_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"out_gw_ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"nat_ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"dns_ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"service_nw": NWSchema(),
								"virtual_nw": NWSchema(),
							},
						},
					},
					"city_geo_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"nat_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"zia_gateway": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"zpa_broker": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"build_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_upgrade_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"upgrade_status": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"upgrade_start_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"upgrade_end_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func MergeSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	final := make(map[string]*schema.Schema)
	for _, s := range schemas {
		for k, v := range s {
			final[k] = v
		}
	}
	return final
}

func getLocationManagementCountries() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateLocationManagementCountries(),
		Description:  "Supported Countries",
		Optional:     true,
		Computed:     true,
	}
}

func getLocationManagementTimeZones() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateLocationManagementTimeZones(),
		Description:  "Timezone of the location. If not specified, it defaults to GMT.",
		Optional:     true,
		Computed:     true,
	}
}

func setIDsSchemaTypeCustom(maxItems *int, desc string) *schema.Schema {
	ids := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
	if maxItems != nil && *maxItems > 0 {
		ids.MaxItems = *maxItems
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": ids,
			},
		},
	}
}

func setIdNameSchemaCustom(maxItems int, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: description,
		MaxItems:    maxItems,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The unique identifier for the resource.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of the resource.",
				},
			},
		},
	}
}

func setIDSchemaCustom(maxItems int, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: description,
		MaxItems:    maxItems,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The unique identifier for the resource.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of the resource.",
				},
			},
		},
	}
}

func setExtIDNameSchemaCustom(maxItems *int, description string) *schema.Schema {
	schema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the application segment.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "External ID of the application segment.",
				},
			},
		},
	}

	if maxItems != nil && *maxItems > 0 {
		schema.MaxItems = *maxItems
	}

	return schema
}

func getISOCountryCodes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateISOCountryCodes,
		},
		Optional: true,
		Computed: true,
	}
}

func flattenIDExtensionsListIDs(list []common.IDNameExtensions) []interface{} {
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

func flattenIDNameSet(idName *common.CommonIDName) []interface{} {
	idNameSet := make([]interface{}, 0)
	if idName != nil {
		idNameSet = append(idNameSet, map[string]interface{}{
			"id":   idName.ID,
			"name": idName.Name,
		})
	}
	return idNameSet
}

func flattenIDSet(idName *common.CommonIDName) []interface{} {
	idNameSet := make([]interface{}, 0)
	if idName != nil {
		idNameSet = append(idNameSet, map[string]interface{}{
			"id": idName.ID,
		})
	}
	return idNameSet
}

func flattenIDName(workloadGroups []common.CommonIDName) []interface{} {
	if workloadGroups == nil {
		return nil
	}

	wgList := make([]interface{}, len(workloadGroups))
	for i, wg := range workloadGroups {
		wgMap := make(map[string]interface{})
		wgMap["id"] = wg.ID
		wgMap["name"] = wg.Name
		wgList[i] = wgMap
	}

	return wgList
}

func expandIDNameExtensionsSet(d *schema.ResourceData, key string) []common.IDNameExtensions {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDNameExtensions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
					result = append(result, common.IDNameExtensions{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.IDNameExtensions{}
}

func expandZPAApplicationSegmentSet(d *schema.ResourceData, key string) []common.ZPAApplicationSegments {
	setInterface, ok := d.GetOk(key)
	if !ok {
		return []common.ZPAApplicationSegments{}
	}

	set := setInterface.(*schema.Set)
	var result []common.ZPAApplicationSegments
	for _, item := range set.List() {
		itemMap, _ := item.(map[string]interface{})
		if itemMap != nil && itemMap["id"] != nil {
			idSet := itemMap["id"].(*schema.Set)
			for _, id := range idSet.List() {
				result = append(result, common.ZPAApplicationSegments{
					ID: id.(int),
				})
			}
		}
	}
	return result
}

func expandZPAApplicationSegmentGroupSet(d *schema.ResourceData, key string) []common.ZPAApplicationSegmentGroups {
	setInterface, ok := d.GetOk(key)
	if !ok {
		return []common.ZPAApplicationSegmentGroups{}
	}

	set := setInterface.(*schema.Set)
	var result []common.ZPAApplicationSegmentGroups
	for _, item := range set.List() {
		itemMap, _ := item.(map[string]interface{})
		if itemMap != nil && itemMap["id"] != nil {
			idSet := itemMap["id"].(*schema.Set)
			for _, id := range idSet.List() {
				result = append(result, common.ZPAApplicationSegmentGroups{
					ID: id.(int),
				})
			}
		}
	}
	return result
}

// Common expand function to support Workload Groups across other resources
func expandWorkloadGroupsIDName(d *schema.ResourceData, key string) []common.CommonIDName {
	// Retrieve the set from the resource data
	if v, ok := d.GetOk(key); ok {
		workloadGroupsSet := v.(*schema.Set)
		// Initialize the slice to hold the expanded workload groups
		workloadGroups := make([]common.CommonIDName, 0, workloadGroupsSet.Len())

		// Iterate over the set and construct the slice of common.IDName
		for _, wgMapInterface := range workloadGroupsSet.List() {
			wgMap := wgMapInterface.(map[string]interface{})
			wg := common.CommonIDName{
				ID:   wgMap["id"].(int),
				Name: wgMap["name"].(string),
			}
			workloadGroups = append(workloadGroups, wg)
		}

		return workloadGroups
	}

	// Return an empty slice if the key is not set
	return []common.CommonIDName{}
}

// expandIDNameSet takes a Terraform set as input and returns a pointer to a common.IDName struct.
func expandIDNameSet(d *schema.ResourceData, key string) *common.CommonIDName {
	idNameList, ok := d.Get(key).(*schema.Set)
	if !ok || idNameList.Len() == 0 {
		return nil
	}

	// Assuming each set can only have one item as per your JSON structure.
	// If it can have multiple, this needs to be adjusted accordingly.
	for _, v := range idNameList.List() {
		item := v.(map[string]interface{})
		return &common.CommonIDName{
			ID:   item["id"].(int),
			Name: item["name"].(string),
		}
	}

	return nil
}

func expandIDSet(d *schema.ResourceData, key string) *common.CommonIDName {
	idNameList, ok := d.Get(key).(*schema.Set)
	if !ok || idNameList.Len() == 0 {
		return nil
	}

	// Assuming each set can only have one item as per your JSON structure.
	// If it can have multiple, this needs to be adjusted accordingly.
	for _, v := range idNameList.List() {
		item := v.(map[string]interface{})
		return &common.CommonIDName{
			ID: item["id"].(int),
		}
	}

	return nil
}

func currentOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := forwarding_rules.GetAll(ctx, service)
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

func flattenIDExtensionsList(idNameExtension *common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, 0)
	if idNameExtension != nil && (idNameExtension.ID != 0 || idNameExtension.Name != "") {
		flattenedList = append(flattenedList, map[string]interface{}{
			"id":         idNameExtension.ID,
			"name":       idNameExtension.Name,
			"extensions": idNameExtension.Extensions,
		})
	}
	return flattenedList
}

func flattenCustomIDNameSet(customID *common.CommonIDNameExternalID) []interface{} {
	if customID == nil || customID.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":   customID.ID,
			"name": customID.Name,
		},
	}
}

func flattenCommonIDNameExternalIDSet(list []common.CommonIDNameExternalID) []interface{} {
	var flattenedList []interface{}
	for _, item := range list {
		m := map[string]interface{}{
			"id":               item.ID,
			"name":             item.Name,
			"is_name_l10n_tag": item.IsNameL10nTag,
			"extensions":       item.Extensions,
			"deleted":          item.Deleted,
			"external_id":      item.ExternalID,
			"association_time": item.AssociationTime,
		}
		flattenedList = append(flattenedList, m)
	}
	return flattenedList
}

func expandCommonIDNameExternalIDSet(l interface{}) []common.CommonIDNameExternalID {
	if l == nil {
		return nil
	}
	setObj := l.(*schema.Set)
	var list []common.CommonIDNameExternalID
	for _, itemObj := range setObj.List() {
		uuidNameObj, _ := itemObj.(map[string]interface{})
		if uuidNameObj != nil {
			if idObj, ok := uuidNameObj["id"]; ok {
				id, idIsOk := idObj.(int)
				if idIsOk {
					item := common.CommonIDNameExternalID{
						ID: id,
					}
					if name, ok := uuidNameObj["name"].(string); ok {
						item.Name = name
					}
					if isNameL10nTag, ok := uuidNameObj["is_name_l10n_tag"].(bool); ok {
						item.IsNameL10nTag = isNameL10nTag
					}
					if extensions, ok := uuidNameObj["extensions"].(map[string]interface{}); ok {
						item.Extensions = extensions
					}
					if deleted, ok := uuidNameObj["deleted"].(bool); ok {
						item.Deleted = deleted
					}
					if externalID, ok := uuidNameObj["external_id"].(string); ok {
						item.ExternalID = externalID
					}
					if associationTime, ok := uuidNameObj["association_time"].(int); ok {
						item.AssociationTime = associationTime
					}
					list = append(list, item)
				}
			}
		}
	}
	return list
}

// For ListIdsSchema() - extracts only IDs and returns format: [{"id": [1,2,3]}]
func flattenCommonIDNameExternalIDToListIds(list []common.CommonIDNameExternalID) []map[string]interface{} {
	if len(list) == 0 {
		return nil
	}
	result := []map[string]interface{}{}
	ids := []int{}
	for _, item := range list {
		ids = append(ids, item.ID)
	}
	result = append(result, map[string]interface{}{
		"id": ids,
	})
	return result
}

// For ListIdsSchema() - expands from format: [{"id": [1,2,3]}] to []CommonIDNameExternalID
func expandListIdsToCommonIDNameExternalID(l interface{}) []common.CommonIDNameExternalID {
	if l == nil {
		return nil
	}
	setObj := l.(*schema.Set)
	var list []common.CommonIDNameExternalID
	for _, itemObj := range setObj.List() {
		uuidNameObj, _ := itemObj.(map[string]interface{})
		if uuidNameObj != nil {
			if idsObj, ok := uuidNameObj["id"]; ok {
				if idsList, ok := idsObj.([]interface{}); ok {
					for _, idObj := range idsList {
						if id, ok := idObj.(int); ok {
							list = append(list, common.CommonIDNameExternalID{
								ID: id,
							})
						}
					}
				}
			}
		}
	}
	return list
}

func getCloudFirewallNwServicesTag() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		ValidateDiagFunc: validateCloudFirewallNwServicesTag(),
		Optional:         true,
		Computed:         true,
	}
}

func resourceNetworkPortsSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"end": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
			},
		},
	}
}

func dataNetworkPortsSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Start of port range",
				},
				"end": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "End of port range",
				},
			},
		},
	}
}

func flattenNetwordPorts(ports []networkservices.NetworkPorts) []interface{} {
	portsObj := make([]interface{}, len(ports))
	for i, val := range ports {
		portsObj[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}
	return portsObj
}

func expandNetworkPorts(d *schema.ResourceData, key string) []networkservices.NetworkPorts {
	var ports []networkservices.NetworkPorts
	if portsInterface, ok := d.GetOk(key); ok {
		portSet, ok := portsInterface.(*schema.Set)
		if !ok {
			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
			return ports
		}
		ports = make([]networkservices.NetworkPorts, len(portSet.List()))
		for i, val := range portSet.List() {
			portItem := val.(map[string]interface{})
			ports[i] = networkservices.NetworkPorts{
				Start: portItem["start"].(int),
				End:   portItem["end"].(int),
			}
		}
	}
	return ports
}
