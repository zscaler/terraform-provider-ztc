---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZTC: location_template"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-location-templates
  API documentation https://help.zscaler.com/cloud-branch-connector/location-templates
  Creates and manages Location Templates.
---

# ztc_location_template (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/location-management#/locationTemplate-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-location-templates)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/location-templates)

Use the **ztc_location_template**  resource allows the creation and management of location templates available in the Zscaler Cloud and Branch Connector Portal. This resource can then be associated with ZTC locations and provisioning URLs.

## Example Usage - Create Location Template

```hcl
resource "ztc_location_template" "this" {
  name = "testAcc_location_template"
  desc = "Location Template Test"
  template {
    template_prefix           = "testAcc-tf"
    aup_timeout_in_days       = 0
    auth_required             = false
    ips_control               = true
    ofw_enabled               = true
    xff_forward_enabled       = true
    enforce_bandwidth_control = true
    up_bandwidth              = 10
    dn_bandwidth              = 10
  }
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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_location_template** can be imported by using `<TEMPLATE_ID>` or `<TEMPLATE_NAME>` as the import ID.

For example:

```shell
terraform import ztc_location_template.example <rule_id>
```

or

```shell
terraform import ztc_location_template.example <rule_name>
```
