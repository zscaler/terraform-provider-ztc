---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTC: traffic_forwarding_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-forwarding-rules
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Get information about Traffic Forwarding Rules.
---

# ztc_traffic_forwarding_rule (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-forwarding-rules)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-rules)

Use the **ztc_traffic_forwarding_rule** data source to get information about traffic forwarding rules available in the Zscaler Cloud and Branch Connector Portal.

## Example Usage - Retrieve by Name

```hcl
data "ztc_traffic_forwarding_rule" "example" {
    name = "example_forwarding_rule"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_traffic_forwarding_rule" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the forwarding rule to be exported.
* `id` - (Optional) The ID of the forwarding rule resource.
* `type` - (Optional) The rule type selected from the available options.
* `access_control` - (Optional) Access permission available for the current user to the rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) A unique identifier assigned to the forwarding rule.
* `name` - (String) The name of the forwarding rule.

### Optional

* `description` - (String) Additional information about the forwarding rule.
* `forward_method` - (String) The type of traffic forwarding method selected from the available options (e.g., DIRECT, ZIA, ECZPA, DROP, LOCAL_SWITCH).
* `state` - (String) Indicates whether the forwarding rule is enabled or disabled.
* `order` - (Number) The order of execution for the forwarding rule order.
* `rank` - (Number) Admin rank assigned to the forwarding rule.
* `type` - (String) The rule type (e.g., FIREWALL, DNS, DNAT, SNAT, FORWARDING, INTRUSION_PREVENTION, EC_DNS, EC_RDR, EC_SELF, DNS_RESPONSE).
* `src_ips` - (List of String) User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
* `dest_addresses` - (List of String) List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses.
* `dest_ip_categories` - (List of String) List of destination IP categories to which the rule applies.
* `dest_countries` - (List of String) Destination countries for which the rule is applicable.
* `res_categories` - (List of String) List of destination domain categories to which the rule applies.
* `wan_selection` - (String) WAN selection (only applicable when configuring a hardware device deployed in gateway mode).
* `source_ip_group_exclusion` - (Boolean) Source IP groups that must be excluded from the rule application.
* `ec_groups` - (List of Object) Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
* `locations` - (List of Object) Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `location_groups` - (List of Object) Name-ID pairs of the location groups to which the forwarding rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `src_ip_groups` - (List of Object) Source IP address groups for which the rule is applicable.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `dest_ip_groups` - (List of Object) User-defined destination IP address groups to which the rule is applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `nw_services` - (List of Object) User-defined network services to which the rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `nw_service_groups` - (List of Object) User-defined network service groups to which the rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `app_service_groups` - (List of Object) List of application service groups.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `zpa_application_segments` - (List of Object) List of ZPA Application Segments for which this rule is applicable (used for ECZPA forwarding method).
  * `id` - (Number) Application segment identifier.
* `zpa_application_segment_groups` - (List of Object) List of ZPA Application Segment Groups for which this rule is applicable (used for ECZPA forwarding method).
  * `id` - (Number) Application segment group identifier.
* `proxy_gateway` - (List of Object) The proxy gateway for which the rule is applicable (for Proxy Chaining forwarding method).
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.
* `zpa_gateway` - (List of Object) The ZPA Server Group for which this rule is applicable (for ZPA forwarding method).
  * `id` - (Number) Server group identifier.
  * `name` - (String) Server group name.
* `src_workload_groups` - (List of Object) The list of preconfigured workload groups to which the policy must be applied.
  * `id` - (Number) Workload group identifier.
  * `name` - (String) Workload group name.
  * `extensions` - (Map of String) Extensions field.