---
subcategory: "Provisioning"
layout: "zscaler"
page_title: "ZTC: provisioning_url"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-provisioning-urls
  API documentation https://help.zscaler.com/cloud-branch-connector/provisioning-urls
  Creates and manages Provisioning URLs.
---

# ztc_provisioning_url (Resource)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/provisioning#/provUrl-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-provisioning-urls)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/provisioning-urls)

Use the **ztc_provisioning_url** resource allows the creation and management of provisioning URLs available in the Zscaler Cloud and Branch Connector Portal.

## Example Usage - Create Provisioning URL

```hcl
resource "ztc_location_template" "this" {
  name = "testAcc_location_template${random_id.test.hex}"
  desc = "Location Template Test"
  template {
    template_prefix           = "testAcc-tf"
    aup_timeout_in_days       = 1
    auth_required             = false
    caution_enabled           = false
    aup_enabled               = true
    ips_control               = true
    ofw_enabled               = true
    xff_forward_enabled       = true
    enforce_bandwidth_control = true
    up_bandwidth              = 10000
    dn_bandwidth              = 10000
  }
}

resource "ztc_provisioning_url" "example" {
  name          = "testAcc_provisioning_url"
  desc          = "Updated Provisioning URL Test"
  prov_url_type = "CLOUD"
  prov_url_data {
    form_factor         = "SMALL"
    cloud_provider_type = "AWS"
    location_template {
      id = ztc_location_template.this.id
    }
  }
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
* `prov_url_type` - (String) Type of the provisioning URL. Supported values: `ONPREM`, `CLOUD`.
* `status` - (String) Status of the provisioning URL.
* `last_mod_time` - (Number) Last modification timestamp.
* `prov_url_data` - (List of Object) Provisioning URL data details.
  * `cloud_provider_type` - (String) Cloud provider type (AWS, AZURE, GCP).
  * `form_factor` - (String) Form factor (SMALL, MEDIUM, LARGE).
  * `location_template` - (List of Object) Location template details. Includes all attributes from the location_template data source.
    * `id` - (Number) Cloud provider identifier.edge_connector_group data source.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZTC configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**ztc_provisioning_url** can be imported by using `<PROV_URL_ID>` or `<PROV_URL_NAME>` as the import ID.

For example:

```shell
terraform import ztc_provisioning_url.example <rule_id>
```

or

```shell
terraform import ztc_provisioning_url.example <rule_name>
```
