---
subcategory: "DNS Gateway"
layout: "zscaler"
page_title: "ZTC: ztc_dns_gateway (Data Source)"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway
  API documentation https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/dns-gateway/ec-dns-gateway-z-resource-add-dns-gateway
  Get information about DNS Forwarding Gateways.
---

# ztc_dns_gateway (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways#/gateways-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/configuring-dns-gateway)
* [API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zcloudconnector/dns-gateway/ec-dns-gateway-z-resource-add-dns-gateway)

Use the **ztc_dns_gateway** data source to get information about a DNS Gateway configuration in the Zscaler Zero Trust Cloud (ZTC) platform. This data source can be used to reference DNS Gateway resources in other Terraform resources.

## Example Usage

```hcl
# Retrieve by name
data "ztc_dns_gateway" "example" {
  name = "Example DNS Gateway"
}

# Retrieve by ID
data "ztc_dns_gateway" "example_by_id" {
  id = 12345
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Optional, Number) The unique identifier for the DNS Gateway.
- `name` - (Optional, String) The name of the DNS Gateway.

~> **NOTE:** One of `id` or `name` must be provided.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `dns_gateway_type` - (String) Type of the DNS Gateway.
- `ec_dns_gateway_options_primary` - (String) Primary DNS gateway option for Edge Connector.
- `ec_dns_gateway_options_secondary` - (String) Secondary DNS gateway option for Edge Connector.
- `failure_behavior` - (String) Defines what happens if the DNS server is unreachable.
- `primary_ip` - (String) IP address of the primary custom DNS server.
- `secondary_ip` - (String) IP address of the secondary custom DNS server.
- `last_modified_time` - (Number) Timestamp when the DNS Gateway was last modified.
- `last_modified_by` - (Set) Details of the user who last modified the DNS Gateway.
  - `id` - (Number) Identifier of the user.
  - `name` - (String) Name of the user.
