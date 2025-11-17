---
subcategory: "Forwarding Gateways"
layout: "zscaler"
page_title: "ZTC: dns_forwarding_gateway"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-gateways
  Creates and manages DNS Forwarding Gateways.
---

# ztc_dns_forwarding_gateway (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways#/gateways-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways)

Use the **ztc_dns_forwarding_gateway** resource allows the creation and management of DNS Forwarding gateways available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTC traffic forwarding rule.

## Example Usage - DNS Gateway Primary and Secondary DNS Address

```hcl
resource "ztc_dns_forwarding_gateway" "ztc_dns01" {
  name             = "DNS_FW_GW01"
  primary_ip       = "4.4.4.4"
  secondary_ip     = "8.8.8.8"
  failure_behavior = "FAIL_ALLOW_IGNORE_DNAT"
}
```

## Example Usage - DNS Gateway ECDNS Options Primary and Secondary

```hcl
resource "ztc_dns_forwarding_gateway" "ztc_dns02" {
  name                             = "DNS_FW_GW02"
  ec_dns_gateway_options_primary   = "WAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "WAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_ALLOW_IGNORE_DNAT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the forwarding gateway to be exported.
* `id` - (Optional) The ID of the forwarding gateway resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) A unique identifier assigned to the forwarding gateway.
* `name` - (String) The name of the Forwarding Gateway.

### Optional

* `primary_ip` - (String) IP address of the primary custom DNS server.
* `secondary_ip` - (String) IP address of the secondary custom DNS server.
* `ec_dns_gateway_options_primary` - (String) IP address of the primary LAN DNS Server. Supported Values: `LAN_PRI_DNS_AS_PRI`, and `LAN_SEC_DNS_AS_SEC`
* `ec_dns_gateway_options_secondary` - (String) IP address of the secondary LAN DNS Server. Supported Values: `LAN_PRI_DNS_AS_PRI`, and `LAN_SEC_DNS_AS_SEC`
* `failure_behavior` - (String) Choose what happens if the DNS server is unreachable. Supported Values: `FAIL_RET_ERR`, and `FAIL_ALLOW_IGNORE_DNAT`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_dns_forwarding_gateway** can be imported by using `<GATEWAY_ID>` or `<GATEWAY_NAME>` as the import ID.

For example:

```shell
terraform import ztc_dns_forwarding_gateway.example <rule_id>
```

or

```shell
terraform import ztc_dns_forwarding_gateway.example <rule_name>
```