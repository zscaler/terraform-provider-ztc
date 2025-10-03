---
subcategory: "Forwarding Gateways"
layout: "zscaler"
page_title: "ZTW: zia_forwarding_gateway"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-traffic-forwarding
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-gateways
  Creates and manages ZIA and Log Control Forwarding Gateways.
---

# ztw_zia_forwarding_gateway (Resource)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-traffic-forwarding)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways)

Use the **ztw_zia_forwarding_gateway** resource allows the creation and management of ZIA and Log Forwarding gateways available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTW traffic forwarding rule.

## Example Usage - ZIA Primary and Secondary Type AUTO

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_gw01" {
  name           = "ZTW_GW01"
  description    = "Example Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ZIA"
}
```

## Example Usage - ZIA Primary and Secondary Type DC

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_gw02" {
  name             = "ZTW_GW02"
  description      = "Example Forwarding Gateway 2"
  fail_closed      = true
  primary_type     = "DC"
  secondary_type   = "DC"
  manual_primary   = "zrh1.svpn.zscalerbeta.net"
  manual_secondary = "syseng.svpn.zscalerbeta.net"
  type             = "ZIA"
}
```

## Example Usage - ZIA Primary and Secondary Type MANUAL_OVERRIDE

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_gw03" {
  name             = "ZTW_GW03"
  description      = "Example Forwarding Gateway 3"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ZIA"
}
```

## Example Usage - Log Forwarding Gateway Primary and Secondary Type AUTO

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_log01" {
  name           = "LOG_FW_GW01"
  description    = "Example Log Forwarding Gateway 1"
  fail_closed    = true
  primary_type   = "AUTO"
  secondary_type = "AUTO"
  type           = "ECSELF"
}
```

## Example Usage - Log Forwarding Gateway Primary and Secondary Type DC

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_log02" {
  name           = "LOG_FW_GW02"
  description    = "Example Log Forwarding Gateway 2"
  fail_closed      = true
  primary_type     = "DC"
  secondary_type   = "DC"
  manual_primary   = "zrh1.svpn.zscalerbeta.net"
  manual_secondary = "syseng.svpn.zscalerbeta.net"
  type             = "ECSELF"
}
```

## Example Usage - Log Forwarding Gateway Primary and Secondary Type MANUAL_OVERRIDE

```hcl
resource "ztw_zia_forwarding_gateway" "ztw_log03" {
  name           = "LOG_FW_GW03"
  description    = "Example Log Forwarding Gateway 3"
  fail_closed      = true
  primary_type     = "MANUAL_OVERRIDE"
  secondary_type   = "MANUAL_OVERRIDE"
  manual_primary   = "1.1.1.1"
  manual_secondary = "2.2.2.2"
  type             = "ECSELF"
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

* `description` - (String) Additional details about the Forwarding Gateway.
* `type` - (String) Type of the gateway. Supported types are `ZIA`. Use `ECSELF`for (Log and Control gateway).
* `fail_closed` - (Boolean) A true value indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable. A false value indicates that traffic must be allowed.
* `manual_primary` - (String) Specifies the primary proxy through which traffic must be forwarded. Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name.
* `manual_secondary` - (String) Specifies the secondary proxy through which traffic must be forwarded. Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name.
* `primary_type` - (String) Type of the primary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through a data center, or manual proxy (IP) that forwards traffic to a specific IP address or domain.
* `secondary_type` - (String) Type of the secondary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through a data center, or manual proxy (IP) that forwards traffic to a specific IP address or domain.

* `subcloud_primary` - (List of Object) If a manual (DC) primary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified data center. This allows for more granular control over which subcloud handles the primary traffic forwarding.
  * `id` - (Number) Identifier that uniquely identifies the subcloud entity.

* `subcloud_secondary` - (List of Object) If a manual (DC) secondary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified data center. This allows for more granular control over which subcloud handles the secondary traffic forwarding.
  * `id` - (Number) Identifier that uniquely identifies the subcloud entity.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTW configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztw_zia_forwarding_gateway** can be imported by using `<GATEWAY_ID>` or `<GATEWAY_NAME>` as the import ID.

For example:

```shell
terraform import ztw_zia_forwarding_gateway.example <rule_id>
```

or

```shell
terraform import ztw_zia_forwarding_gateway.example <rule_name>
```