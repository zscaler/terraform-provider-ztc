---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTC: traffic_forwarding_dns_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-dns-policies
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Get information about Traffic Forwarding DNS rules
---

# ztc_traffic_forwarding_dns_rule (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-dns-policies)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/about-dns-policies)

Use the **ztc_traffic_forwarding_dns_rule** data source to get information about forwarding dns rule available in the Zscaler Cloud and Branch Connector Portal.

## Example Usage - Retrieve by Name

```hcl
data "ztc_traffic_forwarding_dns_rule" "example" {
    name = "example_forwarding_gateway"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_traffic_forwarding_dns_rule" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the dns forwarding gateway to be exported.
* `id` - (Optional) The ID of the dns forwarding gateway resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) A unique identifier assigned to the dns forwarding gateway.
* `name` - (String) The name of the dns Forwarding Gateway.

### Optional

* `description` - (String) Additional information about the forwarding rule.
* `state` - (String) Indicates whether the forwarding rule is enabled or disabled.
* `order` - (Number) The order of execution for the forwarding rule order.
* `rank` - (Number) Admin rank assigned to the forwarding rule.
* `type` - (String) The rule type. Default value `EC_DNS`
* `action` - (String) The rule type. Returned values: `ALLOW`, `BLOCK`, `REDIR_REQ`, `REDIR_ZPA`, 
* `src_ips` - (List of String) User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
* `dest_addresses` - (List of String) List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses.
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
* `dns_gateway` - (List of Object) The dns gateway for which the rule is applicable. Applicable only when action is `REDIR_REQ`
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.
* `zpa_ip_group` - (List of Object) The ip pool group for which the rule is applicable. Applicable only when action is `REDIR_ZPA`
  * `id` - (Number) The ID of the IP pool group resource.
  * `name` - (String) The name of the IP pool group resource.