---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTC: traffic_forwarding_log_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-log-and-control-forwarding
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Creates and manages Log and Control Forwarding rules
---

# ztc_traffic_forwarding_log_rule (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/policy-management#/ecRules/ecRdr-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-log-and-control-forwarding)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-rules)

Use the **ztc_traffic_forwarding_log_rule** resource allows the creation and management of forwarding log rule available in the Zscaler Cloud and Branch Connector Portal.

**NOTE**: Resource available only via OneAPI

## Example Usage - DNS Forwarding Group - Action REDIR_REQ

```hcl
resource "ztc_zia_forwarding_gateway" "this" {
  name             = "ZTC_GW01"
  description      = "Example Forwarding Gateway 1"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ZIA"
}

data "ztc_location_management" "this" {
    name = "SJC_01"
}

resource "ztc_traffic_forwarding_log_rule" "this" {
  name           = "Log_Rule01"
  description    = "Log_Rule01"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  forward_method = "ECSELF"

  locations {
    id = [data.ztc_location_management.this.id]
  }
  proxy_gateway {
    id   = ztc_zia_forwarding_gateway.this.id
    name = tw_zia_forwarding_gateway.this.name
  }
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
* `forward_method` - (String) The type of traffic forwarding method selected from the available options. Supported value: `ECSELF`
* `locations` - (List of Object) Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `ec_groups` - (List of Object) Name-ID pairs of Cloud & Branch Connector Groups
  * `id` - (Number) Identifier that uniquely identifies an entity.
* `proxy_gateway` - (List of Object) The proxy gateway for which the rule is applicable. 
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_zia_forwarding_gateway** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import ztc_zia_forwarding_gateway.example <rule_id>
```

or

```shell
terraform import ztc_zia_forwarding_gateway.example <rule_name>
```