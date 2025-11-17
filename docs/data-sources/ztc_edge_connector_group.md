---
subcategory: "Cloud Connector Groups"
layout: "zscaler"
page_title: "ZTC: edge_connector_group"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-cloud-connector-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/cloud-branch-connector-groups#/ecgroup-get
  Get information about Cloud and Branch Connector Groups.
---

# ztc_edge_connector_group (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/cloud-branch-connector-groups#/ecgroup-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-cloud-connector-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/cloud-branch-connector-groups#/ecgroup-get)

Use the **ztc_edge_connector_group** data source to get information about Cloud and Branch Connector Groups available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTC traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztc_edge_connector_group" "example" {
    name = "example"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_edge_connector_group" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the edge connector group to be exported.
* `id` - (Optional) The ID of the edge connector group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) The unique identifier of the edge connector group.
* `name` - (String) The name of the edge connector group.

### Optional

* `desc` - (String) Description of the edge connector group.
* `deploy_type` - (String) Deployment type of the edge connector group.
* `status` - (String) Status of the edge connector group.
* `platform` - (String) Platform on which the edge connector is deployed.
* `aws_availability_zone` - (String) AWS availability zone for the edge connector group.
* `azure_availability_zone` - (String) Azure availability zone for the edge connector group.
* `max_ec_count` - (Number) Maximum number of edge connectors in the group.
* `tunnel_mode` - (String) Tunnel mode configuration for the edge connector group.
* `location` - (List of Object) Location associated with the edge connector group.
  * `id` - (Number) Identifier that uniquely identifies the location.
  * `name` - (String) The configured name of the location.
  * `is_name_l10n_tag` - (Boolean) Localization tag indicator.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if the location is deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.
* `prov_template` - (List of Object) Provisioning template associated with the edge connector group.
  * `id` - (Number) Identifier that uniquely identifies the template.
  * `name` - (String) The configured name of the template.
  * `is_name_l10n_tag` - (Boolean) Localization tag indicator.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if the template is deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.
* `ec_vms` - (List of Object) List of edge connector VMs in the group.
  * `id` - (Number) The unique identifier of the EC VM.
  * `name` - (String) The name of the EC VM.
  * `form_factor` - (String) Form factor of the EC VM.
  * `city_geo_id` - (Number) City geographical identifier.
  * `nat_ip` - (String) NAT IP address.
  * `zia_gateway` - (String) ZIA gateway.
  * `zpa_broker` - (String) ZPA broker.
  * `build_version` - (String) Build version.
  * `last_upgrade_time` - (Number) Last upgrade timestamp.
  * `upgrade_status` - (Number) Upgrade status code.
  * `upgrade_start_time` - (Number) Upgrade start timestamp.
  * `upgrade_end_time` - (Number) Upgrade end timestamp.
  * `management_nw` - (List of Object) Management network configuration.
    * `id` - (Number) Network identifier.
    * `ip_start` - (String) Starting IP address of the network range.
    * `ip_end` - (String) Ending IP address of the network range.
    * `netmask` - (String) Network mask.
    * `default_gateway` - (String) Default gateway IP address.
    * `nw_type` - (String) Network type.
    * `dns` - (List of Object) DNS configuration.
      * `id` - (Number) DNS identifier.
      * `ips` - (List of String) DNS server IP addresses.
      * `dns_type` - (String) DNS type.
  * `ec_instances` - (List of Object) Edge connector instances.
    * `ec_instance_type` - (String) Instance type.
    * `out_gw_ip` - (String) Outbound gateway IP.
    * `nat_ip` - (String) NAT IP address.
    * `dns_ip` - (String) DNS IP address.
    * `service_nw` - (List of Object) Service network configuration (same structure as management_nw).
    * `virtual_nw` - (List of Object) Virtual network configuration (same structure as management_nw).