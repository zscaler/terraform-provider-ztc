---
subcategory: "Forwarding Gateways"
layout: "zscaler"
page_title: "ZTC: dns_forwarding_gateway"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-gateways
  Get information about DNS Forwarding Gateways.
---

# ztc_dns_forwarding_gateway (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways)

Use the **ztc_dns_forwarding_gateway** data source to get information about dns forwarding gateways available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTC traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztc_dns_forwarding_gateway" "example" {
    name = "example_forwarding_gateway"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_dns_forwarding_gateway" "example" {
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

* `primary_ip` - (String) IP address of the primary custom DNS server.
* `secondary_ip` - (String) IP address of the secondary custom DNS server.
* `ec_dns_gateway_options_primary` - (String) IP address of the primary LAN DNS Server. Supported Values: `LAN_PRI_DNS_AS_PRI`, and `LAN_SEC_DNS_AS_SEC`
* `ec_dns_gateway_options_secondary` - (String) IP address of the secondary LAN DNS Server. Supported Values: `LAN_PRI_DNS_AS_PRI`, and `LAN_SEC_DNS_AS_SEC`
* `failure_behavior` - (String) Choose what happens if the DNS server is unreachable. Supported Values: `FAIL_RET_ERR`, and `FAIL_ALLOW_IGNORE_DNAT`
