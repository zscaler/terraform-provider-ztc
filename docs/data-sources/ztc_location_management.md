---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZTC: location_management"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-locations
  API documentation https://help.zscaler.com/cloud-branch-connector/location-management
  Get information about Location Management.
---

# ztc_location_management (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-locations)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/location-management)

Use the **ztc_location_management** data source to get information about locations available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with with ZTC policies and rules.

## Example Usage - Retrieve by Name

```hcl
data "ztc_location_management" "example" {
    name = "example_location"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_location_management" "example" {
  id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the location to be exported.
* `id` - (Optional) The ID of the location resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) The unique identifier of the location.
* `name` - (String) The name of the location.

### Optional

* `description` - (String) Additional information about the location.
* `non_editable` - (Boolean) Indicates whether the location can be edited.
* `parent_id` - (Number) Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID.
* `enforce_bandwidth_control` - (Boolean) Enable to specify the maximum bandwidth limits for download (Mbps) and upload (Mbps).
* `up_bandwidth` - (Number) Upload bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.
* `dn_bandwidth` - (Number) Download bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.
* `country` - (String) Country of the location.
* `state` - (String) State of the location.
* `language` - (String) Language of the location.
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT.
* `auth_required` - (Boolean) Indicates whether to enforce authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
* `xff_forward_enabled` - (Boolean) Enable XFF Forwarding for a location. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header. Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location.
* `ec_location` - (Boolean) Indicates whether this is a Cloud or Branch Connector location (true) or a generic location (false).
* `ofw_enabled` - (Boolean) Indicates whether to enable firewall for this location.
* `ips_control` - (Boolean) Indicates whether to enable IPS for this location.
* `aup_enabled` - (Boolean) Indicates whether to enable Acceptable Use Policy (AUP) for this location.
* `caution_enabled` - (Boolean) Indicates whether to enable Caution for this location.
* `exclude_from_dynamic_groups` - (Boolean) Indicates whether to exclude this location from dynamic location groups when created.
* `exclude_from_manual_groups` - (Boolean) Indicates whether to exclude this location from manual location groups when created.
* `public_cloud_account_id` - (List of Object) AWS/Azure subscription ID associated with this location.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `is_name_l10n_tag` - (Boolean) Indicates the external ID. Applicable only when this reference is of an external entity.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) General purpose field.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.