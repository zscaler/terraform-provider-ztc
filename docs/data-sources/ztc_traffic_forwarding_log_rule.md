---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTC: traffic_forwarding_log_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-log-and-control-forwarding
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Get information about Log and Control Forwarding rules
---

# ztc_traffic_forwarding_log_rule (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/policy-management#/ecRules/ecRdr-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-log-and-control-forwarding)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-rules)

Use the **ztc_traffic_forwarding_log_rule** data source to get information about forwarding log rule available in the Zscaler Cloud and Branch Connector Portal.

**NOTE**: Resource available only via OneAPI

## Example Usage - Retrieve by Name

```hcl
data "ztc_traffic_forwarding_log_rule" "example" {
    name = "Rule01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_traffic_forwarding_log_rule" "example" {
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
* `type` - (String) The rule type.
* `forward_method` - (String) The type of traffic forwarding method selected from the available options
* `default_rule` - (String) Indicates whether the forwarding rule is a default rule
* `locations` - (List of Object) Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `ec_groups` - (List of Object) Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `proxy_gateway` - (List of Object) The proxy gateway for which the rule is applicable. 
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.