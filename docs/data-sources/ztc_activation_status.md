---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZTC: activation_status"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-activation
  API documentation https://help.zscaler.com/cloud-branch-connector/activation-status
  Get information about Activation Status.
---

# ztc_activation_status (Data Source)

[![General Availability](https://img.shields.io/badge/Lifecycle%20Stage-General%20Availability-%2345c6e8)](https://help.zscaler.com/cloud-branch-connector/activation#/ecAdminActivateStatus-get)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-activation)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/activation-status)

Use the **ztc_activation_status** data source to get information about the current activation status in the Zscaler Cloud and Branch Connector Portal. This data source provides details about organization and admin activation status.

## Example Usage

```hcl
data "ztc_activation_status" "this" {}

output "org_status" {
  value = data.ztc_activation_status.this.org_edit_status
}
```

## Argument Reference

No arguments are required for this data source.

## Attribute Reference

The following attributes are exported:

* `org_edit_status` - (String) Organization policy edit status.
* `org_last_activate_status` - (String) Organization policy last activation status.
* `admin_activate_status` - (String) Admin activation status.
* `admin_status_map` - (Map of String) Admin status details.