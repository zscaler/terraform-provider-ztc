---
subcategory: "DNS Gateway"
layout: "zscaler"
page_title: "ZTC: ztc_dns_gateway"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway
  API documentation https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/dns-gateway/ec-dns-gateway-z-resource-add-dns-gateway
  Get information about DNS Forwarding Gateways.
---

# ztc_dns_gateway (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways#/gateways-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway)
* [API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/dns-gateway/ec-dns-gateway-z-resource-add-dns-gateway)

The **ztc_dns_gateway** resource allows you to create and manage DNS Gateway configurations in the Zscaler Zero Trust Cloud (ZTC) platform.

## Example Usage - EC DNS Gateway

```hcl
resource "ztc_dns_gateway" "example" {
  name                             = "Example DNS Gateway"
  ec_dns_gateway_options_primary   = "LAN_PRI_DNS_AS_PRI"
  ec_dns_gateway_options_secondary = "LAN_SEC_DNS_AS_SEC"
  failure_behavior                 = "FAIL_RET_ERR"
}
```

## Example Usage - Custom DNS Server

```hcl
resource "ztc_dns_gateway" "example" {
  name                             = "Example DNS Gateway"
  primary_ip                       = "1.1.1.1"
  secondary_ip                     = "2.2.2.2"
  failure_behavior                 = "FAIL_RET_ERR"
}
```

## Argument Reference

### Required

- `name` - (String) The name of the DNS Gateway.

### Optional

- `dns_gateway_type` - (String) Type of the DNS Gateway. Supported value: `EC_DNS_GW`.
- `ec_dns_gateway_options_primary` - (String) Primary DNS gateway option for Edge Connector. Supported values: `LAN_PRI_DNS_AS_PRI`, `LAN_SEC_DNS_AS_SEC`, `WAN_PRI_DNS_AS_PRI`, `WAN_SEC_DNS_AS_SEC`.
- `ec_dns_gateway_options_secondary` - (String) Secondary DNS gateway option for Edge Connector. Supported values: `LAN_PRI_DNS_AS_PRI`, `LAN_SEC_DNS_AS_SEC`, `WAN_PRI_DNS_AS_PRI`, `WAN_SEC_DNS_AS_SEC`.
- `failure_behavior` - (String) Choose what happens if the DNS server is unreachable. Supported values: `FAIL_RET_ERR`, `"FAIL_ALLOW_IGNORE_DNAT"`.
- `primary_ip` - (String) IP address of the primary custom DNS server.
- `secondary_ip` - (String) IP address of the secondary custom DNS server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the DNS Gateway.
- `gateway_id` - (Number) The numeric identifier assigned by the API.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_dns_gateway** can be imported by using `<GATEWAY_ID>` or `<GATEWAY_NAME>` as the import ID.

For example:

```shell
terraform import ztc_dns_gateway.example <gateway_id>
```

or

```shell
terraform import ztc_dns_gateway.example <gateway_name>
```