---
subcategory: "Provisioning"
layout: "zscaler"
page_title: "ZTC: provisioning_url"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-provisioning-urls
  API documentation https://help.zscaler.com/cloud-branch-connector/provisioning-urls
  Get information about Provisioning URLs.
---

# ztc_provisioning_url (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/provisioning#/provUrl-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-provisioning-urls)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/provisioning-urls)

Use the **ztc_provisioning_url** data source to get information about provisioning URLs available in the Zscaler Cloud and Branch Connector Portal. This data source is used to retrieve provisioning information for edge connectors.

## Example Usage - Retrieve by Name

```hcl
data "ztc_provisioning_url" "example" {
    name = "example_provisioning_url"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztc_provisioning_url" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the provisioning URL to be exported.
* `id` - (Optional) The ID of the provisioning URL resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) The unique identifier of the provisioning URL.
* `name` - (String) The name of the provisioning URL.

### Optional

* `desc` - (String) Description of the provisioning URL.
* `prov_url` - (String) The actual provisioning URL.
* `prov_url_type` - (String) Type of the provisioning URL (e.g., ONPREM, CLOUD).
* `status` - (String) Status of the provisioning URL.
* `last_mod_time` - (Number) Last modification timestamp.
* `prov_url_data` - (List of Object) Provisioning URL data details.
  * `zs_cloud_domain` - (String) Zscaler cloud domain.
  * `org_id` - (Number) Organization ID.
  * `config_server` - (String) Configuration server URL.
  * `registration_server` - (String) Registration server URL.
  * `api_server` - (String) API server URL.
  * `pac_server` - (String) PAC server URL.
  * `cloud_provider_type` - (String) Cloud provider type (AWS, AZURE, GCP).
  * `form_factor` - (String) Form factor (SMALL, MEDIUM, LARGE).
  * `hypervisors` - (String) Hypervisor type.
  * `location_template` - (List of Object) Location template details. Includes all attributes from the location_template data source.
  * `cloud_provider` - (List of Object) Cloud provider information.
    * `id` - (Number) Cloud provider identifier.
    * `name` - (String) Cloud provider name.
    * `is_name_l10n_tag` - (Boolean) Localization tag indicator.
    * `extensions` - (Map of String) General purpose field.
    * `deleted` - (Boolean) Indicates if deleted.
    * `external_id` - (String) External identifier.
    * `association_time` - (Number) Association timestamp.
  * `location` - (List of Object) Location information (same structure as cloud_provider).
  * `bc_group` - (List of Object) Branch Connector group information. Includes all attributes from the edge_connector_group data source.
* `used_in_ec_groups` - (List of Object) Edge Connector groups using this provisioning URL.
  * `id` - (Number) Group identifier.
  * `name` - (String) Group name.
  * `is_name_l10n_tag` - (Boolean) Localization tag indicator.
  * `extensions` - (Map of String) General purpose field.
  * `deleted` - (Boolean) Indicates if deleted.
  * `external_id` - (String) External identifier.
  * `association_time` - (Number) Association timestamp.
* `last_mod_uid` - (List of Object) Last modifier information (same structure as used_in_ec_groups).