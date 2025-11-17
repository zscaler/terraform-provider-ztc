---
subcategory: "Forwarding Gateways"
layout: "zscaler"
page_title: "ZTC: zia_forwarding_gateway"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-forwarding-gateways
  API documentation https://help.zscaler.com/cloud-branch-connector/forwarding-gateways
  Get information about Forwarding Gateways.
---

# ztc_zia_forwarding_gateway (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways#/gateways-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-forwarding-gateways)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/forwarding-gateways)

Use the **ztc_zia_forwarding_gateway** data source to get information about forwarding gateways available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTC traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztc_zia_forwarding_gateway" "example" {
    name = "example_forwarding_gateway"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_zia_forwarding_gateway" "example" {
    id = 5458452
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
* `type` - (String) Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
* `fail_closed` - (Boolean) A true value indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable. A false value indicates that traffic must be allowed.
* `manual_primary` - (String) Specifies the primary proxy through which traffic must be forwarded. Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name.
* `manual_secondary` - (String) Specifies the secondary proxy through which traffic must be forwarded. Depending on the proxy forwarding type specified (AUTODC), this field includes a preconfigured data center, or a specified IP address or domain name.
* `primary_type` - (String) Type of the primary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through a data center, or manual proxy (IP) that forwards traffic to a specific IP address or domain.
* `secondary_type` - (String) Type of the secondary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through a data center, or manual proxy (IP) that forwards traffic to a specific IP address or domain.
* `last_modified_time` - (Number) Timestamp when the forwarding gateway was last modified.
* `subcloud_primary` - (List of Object) If a manual (DC) primary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified data center. This allows for more granular control over which subcloud handles the primary traffic forwarding.
  * `id` - (Number) Identifier that uniquely identifies the subcloud entity.
  * `name` - (String) The configured name of the subcloud entity.
  * `is_name_l10n_tag` - (Number) Indicates the external ID. Applicable only when this reference is of an external entity.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if the entity is deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.
* `subcloud_secondary` - (List of Object) If a manual (DC) secondary proxy is used and if the organization has subclouds associated, you can specify a subcloud using this field for the specified data center. This allows for more granular control over which subcloud handles the secondary traffic forwarding.
  * `id` - (Number) Identifier that uniquely identifies the subcloud entity.
  * `name` - (String) The configured name of the subcloud entity.
  * `is_name_l10n_tag` - (Number) Indicates the external ID. Applicable only when this reference is of an external entity.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if the entity is deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.
* `last_modified_by` - (List of Object) User who last modified the forwarding gateway.
  * `id` - (Number) Identifier that uniquely identifies the user.
  * `name` - (String) The configured name of the user.
  * `is_name_l10n_tag` - (Number) Indicates the external ID. Applicable only when this reference is of an external entity.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if the entity is deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.