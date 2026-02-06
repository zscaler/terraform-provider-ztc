---
subcategory: "Policy Management"
layout: "zscaler"
page_title: "ZTC: traffic_forwarding_dns_rule"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-dns-policies
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-rules
  Creates and manages Traffic Forwarding DNS rules
---

# ztc_traffic_forwarding_dns_rule (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/policy-management#/ecRules/ecRdr-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-dns-policies)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/about-dns-policies)

Use the **ztc_traffic_forwarding_dns_rule** resource allows the creation and management of traffic forwarding DNS rule in the Zscaler Cloud and Branch Connector Portal.

**NOTE**: Resource available only via OneAPI

## Example Usage - DNS Forwarding Group - Action REDIR_REQ

```hcl
resource "ztc_dns_forwarding_gateway" "ztc_dns01" {
  name                             = "DNS_FW_GW01"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}

resource "ztc_ip_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
}

resource "zia_ip_source_groups" "this" {
  name         = "example1"
  description  = "example1"
  ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}

data "ztc_location_management" "this" {
    name = "SJC_01"
}

resource "ztc_traffic_forwarding_dns_rule" "this" {
  name              = "DNS_Rule01"
  description       = "DNS_Rule01"
  order             = 1
  rank              = 7
  state             = "ENABLED"
  action            = "REDIR_REQ"
  src_ips           = ["192.168.200.200"]
  dest_addresses    = ["server1.acme.com"]
  src_ip_groups {
    id = [zia_ip_source_groups.this.id]
  }
  dest_ip_groups {
    id = [ztc_ip_destination_groups.dstn_fqdn.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  dns_gateway {
    id   = ztc_dns_forwarding_gateway.this.id
    name = ztc_dns_forwarding_gateway.this.name
  }
}
```

## Example Usage - DNS Forwarding Group - Action REDIR_ZPA

```hcl
resource "ztc_dns_forwarding_gateway" "ztc_dns01" {
  name                             = "DNS_FW_GW01"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}

resource "ztc_ip_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses   = ["test1.acme.com", "test2.acme.com", "test3.acme.com"]
}

resource "zia_ip_source_groups" "this" {
  name         = "example1"
  description  = "example1"
  ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}

resource "ztc_ip_pool_groups" "this" {
  name        = "Example IP Pool"
  description = "Updated Example IP Pool for testing"
  ip_addresses = [
    "10.0.0.0/24"
  ]
}

data "ztc_location_management" "this" {
    name = "SJC_01"
}

resource "ztc_traffic_forwarding_dns_rule" "this" {
  name              = "DNS_Rule01"
  description       = "DNS_Rule01"
  order             = 1
  rank              = 7
  state             = "ENABLED"
  action            = "REDIR_ZPA"
  src_ips           = ["192.168.200.200"]
  dest_addresses    = ["server1.acme.com"]
  src_ip_groups {
    id = [zia_ip_source_groups.this.id]
  }
  locations {
    id = [data.ztc_location_management.this.id]
  }
  zpa_ip_group {
    id   = ztc_ip_pool_groups.this.id
    name = ztc_ip_pool_groups.this.name
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
* `action` - (String) The rule type. Supported values: `ALLOW`, `BLOCK`, `REDIR_REQ`, `REDIR_ZPA`, 
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
* `dest_ip_groups` - (List of Object) User-defined destination IP address groups to which the rule is applied. Not supported when action is `REDIR_ZPA`
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Extensions field.
* `dns_gateway` - (List of Object) The dns gateway for which the rule is applicable. Applicable only when action is `REDIR_REQ`
  * `id` - (Number) Gateway identifier.
  * `name` - (String) Gateway name.
* `zpa_ip_group` - (List of Object) The ip pool group for which the rule is applicable. Applicable only when action is `REDIR_ZPA`
  * `id` - (Number) The ID of the IP pool group resource.
  * `name` - (String) The name of the IP pool group resource.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_traffic_forwarding_dns_rule** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import ztc_traffic_forwarding_dns_rule.example <rule_id>
```

or

```shell
terraform import ztc_traffic_forwarding_dns_rule.example <rule_name>
```