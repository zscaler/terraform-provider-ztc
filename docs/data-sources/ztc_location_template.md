---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZTC: location_template"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-location-templates
  API documentation https://help.zscaler.com/cloud-branch-connector/location-templates
  Get information about Location Templates.
---

# ztc_location_template (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-location-templates)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/location-templates)

Use the **ztc_location_template** data source to get information about location templates available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with with ZTC locations and provisioning URLs.

## Example Usage - Retrieve by Name

```hcl
data "ztc_location_template" "example" {
    name = "example_location_template"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_location_template" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the location template to be exported.
* `id` - (Optional) The ID of the location template resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of Cloud & Branch Connector location template.
* `name` - (String) Name of Cloud & Branch Connector location template.

### Optional

* `desc` - (String) Description of Cloud & Branch Connector location template.
* `editable` - (Boolean) Whether Cloud & Branch Connector location template is editable.
* `last_mod_time` - (Number) Last time Cloud & Branch Connector location template was modified.
* `template` - (List of Object) Template configuration details.
  * `template_prefix` - (String) Prefix of Cloud & Branch Connector location template.
  * `xff_forward_enabled` - (Boolean) Enable XFF Forwarding for a location. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header. Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location.
  * `auth_required` - (Boolean) Indicates whether to enforce authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
  * `caution_enabled` - (Boolean) Indicates whether to enable Caution for this location.
  * `aup_enabled` - (Boolean) Indicates whether to enable Acceptable Use Policy (AUP) for this location.
  * `aup_timeout_in_days` - (Number) Number of days for AUP timeout.
  * `ofw_enabled` - (Boolean) Indicates whether to enable firewall for this location.
  * `ips_control` - (Boolean) Indicates whether to enable IPS for this location.
  * `enforce_bandwidth_control` - (Boolean) Enable to specify the maximum bandwidth limits for download (Mbps) and upload (Mbps).
  * `up_bandwidth` - (Number) Upload bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.
  * `dn_bandwidth` - (Number) Download bandwidth in Kbps. The value 0 implies no Bandwidth Control enforcement.
* `last_mod_uid` - (List of Object) Last modifier information.
  * `id` - (Number) Identifier that uniquely identifies the user.
  * `name` - (String) The configured name of the user.
  * `is_name_l10n_tag` - (Boolean) Localization tag indicator.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.