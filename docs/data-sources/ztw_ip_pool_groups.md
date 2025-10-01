---
subcategory: "Policy Resources"
layout: "zscaler"
page_title: "ZTW: ip_pool_groups"
description: |-
  Official documentation https://help.zscaler.com/cloud-branch-connector/about-ip-pool-groups
  API documentation https://help.zscaler.com/cloud-branch-connector/ip-pool-groups
  Get information about IP Pool Groups.
---

# ztw_ip_pool_groups (Data Source)

* [Official documentation](https://help.zscaler.com/cloud-branch-connector/about-ip-pool-groups)
* [API documentation](https://help.zscaler.com/cloud-branch-connector/ip-pool-groups)

Use the **ztw_ip_pool_groups** data source to get information about IP Pool Groups available in the Zscaler Cloud and Branch Connector Portal. This data source can then be associated with ZTW traffic forwarding rule.

## Example Usage - Retrieve by Name

```hcl
data "ztw_ip_pool_groups" "example" {
    name = "example_ip_pool_group"
}
```

## Example Usage - Retrieve by ID

```hcl
data "ztw_ip_pool_groups" "example" {
    id = 5458452
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IP pool group to be exported.
* `id` - (Optional) The ID of the IP pool group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Required

* `id` - (Number) ID of the IP address group or IP pool.
* `name` - (String) Name of the IP address group or IP pool.

### Optional

* `description` - (String) Description of the IP group or IP pool.
* `ip_addresses` - (List of String) IP addresses included in the IP group or IP pool.
* `creator_context` - (String) Indicates that the IP group or IP pool is created in Cloud & Branch Connector (EC) (only applicable value).
* `is_non_editable` - (Boolean) Indicates whether the group is view-only (true) or editable (false).
* `extranet_ip_pool` - (Boolean) Indicates whether this is an extranet IP pool.
* `is_predefined` - (Boolean) Indicates whether the IP pool group is predefined.